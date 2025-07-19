package backup

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/backup"
	"ygpt/server/service"
	"ygpt/server/utils"
	"ygpt/server/utils/minio"
)

type ImageApi struct {
}

//	type getBackupImageReq struct {
//		PageIndex int64 `json:"pageIndex"`
//		PageSize  int64 `json:"pageSize"`
//	}
//
// // 查询套餐列表
//
//	func (p *ImageApi) GetBackupImageList(c *gin.Context) {
//		var r getBackupImageReq
//		err := c.ShouldBindJSON(&r)
//		if err != nil {
//			response.FailWithMessage(err, c)
//			return
//		}
//		s, err := global.GetCloudClientSession()
//		if err != nil {
//			response.FailWithMessage(err, c)
//		}
//		params := jsonutils.NewDict()
//		params.Set("limit", jsonutils.NewInt(r.PageSize))
//		params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
//		result, _ := image.Images.List(s, params)
//
//		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
//	}
//
//	type backupImageBaseReq struct {
//		ID string `json:"ID"`
//	}
//
//
//	func (p *ImageApi) AddBackupImage(c *gin.Context) {
//		var r backupImageBaseReq
//		err := c.ShouldBindJSON(&r)
//		if err != nil {
//			response.FailWithMessage(err, c)
//			return
//		}
//		s, err := global.GetCloudClientSession()
//		if err != nil {
//			response.FailWithMessage(err, c)
//		}
//		params := jsonutils.NewDict()
//		params.Set("id", jsonutils.NewString(r.ID))
//
//		result, err := image.Images.Create(s, params)
//		if err != nil {
//			response.FailWithMessage("关机失败", c)
//		}
//		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "关机成功", c)
//	}
//
//	type backupImageDownloadReq struct {
//		Format string `json:"format"`
//		ID     string `json:"ID"`
//	}
//
//
//	func (p *ImageApi) DownloadBackupImage(c *gin.Context) {
//		var r backupImageDownloadReq
//		err := c.ShouldBindJSON(&r)
//		if err != nil {
//			response.FailWithMessage(err, c)
//			return
//		}
//		s, err := global.GetCloudClientSession()
//		if err != nil {
//			response.FailWithMessage(err, c)
//		}
//		params := jsonutils.NewDict()
//		params.Set("id", jsonutils.NewString(r.ID))
//
//		result, err := image.Images.List(s, params)
//		if err != nil {
//			response.FailWithMessage("关机失败", c)
//		}
//		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "关机成功", c)
//	}
//
// 添加
func (p *ImageApi) AddBackImage(c *gin.Context) {
	var r backup.BackupImage
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if r.GenerateWay == 1 {
		r.SourceID = int(utils.GetUserID(c))
	}
	serviceReturn, err := ImageService.AddImageService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

var systemConfigService = service.ServiceGroupApp.YunGuanServiceGroup.SystemConfigService

// 添加
func (p *ImageApi) AddBackImageFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	id, _ := strconv.Atoi(c.PostForm("id"))
	data := backup.BackupImage{
		GVA_MODEL: global.GVA_MODEL{
			ID: uint(id),
		},
		Size: uint32(file.Size),
	}
	//获取文件名和文件后缀
	fileExt := filepath.Ext(file.Filename)
	image, err := ImageService.GetImageByIdService(id)

	fileName := image.Name + "-" + c.PostForm("id") + fileExt
	data.FileName = fileName
	err = ImageService.UpdateImageService(data)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败", c)
		return
	}

	// 上传文件到minio存储桶中
	{
		config, err := systemConfigService.GetSystemConfigByKeyService("image")
		var imageConfig global.MinioData
		json.Unmarshal([]byte(config.Value), &imageConfig)
		fmt.Println(imageConfig)

		if file.Filename != "" {
			err = minio.UploadFile(imageConfig.Endpoint, imageConfig.AccessKey, imageConfig.SecretKey, imageConfig.BucketName, fileName, file)
			if err != nil {
				response.FailWithMessage(err.Error()+"上传文件失败", c)
				return
			}
		}
		if err != nil {
			global.GVA_LOG.Error("添加失败!", zap.Error(err))
			response.FailWithMessage("添加失败", c)
			return
		}
	}

	response.OkWithMessage("添加成功", c)
}

// 修改
func (p *ImageApi) UpdateBackImage(c *gin.Context) {
	var r backup.BackupImage
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = ImageService.UpdateImageService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix(), "id": r.ID}, "修改成功", c)
}

// 查询
func (p *ImageApi) QueryBackImage(c *gin.Context) {
	var pageInfo backup.GetImageListReq
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if pageInfo.SourceType == 2 {
		pageInfo.SourceID = int(utils.GetUserID(c))
	}
	list, total, err := ImageService.GetImageListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 查询
func (p *ImageApi) QueryBackImageAll(c *gin.Context) {
	list, total, err := ImageService.GetImageAllService()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}

// 删除
func (p *ImageApi) DeleteBackImage(c *gin.Context) {
	var r backup.BackupImage
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//删除minio中的文件
	{
		image, err := ImageService.GetImageByIdService(int(r.ID))
		config, err := systemConfigService.GetSystemConfigByKeyService("image")
		var imageConfig global.MinioData
		json.Unmarshal([]byte(config.Value), &imageConfig)

		fileName := image.FileName
		//删除minio存储桶中的文件
		err = minio.RemoveFile(imageConfig.Endpoint, imageConfig.AccessKey, imageConfig.SecretKey, imageConfig.BucketName, fileName)
		if err != nil {
			response.FailWithMessage(err.Error()+"删除文件失败", c)
			return
		}
	}

	err = ImageService.DeleteImageService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
func (p *ImageApi) DownloadBackImage(ctx *gin.Context) {
	imageID, _ := strconv.Atoi(ctx.Query("id"))
	image, err := ImageService.GetImageByIdService(imageID)
	config, err := systemConfigService.GetSystemConfigByKeyService("image")
	var imageConfig global.MinioData
	json.Unmarshal([]byte(config.Value), &imageConfig)
	fileName := image.FileName

	reader, err := minio.DownloadFile(imageConfig.Endpoint, imageConfig.AccessKey, imageConfig.SecretKey, imageConfig.BucketName, fileName)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	tempFile, err := os.CreateTemp(os.TempDir(), "download-")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	defer os.Remove(tempFile.Name())
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	ctx.File(tempFile.Name())
}
