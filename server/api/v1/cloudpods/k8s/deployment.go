package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	data "ygpt/server/model/yunguan/k8s"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/service"
	"ygpt/server/utils"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
	"yunion.io/x/onecloud/pkg/mcclient/modules/webconsole"
)

type DeploymentAPI struct {
}

type containerCreateReq struct {
	Name        string               `json:"name"`
	Cluster     string               `json:"cluster"`
	Namespace   string               `json:"namespace"`
	Replicas    int                  `json:"replicas"`
	Labels      jsonutils.JSONObject `json:"labels"`
	Annotations jsonutils.JSONObject `json:"annotations"`
	Template    Template             `json:"template"`
	Service     Service              `json:"service"`
}

type Template struct {
	Spec Spec `json:"spec"`
}

type Volume struct {
	Name                  string `json:"name"`
	PersistentVolumeClaim struct {
		ClaimName string `json:"claimName"`
	} `json:"persistentVolumeClaim"`
}

type Spec struct {
	Containers    []Container `json:"containers"`
	Volumes       []Volume    `json:"volumes"`
	RestartPolicy string      `json:"restartPolicy"`
}
type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}
type Container struct {
	Name            string          `json:"name"`
	Image           string          `json:"image"`
	SecurityContext SecurityContext `json:"securityContext"`
	VolumeMounts    []VolumeMount   `json:"volumeMounts"`
	Resources       Resources       `json:"resources"`
	Command         []string        `json:"command"`
	Args            []string        `json:"args"`
}

type SecurityContext struct {
	Privileged bool `json:"privileged"`
	RunAsUser  int  `json:"runAsUser"`
	RunAsGroup int  `json:"runAsGroup"`
	FsGroup    int  `json:"fsGroup"`
}

type Resources struct {
	Requests Requests `json:"requests"`
	Limits   Limits   `json:"limits"`
}

type Requests struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	EphemeralStorage string `json:"ephemeral-storage"`
	NvidiaGpu        string `json:"nvidia.com/gpu"`
}

type Limits struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	EphemeralStorage string `json:"ephemeral-storage"`
	NvidiaGpu        string `json:"nvidia.com/gpu"`
}

type Service struct {
	IsExternal   bool          `json:"isExternal"`
	PortMappings []PortMapping `json:"portMappings"`
	Type         string        `json:"type"`
}

type PortMapping struct {
	Port       int `json:"port"`       //服务端口，服务暴露给外部访问的端口
	TargetPort int `json:"targetPort"` //目标端口，Pod内部容器的端口
	//NodePort   int    `json:"nodePort"`   //Node端口，每个Node转发的端口，区间范围 30000-32767
	Protocol string `json:"protocol"` //协议类型，TCP/UDP
}

type ContainerEvs struct {
	ContainerNas   string `json:"containerNas"`   //容器Nas名称，新建云硬盘时需要传
	Size           string `json:"size"`           //大小，新建云硬盘时需要传
	ContainerEvsId string `json:"containerEvsId"` //容器云硬盘Id，选择已有的云硬盘时需要传
	MountPath      string `json:"mountPath"`      //挂载点，必传
}

type deploymentCreateReq struct {
	Name         string         `json:"name"`         //名称
	ClusterID    string         `json:"clusterId"`    //集群ID
	ImageUrl     string         `json:"imageUrl"`     //容器镜像地址，例如10.10.100.202/10516/nginx:lastest
	Cpu          int            `json:"cpu"`          //单位:核
	Memory       int            `json:"memory"`       //单位:GB
	Gpu          int            `json:"gpu"`          //单位:卡
	Vendors      string         `json:"vendors"`      //gpu厂商，默认为nvidia
	GpuModel     string         `json:"gpuModel"`     //gpu型号
	PortMappings []PortMapping  `json:"portMappings"` //端口映射
	RootPasswd   string         `json:"rootPasswd"`   //root密码
	ContainerEvs []ContainerEvs `json:"containerEvs"` //容器云硬盘
}
type monitorDataReq struct {
	Interval    int    `json:"interval"`
	Duration    int    `json:"duration"`
	K8sName     string `json:"k8sName"`
	Name        string `json:"name"`
	MonitorType string `json:"monitorType"`
}

