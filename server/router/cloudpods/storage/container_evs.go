package storage

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ContainerEvsRouter struct {
}

func (s *ContainerEvsRouter) InitContainerEvsRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/storage/containerEvs")
	containerEvsApi := v1.ApiGroupApp.CloudApiGroup.ContainerEvsApi

	{
		dataRouter.POST("list", containerEvsApi.GetContainerEvsList)            // 获取容器云硬盘列表
		dataRouter.POST("create", containerEvsApi.CreateContainerEvs)           // 创建容器云硬盘
		dataRouter.POST("delete", containerEvsApi.DeleteContainerEvs)           // 删除容器云硬盘
		dataRouter.POST("deleteByIds", containerEvsApi.DeleteContainerEvsByIds) // 删除容器云硬盘
	}
}
