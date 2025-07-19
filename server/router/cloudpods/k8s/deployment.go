package k8s

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type K8sDeploymentRouter struct {
}

func (*K8sDeploymentRouter) InitK8sDeploymentRouter(Router *gin.RouterGroup) {
	deploymentGroup := Router.Group("cloud/k8s/deployment")
	deploymentApi := v1.ApiGroupApp.CloudApiGroup.K8sAPiGroup.DeploymentAPI
	{
		deploymentGroup.POST("list", deploymentApi.DeploymentList)                 // 获取deployment列表
		deploymentGroup.GET("webConsole", deploymentApi.GetWebConsoleSSH)          // 远程终端访问
		deploymentGroup.POST("create", deploymentApi.DeploymentCreate)             // 创建deployment
		deploymentGroup.DELETE("delete", deploymentApi.DeploymentDelete)           // 删除deployment
		deploymentGroup.POST("sshInfo", deploymentApi.GetDeploymentSshInfo)        //获取容器的ssh信息
		deploymentGroup.GET("pods", deploymentApi.GetDeploymentPods)               //
		deploymentGroup.POST("getK8sNodeList", deploymentApi.GetK8sNodeList)       //获取集群节点列表
		deploymentGroup.POST("modK8sNodeStatus", deploymentApi.ModK8sNodeStatus)   //修改集群相关节点调度状态
		deploymentGroup.POST("deploymentMonitor", deploymentApi.DeploymentMonitor) //监控容器
	}
}
