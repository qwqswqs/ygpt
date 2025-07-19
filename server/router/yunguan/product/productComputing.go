package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductComputingRouter struct {
}

func (s *ProductComputingRouter) InitProductComputingRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("product/computing")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductComputingApi

	{
		dataRouter.POST("addComputingInfo", supplyInfoApi.AddComputingInfo)             // 添加数据
		dataRouter.PUT("updateComputingInfo", supplyInfoApi.UpdateComputingInfo)        // 修改数据
		dataRouter.POST("syncedComputingInfo", supplyInfoApi.SyncedComputingInfo)       // 同步单个数据
		dataRouter.POST("syncedAllComputingInfo", supplyInfoApi.SyncedAllComputingInfo) // 同步所有数据
		dataRouter.DELETE("deleteComputingInfo", supplyInfoApi.DeleteComputingInfo)     // 删除数据
		dataRouter.POST("queryComputingInfo", supplyInfoApi.QueryComputingInfo)         // 获取数据列表
	}
}
