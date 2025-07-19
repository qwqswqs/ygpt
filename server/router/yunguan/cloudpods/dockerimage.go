package cloudpods

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type DockerImageRouter struct {
}

func (s *SystemImageRouter) InitDockerImageRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/dockerimage")
	imageApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.ImageApi

	{
		dataRouter.GET("list", imageApi.GetDockerImage) // 列出系统镜像
	}
}
