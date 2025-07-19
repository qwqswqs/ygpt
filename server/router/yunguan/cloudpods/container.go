package cloudpods

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ContainerRouter struct {
}

func (c *ContainerRouter) InitContainerRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/container")
	containerApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.ContainerApi

	{
		dataRouter.POST("create", containerApi.HostCreate)        // 创建容器
		dataRouter.POST("connect", containerApi.GetWebConsoleSSH) // 获取容器控制台
		dataRouter.POST("delete", containerApi.HostDelete)        // 删除容器
	}
}
