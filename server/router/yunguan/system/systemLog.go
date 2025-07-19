package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemLogRouter struct{}

func (s *SystemLogRouter) InitSystemLogRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/log")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemLogApi
	{
		dataRouter.DELETE("deleteSysLog", modelApi.DeleteSystemLog) // 删除数据
		dataRouter.POST("getSysLogList", modelApi.GetSystemLogList) // 获取数据列表
	}
}