// 监控容器
func (*DeploymentAPI) DeploymentMonitor(c *gin.Context) {
	var r monitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	var k8sData data.K8sModel
	var queryData string
	global.GVA_DB.Where("name = ?", r.K8sName).First(&k8sData)
	switch r.MonitorType {
	//cpu使用率
	case "cpu":
		queryData = "container_cpu_usage_seconds_total"
		break
	//内存使用率
	case "memory":
		queryData = "container_spec_memory_limit_bytes"
		break
	//网络接收速率
	case "networkRecevice":
		queryData = "container_network_receive_bytes_total"
		break

	}
	query := fmt.Sprintf(`sum by (namespace, pod, container) (rate(%s{namespace="ygpt", container="%s"}[%dh]))`, queryData, r.Name, r.Duration)
	encodedQuery := url.QueryEscape(query)
	now := time.Now()
	stepData := r.Interval * 60
	// 十分钟前
	tenMinutesAgo := now.Add(-time.Duration(r.Duration) * time.Hour)
	url := fmt.Sprintf("http://%s:%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d", k8sData.ApiIP, k8sData.Endpoint,
		encodedQuery, tenMinutesAgo.Unix(), now.Unix(), stepData)
	fmt.Println("url:" + url)
	// 发起 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		response.FailWithMessage("监控错误", c)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage("监控错误", c)
		return
	}
	fmt.Println(string(body))
	var data interface{}

	// 解析 JSON 到 data
	err = json.Unmarshal(body, &data)
	if err != nil {
		response.FailWithMessage("监控错误", c)
		return
	}
	response.OkWithDetailed(data, "监控成功", c)
}
func (*DeploymentAPI) DeploymentCreate(c *gin.Context) {
	var deployment jsonutils.JSONObject
	var r deploymentCreateReq
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

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("cluster", jsonutils.NewString(r.ClusterID))
	params.Set("limit", jsonutils.NewInt(0))

	namespaceId := utils.GetNamespaceId(r.ClusterID)
	if namespaceId == "" {
		response.FailWithMessage("未找到云管平台命名空间", c)
		return
	}

	var cr containerCreateReq
	var container Container

	cr.Name = r.Name
	cr.Labels = jsonutils.NewDict()
	cr.Annotations = jsonutils.NewDict()
	cr.Cluster = r.ClusterID
	cr.Namespace = namespaceId
	cr.Replicas = 1
	cr.Service.IsExternal = false
	cr.Service.PortMappings = []PortMapping{
		{
			Port:       22,
			TargetPort: 22,
			Protocol:   "TCP",
		},
	}
	cr.Service.Type = "NodePort"

	container.Name = r.Name
	container.Image = r.ImageUrl
	container.Resources.Requests.CPU = strconv.Itoa(r.Cpu)
	container.Resources.Requests.Memory = strconv.Itoa(r.Memory) + "Gi"
	container.Resources.Limits.CPU = strconv.Itoa(r.Cpu)
	container.Resources.Limits.Memory = strconv.Itoa(r.Memory) + "Gi"

	switch strings.ToUpper(r.Vendors) {
	case "NVIDIA":
		container.Resources.Requests.NvidiaGpu = strconv.Itoa(r.Gpu)
		container.Resources.Limits.NvidiaGpu = strconv.Itoa(r.Gpu)
	case "":
		break
	default:
		response.FailWithMessage("暂不支持该gpu厂商", c)
		return
	}

	container.SecurityContext.Privileged = false
	container.Command = []string{"/bin/sh", "-c"}

	if r.RootPasswd != "" {
		container.Args = []string{
			fmt.Sprintf("echo 'root:%s' | chpasswd && (service ssh start || true) && (systemctl start ssh || true) && (service agent start || true) && (systemctl start agent || true) && (sleep infinity || true)", r.RootPasswd),
		}
	} else {
		container.Args = []string{
			fmt.Sprintf("echo 'root:%s' | chpasswd && (service ssh start || true) && (systemctl start ssh || true) && (service agent start || true) && (systemctl start agent || true) && (sleep infinity || true)", "9ijnBHU*@123"),
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(r.ContainerEvs)) // 创建错误通道

	for _, containerEvs := range r.ContainerEvs {
		if containerEvs.MountPath == "" {
			response.FailWithMessage("参数错误", c)
			return
		}

		wg.Add(1)
		go func(containerEvs ContainerEvs) {
			defer wg.Done()
			if containerEvs.ContainerEvsId != "" {
				getParams := jsonutils.NewDict()
				getParams.Set("cluster", jsonutils.NewString(r.ClusterID))
				getParams.Set("namespace", jsonutils.NewString(namespaceId))
				persistentVolumeClaim, err := k8s.PersistentVolumeClaims.GetById(s, containerEvs.ContainerEvsId, getParams)
				if err != nil {
					errChan <- err // 发送错误到通道
					return
				}
				status, _ := persistentVolumeClaim.GetString("status")
				if status != "Bound" {
					errChan <- fmt.Errorf("云硬盘状态异常: %s", containerEvs.ContainerEvsId)
					return
				}

				var volume Volume
				volume.Name = containerEvs.ContainerEvsId
				volume.PersistentVolumeClaim.ClaimName = containerEvs.ContainerEvsId

				var volumeMount VolumeMount
				volumeMount.Name = containerEvs.ContainerEvsId
				volumeMount.MountPath = containerEvs.MountPath

				cr.Template.Spec.Volumes = append(cr.Template.Spec.Volumes, volume)
				container.VolumeMounts = append(container.VolumeMounts, volumeMount)

				return
			}

			if containerEvs.Size != "" {
				uuidStr := uuid.NewString()

				createParams := jsonutils.NewDict()
				createParams.Set("cluster", jsonutils.NewString(r.ClusterID))
				createParams.Set("namespace", jsonutils.NewString(namespaceId))
				if len(uuidStr) >= 8 {
					createParams.Set("name", jsonutils.NewString(r.Name+"-"+uuidStr[:7]))
				} else {
					createParams.Set("name", jsonutils.NewString(r.Name+"-"+uuidStr))
				}

				createParams.Set("size", jsonutils.NewString(containerEvs.Size))
				createParams.Set("storageClass", jsonutils.NewString(containerEvs.ContainerNas))

				persistentVolumeClaim, err := k8s.PersistentVolumeClaims.Create(s, createParams)
				if err != nil {
					errChan <- err // 发送错误到通道
					return
				}

				id, err := persistentVolumeClaim.GetString("id")
				if err != nil {
					errChan <- err // 发送错误到通道
					return
				}

				status, _ := persistentVolumeClaim.GetString("status")
				if status != "Bound" {
					getParams := jsonutils.NewDict()
					getParams.Set("cluster", jsonutils.NewString(r.ClusterID))
					getParams.Set("namespace", jsonutils.NewString(namespaceId))

					size := 10
					for i := 0; i < size; i++ {
						if i < size-1 {
							time.Sleep(500 * time.Millisecond)
						}
						persistentVolumeClaim, err = k8s.PersistentVolumeClaims.GetById(s, id, getParams)
						if err != nil {
							errChan <- err // 发送错误到通道
							return
						}
						status, _ = persistentVolumeClaim.GetString("status")
						if status == "Bound" {
							break
						}
					}
				}

				if status != "Bound" {
					errChan <- fmt.Errorf("云硬盘创建错误: %s", containerEvs.Size)
					return
				}

				var volume Volume
				volume.Name = id
				volume.PersistentVolumeClaim.ClaimName = id

				var volumeMount VolumeMount
				volumeMount.Name = id
				volumeMount.MountPath = containerEvs.MountPath

				cr.Template.Spec.Volumes = append(cr.Template.Spec.Volumes, volume)
				container.VolumeMounts = append(container.VolumeMounts, volumeMount)
			}
		}(containerEvs)
	}

	wg.Wait()      // 等待所有协程完成
	close(errChan) // 关闭错误通道

	// 检查是否有错误发生
	for err := range errChan {
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	cr.Template.Spec.Containers = append(cr.Template.Spec.Containers, container)
	cr.Template.Spec.RestartPolicy = "Always"

	params1 := jsonutils.Marshal(cr)
	deployment, err = k8s.Deployments.Create(s, params1)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(deployment).PrettyString(), "创建成功", c)
}

type ContainerConnectReq struct {
	Id string `json:"id" form:"id"`
}

func (*DeploymentAPI) GetWebConsoleSSH(c *gin.Context) {
	var r ContainerConnectReq

	var result jsonutils.JSONObject

	err := c.ShouldBindQuery(&r)
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

type ContainerDeleteReq struct {
	ID string `json:"id"`
}

func (*DeploymentAPI) DeploymentDelete(c *gin.Context) {
	var r ContainerDeleteReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	deployment, err := k8s.Deployments.Get(s, r.ID, nil)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	clusterId, err := deployment.GetString("cluster_id")
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	namespace, err := deployment.GetString("namespace_id")
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	name, err := deployment.GetString("name")
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("cluster", jsonutils.NewString(clusterId))
	params.Set("namespace", jsonutils.NewString("ygpt"))
	params.Set("owner_name", jsonutils.NewString(name))
	params.Set("owner_kind", jsonutils.NewString("Deployment"))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	services, err := k8s.Services.List(s, params)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}

	// 删除deployment的服务
	for _, item := range services.Data {
		serviceId, err := item.GetString("id")
		if err != nil {
			response.FailWithMessage("删除容器失败:"+err.Error(), c)
			return
		}

		params1 := jsonutils.NewDict()
		params1.Set("cluster", jsonutils.NewString(clusterId))
		params1.Set("namespace", jsonutils.NewString(namespace))
		params1.Set("id", jsonutils.NewString(serviceId))

		_, err = k8s.Services.Delete(s, serviceId, params)
		if err != nil {
			response.FailWithMessage("删除容器失败:"+err.Error(), c)
			return
		}
	}

	body := jsonutils.NewDict()
	body.Set("cluster", jsonutils.NewString(clusterId))
	body.Set("namespace", jsonutils.NewString(namespace))
	body.Set("id", jsonutils.NewString(r.ID))

	_, err = k8s.Deployments.Delete(s, r.ID, body)
	if err != nil {
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}
	var model renter.RenterRes
	err = global.GVA_DB.Where("resource_id = ?", r.ID).First(&model).Error
	if err != nil {
		response.OkWithMessage("删除容器成功", c)
		return
	}
	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", r.ID).Updates(&renter.RenterRes{
		EndTime: time.Now(),
	}).Error
	if err != nil {
		response.OkWithMessage("删除容器成功", c)
		return
	}
	response.OkWithMessage("删除容器成功", c)
}

