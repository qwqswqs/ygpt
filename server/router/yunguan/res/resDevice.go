package res

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ResDeviceRouter struct {
}

func (s *ResDeviceRouter) InitResDeviceRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("res/device")
	diskApi := v1.ApiGroupApp.YunGuanApiGroup.ResApi.ResDeviceApi

	{
		dataRouter.POST("addDevice", diskApi.AddResDevice)                 // 添加数据
		dataRouter.PUT("updateDevice", diskApi.UpdateResDevice)            // 修改数据
		dataRouter.DELETE("deleteDevice", diskApi.DeleteResDevice)         // 删除数据
		dataRouter.POST("queryDevice", diskApi.QueryResDevice)             // 获取数据列表
		dataRouter.POST("queryAllServer", diskApi.QueryAllServerResDevice) // 获取数据列表
	}
}
