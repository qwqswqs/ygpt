package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type systemOperateRouter struct{}

func (s *systemOperateRouter) InitSystemOperateRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/operate")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemOperateApi
	{
		dataRouter.DELETE("deleteSysOperate", modelApi.DeleteSystemOperate) // 删除数据
		dataRouter.POST("getSysOperateList", modelApi.GetSystemOperateList) // 获取数据列表
	}
}
