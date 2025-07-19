package config

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type AlertConfigRouter struct {
}

func (s *LicenseConfigRouter) InitAlertConfigRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("config/alertConfig")
	backApi := v1.ApiGroupApp.YunGuanApiGroup.ConfigApi.AlertConfigApi

	{
		dataRouter.POST("addAlertConfig", backApi.AddAlertConfig)       // 添加数据
		dataRouter.POST("updateAlertConfig", backApi.UpdateAlertConfig) // 修改数据
		dataRouter.POST("deleteAlertConfig", backApi.DeleteAlertConfig) // 修改数据
		dataRouter.POST("queryAlertConfig", backApi.QueryAlertConfig)   // 修改数据
	}
}
