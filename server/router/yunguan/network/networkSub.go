package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkSubRouter struct {
}

func (s *NetworkSubRouter) InitNetworkSubRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("network/sub")
	natApi := v1.ApiGroupApp.YunGuanApiGroup.NetworkApi.NetworkSubApi

	{
		dataRouter.POST("addSub", natApi.AddNetworkSub)           // 添加数据
		dataRouter.PUT("updateSub", natApi.UpdateNetworkSub)      // 修改数据
		dataRouter.DELETE("deleteSub", natApi.DeleteNetworkSub)   // 删除数据
		dataRouter.POST("querySub", natApi.QueryNetworkSub)       // 获取数据列表
		dataRouter.POST("queryAllSub", natApi.QueryAllNetworkSub) // 获取数据列表
	}
}
