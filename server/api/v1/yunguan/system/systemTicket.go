package system

import (
	"encoding/json"
	"mime/multipart"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/utils"
	"ygpt/server/utils/minio"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemTicketApi struct {
}

// 添加系统运行工具
func (d *SystemTicketApi) AddSystemTicket(c *gin.Context) {
	var r system.SystemTicket
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r.CreatePerson = int(utils.GetUserID(c))
	serviceReturn, err := systemTicketService.AddSystemTicketService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改系统运行工具
func (d *SystemTicketApi) UpdateSystemTicket(c *gin.Context) {

	var model2 system.UpdateAndUploadReq
	err := c.ShouldBindJSON(&model2)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var model system.SystemTicket
	model.GVA_MODEL = model2.GVA_MODEL
	model.TicketType = model2.TicketType
	model.RenterID = model2.RenterID
	model.Type = model2.Type
	model.CreatePerson = model2.CreatePerson
	model.Status = model2.Status
	model.IsThird = model2.IsThird
	model.IsQuotation = model2.IsQuotation
	model.Quotation = model2.Quotation
	model.OrderId = model2.OrderId
	model.Description = model2.Description
	model.AssignPerson = model2.AssignPerson
	model.AssignTime = model2.AssignTime
	model.HandelTime = model2.HandelTime
	model.HandlePerson = model2.HandlePerson
	model.HandelCondition = model2.HandelCondition
	//将合同和发票信息发送到算力平台
	var quote float64
	if model2.Quotation == nil {
		quote = 0.0
	} else {
		quote = *model2.Quotation
	}
	if model2.IsQuotation == 1 || model2.TicketType == 4 {
		err = service.SuanLiServiceGroupApp.TicketService.QuoteRequest(int(model2.ID), model2.Contract, model2.Invoice, quote)
		if err != nil {
			global.GVA_LOG.Error("更新失败!", zap.Error(err))
			response.FailWithMessage("更新失败", c)
			return
		}
	}
	var ticket system.SystemTicket
	if err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id = ?", model2.ID).First(&ticket).Error; err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	if ticket.TicketType == 2 {
		var renterRes []renter.RenterRes
		if err = global.GVA_DB.Model(&renter.RenterRes{}).Where("order_id = ?", ticket.OrderId).Find(&renterRes).Error; err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Error(err))
			response.FailWithMessage("获取失败", c)
			return
		}
		for _, v := range renterRes {
			err = service.SuanLiServiceGroupApp.ResourceService.AllocateResponse(nil, v.OrderComputingId, v.ResourceID, true)
			if err != nil {
				global.GVA_LOG.Error("更新失败!", zap.Error(err))
				response.FailWithMessage("更新失败", c)
				return
			}
			s, err := global.GetCloudClientSession()
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			params := jsonutils.NewDict()
			params.Set("name", jsonutils.NewString(global.GenerateUniqueName(v.RenterID)))
			_, _ = compute.Servers.Update(s, v.ResourceID, params)
		}
	}
	err = systemTicketService.UpdateSystemTicketService(model)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 分配
func (d *SystemTicketApi) AssignSystemTicket(c *gin.Context) {
	var model system.SystemTicket
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	model.AssignPerson = int(utils.GetUserID(c))
	model.AssignTime = time.Now()
	model.Status = 2
	err = systemTicketService.UpdateSystemTicketService(model)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 获取系统运行工具列表
func (b *SystemTicketApi) GetSystemTicketList(c *gin.Context) {
	var pageInfo system.GetTicketReq
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, totalNum, unAssignNum, unHandleNum, handleNum, err := systemTicketService.GetSystemTicketListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{
		"list":        list,
		"total":       total,
		"page":        pageInfo.Page,
		"pageSize":    pageInfo.PageSize,
		"totalNum":    totalNum,
		"unAssignNum": unAssignNum,
		"unHandleNum": unHandleNum,
		"handleNum":   handleNum,
	}, "获取成功", c)
}

