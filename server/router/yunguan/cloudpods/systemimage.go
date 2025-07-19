package cloudpods

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemImageRouter struct {
}

func (s *SystemImageRouter) InitSystemImageRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/systemimage")
	imageApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.ImageApi

	{
		dataRouter.GET("list", imageApi.GetSystemImage) // 列出系统镜像
	}
}
