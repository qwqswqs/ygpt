package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemPlatformConnectRouter struct{}

func (s *SystemPlatformConnectRouter) InitSystemPlatformConnectRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/connect")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.PlatformConnectApi
	{
		dataRouter.GET("getIsConnected", modelApi.GetIsConnected)
		dataRouter.POST("updateConnectConfigAndConnect", modelApi.UpdateConnectConfigAndConnect)
	}
}
