package cloudpods

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type JumpRouter struct {
}

func (j *JumpRouter) InitJumpRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods")
	jumpApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.JumpApi

	{
		dataRouter.GET("jump", jumpApi.Jump) // 获取跳转url
	}
}