type deploymentListReq struct {
	ClusterID string `json:"cluster" form:"clusterID"`
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	TimeDesc  string `json:"timeDesc"`
}

func (*DeploymentAPI) DeploymentList(c *gin.Context) {
	var r deploymentListReq
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

	namespaceId := utils.GetNamespaceId(r.ClusterID)
	if namespaceId == "" {
		response.FailWithMessage("未找到云管平台命名空间", c)
		return
	}

	params := jsonutils.NewDict()
	params.Set("namespace", jsonutils.NewString(namespaceId))
	params.Set("cluster", jsonutils.NewString(r.ClusterID))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	if r.TimeDesc != "" {
		params.Set("order_by", jsonutils.NewString("creationTimestamp"))
		params.Set("order", jsonutils.NewString(r.TimeDesc))
	}
	result, err := k8s.Deployments.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	renterIds := make([]int, 0)
	list := jsonutils.Marshal(result).Interface()
	rawMap, _ := list.(map[string]interface{})
	if arr, ok := rawMap["data"].([]interface{}); ok {
		for _, item := range arr {
			if obj, i := item.(map[string]interface{}); i {
				if id, j := obj["id"].(string); j {
					fmt.Println("id:", id)
					serviceReturn, _ := service.ServiceGroupApp.YunGuanServiceGroup.RenterServiceGroup.RenterResService.QueryRenterResInfo(renter.GetRenterResReqByResourceID{ResourceID: id})
					renterIds = append(renterIds, serviceReturn.RenterID)
				}
			}
		}
	}
	for i, item := range rawMap["data"].([]interface{}) {
		if obj, ok := item.(map[string]interface{}); ok {
			// 将 renterId 添加到当前对象中
			obj["renterId"] = renterIds[i]
		} else {
			fmt.Printf("Item at index %d is not a map[string]interface{}\n", i)
		}
	}
	response.OkWithDetailed(rawMap, "获取成功", c)
}

