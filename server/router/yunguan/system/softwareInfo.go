package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemSoftwareRouter struct{}

func (s *SystemSoftwareRouter) InitSystemSoftwareRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/software")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemSoftwareApi
	{
		dataRouter.POST("addSysSoftware", modelApi.AddSystemSoftware)         // 添加数据
		dataRouter.PUT("updateSysSoftware", modelApi.UpdateSystemSoftware)    // 添加数据
		dataRouter.DELETE("deleteSysSoftware", modelApi.DeleteSystemSoftware) // 删除数据
		dataRouter.POST("getSysSoftwareList", modelApi.GetSystemSoftwareList) // 获取数据列表
		dataRouter.GET("startScan", modelApi.StartScan)                       //启动病毒扫描
		dataRouter.GET("stopScan", modelApi.StopScan)                         //停止病毒扫描
		//dataRouter.GET("uninstallClamAV", modelApi.UninstallClamAV)          //卸载clamav
		dataRouter.GET("systemCheck", modelApi.SystemCheck) //一键检测
	}
}
