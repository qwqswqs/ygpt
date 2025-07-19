package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemConfigRouter struct{}

func (s *SystemConfigRouter) InitSystemConfigRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/config")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemConfigApi
	{
		dataRouter.POST("addSysConfig", modelApi.AddSystemConfig)         // 添加数据
		dataRouter.PUT("updateSysConfig", modelApi.UpdateSystemConfig)    // 添加数据
		dataRouter.DELETE("deleteSysConfig", modelApi.DeleteSystemConfig) // 删除数据
		dataRouter.POST("getSysConfigList", modelApi.GetSystemConfigList) // 获取数据列表
	}
}