type GetDeploymentPodsReq struct {
	Scope          string `json:"scope" form:"scope"`
	ShowFailReason bool   `json:"show_fail_reason" form:"show_fail_reason"`
	OwnerKind      string `json:"owner_kind" form:"owner_kind"`
	OwnerName      string `json:"owner_name" form:"owner_name"`
	Namespace      string `json:"namespace" form:"namespace"`
	Cluster        string `json:"cluster" form:"cluster"`
	Details        bool   `json:"details" form:"details"`
	Limit          int    `json:"limit" form:"limit"`
	Page           int    `json:"page" form:"page"`
	PageSize       int    `json:"pageSize" form:"pageSize"`
	Offset         int    `json:"offset" form:"offset"`
}

// GetDeploymentPods 获取Deployment关联的Pods
func (*DeploymentAPI) GetDeploymentPods(c *gin.Context) {
	var req GetDeploymentPodsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	req.Limit = req.PageSize
	req.Offset = (req.Page - 1) * req.PageSize
	result, err := k8s.Pods.List(s, jsonutils.Marshal(&req))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithData(jsonutils.Marshal(result).Interface(), c)
}

type GetDeploymentSshInfoReq struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}

func (*DeploymentAPI) GetDeploymentSshInfo(c *gin.Context) {
	var req GetDeploymentSshInfoReq
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

	param := jsonutils.NewDict()
	param.Set("scope", jsonutils.NewString("system"))
	param.Set("show_fail_reason", jsonutils.NewBool(true))
	param.Set("owner_kind", jsonutils.NewString("Deployment"))
	param.Set("owner_name", jsonutils.NewString(req.Name))
	param.Set("namespace", jsonutils.NewString("ygpt"))
	param.Set("cluster", jsonutils.NewString(req.Cluster))
	param.Set("details", jsonutils.NewBool(true))

	pods, err := k8s.Pods.List(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	if pods.Data == nil || len(pods.Data) == 0 {
		response.FailWithMessage("未找到容器", c)
		return
	}

	pod := pods.Data[0]
	nodeNameInterface := k8s.Pods.Get_Node(pod)
	nodeName, ok := nodeNameInterface.(string)
	if !ok || nodeName == "" {
		response.FailWithMessage("未找到容器所在节点或类型断言失败", c)
		return
	}

	containers, err := pod.GetArray("containers")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if containers == nil || len(containers) == 0 {
		response.FailWithMessage("未找到容器", c)
		return
	}

	container := containers[0]

	args, err := container.GetArray("args")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if args == nil || len(args) == 0 {
		response.FailWithMessage("未找到容器启动参数", c)
		return
	}

	arg := args[0]
	username, password, err := extractUsernameAndPassword(arg.String())
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	services, err := k8s.Services.List(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	if services.Data == nil || len(services.Data) == 0 {
		response.FailWithMessage("未找到容器关联的service", c)
		return
	}

	var nodePort int64
	for _, item := range services.Data {
		serviceType := k8s.Services.Get_Type(item)
		if serviceType != "NodePort" {
			continue
		}

		portObjs, _ := item.GetArray("internalEndpoint", "ports")
		if len(portObjs) != 0 {
			for _, obj := range portObjs {
				port, err := obj.Int("port")
				if err != nil {
					continue
				}
				protocol, err := obj.GetString("protocol")
				if err != nil {
					continue
				}

				if port == 22 && protocol == "TCP" {
					nodePort, err = obj.Int("nodePort")
					if err != nil {
						continue
					}
				}
			}
		}
	}

	node, err := k8s.K8sNodes.GetByName(s, nodeName, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	annotations, err := node.Get("annotations")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	ipv4Address, err := annotations.GetString("projectcalico.org/IPv4Address")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	// 提取纯IP地址
	ip := ""
	if ipv4Address != "" {
		parts := strings.Split(ipv4Address, "/")
		if len(parts) > 0 {
			ip = parts[0]
		}
	}

	if ip == "" {
		response.FailWithMessage("未找到容器所在节点有效的IPv4地址", c)
		return
	}

	response.OkWithData(map[string]interface{}{
		"ip":       ip,
		"port":     nodePort,
		"username": username,
		"password": password,
	}, c)
}

func extractUsernameAndPassword(command string) (string, string, error) {
	// 定义正则表达式模式
	pattern := `echo '([^:]+):([^']+)'\s*\| chpasswd`
	re := regexp.MustCompile(pattern)

	// 查找匹配项
	matches := re.FindStringSubmatch(command)
	if matches == nil {
		return "", "", fmt.Errorf("未找到匹配的账户和密码")
	}

	// 提取账户和密码
	username := matches[1]
	password := matches[2]

	return username, password, nil
}

// 获取集群相关节点
type GetK8SNodeReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	ClusterID string `json:"clusterID"`
}
type ModK8SNodeStatusReq struct {
	Cluster string `json:"cluster"`
	ID      string `json:"id"`
	Status  bool   `json:"status"`
}

// 获取集群相关节点列表
func (*DeploymentAPI) GetK8sNodeList(c *gin.Context) {
	var req GetK8SNodeReq
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

	param := jsonutils.NewDict()
	param.Set("scope", jsonutils.NewString("system"))
	param.Set("show_fail_reason", jsonutils.NewBool(true))
	param.Set("details", jsonutils.NewBool(true))
	param.Set("cluster", jsonutils.NewString(req.ClusterID))
	if req.PageSize != 0 {
		param.Set("limit", jsonutils.NewInt(req.PageSize))
		param.Set("offset", jsonutils.NewInt(req.PageSize*(req.PageIndex-1)))
	}

	pods, err := k8s.K8sNodes.List(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(pods).Interface(), "获取成功", c)
}

// 修改集群相关节点调度状态
func (*DeploymentAPI) ModK8sNodeStatus(c *gin.Context) {
	var req ModK8SNodeStatusReq
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

	var pods jsonutils.JSONObject
	param := jsonutils.NewDict()
	param.Set("cluster", jsonutils.NewString(req.Cluster))
	if req.Status {
		pods, err = k8s.K8sNodes.PerformAction(s, req.ID, "uncordon", param)
	} else {
		pods, err = k8s.K8sNodes.PerformAction(s, req.ID, "cordon", param)
	}

	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(pods).PrettyString(), "修改成功", c)
}
