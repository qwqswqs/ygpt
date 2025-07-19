package cloudpods

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type BareMetalRouter struct {
}

func (c *BareMetalRouter) InitBareMetalRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/bareMetal")
	bareMetalApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.BareMetalApi

	{
		dataRouter.POST("list", bareMetalApi.HostList) // 获取物理机或裸金属
	}
}
