package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkVpcRouter struct {
}

func (s *NetworkVpcRouter) InitNetworkVpcRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/network/vpc")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.NetApiGroup.NetworkVpcApi

	{
		dataRouter.POST("getVpcList", vpcApi.GetVpcList)         // 获取数据列表
		dataRouter.POST("addVpc", vpcApi.AddVpc)                 // 新增
		dataRouter.POST("deleteVpc", vpcApi.DeleteVpc)           // 删除
		dataRouter.POST("deleteVpcByIds", vpcApi.DeleteVpcByIds) // 批量删除
	}
}
