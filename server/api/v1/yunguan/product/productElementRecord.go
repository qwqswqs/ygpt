package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/utils"
)

type ProductElementRecordApi struct {
}

// 添加
func (p *ProductElementRecordApi) AddElementRecordInfo(c *gin.Context) {
	var r product.ProductElementRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	serviceReturn, err := productElementRecordService.AddProductElementRecordService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

// 修改
func (p *ProductElementRecordApi) UpdateElementRecordInfo(c *gin.Context) {
	var r product.ProductElementRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productElementRecordService.UpdateProductElementRecordService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *ProductElementRecordApi) QueryElementRecordInfo(c *gin.Context) {
	var pageInfo request.PageInfo
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
	list, total, err := productElementRecordService.GetElementRecordListService(pageInfo)
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

// 删除
func (p *ProductElementRecordApi) DeleteElementRecordInfo(c *gin.Context) {
	var r product.ProductElementRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productElementRecordService.DeleteProductElementRecordService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
