package k8s

import (
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type PodAPI struct {
}

type ListNameSpaceReq struct {
	Cluster string `json:"cluster" form:"cluster"`
	Limit   int    `json:"limit" form:"limit"`
	Scope   string `json:"scope" form:"scope"`
}

type Namespace struct {
	APIVersion         string            `json:"apiVersion"`
	CanDelete          bool              `json:"can_delete"`
	CanUpdate          bool              `json:"can_update"`
	Cluster            string            `json:"cluster"`
	ClusterID          string            `json:"clusterID"`
	ClusterIDAlt       string            `json:"cluster_id"`
	CreatedAt          string            `json:"created_at"`
	CreationTimestamp  string            `json:"creationTimestamp"`
	Deleted            bool              `json:"deleted"`
	Distribution       string            `json:"distribution"`
	DomainID           string            `json:"domain_id"`
	DomainSrc          string            `json:"domain_src"`
	ExternalID         string            `json:"external_id"`
	Generation         int               `json:"generation"`
	ID                 string            `json:"id"`
	ImportedAt         string            `json:"imported_at"`
	IsEmulated         bool              `json:"is_emulated"`
	Kind               string            `json:"kind"`
	Labels             map[string]string `json:"labels"`
	Name               string            `json:"name"`
	Progress           int               `json:"progress"`
	ProjectDomain      string            `json:"project_domain"`
	ResourceVersion    string            `json:"resourceVersion"`
	ResourceVersionAlt string            `json:"resource_version"`
	Source             string            `json:"source"`
	Status             string            `json:"status"`
	UpdateVersion      int               `json:"update_version"`
	UpdatedAt          string            `json:"updated_at"`
}

// 获取集群的NameSpace
func (p *PodAPI) ListNameSpace(c *gin.Context) {
	var req ListNameSpaceReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	result, err := k8s.Namespaces.List(s,
		jsonutils.Marshal(&req))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	var listNameSpace []Namespace
	for _, item := range result.Data {
		var nameSpace Namespace
		err := json.Unmarshal([]byte(item.String()), &nameSpace)
		if err != nil {
			continue
		}
		if nameSpace.ClusterID == req.Cluster {
			listNameSpace = append(listNameSpace, nameSpace)
		}
	}
	response.OkWithData(&gin.H{
		"list":  listNameSpace,
		"total": result.Total,
	}, c)
}

type ListPodsReq struct {
	Scope          string `json:"scope" form:"scope"`
	ShowFailReason bool   `json:"show_fail_reason" form:"show_fail_reason:"`
	Detail         bool   `json:"detail" form:"detail"`
	Cluster        string `json:"cluster" form:"cluster"`
	Namespace      string `json:"namespace" form:"namespace"`
	Limit          int    `json:"limit" form:"limit"`
}

func (*PodAPI) ListPods(c *gin.Context) {
	var req ListPodsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 创建会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 查询POds
	result, err := k8s.Pods.List(s, jsonutils.Marshal(&req))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 解析返回结果
	var list []string
	for _, item := range result.Data {
		list = append(list, item.PrettyString())
	}
	response.OkWithData(&gin.H{
		"list":  list,
		"total": result.Total,
	}, c)
}

type GetPodInfoReq struct {
	PodId          string `json:"pod_id" form:"pod_id"`
	Scope          string `json:"scope" form:"scope"`
	ShowFailReason bool   `json:"show_fail_reason" form:"show_fail_reason:"`
	Detail         bool   `json:"detail" form:"detail"`
	Cluster        string `json:"cluster" form:"cluster"`
	NameSpace      string `json:"namespace" form:"namespace"`
}

// 获取Pod信息
func (*PodAPI) GetPodInfo(c *gin.Context) {
	var req GetPodInfoReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 创建会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 查询PodInfo
	result, err := k8s.Pods.Get(s, req.PodId, jsonutils.Marshal(&req))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 响应
	response.OkWithData(result.PrettyString(), c)
}

type GetPodYamlCOnfigReq struct {
	PodId string `json:"pod_id" form:"pod_id"`
}

// 获取Pod的yaml配置
func (*PodAPI) GetPodYamlConfigInfo(c *gin.Context) {
	var req GetPodInfoReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 创建会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	result, err := k8s.Pods.GetRaw(s, req.PodId, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithData(result.PrettyString(), c)
}

type UpdatePodRawInfoReq struct {
	PodId     string `json:"pod_id" form:"pod_id"`
	RawString string `json:"raw_string" form:"raw_string"`
	Cluster   string `json:"cluster" form:"cluster"`
	NameSpace string `json:"namespace" form:"namespace"`
}

func (*PodAPI) UpdatePodRawYamlInfo(c *gin.Context) {
	var req UpdatePodRawInfoReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 创建会话
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 更新Pod
	query := jsonutils.NewDict()
	query.Add(jsonutils.NewString(req.Cluster), "cluster")
	query.Add(jsonutils.NewString(req.NameSpace), "namespace")
	rawObj, err := jsonutils.ParseString(req.RawString)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	dict := jsonutils.NewDict()
	dict.UpdateDefault(rawObj)
	result, err := k8s.Pods.UpdateRaw(s, req.PodId, query, dict)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithData(result.PrettyString(), c)
}

type PodMonitorReq struct {
	MetricQuery []struct {
		Model struct {
			Measurement string `json:"measurement"`
			Select      [][]struct {
				Type   string   `json:"type"`
				Params []string `json:"params"`
			} `json:"select"`
			Tags []struct {
				Key      string `json:"key"`
				Value    string `json:"value"`
				Operator string `json:"operator"`
			} `json:"tags"`
			GroupBy []struct {
				Type   string   `json:"type"`
				Params []string `json:"params"`
			} `json:"group_by"`
		} `json:"model"`
	} `json:"metric_query"`
	Scope           string `json:"scope"`
	From            string `json:"from"`
	Interval        string `json:"interval"`
	Unit            bool   `json:"unit"`
	SkipCheckSeries bool   `json:"skip_check_series"`
	Signature       string `json:"signature"`
}

// 监控容器组的信息
func (*PodAPI) PodMonitor() {
	//TODO 获取容器组的监控信息
}

type GetPodEventsReq struct {
	Scope          string `json:"scope" form:"scope"`
	ShowFailReason bool   `json:"show_fail_reason" form:"show_fail_reason:"`
	OwnerName      string `json:"owner_name" form:"owner_name"`
	OwnerKind      string `json:"owner_kind"`
}

// 获取容器组的事件
func (*PodAPI) GetPodEvents(c *gin.Context) {
	var req GetPodEventsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return

	}
	result, err := k8s.Events.List(s, jsonutils.Marshal(&req))
	if err != nil {
		response.FailWithMessage(err, c)
	}
	var events []string
	for _, item := range result.Data {
		events = append(events, item.PrettyString())
	}
	response.OkWithData(&gin.H{
		"list":  events,
		"total": result.Total,
	}, c)
}

//type containerCreateReq struct {
//	Name        string               `json:"name"`
//	Cluster     string               `json:"cluster"`
//	Namespace   string               `json:"namespace"`
//	Replicas    int                  `json:"replicas"`
//	Labels      jsonutils.JSONObject `json:"labels"`
//	Annotations jsonutils.JSONObject `json:"annotations"`
//	Template    Template             `json:"template"`
//	Service     Service              `json:"service"`
//}
//
//type Template struct {
//	Spec Spec `json:"spec"`
//}
//
//type Spec struct {
//	Containers    []Container `json:"containers"`
//	RestartPolicy string      `json:"restartPolicy"`
//}
//
//type Container struct {
//	Name            string          `json:"name"`
//	Image           string          `json:"image"`
//	SecurityContext SecurityContext `json:"securityContext"`
//	Resources       Resources       `json:"resources"`
//}
//
//type SecurityContext struct {
//	Privileged bool `json:"privileged"`
//}
//
//type Resources struct {
//	Requests Requests `json:"requests"`
//	Limits   Limits   `json:"limits"`
//}
//
//type Requests struct {
//	CPU              string `json:"cpu"`
//	Memory           string `json:"memory"`
//	EphemeralStorage string `json:"ephemeral-storage"`
//	NvidiaGpu        string `json:"nvidia.com/gpu"`
//}
//
//type Limits struct {
//	CPU              string `json:"cpu"`
//	Memory           string `json:"memory"`
//	EphemeralStorage string `json:"ephemeral-storage"`
//	NvidiaGpu        string `json:"nvidia.com/gpu"`
//}
//
//type Service struct {
//	IsExternal   bool          `json:"isExternal"`
//	PortMappings []PortMapping `json:"portMappings"`
//	Type         string        `json:"type"`
//}
//
//type PortMapping struct {
//	Port       int    `json:"port"`       //服务端口，服务暴露给外部访问的端口
//	TargetPort int    `json:"targetPort"` //目标端口，Pod内部容器的端口
//	NodePort   int    `json:"nodePort"`   //Node端口，每个Node转发的端口，区间范围 30000-32767
//	Protocol   string `json:"protocol"`   //协议类型，TCP/UDP
//}
//
//type podCreateReq struct {
//	Name         string        `json:"name"`         //名称
//	ClusterID    string        `json:"clusterId"`    //集群ID
//	ImageUrl     string        `json:"imageUrl"`     //容器镜像地址，例如10.10.100.202/10516/nginx:lastest
//	Cpu          int           `json:"cpu"`          //单位:核
//	Memory       int           `json:"memory"`       //单位:GB
//	Gpu          int           `json:"gpu"`          //单位:卡
//	Vendors      string        `json:"vendors"`      //gpu厂商，默认为nvidia
//	GpuModel     string        `json:"gpuModel"`     //gpu型号
//	PortMappings []PortMapping `json:"portMappings"` //端口映射
//}
//
//func (*PodAPI) PodCreate(c *gin.Context) {
//	var deployment jsonutils.JSONObject
//	var r podCreateReq
//	err := c.ShouldBindJSON(&r)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	s, err := global.GetCloudClientSession()
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	params := jsonutils.NewDict()
//	params.Set("scope", jsonutils.NewString("system"))
//	params.Set("cluster", jsonutils.NewString(r.ClusterID))
//	params.Set("limit", jsonutils.NewInt(0))
//
//	namespace, err := k8s.Namespaces.GetByName(s, "default", nil)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	namespaceId, err := namespace.GetString("id")
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	var cr containerCreateReq
//
//	var container Container
//
//	//uniqueName := global.GenerateUniqueName()
//	cr.Name = r.Name
//	cr.Labels = jsonutils.NewDict()
//	cr.Annotations = jsonutils.NewDict()
//	cr.Cluster = r.ClusterID
//	cr.Namespace = namespaceId
//	cr.Replicas = 1
//	cr.Service.IsExternal = false
//	cr.Service.PortMappings = r.PortMappings
//	container.Name = r.Name
//
//	container.Image = r.ImageUrl
//	container.Resources.Requests.CPU = strconv.Itoa(r.Cpu)
//	container.Resources.Requests.Memory = strconv.Itoa(r.Memory) + "Gi"
//
//	container.Resources.Limits.CPU = strconv.Itoa(r.Cpu)
//	container.Resources.Limits.Memory = strconv.Itoa(r.Cpu) + "Gi"
//
//	switch r.Vendors {
//	case "nvidia":
//		container.Resources.Requests.NvidiaGpu = strconv.Itoa(r.Gpu)
//		container.Resources.Limits.NvidiaGpu = strconv.Itoa(r.Gpu)
//	case "":
//		break
//	default:
//		response.FailWithMessage("暂不支持该gpu厂商", c)
//		return
//	}
//
//	container.SecurityContext.Privileged = false
//
//	cr.Template.Spec.Containers = append(cr.Template.Spec.Containers, container)
//	cr.Template.Spec.RestartPolicy = "Always"
//
//	params1 := jsonutils.Marshal(cr)
//	deployment, err = k8s.Deployments.Create(s, params1)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	response.OkWithDetailed(jsonutils.Marshal(deployment).PrettyString(), "创建成功", c)
//}
//
//type ContainerConnectReq struct {
//	Id string `json:"id"`
//}
//
//func (*PodAPI) GetWebConsoleSSH(c *gin.Context) {
//	var r ContainerConnectReq
//
//	var result jsonutils.JSONObject
//
//	err := c.ShouldBindJSON(&r)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	s, err := global.GetCloudClientSession()
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	deployment, err := k8s.Deployments.Get(s, r.Id, nil)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	ownerName, err := deployment.GetString("name")
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	namespace, err := deployment.GetString("namespace")
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	clusterId, err := deployment.GetString("cluster_id")
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//
//	podListParams := jsonutils.NewDict()
//	podListParams.Set("cluster", jsonutils.NewString(clusterId))
//	podListParams.Set("namespace", jsonutils.NewString(namespace))
//	podListParams.Set("details", jsonutils.NewBool(true))
//	podListParams.Set("owner_name", jsonutils.NewString(ownerName))
//	podListParams.Set("scope", jsonutils.NewString("system"))
//	podListParams.Set("show_fail_reason", jsonutils.NewBool(true))
//	podListParams.Set("owner_kind", jsonutils.NewString("Deployment"))
//
//	podsList, err := k8s.Pods.List(s, podListParams)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	for _, item := range podsList.Data {
//		clusterId, err = item.GetString("cluster")
//		if err != nil {
//			continue
//		}
//		namespace, err = item.GetString("namespace")
//		if err != nil {
//			continue
//		}
//		name, err := item.GetString("name")
//		if err != nil {
//			continue
//		}
//		containers, err := item.GetArray("containers")
//		if err != nil {
//			continue
//		}
//		if len(containers) == 0 {
//			continue
//		}
//		container := containers[0]
//
//		containerName, err := container.GetString("name")
//		if err != nil {
//			continue
//		}
//
//		params := jsonutils.NewDict()
//		params.Set("cluster", jsonutils.NewString(clusterId))
//		params.Set("namespace", jsonutils.NewString(namespace))
//		params.Set("name", jsonutils.NewString(name))
//		params.Set("container", jsonutils.NewString(containerName))
//
//		result, err = webconsole.WebConsole.DoK8sShellConnect(s, name, params)
//		if err != nil {
//			continue
//		}
//		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "容器远程连接成功", c)
//		return
//	}
//
//	response.FailWithMessage("容器远程连接失败:not found", c)
//}
//
//type ContainerDeleteReq struct {
//	ID string `json:"id"`
//}
//
//func (*PodAPI) PodDelete(c *gin.Context) {
//	var r ContainerDeleteReq
//	err := c.ShouldBindJSON(&r)
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//
//	s, err := global.GetCloudClientSession()
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//
//	deployment, err := k8s.Deployments.Get(s, r.ID, nil)
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//
//	clusterId, err := deployment.GetString("cluster_id")
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//
//	namespace, err := deployment.GetString("namespace")
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//
//	body := jsonutils.NewDict()
//	body.Set("cluster", jsonutils.NewString(clusterId))
//	body.Set("namespace", jsonutils.NewString(namespace))
//	body.Set("id", jsonutils.NewString(r.ID))
//
//	_, err = k8s.Deployments.Delete(s, r.ID, body)
//	if err != nil {
//		response.FailWithMessage("删除容器失败:"+err.Error(), c)
//		return
//	}
//	response.OkWithMessage("删除容器成功", c)
//}
