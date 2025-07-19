package compute

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type ScalingConfigRouter struct {
}

func (s *ScalingConfigRouter) InitScalingConfigRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/scaling/config")
	ipApi := v1.ApiGroupApp.CloudApiGroup.ScalingConfigApi
	{
		dataRouter.POST("list", ipApi.List)
		dataRouter.POST("create", ipApi.Create)
		dataRouter.POST("delete", ipApi.Delete)
		dataRouter.POST("deleteByIds", ipApi.DeleteByIds)
	}
}
