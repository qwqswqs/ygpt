package cloudpods

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
	"yunion.io/x/onecloud/pkg/mcclient/modules/webconsole"
)

type ContainerApi struct {
}

type ContainerDeleteReq struct {
	ID string `json:"id"`
}

type ContainerCreateReq struct {
	Name string `json:"name"`
	//Tag    string `json:"tag"`
	Cpu    int `json:"cpu"`
	Memory int `json:"memory"`
}

// 创建容器
func (con *ContainerApi) HostCreate(c *gin.Context) {
	//var r ContainerCreateReq
	//err := c.ShouldBindJSON(&r)
	//if err != nil {
	//	response.FailWithMessage("创建容器失败:"+err.Error(), c)
	//	return
	//}
	//
	//deployment, err := ContainerService.HostCreate(r.Cpu, r.Memory, r.Name)
	//if err != nil {
	//	response.FailWithMessage("创建容器失败:"+err.Error(), c)
	//	return
	//}

	//response.OkWithDetailed(jsonutils.Marshal(deployment).PrettyString(), "获取成功", c)
	response.FailWithMessage("该接口已经隐藏!", c)
}

func (con *ContainerApi) HostDelete(c *gin.Context) {
	var r ContainerDeleteReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	err = ContainerService.HostDelete(r.ID)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除容器成功", c)
}

// 关机
func (con *ContainerApi) HostShutDown(c *gin.Context) {
}

// 开机
func (con *ContainerApi) HostStart(c *gin.Context) {
}

// 重启
func (con *ContainerApi) HostRestart(c *gin.Context) {
}

// 挂起
func (con *ContainerApi) HostSuspend(c *gin.Context) {
}

// 恢复
func (con *ContainerApi) HostResume(c *gin.Context) {
}

type ContainerConnectReq struct {
	Id string `json:"id"`
}

// 容器远程连接
func (con *ContainerApi) GetWebConsoleSSH(c *gin.Context) {
	var r ContainerConnectReq

	var result jsonutils.JSONObject

	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	deployment, err := k8s.Deployments.Get(s, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	ownerName, err := deployment.GetString("name")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	namespace, err := deployment.GetString("namespace")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	clusterId, err := deployment.GetString("cluster_id")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	podListParams := jsonutils.NewDict()
	podListParams.Set("cluster", jsonutils.NewString(clusterId))
	podListParams.Set("namespace", jsonutils.NewString(namespace))
	podListParams.Set("details", jsonutils.NewBool(true))
	podListParams.Set("owner_name", jsonutils.NewString(ownerName))
	podListParams.Set("scope", jsonutils.NewString("system"))
	podListParams.Set("show_fail_reason", jsonutils.NewBool(true))
	podListParams.Set("owner_kind", jsonutils.NewString("Deployment"))

	podsList, err := k8s.Pods.List(s, podListParams)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	for _, item := range podsList.Data {
		clusterId, err = item.GetString("cluster")
		if err != nil {
			continue
		}
		namespace, err = item.GetString("namespace")
		if err != nil {
			continue
		}
		name, err := item.GetString("name")
		if err != nil {
			continue
		}
		containers, err := item.GetArray("containers")
		if err != nil {
			continue
		}
		if len(containers) == 0 {
			continue
		}
		container := containers[0]

		containerName, err := container.GetString("name")
		if err != nil {
			continue
		}

		params := jsonutils.NewDict()
		params.Set("cluster", jsonutils.NewString(clusterId))
		params.Set("namespace", jsonutils.NewString(namespace))
		params.Set("name", jsonutils.NewString(name))
		params.Set("container", jsonutils.NewString(containerName))

		result, err = webconsole.WebConsole.DoK8sShellConnect(s, name, params)
		if err != nil {
			continue
		}
		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "容器远程连接成功", c)
		return
	}

	response.FailWithMessage("容器远程连接失败:not found", c)
}

// 容器监控
func (con *ContainerApi) HostMonitor(c *gin.Context) {
}

// 获取容器列表
func (con *ContainerApi) HostGet(c *gin.Context) {
}

// 创建硬盘快照
func (con *ContainerApi) HostCreateDiskSnap(c *gin.Context) {
}
