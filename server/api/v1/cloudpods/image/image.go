package image

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
	"net/http"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/appsrv"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules/image"
	"yunion.io/x/pkg/util/stringutils"
	"yunion.io/x/pkg/utils"
)

type ImageApi struct {
}

type getBackupImageReq struct {
	PageIndex     int64  `json:"pageIndex"`
	PageSize      int64  `json:"pageSize"`
	Description   string `json:"description"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	OsArch        string `json:"osArch"`
	Distributions string `json:"distributions"`
	SizeDesc      string `json:"sizeDesc"`
	TimeDesc      string `json:"timeDesc"`
}
type getImageByNameReq struct {
	Name string `json:"name"`
}
type getBaseImageReq struct {
	ID string `json:"id"`
}

// 恢复回收镜像
func (p *ImageApi) RecoverRecycleImage(c *gin.Context) {
	var r getBaseImageReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := image.Images.PerformAction(s, r.ID, "cancel-delete", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "清除成功", c)
}

// 清除回收镜像
func (p *ImageApi) ClearRecycleImage(c *gin.Context) {
	var r getBaseImageReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	params.Set("override_pending_delete", jsonutils.NewBool(true))
	result, err := image.Images.DeleteWithParam(s, r.ID, params, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "清除成功", c)
}

// 获取回收站列表
func (p *ImageApi) GetRecycleImageList(c *gin.Context) {
	var r getBackupImageReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	filters := make([]string, 0)
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("pending_delete", jsonutils.NewBool(true))
	params.Set("is_guest_image", jsonutils.NewBool(false))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.TimeDesc != "" {
		params.Set("order_by", jsonutils.NewString("auto_delete_at"))
		params.Set("order", jsonutils.NewString(r.TimeDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := image.Images.List(s, params)

	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 获取镜像列表
func (p *ImageApi) GetImageList(c *gin.Context) {
	var r getBackupImageReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	filters := make([]string, 0)
	params := jsonutils.NewDict()
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.Description != "" {
		filters = append(filters, "description.contains('"+r.Description+"')")
	}
	if r.OsArch != "" {
		params.Set("os_arch", jsonutils.NewString(r.OsArch))
	}
	if r.Distributions != "" {
		params.Set("distributions", jsonutils.NewString(r.Distributions))
	}
	if r.SizeDesc != "" {
		params.Set("order_by", jsonutils.NewString("size"))
		params.Set("order", jsonutils.NewString(r.SizeDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := image.Images.List(s, params)

	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 获取单个镜像
func (p *ImageApi) GetImageByName(c *gin.Context) {
	var r getImageByNameReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("scope", jsonutils.NewString("system"))
	result, err := image.Images.GetByName(s, r.Name, params)
	if err != nil {
		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取失败", c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

type ImageUploadReq struct {
	Name       string                `form:"name" json:"name"`
	OsArch     string                `form:"osArch" json:"osArch"`
	OsType     string                `form:"osType" json:"osType"`
	ImageFile  *multipart.FileHeader `form:"image" json:"image"`
	UploadType string                `form:"uploadType" json:"uploadType"` // url file
	Url        string                `form:"url" json:"url"`
}

//func (*ImageApi)

func (*ImageApi) BetterUpImage(c *gin.Context) {
	params, query, _ := appsrv.FetchEnv(context.Background(), c.Writer, c.Request)
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// input params
	imageId, ok := params["<image_id>"]
	if !ok || len(imageId) == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}
	// 是否直接下载
	format, _ := query.GetString("format")
	err = imageDownloadValidateStatus(s, imageId, format)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	direct, _ := query.Bool("direct")
	if !direct {
		ret, err := imageDownloadUrl(imageId, format)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		appsrv.SendJSON(c.Writer, ret)
		return
	}

	// // 直接下载镜像
	// imageDownload(ctx, w, s, imageId, format)
}

func imageDownloadUrl(id string, format string) (jsonutils.JSONObject, error) {
	// 加密下载url
	expired := time.Now().Add(24 * time.Hour)
	imageInfo := jsonutils.NewDict()
	imageInfo.Set("id", jsonutils.NewString(id))
	imageInfo.Set("format", jsonutils.NewString(format))
	imageInfo.Set("expired", jsonutils.NewInt(expired.Unix()))
	token, err := utils.EncryptAESBase64Url(stringutils.UUID4(), imageInfo.String())
	if err != nil {
		return nil, err
	}
	ret := jsonutils.NewDict()
	ret.Set("signature", jsonutils.NewString(token))
	return ret, nil
}

func imageDownloadValidateStatus(s *mcclient.ClientSession, id string, format string) error {
	var imageInfo jsonutils.JSONObject
	var err error
	if len(format) == 0 {
		image.Images.Get(s, id, nil)
		imageInfo, err = image.Images.Get(s, id, nil)
		return err
	} else {
		resp, e := image.Images.GetSpecific(s, id, "subformats", nil)
		if e != nil {
			return errors.Wrap(e, "images.get.subformats")
		}

		images, _ := resp.(*jsonutils.JSONArray).GetArray()
		for i := range images {
			if f, _ := images[i].GetString("format"); f == format {
				imageInfo = images[i]
			}
		}
	}

	status, _ := imageInfo.GetString("status")
	if status != "active" {
		return errors.New("image is not in status 'active'")
	}

	return nil
}

// 上传镜像
func (p *ImageApi) UploadImage(c *gin.Context) {
	var req ImageUploadReq
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("请求已接收，正在处理上传...", c)

	go func() {
		buff, err := req.ImageFile.Open()
		if req.UploadType == "file" {
			if err != nil {
				response.FailWithMessage("上传失败", c)
				return
			}
		}
		s, err := global.GetCloudClientSession()
		if err != nil {
			response.FailWithMessage("上传失败", c)
			return
		}
		imageParams := jsonutils.NewDict()
		prop := jsonutils.NewDict()
		prop.Set("os_type", jsonutils.NewString(req.OsType))
		prop.Set("os_version", jsonutils.NewString(""))
		prop.Set("os_arch", jsonutils.NewString(req.OsArch))
		imageParams.Set("name", jsonutils.NewString(req.Name))
		imageParams.Set("domain_id", jsonutils.NewString("default"))
		imageParams.Set("os_arch", jsonutils.NewString(req.OsArch))
		imageParams.Set("os_version", jsonutils.NewString(""))
		imageParams.Set("os_type", jsonutils.NewString(req.OsType))
		imageParams.Set("project_id", jsonutils.NewString(s.GetProjectId()))

		if req.UploadType == "file" {
			imageParams.Set("image_size", jsonutils.NewInt(req.ImageFile.Size))
			_, err = image.Images.Upload(s, imageParams, buff, req.ImageFile.Size)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
		} else {
			imageParams.Set("copy_from", jsonutils.NewString(req.Url))
			imageParams.Set("properties", prop)
			_, err = image.Images.Create(s, imageParams)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
		}
	}()
}

type DeleteImageCloudObj struct {
	ImageId       string `json:"id"`
	DisableDelete bool   `json:"disable_delete"`
	Protected     bool   `json:"protected"`
}
type DeleteImageByIds struct {
	Ids           []string `json:"ids"`
	DisableDelete bool     `json:"disable_delete"`
	Protected     bool     `json:"protected"`
}

// 删除镜像
func (p *ImageApi) DeleteImage(c *gin.Context) {
	var req DeleteImageCloudObj
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	var cloudReq = DeleteImageCloudObj{
		ImageId:       req.ImageId,
		DisableDelete: false,
		Protected:     false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	//取消镜像的删除保护
	_, err = image.Images.Update(s, req.ImageId, jsonutils.Marshal(cloudReq))
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	//进行删除
	result, err := image.Images.Delete(s, req.ImageId, nil)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	} else {
		response.OkWithData(result, c)
	}
}

// 批量删除镜像
func (p *ImageApi) DeleteImageByIds(c *gin.Context) {
	var req DeleteImageByIds
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	for _, id := range req.Ids {
		var cloudReq = DeleteImageCloudObj{
			ImageId:       id,
			DisableDelete: false,
			Protected:     false,
		}
		//取消镜像的删除保护
		_, err = image.Images.Update(s, id, jsonutils.Marshal(cloudReq))
		if err != nil {
			response.FailWithMessage("删除失败", c)
			return
		}
		//进行删除
		_, err := image.Images.Delete(s, id, nil)
		if err != nil {
			response.FailWithMessage("删除失败", c)
			return
		}
	}
	response.OkWithData("删除成功", c)
}

type UpdateImageReq struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	OsArch     string          `json:"os_arch"`
	Properties ImageProperties `json:"properties"`
	MinDisk    string          `json:"min-disk"`
}

// ImageProperties 定义镜像属性的结构体
type ImageProperties struct {
	DiskDriver  string `json:"disk_driver"`
	NetDriver   string `json:"net_driver"`
	OsArch      string `json:"os_arch"`
	UefiSupport string `json:"uefi_support"`
	VdiProtocol string `json:"vdi_protocol"`
}

// 镜像信息修改
func (*ImageApi) UpdateImageInfo(c *gin.Context) {
	var req UpdateImageReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 获取会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	resulr, err := image.Images.Update(s, req.Id, jsonutils.Marshal(req))
	if err != nil {
		response.FailWithMessage("更新失败", c)
		return
	} else {
		response.OkWithData(resulr, c)
	}
}

type DownloadImageReq struct {
	ImageId string `json:"image_id" form:"image_id"`
	Format  string `json:"format" form:"format"`
}

// 下载镜像
func (*ImageApi) DownloadImage(c *gin.Context) {
	var req DownloadImageReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 获取会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("下载失败", c)
		return
	}
	//进行下载
	meta, body, size, err := image.Images.Download2(s, req.ImageId, req.Format, false)
	if err != nil {
		response.FailWithMessage("下载失败", c)
		return
	}

	name, _ := meta.GetString("name")
	if len(name) == 0 {
		format, _ := meta.GetString("disk_format")
		name = "os_image"
		if len(format) > 0 {
			name += "." + format
		}
	}

	hdr := http.Header{}
	hdr.Set("Content-Type", "application/octet-stream")
	hdr.Set("Content-Disposition", fmt.Sprintf("Attachment; filename=%s", name))
	appsrv.SendStream(c.Writer, false, hdr, body, size)
}

type imageTagReq struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 镜像标签修改
func (*ImageApi) UpdateImageTag(c *gin.Context) {
	var req imageTagReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 获取会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	resulr, err := image.Images.Update(s, req.ID, jsonutils.Marshal(req))
	if err != nil {
		response.FailWithMessage("更新失败", c)
		return
	} else {
		response.OkWithData(resulr, c)
	}
}
