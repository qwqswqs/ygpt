package udp

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type MonitorRouter struct {
}

func (s *MonitorRouter) InitUdpMonitorRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("udp/instance")
	instanceApi := v1.ApiGroupApp.ComputeApiGroup.UdpApi.InstanceApi
	{

		dataRouter.POST("install", instanceApi.Install)       // 安装软件
		dataRouter.POST("download", instanceApi.Download)     // 下载数据集、模型
		dataRouter.POST("addTask", instanceApi.AddTask)       // 添加任务
		dataRouter.POST("listTask", instanceApi.ListTask)     //获取任务列表
		dataRouter.POST("manageTask", instanceApi.ManageTask) //获取任务列表
	}
}
