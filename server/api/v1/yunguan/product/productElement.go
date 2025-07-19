package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/system"
	wsService "ygpt/server/service/compute/websocket/service"
	"ygpt/server/utils"
	"ygpt/server/utils/minio"
)

type ProductElemApi struct {
}

// AddInfo 添加
func (i *ProductElemApi) AddInfo(c *gin.Context) {
	filepathUrl := c.PostForm("filePath")
	fileName := c.PostForm("fileName")
	fileSizeStr := c.PostForm("fileSize")      // 获取 fileSize 字符串
	fileSize, err := strconv.Atoi(fileSizeStr) // 转换为 int 类型
	name := c.PostForm("name")
	description := c.PostForm("description")
	startTime := c.PostForm("startTime")
	endTime := c.PostForm("endTime")
	openType, _ := strconv.Atoi(c.PostForm("openType"))
	elemType, _ := strconv.Atoi(c.PostForm("type"))
	source, _ := strconv.Atoi(c.PostForm("source"))
	usage := c.PostFormArray("usage[]")
	fmt.Println(usage)
	priceStr := c.PostForm("price")
	price, _ := strconv.ParseFloat(priceStr, 64)

	//md5 := strings.TrimSuffix(file.Filename, ext)
	//md5 = utils.MD5V([]byte(md5))
	data := product.ProductElementInfo{
		ElementInfoID: utils.GenerateUUID(),
		Name:          name,
		Description:   description,
		Usage:         usage,
		UploadTime:    time.Now(),
		UploadBy:      utils.GetUserID(c),
		Type:          elemType,
		OpenType:      openType,
		Source:        source,
		MD5:           "",
		FilePath:      filepathUrl,
		FileSize:      fileSize,
		FileName:      fileName,
		Price:         price,
		StartTime:     utils.TimeStringToGoTime(startTime),
		EndTime:       utils.TimeStringToGoTime(endTime),
		Endpoint:      global.GVA_CONFIG.Minio.EndPoint,
		AccessKey:     global.GVA_CONFIG.Minio.AccessKey,
		SecretKey:     global.GVA_CONFIG.Minio.SecretKey,
		BucketName:    global.GVA_CONFIG.Minio.BucketName,
	}
	serviceReturn, err := productElemService.AddElemInfoService(data)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功，文件上传中", c)
}

