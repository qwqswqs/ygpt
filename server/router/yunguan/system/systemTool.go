package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemToolRouter struct{}

func (s *SystemToolRouter) InitSystemToolRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/tool")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemToolApi
	{
		dataRouter.POST("addSysTool", modelApi.AddSystemTool)         // 添加数据
		dataRouter.PUT("updateSysTool", modelApi.UpdateSystemTool)    // 添加数据
		dataRouter.DELETE("deleteSysTool", modelApi.DeleteSystemTool) // 删除数据
		dataRouter.POST("getSysToolList", modelApi.GetSystemToolList) // 获取数据列表
	}
}
