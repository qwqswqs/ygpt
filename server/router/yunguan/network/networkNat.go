package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkNatRouter struct {
}

func (s *NetworkNatRouter) InitNetworkNatRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("network/nat")
	natApi := v1.ApiGroupApp.YunGuanApiGroup.NetworkApi.NetworkNatApi

	{
		dataRouter.POST("addNat", natApi.AddNetworkNat)         // 添加数据
		dataRouter.PUT("updateNat", natApi.UpdateNetworkNat)    // 修改数据
		dataRouter.DELETE("deleteNat", natApi.DeleteNetworkNat) // 删除数据
		dataRouter.POST("queryNat", natApi.QueryNetworkNat)     // 获取数据列表
	}
}