// UpdateInfo 修改
func (i *ProductElemApi) UpdateInfo(c *gin.Context) {
	var info product.ProductElementInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productElemService.UpdateElemInfoService(info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// UploadFile
func (i *ProductElemApi) UploadFile(c *gin.Context) {

	//noSave标志是否在框架自带的数据库保存记录，不需要
	//noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	} // 创建文件 defer 关闭
	endpoint := global.GVA_CONFIG.Minio.EndPoint
	accessKey := global.GVA_CONFIG.Minio.AccessKey
	secretKey := global.GVA_CONFIG.Minio.SecretKey
	bucketName := global.GVA_CONFIG.Minio.BucketName

	//生成唯一文件名
	ext := filepath.Ext(header.Filename)
	this := utils.GenerateUUID()
	onlyName := strings.TrimSuffix(header.Filename, ext)              // 获取去掉扩展名的文件名
	fileName := fmt.Sprintf("%s_%s%s%s", onlyName, "AI产品", this, ext) // 组合新的文件名
	go func() {
		err = minio.UploadFile(endpoint, accessKey, secretKey, bucketName, fileName, header)
		if err != nil {
			global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
			return
		}
	}()
	file := multipart.FileHeader{
		Filename: fileName,
		Size:     header.Size,
	}
	response.OkWithDetailed(gin.H{"file": file}, "上传成功", c)
}

// 与算力圈同步
func (i *ProductElemApi) SyncInfo(c *gin.Context) {
	var info product.ProductElementInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		global.GVA_LOG.Error("同步失败!"+err.Error(), zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	err = wsService.SuanLiServiceGroupApp.ElementRelease(int(info.ID))
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败"+err.Error(), c)
		return
	}
	info.Status = 2
	info.IsSynced = 2
	err = productElemService.UpdateElemInfoService(info)
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "同步成功", c)
}

// UpdateFileInfo 修改
func (i *ProductElemApi) UpdateFileInfo(c *gin.Context) {
	file, err := c.FormFile("file")

	// //删除minio存储桶中的原文件
	id, _ := strconv.Atoi(c.PostForm("id"))
	info, _ := productElemService.QueryInfoById(id)
	err = minio.RemoveFile(info.Endpoint, info.AccessKey, info.SecretKey, info.BucketName, info.FileName)
	if err != nil {
		response.FailWithMessage(err.Error()+"删除文件失败", c)
		return
	}

	//上传新文件到minio存储桶中

	typeName := ""
	switch info.Type {
	case 1:
		typeName = "模型"
		break
	case 2:
		typeName = "算法"
		break
	case 3:
		typeName = "数据集"
		break
	}
	this := utils.GenerateUUID()
	fileName := file.Filename + "," + typeName + "-" + this
	err = minio.UploadFile(info.Endpoint, info.AccessKey, info.SecretKey, info.BucketName, fileName, file)
	if err != nil {
		response.FailWithMessage(err.Error()+"上传文件失败", c)
		return
	}

	ext := filepath.Ext(file.Filename)
	md5 := strings.TrimSuffix(file.Filename, ext)
	md5 = utils.MD5V([]byte(md5))
	data := product.ProductElementInfo{
		GVA_MODEL: global.GVA_MODEL{
			ID: uint(id),
		},
		MD5:        md5,
		FileName:   fileName,
		FileSize:   int(file.Size),
		UploadTime: time.Now(),
	}
	err = productElemService.UpdateElemInfoService(data)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// QueryInfo 查询节点的元素
func (i *ProductElemApi) QueryInfo(c *gin.Context) {

	var infoReq product.GetElemInfoReq
	err := c.ShouldBindJSON(&infoReq)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = utils.Verify(infoReq, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := productElemService.QueryElemInfoService(infoReq)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     infoReq.PageIndex,
		PageSize: infoReq.PageSize,
	}, "获取成功", c)
}

// 获取租户的元素
func (i *ProductElemApi) QueryUserInfo(c *gin.Context) {

	var infoReq product.GetElemInfoReq
	err := c.ShouldBindJSON(&infoReq)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = utils.Verify(infoReq, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	infoReq.UserID = utils.GetUserID(c)
	list, total, err := productElemService.QueryElemInfoService(infoReq)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     infoReq.PageIndex,
		PageSize: infoReq.PageSize,
	}, "获取成功", c)
}

// 获取共享的元素
func (i *ProductElemApi) QueryShareInfo(c *gin.Context) {

	var infoReq product.GetElemInfoReq
	err := c.ShouldBindJSON(&infoReq)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = utils.Verify(infoReq, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := productElemService.QueryShareElemInfoService(infoReq)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     infoReq.PageIndex,
		PageSize: infoReq.PageSize,
	}, "获取成功", c)
}

