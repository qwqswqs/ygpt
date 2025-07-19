package baseRes

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type BaseHostRouter struct {
}

func (s *BaseHostRouter) InitBaseHostRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/baseRes/baseHost")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.BaseResApiGroup.BaseHostApi

	{
		dataRouter.POST("getBaseHostList", vpcApi.GetBaseHostList)                     // 获取宿主机列表
		dataRouter.POST("getBaseBareHostList", vpcApi.GetBaseBareHostList)             // 获取物理机列表
		dataRouter.POST("getImpiInfo", vpcApi.GetImpiInfo)                             // 获取IMPI信息
		dataRouter.POST("getLoginInfo", vpcApi.GetLoginInfo)                           // 获取物理机登录信息
		dataRouter.POST("updateBareHostStatus", vpcApi.UpdateBareHostStatus)           // 修改物理机状态
		dataRouter.POST("updateBareHostDescription", vpcApi.UpdateBareHostDescription) // 修改物理机备注
		dataRouter.POST("deleteBareHost", vpcApi.DeleteBareHost)                       // 删除物理机
		dataRouter.POST("deleteBareHostByIds", vpcApi.DeleteBareHostByIds)             // 批量删除物理机
		dataRouter.POST("addBareHost", vpcApi.AddBareHost)                             // 创建物理机
		dataRouter.POST("getWebConsoleSSH", vpcApi.GetWebConsoleSSH)                   // 远程连接物理机

		dataRouter.POST("queryCpu", vpcApi.GetHostMonitorCPUData)             // 获取监控CPU数据
		dataRouter.POST("queryMem", vpcApi.GetHostMonitorMemData)             // 获取监控内存数据
		dataRouter.POST("queryDisk", vpcApi.GetHostMonitorDiskData)           // 获取监控磁盘使用率
		dataRouter.POST("queryDiskRead", vpcApi.GetHostMonitorDiskReadData)   // 获取监控磁盘读速率
		dataRouter.POST("queryDiskWrite", vpcApi.GetHostMonitorDiskWriteData) // 获取监控磁盘写速率
		dataRouter.POST("queryBpsRecv", vpcApi.GetHostMonitorBpsRecvData)     // 获取监控网络入带宽
		dataRouter.POST("queryBpsSent", vpcApi.GetHostMonitorBpsSentData)     // 获取监控网络出带宽
	}
}
