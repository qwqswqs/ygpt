package storage

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type VmEvsRouter struct {
}

func (s *VmEvsRouter) InitVmEvsRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/storage/vmEvs")
	vmEvsApi := v1.ApiGroupApp.CloudApiGroup.VmEvsApi

	{
		dataRouter.POST("list", vmEvsApi.GetVmEvsList)                 // 获取虚拟机云硬盘列表
		dataRouter.POST("create", vmEvsApi.CreateVmEvs)                // 创建虚拟机云硬盘
		dataRouter.POST("delete", vmEvsApi.DeleteVmEvs)                // 删除虚拟机云硬盘
		dataRouter.POST("deleteByIds", vmEvsApi.DeleteVmEvsByIds)      // 删除虚拟机云硬盘
		dataRouter.POST("resize", vmEvsApi.ResizeVmEvs)                // 扩容虚拟机云硬盘
		dataRouter.POST("attach", vmEvsApi.AttachVmEvs)                // 挂载虚拟机云硬盘
		dataRouter.POST("detach", vmEvsApi.DetachVmEvs)                // 卸载虚拟机云硬盘
		dataRouter.POST("getAttachableVms", vmEvsApi.GetAttachableVms) // 获取可挂载的虚拟机

		dataRouter.POST("getRecycleVmEvsList", vmEvsApi.GetRecycleVmEvsList) // 获取回收站硬盘列表
		dataRouter.POST("clearRecycleVmEvs", vmEvsApi.ClearRecycleVmEvs)     // 获取回收站硬盘列表
		dataRouter.POST("recoverRecycleVmEvs", vmEvsApi.RecoverRecycleVmEvs) // 获取回收站硬盘列表
	}
}