// 获取共享的元素
func (i *ProductElemApi) QueryAllShareInfo(c *gin.Context) {

	list, total, err := productElemService.QueryAllShareElemInfoService()
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

// DeleteInfo 删除
func (i *ProductElemApi) DeleteInfo(c *gin.Context) {
	var info product.ProductElementInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if info.IsSynced == 2 {
		err = wsService.SuanLiServiceGroupApp.SyncProductStatusChangeToPlatform(int(info.ID), 2, 5)
		if err != nil {
			global.GVA_LOG.Error("删除失败!无法同步到算力平台", zap.Error(err))
			response.FailWithMessage("删除失败!无法同步到算力平台", c)
			return
		}
	}
	//删除minio存储桶中的文件
	file, _ := productElemService.QueryInfoById(int(info.ID))
	err = minio.RemoveFile(file.Endpoint, file.AccessKey, file.SecretKey, file.BucketName, file.FileName)
	if err != nil {
		global.GVA_LOG.Error("从minio中删除文件失败!", zap.Error(err))
	}
	err = productElemService.DeleteElemInfoService(info.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

// GetAllAvailableModel 获取所有上架的模型元素
func (i *ProductElemApi) GetAllAvailableModel(c *gin.Context) {
	var modelList []product.ProductElementInfo
	err := global.GVA_DB.Model(&product.ProductElementInfo{}).Where("audit_status = 1 and type = 1 and is_synced = 2").Find(&modelList).Error
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithData(gin.H{"list": modelList}, c)
}

// GetAllAvailableDataSet 获取所有上架的数据集元素
func (i *ProductElemApi) GetAllAvailableDataSet(c *gin.Context) {
	var modelList []product.ProductElementInfo
	err := global.GVA_DB.Model(&product.ProductElementInfo{}).Where("audit_status = 1 and type = 2 and is_synced = 2").Find(&modelList).Error
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithData(gin.H{"list": modelList}, c)
}
func (i *ProductElemApi) CreateProduct(c *gin.Context) {
	var req product.ProductCreateV2Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误!", c)
		return
	}
	//代表是数据集和模型
	var pro product.ProductElem
	err = copier.Copy(&pro, &req)
	pro.CreateTime = time.Now()
	pro.CreateBy = utils.GetUserName(c)
	pro.CreateUserId = utils.GetUserIDtoInt(c)
	pro.UserId = utils.GetUserIDtoI64(c)
	pro.Parameters = req.Parameter
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Create(&pro).Error; err != nil {
			global.GVA_LOG.Error("发布产品错误!", zap.Error(err))
			return err
		}
		//上传文件到minio
		if req.FileName != nil && *req.FileName != "" {
			var productFile product.ProductFile
			productFile.CreateBy = utils.GetUserName(c)
			productFile.CreateTime = utils.Pointer(time.Now())
			productFile.CreateUserId = utils.GetUserIDtoIntPointer(c)
			productFile.IsDeleted = 0
			productFile.Name = *req.FileName
			productFile.Description = req.Description
			productFile.Size = utils.Pointer(req.Size)
			if req.Type == 5 {
				productFile.Type = utils.Pointer(3)
			} else if req.Type == 6 {
				productFile.Type = utils.Pointer(4)
			}
			//productFile.Key = *req.Url
			productFile.ProductId = utils.Pointer(int(pro.Id))
			if err = global.GVA_DB.Model(&product.ProductFile{}).Create(&productFile).Error; err != nil {
				global.GVA_LOG.Error("发布产品错误!", zap.Error(err))
				return err
			}
			//提供了下载链接
		} else if req.Url != nil && *req.Url != "" {
			var productFile product.ProductFile
			productFile.CreateBy = utils.GetUserName(c)
			productFile.CreateTime = utils.Pointer(time.Now())
			productFile.CreateUserId = utils.GetUserIDtoIntPointer(c)
			productFile.IsDeleted = 0
			productFile.Name = ""
			//如果只有链接，没有文件名，说明用户只提供了下载链接的文件形式
			productFile.Url = *req.Url
			productFile.Description = req.Description
			productFile.Size = utils.Pointer(int64(0))
			if req.Type == 5 {
				productFile.Type = utils.Pointer(3)
			} else if req.Type == 6 {
				productFile.Type = utils.Pointer(4)
			}
			//productFile.Key = *req.Url
			productFile.ProductId = utils.Pointer(int(pro.Id))
			if err = global.GVA_DB.Model(&product.ProductFile{}).Create(&productFile).Error; err != nil {
				global.GVA_LOG.Error("发布产品错误!", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		if err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Create(&pro).Error; err != nil {
			global.GVA_LOG.Error("发布产品错误!", zap.Error(err))
			response.FailWithMessage("发布产品错误!", c)
			return
		}
	}
	//返回productID 提醒用户上传模型和数据集
	response.OkWithDetailed(map[string]interface{}{
		"productId": pro.Id,
	}, "发布产品成功", c)
	return
}
func (i *ProductElemApi) UpdateProduct(c *gin.Context) {
	var req product.ProductElem
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	fmt.Printf("%+v\n", req)
	req.UpdateBy = utils.GetUserName(c)
	req.UpdateTime = time.Now()
	req.UpdateUserId = utils.GetUserIDtoInt(c)
	if req.LeaseStartTime.IsZero() || req.LeaseEndTime.IsZero() {
		// 默认信息开始时间是此刻
		req.LeaseStartTime = time.Now()
		// 默认信息有效期是一月
		req.LeaseEndTime = time.Now().AddDate(0, 1, 0)
	}
	if req.Price == 0 {
		global.GVA_DB.Debug().Model(&product.ProductElem{}).Select("price").Where("id = ?", req.Id).Updates(&product.ProductElem{Price: 0})
	}
	err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Where("id = ?", req.Id).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	if req.AuditStatus == 0 {
		global.GVA_DB.Debug().Model(&product.ProductElem{}).Select("audit_status").Where("id = ?", req.Id).Updates(&product.ProductElem{AuditStatus: 0})
	}
	//产品更新后要重新同步
	err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Where("id = ?", req.Id).Updates(map[string]interface{}{
		"IsSynced": 0,
	}).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (i *ProductElemApi) GetEleListInfo(c *gin.Context) {
	var req product.ProductSearch3
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// todo 优化项 过于复杂的sql建议下沉
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	var eleTypeList []int
	if req.EleType != nil {
		global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ?", req.EleType).Pluck("id", &eleTypeList)
	}
	db := global.GVA_DB.Model(&product.ProductElem{}).Where("is_deleted  =  0 and type = 2")
	{
		//这一坨请求参数处理
		if req.UserId != nil && *req.UserId > 0 {
			db.Where("user_id  = ?", *req.UserId)
		}
		if req.SubType != nil && *req.SubType > 0 {
			db.Where("subtype = ?", *req.SubType)
		}
		if req.TypeId != nil && *req.TypeId > 0 {
			db.Where("type_id  = ?", *req.TypeId)
		}
		if req.Keyword != "" {
			db.Where("name  LIKE ?", "%"+req.Keyword+"%")
		}
		if req.Source != nil {
			db.Where("source  = ?", *req.Source)
		}
		if req.Status != nil {
			db.Where("Status  = ?", *req.Status)
		}
		if req.IsSynced != nil {
			db.Where("is_synced  = ?", *req.IsSynced)
		}
		if req.ReNodeCode != nil && *req.ReNodeCode != "" {
			db.Where("re_node_code = ?", *req.ReNodeCode)
		}
		if req.EleType != nil {
			db.Where("type_id IN ?", eleTypeList)
		}
		if req.PriceDesc == "desc" {
			db.Order("price desc")
		} else {
			db.Order("price asc")
		}
	}

	var count int64
	err = db.Debug().Count(&count).Error
	if limit > 0 {
		db.Limit(limit).Offset(offset).Order("id desc")
	}
	var list []product.ProductElem
	err = db.Debug().Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("查询错误", zap.Error(err))
		response.FailWithMessage("查询错误", c)
		return
	}
	type Ret struct {
		Product      product.ProductElem   `json:"product"`
		ProductFiles []product.ProductFile `json:"productFiles"`
		SysType      system.ProductSysType `json:"sysType"`
		Download     int64                 `json:"download"`
	}
	var Rett []Ret
	for _, v := range list {
		var e Ret
		var sysType system.ProductSysType
		global.GVA_DB.Model(&system.ProductSysType{}).
			Where("is_deleted  = 0 and id = ?", v.TypeId).First(&sysType)
		e.Product = v
		e.SysType = sysType
		//
		var downCount int64 = 0
		//文件列表

		files := make([]product.ProductFile, 0)
		global.GVA_DB.Model(&product.ProductFile{}).Where("product_id =  ?  and is_deleted  = 0", v.Id).Find(&files)
		e.Download = downCount
		e.ProductFiles = files
		Rett = append(Rett, e)
	}
	response.OkWithDetailed(response.PageResult{
		List:     Rett,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "查询成功", c)
}

func (i *ProductElemApi) DeleteProduct(c *gin.Context) {
	id := c.Query("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	//isForce := c.Query("isForce")
	//if isForce != "true" {
	//	// 请求云管平台删除产品
	//	err = ConnMap.ProductStatusChangeRequest(productId, 5)
	//	if err != nil {
	//		global.GVA_LOG.Error("删除失败!", zap.Error(err))
	//		response.WSFail(c)
	//		return
	//	}
	//}
	deleteMap := map[string]interface{}{
		//"is_deleted":     constant.CommonRecordDeleted,
		"is_deleted":     1,
		"update_by":      utils.GetUserName(c),
		"update_user_id": utils.GetUserIDtoInt(c),
		"update_time":    time.Now(),
	}
	if err = global.GVA_DB.Model(&product.ProductElem{}).Where("id = ?", productId).Updates(deleteMap).Error; err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (i *ProductElemApi) DeleteProductByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idsInt := utils.StrslicetoInt(ids)
	//isForce := c.Query("isForce")
	//errHistory := ""
	//var SendErr error
	//if isForce != "true" {
	//	// 批量请求云管平台删除产品
	//	j := 0
	//	for _, v := range idsInt {
	//		SendErr = ConnMap.ProductStatusChangeRequest(v, 5)
	//		if SendErr != nil {
	//			var product admin.Product
	//			if er := global.GVA_DB.Model(&admin.Product{}).Where("id = ?", v).First(&product).Error; er != nil {
	//				global.GVA_LOG.Error("查找产品名失败!", zap.Error(er))
	//			}
	//			global.GVA_LOG.Error("删除失败!", zap.Error(SendErr))
	//			errHistory += product.Name + SendErr.Error() + ","
	//		} else {
	//			// 逻辑上，请求成功保留，失败则从切片中删除该id
	//			//将请求成功的id前移
	//			idsInt[j] = v
	//			j++
	//		}
	//	}
	//	// 截断多余部分
	//	idsInt = idsInt[:j]
	//}
	var productV2 product.ProductElem
	productV2.UpdateTime = time.Now()
	productV2.UpdateBy = utils.GetUserName(c)
	productV2.UpdateUserId = utils.GetUserIDtoInt(c)
	//productV2.IsDeleted = constant.CommonRecordDeleted
	productV2.IsDeleted = 1
	if err := global.GVA_DB.Model(&product.ProductElem{}).Where("id  in ? ", idsInt).Updates(&productV2).Error; err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败!", c)
		return
	}
	response.OkWithMessage("批量删除完成", c)
}
func (i *ProductElemApi) GetUserSupplyList(c *gin.Context) {
	var req struct {
		request.PageInfo
		//SubType 1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他
		SubType *int `json:"subType"  form:"subType"`
		//Status 审核状态（-1=驳回，0=待审核*，1=通过）
		Status  *int `json:"status"  form:"status"`
		EleType *int `json:"eleType" form:"eleType"`
		//Source  来源 0=平台 1=用户 2=节点
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	var eleTypeList []int
	if req.EleType != nil {
		global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ?", req.EleType).Pluck("id", &eleTypeList)
	}
	db := global.GVA_DB.Model(&product.ProductElem{}).Debug().Where("is_deleted  = 0 and type = 2 and  user_id = ? ", utils.GetUserIDtoI64(c))
	var count int64
	var list []product.ProductElem
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	if req.SubType != nil {
		db.Where("subtype = ?", *req.SubType)
	}
	if req.EleType != nil {
		db.Where("type_id IN ?", eleTypeList)
	}
	err = db.Debug().Count(&count).Error

	if limit > 0 {
		db.Debug().Limit(offset).Offset(offset).Order("id desc ")
	}
	err = db.Debug().Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	type Ret struct {
		Product      product.ProductElem   `json:"product"`
		ProductFiles []product.ProductFile `json:"productFiles"`
		SysType      system.ProductSysType `json:"sysType"`
		Download     int64                 `json:"download"`
	}
	var Rett []Ret
	for _, v := range list {
		var e Ret
		var sysType system.ProductSysType
		global.GVA_DB.Model(&system.ProductSysType{}).
			Where("is_deleted  = 0 and id = ?", v.TypeId).First(&sysType)
		e.Product = v
		e.SysType = sysType
		//
		var downCount int64 = 0
		//文件列表

		files := make([]product.ProductFile, 0)
		global.GVA_DB.Model(&product.ProductFile{}).Where("product_id =  ?  and is_deleted  = 0", v.Id).Find(&files)
		e.Download = downCount
		e.ProductFiles = files
		Rett = append(Rett, e)
	}
	response.OkWithDetailed(response.PageResult{
		List:     Rett,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "查询成功", c)
}
