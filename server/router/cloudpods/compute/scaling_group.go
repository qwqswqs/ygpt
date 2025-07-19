package compute

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type ScalingGroupRouter struct {
}

func (s *ScalingGroupRouter) InitScalingGroupRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/scaling/group")
	ipApi := v1.ApiGroupApp.CloudApiGroup.ScalingGroupApi
	{
		dataRouter.POST("list", ipApi.List)
		dataRouter.POST("create", ipApi.Create)
		dataRouter.POST("delete", ipApi.Delete)
		dataRouter.POST("deleteByIds", ipApi.DeleteByIds)
		dataRouter.POST("enable", ipApi.Enable)
		dataRouter.POST("disable", ipApi.Disable)
		dataRouter.POST("detach", ipApi.Detach)
	}
}
