package network

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type NetEipRouter struct {
}

func (s *NetEipRouter) InitNetworkEipRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/network/eip")
	ipApi := v1.ApiGroupApp.CloudApiGroup.NetApiGroup.NetworkEipApi

	{
		dataRouter.POST("getNetEipList", ipApi.GetEipList)            // 获取eip列表
		dataRouter.POST("addNetEip", ipApi.AddEip)                    // 新增eip
		dataRouter.POST("deleteNetEip", ipApi.DeleteEip)              // 删除eip
		dataRouter.POST("deleteNetEipByIds", ipApi.DeleteEipByIds)    // 删除eip
		dataRouter.POST("bindEipServer", ipApi.BindEipServer)         // 绑定
		dataRouter.POST("unbindEipServer", ipApi.UnbindEipServer)     // 解绑
		dataRouter.POST("getBindServerList", ipApi.GetBindServerList) // 获取可绑定云主机
		dataRouter.POST("modEipBandWidth", ipApi.ModEipBandWidth)     // 修改带宽
	}
}
