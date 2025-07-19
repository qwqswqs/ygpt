package compute

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type ServersRouter struct {
}

func (s *ServersRouter) InitServersRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/compute/servers")
	ipApi := v1.ApiGroupApp.CloudApiGroup.ComputeServersApi
	{
		dataRouter.POST("getServerList", ipApi.GetServerList)       // 获取云主机列表
		dataRouter.POST("getWebConsoleSSH", ipApi.GetWebConsoleSSH) // 获取WebSSH连接
		dataRouter.POST("hostShutDown", ipApi.HostShutDown)         // 关机
		dataRouter.POST("hostStart", ipApi.HostStart)               // 开机
		dataRouter.POST("hostRestart", ipApi.HostRestart)           // 重启
		dataRouter.POST("hostSuspend", ipApi.HostSuspend)           // 挂起
		dataRouter.POST("hostResume", ipApi.HostResume)             // 恢复
		dataRouter.POST("hostCreate", ipApi.HostCreate)             // 云主机创建
		dataRouter.POST("hostDelete", ipApi.HostDelete)             // 删除
		dataRouter.POST("hostDeleteByIds", ipApi.HostDeleteByIds)   // 批量删除
		dataRouter.POST("hostDetails", ipApi.HostDetails)           // 获取主机详情
		dataRouter.POST("getGpuInfo", ipApi.GetGpuInfo)             // 获取GPU信息
		dataRouter.POST("manageGpu", ipApi.ManageGpu)               // 管理GPU
		dataRouter.POST("migrateForecast", ipApi.MigrateForecast)   // 获取云主机可迁移的宿主机列表
		dataRouter.POST("migrate", ipApi.Migrate)                   // 云主机迁移

		dataRouter.POST("getBaremetalList", ipApi.GetBareServerList) // 获取裸金属列表
		dataRouter.POST("baremetalCreate", ipApi.BaremetalCreate)    //裸金属创建

		dataRouter.POST("queryCpu", ipApi.GetHostMonitorCPUData)                   // 获取监控CPU数据
		dataRouter.POST("queryMem", ipApi.GetHostMonitorMemData)                   // 获取监控内存数据
		dataRouter.POST("queryDisk", ipApi.GetHostMonitorDiskData)                 // 获取监控磁盘使用率
		dataRouter.POST("queryDiskRead", ipApi.GetHostMonitorDiskReadData)         // 获取监控磁盘读速率
		dataRouter.POST("queryDiskWrite", ipApi.GetHostMonitorDiskWriteData)       // 获取监控磁盘写速率
		dataRouter.POST("queryBpsRecv", ipApi.GetHostMonitorBpsRecvData)           // 获取监控网络入带宽
		dataRouter.POST("queryBpsSent", ipApi.GetHostMonitorBpsSentData)           // 获取监控网络出带宽
		dataRouter.POST("queryPpsRecv", ipApi.GetHostMonitorPpsRecvData)           // 获取监控网络收包数
		dataRouter.POST("queryPpsSent", ipApi.GetHostMonitorPpsSentData)           // 获取监控网络发包数
		dataRouter.POST("queryGpu", ipApi.GetHostMonitorGPUData)                   // 获取GPU使用率
		dataRouter.POST("queryGpuMem", ipApi.GetHostMonitorGPUMemData)             // 获取GPU显存使用率
		dataRouter.POST("queryGpuTemperature", ipApi.GetHostMonitorGpuTemperature) // 获取Gpu温度

		dataRouter.POST("getResSkuList", ipApi.GetResSkuList) // 查询套餐列表
		dataRouter.POST("getLoginInfo", ipApi.GetLoginInfo)   // 获取账号密码
		dataRouter.POST("resetPassword", ipApi.ResetPassword) // 重置账号密码

		dataRouter.POST("getGpuList", ipApi.GetGpuList)                       // 获取可使用的GPU信息
		dataRouter.POST("getRecycleServerList", ipApi.GetRecycleServerList)   // 获取回收站列表
		dataRouter.POST("clearRecycleHost", ipApi.ClearRecycleHost)           // 清除回收主机
		dataRouter.POST("clearRecycleHostByIds", ipApi.ClearRecycleHostByIds) // 清除回收主机
		dataRouter.POST("recoverRecycleHost", ipApi.RecoverRecycleHost)       // 恢复回收主机

		dataRouter.POST("updateName", ipApi.UpdateName) // 修改名称

		dataRouter.POST("forecastAgent", ipApi.ForecastAgent)                     //探测是否有agent
		dataRouter.POST("forecastAgentCanInstall", ipApi.ForecastAgentCanInstall) //探测是否具备安装agent的条件
		dataRouter.POST("setAgentCanInstall", ipApi.SetAgentCanInstall)           //输入ssh信息设置其能安装agent
		dataRouter.POST("installAgent", ipApi.InstallAgent)                       //安装agent
	}
}
