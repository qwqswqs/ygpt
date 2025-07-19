package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemCacheRouter struct{}

func (s *SystemToolRouter) InitSystemCacheRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/cache")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemCacheApi
	{
		dataRouter.POST("addSysCache", modelApi.AddSystemCache)         // 添加数据
		dataRouter.PUT("updateSysCache", modelApi.UpdateSystemCache)    // 添加数据
		dataRouter.DELETE("deleteSysCache", modelApi.DeleteSystemCache) // 删除数据
		dataRouter.GET("getSysCacheList", modelApi.GetSystemCacheList)  // 获取数据列表
	}
}
