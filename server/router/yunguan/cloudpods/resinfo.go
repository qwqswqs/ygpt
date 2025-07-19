package cloudpods

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type CloudpodsResInfoRouter struct {
}

func (s *CloudpodsResInfoRouter) InitCloudpodsResInfoRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/resinfo")
	resInfoApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.ResInfoApi

	{
		dataRouter.GET("", resInfoApi.GetResInfo) // 列出资源信息
	}
}
