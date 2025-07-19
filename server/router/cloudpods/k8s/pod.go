package k8s

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type K8sPodRouter struct {
}

func (*K8sPodRouter) InitK8sPodRouter(Router *gin.RouterGroup) {
	podGroup := Router.Group("cloud/k8s/pod")
	podApi := v1.ApiGroupApp.CloudApiGroup.K8sAPiGroup.PodApi
	{
		podGroup.GET("listNameSpace", podApi.ListNameSpace)
		podGroup.GET("listPods", podApi.ListPods)
	}
}
