package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/utils"
)

type ProductConfigApi struct {
}

// 添加
func (p *ProductConfigApi) AddConfigInfo(c *gin.Context) {
	var r product.ProductConfig
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	serviceReturn, err := productConfigService.AddProductConfigService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!"+err.Error(), zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

// 修改
func (p *ProductConfigApi) UpdateConfigInfo(c *gin.Context) {
	var r product.ProductConfig
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productConfigService.UpdateProductConfigService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *ProductConfigApi) QueryConfigInfo(c *gin.Context) {
	var pageInfo product.GetProductConfigListReq
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
	list, total, err := productConfigService.GetConfigListService(pageInfo)
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
func (p *ProductConfigApi) QueryAllConfigInfo(c *gin.Context) {
	list, total, err := productConfigService.GetAllConfigListService()
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
func (p *ProductConfigApi) DeleteConfigInfo(c *gin.Context) {
	var r product.ProductConfig
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productConfigService.DeleteProductConfigService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
