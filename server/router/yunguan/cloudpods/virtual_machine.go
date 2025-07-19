package cloudpods

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type VirtualMachineRouter struct {
}

func (v *VirtualMachineRouter) InitVirtualMachineRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloudpods/virtual_machine")
	virtualMachineApi := v1.ApiGroupApp.YunGuanApiGroup.CloudpodsApi.VirtualMachineApi

	{
		dataRouter.POST("create", virtualMachineApi.HostCreate)                 // 创建云主机
		dataRouter.POST("connect", virtualMachineApi.GetWebConsoleSSH)          // 获取云主机控制台
		dataRouter.POST("delete", virtualMachineApi.HostDelete)                 // 删除云主机
		dataRouter.POST("limitBandwidth", virtualMachineApi.HostBandwidthLimit) //限制云主机带宽
		dataRouter.POST("listByIds", virtualMachineApi.HostListByIds)

		dataRouter.POST("hostStart", virtualMachineApi.HostStart)                   // 开机
		dataRouter.POST("hostShutDown", virtualMachineApi.HostShutDown)             // 关机
		dataRouter.POST("hostRestart", virtualMachineApi.HostRestart)               // 重启
		dataRouter.POST("hostSuspend", virtualMachineApi.HostSuspend)               // 挂起
		dataRouter.POST("hostResume", virtualMachineApi.HostResume)                 // 恢复
		dataRouter.POST("createServerSnap", virtualMachineApi.HostCreateServerSnap) // 创建主机快照
		dataRouter.POST("resetServerSnap", virtualMachineApi.HostResetServerSnap)   // 主机回滚快照
	}
}
