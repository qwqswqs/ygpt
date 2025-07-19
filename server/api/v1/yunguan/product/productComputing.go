package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/service/compute/websocket/service/business"
	"ygpt/server/utils"
)

type ProductComputingApi struct {
}

// 添加
func (p *ProductComputingApi) AddComputingInfo(c *gin.Context) {
	var r product.ProductComputing
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	r.AuditStatus = 0
	serviceReturn, err := productComputingService.AddProductComputingService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

// 修改
func (p *ProductComputingApi) UpdateComputingInfo(c *gin.Context) {
	var r product.ProductComputing
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	r.IsSynced = 1
	err = productComputingService.UpdateProductComputingService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

var syncedService = business.SpecUploadService{}

// 同步单个数据
func (p *ProductComputingApi) SyncedComputingInfo(c *gin.Context) {
	var r product.ProductComputing
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = syncedService.Push(int(r.ID))
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败"+err.Error(), c)
		return
	}
	err = productComputingService.UpdateProductComputingService(r)
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "同步成功", c)
}

// 同步所有数据
func (p *ProductComputingApi) SyncedAllComputingInfo(c *gin.Context) {
	err := syncedService.SyncCloudComputeInfo()
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	err = productComputingService.SyncedAllProductComputingService()
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "同步成功", c)
}

// 查询
func (p *ProductComputingApi) QueryComputingInfo(c *gin.Context) {
	var pageInfo product.GetProductComputingReq
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
	list, isSyncedNum, total, err := productComputingService.GetComputingInfoListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list, "total": total, "isSynced": isSyncedNum, "page": pageInfo.Page, "pageSize": pageInfo.PageSize}, "获取成功", c)
}

// 删除
func (p *ProductComputingApi) DeleteComputingInfo(c *gin.Context) {
	var r product.ProductComputing
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productComputingService.DeleteProductComputingService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	err = service.SuanLiServiceGroupApp.SyncProductStatusChangeToPlatform(int(r.ID), 1, 5)
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
