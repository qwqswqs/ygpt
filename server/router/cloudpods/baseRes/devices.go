package baseRes

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type BaseDeviceRouter struct {
}

func (s *BaseDeviceRouter) InitBaseDeviceRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/baseRes/device")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.BaseResApiGroup.BaseDevicesApi

	{
		dataRouter.POST("getBaremetalGpuList", vpcApi.GetBaremetalGpuList) // 获取裸金属gpu列表
	}
}
