package compute

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ScalingPolicyRouter struct {
}

func (s *ScalingGroupRouter) InitScalingPolicyRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/scaling/policy")
	ipApi := v1.ApiGroupApp.CloudApiGroup.ScalingPolicyApi
	{
		dataRouter.POST("list", ipApi.List)
		dataRouter.POST("create", ipApi.Create)
		dataRouter.POST("delete", ipApi.Delete)
		dataRouter.POST("enable", ipApi.Enable)
		dataRouter.POST("disable", ipApi.Disable)
	}
}
