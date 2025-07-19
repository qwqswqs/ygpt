package config

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type LicenseConfigRouter struct {
}

func (s *LicenseConfigRouter) InitLicenseConfigRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("config/licenseConfig")
	backApi := v1.ApiGroupApp.YunGuanApiGroup.ConfigApi.LicenseConfigApi

	{
		dataRouter.POST("updateLicense", backApi.UpdateLicenseConfig) // 添加数据
		dataRouter.POST("queryLicense", backApi.GetLicenseConfig)     // 修改数据
		dataRouter.POST("loginLicense", backApi.LoginLicenseConfig)   // 修改数据
		dataRouter.POST("updateLicensePwd", backApi.UpdateLicensePwd) // 修改数据
	}
}
