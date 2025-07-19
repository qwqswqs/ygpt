package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemAlertRouter struct{}

func (s *SystemAlertRouter) InitSystemAlertRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/alert")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemAlertApi
	{
		dataRouter.POST("addSystemAlert", modelApi.AddSystemAlert)         // 添加数据
		dataRouter.POST("deleteSystemAlert", modelApi.DeleteSystemAlert)   // 删除数据
		dataRouter.POST("updateSystemAlert", modelApi.UpdateSystemAlert)   // 删除数据
		dataRouter.POST("getSystemAlertList", modelApi.GetSystemAlertList) // 获取数据列表
	}
}
