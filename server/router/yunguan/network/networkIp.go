package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkIpRouter struct {
}

func (s *NetworkIpRouter) InitNetworkIpRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("network/ip")
	ipApi := v1.ApiGroupApp.YunGuanApiGroup.NetworkApi.NetworkIpApi

	{
		dataRouter.POST("addIp", ipApi.AddNetworkIp)           // 添加数据
		dataRouter.PUT("updateIp", ipApi.UpdateNetworkIp)      // 修改数据
		dataRouter.DELETE("deleteIp", ipApi.DeleteNetworkIp)   // 删除数据
		dataRouter.POST("queryIp", ipApi.QueryNetworkIp)       // 获取数据列表
		dataRouter.POST("queryAllIp", ipApi.QueryAllNetworkIp) // 获取数据列表
	}
}
