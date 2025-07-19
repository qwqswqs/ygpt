package storage

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ContainerNasRouter struct {
}

func (s *ContainerEvsRouter) InitContainerNasRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/storage/containerNas")
	containerNasApi := v1.ApiGroupApp.CloudApiGroup.ContainerNasApi

	{
		dataRouter.POST("list", containerNasApi.GetContainerNasList)            // 获取容器Nas列表
		dataRouter.POST("create", containerNasApi.CreateContainerNas)           // 创建容器Nas
		dataRouter.POST("delete", containerNasApi.DeleteContainerNas)           // 删除容器Nas
		dataRouter.POST("deleteByIds", containerNasApi.DeleteContainerNasByIds) // 删除容器Nas

		dataRouter.POST("secretList", containerNasApi.GetCephSecretList)      // 获取保密字典列表
		dataRouter.POST("secretCreate", containerNasApi.CreateCephSecret)     // 创建保密字典
		dataRouter.POST("secretDelete", containerNasApi.DeleteCephSecret)     // 删除保密字典
		dataRouter.POST("getCephClusterId", containerNasApi.GetCephClusterId) // 获取ceph集群的id
		dataRouter.POST("getCephPools", containerNasApi.GetCephPools)         // 获取ceph集群的存储池列表
	}
}
