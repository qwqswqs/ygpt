package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
)

type ProductCategoryApi struct {
}

// 添加
func (p *ProductCategoryApi) AddCategoryInfo(c *gin.Context) {
	var r product.ProductCategory
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	serviceReturn, err := productCategoryService.AddProductCategoryService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!"+err.Error(), zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败!"+err.Error(), c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

// 查询
func (p *ProductCategoryApi) QueryCategoryInfo(c *gin.Context) {
	list, total, err := productCategoryService.GetCategoryListService()
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
func (p *ProductCategoryApi) DeleteCategoryInfo(c *gin.Context) {
	var r product.ProductCategory
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productCategoryService.DeleteProductCategoryService(r)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
