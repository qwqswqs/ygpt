package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkVpcRouter struct {
}

func (s *NetworkVpcRouter) InitNetworkVpcRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("network/vpc")
	vpcApi := v1.ApiGroupApp.YunGuanApiGroup.NetworkApi.NetworkVpcApi

	{
		dataRouter.POST("addVpc", vpcApi.AddNetworkVpc)             // 修改数据
		dataRouter.PUT("updateVpc", vpcApi.UpdateNetworkVpc)        // 修改数据
		dataRouter.DELETE("deleteVpc", vpcApi.DeleteNetworkVpc)     // 删除数据
		dataRouter.POST("getVpcList", vpcApi.QueryNetworkVpc)       // 获取数据列表
		dataRouter.POST("getAllVpcList", vpcApi.QueryAllNetworkVpc) // 获取数据列表
	}
}
