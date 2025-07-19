package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetIpRouter struct {
}

func (s *NetIpRouter) InitNetworkIpRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/network/ip")
	ipApi := v1.ApiGroupApp.CloudApiGroup.NetApiGroup.NetIpApi

	{
		dataRouter.POST("getNetIpList", ipApi.GetNetIPList)         // 获取Ip列表
		dataRouter.POST("addNetIP", ipApi.AddNetIP)                 // 新增Ip
		dataRouter.POST("deleteNetIP", ipApi.DeleteNetIP)           // 删除Ip
		dataRouter.POST("deleteNetIpByIds", ipApi.DeleteNetIpByIds) // 批量删除Ip
		dataRouter.POST("updateNetIP", ipApi.UpdateNetIP)           // 修改Ip
	}
}
