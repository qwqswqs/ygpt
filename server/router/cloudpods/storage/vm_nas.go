package storage

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type VmNasRouter struct {
}

func (s *VmNasRouter) InitVmNasRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/storage/vmNas")
	VmNasApi := v1.ApiGroupApp.CloudApiGroup.VmNasApi

	{
		dataRouter.POST("list", VmNasApi.GetVmNasList)                       // 获取NAS列表
		dataRouter.POST("create", VmNasApi.AddVmNas)                         // 创建NAS存储
		dataRouter.POST("delete", VmNasApi.DeleteVmNas)                      // 删除NAS存储
		dataRouter.POST("deleteByIds", VmNasApi.DeleteVmNasByIds)            // 批量删除NAS存储
		dataRouter.POST("disable", VmNasApi.DisableVmNas)                    // 禁用NAS存储
		dataRouter.POST("enable", VmNasApi.EnableVmNas)                      // 启用NAS存储
		dataRouter.POST("vmNasNoConnHostList", VmNasApi.VmNasNoConnHostList) // NAS存储未关联的宿主机列表
		dataRouter.POST("vmNasConnHostList", VmNasApi.VmNasConnHostList)     // NAS存储关联的宿主机列表
		dataRouter.POST("vmNasDisConnHost", VmNasApi.VmNasDisConnHost)       // 宿主机取消NAS存储关联
		dataRouter.POST("vmNasConnHost", VmNasApi.VmNasConnHost)             // 宿主机关联NAS存储
	}
}
