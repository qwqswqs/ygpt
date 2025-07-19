package baseRes

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type BaseK8SRouter struct {
}

func (s *BaseK8SRouter) InitBaseK8SRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/baseRes/k8s")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.BaseResApiGroup.BaseK8SApi

	{
		dataRouter.POST("getBaseK8SList", vpcApi.GetBaseK8SList)                 // 获取集群列表
		dataRouter.POST("importBaseCluster", vpcApi.ImportBaseCluster)           // 导入集群
		dataRouter.POST("deleteBaseCluster", vpcApi.DeleteBaseCluster)           // 删除集群
		dataRouter.POST("deleteBaseClusterByIds", vpcApi.DeleteBaseClusterByIds) // 删除集群
		dataRouter.POST("updateBaseCluster", vpcApi.UpdateBaseCluster)           // 修改集群端口号
	}
}
