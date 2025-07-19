package common

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type CloudCommonLogRouter struct{}

func (s *CloudCommonLogRouter) InitCloudCommonLogRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/common/log")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.CloudCommonApiGroup.CloudCommonLogApi

	{
		dataRouter.POST("getLogsList", vpcApi.GetLogsList) // 获取日志列表
	}
}
