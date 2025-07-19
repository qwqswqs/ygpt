package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductCategoryRouter struct {
}

func (s *ProductCategoryRouter) InitProductCategoryRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("product/Category")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductCategoryApi

	{
		dataRouter.POST("addCategoryInfo", supplyInfoApi.AddCategoryInfo)         // 添加数据
		dataRouter.DELETE("deleteCategoryInfo", supplyInfoApi.DeleteCategoryInfo) // 删除数据
		dataRouter.POST("queryCategoryInfo", supplyInfoApi.QueryCategoryInfo)     // 获取数据列表
	}
}