// 删除系统运行工具
func (b *SystemTicketApi) DeleteSystemTicket(c *gin.Context) {
	var model system.SystemLog
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemTicketService.DeleteSystemTicketService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

// UploadContract 工单合同上传
func (*SystemTicketApi) UploadContract(c *gin.Context) {
	var req struct {
		File   *multipart.FileHeader `json:"contract_file" form:"contract_file"` //合同文件
		Id     uint                  `json:"id" form:"id" `                      // 工单Id
		Remark string                `json:"remark" form:"remark"`               //合同描述
	}
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//上传文件
	uuid := utils.GenerateUUID()
	minioConfig := global.GVA_CONFIG.Minio
	err = minio.UploadFile(minioConfig.EndPoint, minioConfig.AccessKey, minioConfig.SecretKey, minioConfig.BucketName, uuid, req.File)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 数据库信息更改
	info, err := json.Marshal(map[string]interface{}{
		"objectId": req.Id,
		"uploadAt": time.Now().Unix(),
		"remark":   req.Remark,
		"fileName": req.File.Filename,
	})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id = ?", req.Id).Updates(&system.SystemTicket{
		ContractJson: string(info),
	}).Error
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("合同信息上传成功", c)
	//添加合同记录

	//var req struct {
	//	Id               uint   `json:"id"`                 // 工单Id
	//	ContractFileLink string `json:"contract_file_link"` //合同附件地址
	//	Description      string `json:"description"`        // 合同备注描述
	//}
	//err := c.ShouldBindJSON(&req)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//// 处理添加工单合同
	//var contractJson string
	//info, err := json.Marshal(&req)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//contractJson = string(info)
	//var do = &system.SystemTicket{
	//	ContractJson: contractJson,
	//}
	//err = global.GVA_DB.Model(&do).Where("id = ?", req.Id).Updates(&do).Error
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
}

// UploadInvoice 发票上传
func (*SystemTicketApi) UploadInvoice(c *gin.Context) {
	var req struct {
		File   *multipart.FileHeader `json:"file" form:"file"`     //工单文件
		Id     uint                  `json:"id" form:"id" `        // 工单Id
		Remark string                `json:"remark" form:"remark"` //合同描述
	}
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uuid := utils.GenerateUUID()
	//上传文件
	minioConfig := global.GVA_CONFIG.Minio
	err = minio.UploadFile(minioConfig.EndPoint, minioConfig.AccessKey, minioConfig.SecretKey, minioConfig.BucketName, uuid, req.File)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 数据库信息更改
	info, err := json.Marshal(map[string]interface{}{
		"objectId": req.Id,
		"uploadAt": time.Now().Unix(),
		"remark":   req.Remark,
		"fileName": req.File.Filename,
	})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id = ?", req.Id).Updates(&system.SystemTicket{
		InvoiceJson: string(info),
	}).Error
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("发票信息上传成功", c)
	//var req struct {
	//	Id              uint   `json:"id"`
	//	InvoiceFileLink string `json:"invoice_link"`
	//}
	//
	//err := c.ShouldBindQuery(&req)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//info, err := json.Marshal(&req)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//var do = &system.SystemTicket{
	//	InvoiceJson: string(info),
	//}
	//err = global.GVA_DB.Model(&do).Where("id = ?", req.Id).Updates(&do).Error
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//response.OkWithMessage("发票上传成功！", c)
}

// PostQuotation 工单报价设置
func (*SystemTicketApi) PostQuotation(c *gin.Context) {
	var req struct {
		Id        uint    `json:"id"`
		Quotation float64 `json:"quotation"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//更改数据
	err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id = ? and is_quotation =  ?", req.Id, 1).Updates(&system.SystemTicket{
		Quotation: &req.Quotation,
	}).Error
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	////发送报价请求
	//err = service.SuanLiServiceGroupApp.TicketService.QuoteRequest(int(req.Id))
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	response.OkWithMessage("工单报价成功", c)
}

var AssignResourceReq struct {
	ResourceIDs []string `json:"resourceIDs"`
}

func (*SystemTicketApi) AssignResource(c *gin.Context) {

}

type GetResReq struct {
	Type int `json:"type"`
}

func (*SystemTicketApi) GetNoAssignRes(c *gin.Context) {
	var req GetResReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	res, total, err := systemTicketService.GetNoAssignRes(req.Type)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  res,
		Total: total,
	}, "获取成功", c)
}

func (*SystemTicketApi) GetAllRes(c *gin.Context) {
	var req GetResReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	res, total, err := systemTicketService.GetAllRes(req.Type)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  res,
		Total: total,
	}, "获取成功", c)
}
