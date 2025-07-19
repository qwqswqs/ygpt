package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductConfigRouter struct {
}

func (s *ProductConfigRouter) InitProductConfigRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("product/config")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductConfigApi

	{
		dataRouter.POST("addConfigInfo", supplyInfoApi.AddConfigInfo)           // 添加数据
		dataRouter.PUT("updateConfigInfo", supplyInfoApi.UpdateConfigInfo)      // 修改数据
		dataRouter.DELETE("deleteConfigInfo", supplyInfoApi.DeleteConfigInfo)   // 删除数据
		dataRouter.POST("queryConfigInfo", supplyInfoApi.QueryConfigInfo)       // 获取数据列表
		dataRouter.POST("queryAllConfigInfo", supplyInfoApi.QueryAllConfigInfo) // 获取数据列表
	}
}
