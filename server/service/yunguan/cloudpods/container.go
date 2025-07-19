package cloudpods

import (
	"errors"
	"strconv"
	"time"
	"ygpt/server/global"
	"ygpt/server/utils"

	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type ContainerService struct {
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
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	NodePort   int    `json:"nodePort"`
	Protocol   string `json:"protocol"`
}

func (c *ContainerService) HostCreate(imageName string, cpu, memory int, sysDisk int, dataDisk int, gpu int, renterId int) (jsonutils.JSONObject, error) {
	var deployment jsonutils.JSONObject
	var r containerCreateReq
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	clusters, err := k8s.KubeClusters.List(s, params)
	if err != nil {
		return nil, err
	}
	if clusters.Total == 0 || len(clusters.Data) == 0 {
		return nil, errors.New("no running cluster found")
	}
	for _, cluster := range clusters.Data {
		status, err := cluster.GetString("status")
		if err != nil {
			continue
		}
		if status != "running" {
			continue
		}
		clusterId, err := cluster.GetString("id")
		if err != nil {
			continue
		}

		params = jsonutils.NewDict()
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("cluster", jsonutils.NewString(clusterId))

		namespaceId := utils.GetNamespaceId(clusterId)
		if namespaceId == "" {
			return nil, errors.New("no ygpt namespace select")
		}

		var container Container

		uniqueName := global.GenerateUniqueName(renterId)
		r.Name = uniqueName
		r.Labels = jsonutils.NewDict()
		r.Annotations = jsonutils.NewDict()
		r.Cluster = clusterId
		r.Namespace = namespaceId
		r.Replicas = 1
		r.Service.IsExternal = false
		r.Service.PortMappings = []PortMapping{
			{
				Port:       22,
				TargetPort: 22,
				Protocol:   "TCP",
			},
		}
		r.Service.Type = "NodePort"

		container.Name = uniqueName

		imageService := &ImageService{}
		imageUrl, err := imageService.FindDockerImage(imageName)
		if err != nil {
			return nil, err
		}

		container.Image = imageUrl
		container.Resources.Requests.CPU = strconv.Itoa(cpu)
		container.Resources.Requests.Memory = strconv.Itoa(memory) + "Gi"
		//container.Resources.Requests.EphemeralStorage = "2Gi"
		//container.Resources.Requests.NvidiaGpu = "1"

		container.Resources.Limits.CPU = strconv.Itoa(cpu)
		container.Resources.Limits.Memory = strconv.Itoa(memory) + "Gi"
		//container.Resources.Limits.EphemeralStorage = "2Gi"
		//container.Resources.Limits.NvidiaGpu = "1"

		if gpu > 0 {
			container.Resources.Requests.NvidiaGpu = strconv.Itoa(gpu)
			container.Resources.Limits.NvidiaGpu = strconv.Itoa(gpu)
		}

		if sysDisk > 0 || dataDisk > 0 {
			storageClassParams := jsonutils.NewDict()
			storageClassParams.Set("scope", jsonutils.NewString("system"))
			storageClassParams.Set("show_fail_reason", jsonutils.NewBool(true))
			storageClassParams.Set("details", jsonutils.NewBool(true))
			storageClassParams.Set("cluster", jsonutils.NewString(r.Cluster))
			storageClassParams.Set("namespace", jsonutils.NewString(namespaceId))

			storageClasses, err := k8s.Storageclass.List(s, storageClassParams)
			if err != nil {
				return nil, errors.New("未找到容器Nas存储")
			}

			if storageClasses.Total == 0 || len(storageClasses.Data) == 0 {
				return nil, errors.New("未找到容器Nas存储")
			}

			storageClass := storageClasses.Data[0]

			storageClassName, err := storageClass.GetString("name")
			if err != nil {
				return nil, err
			}

			if sysDisk > 0 {
				createParams := jsonutils.NewDict()
				createParams.Set("cluster", jsonutils.NewString(r.Cluster))
				createParams.Set("namespace", jsonutils.NewString(namespaceId))
				createParams.Set("name", jsonutils.NewString(r.Name+"-sys"))
				createParams.Set("size", jsonutils.NewString(strconv.Itoa(sysDisk)+"Gi"))
				createParams.Set("storageClass", jsonutils.NewString(storageClassName))

				persistentVolumeClaim, err := k8s.PersistentVolumeClaims.Create(s, createParams)
				if err != nil {
					return nil, err
				}

				id, err := persistentVolumeClaim.GetString("id")
				if err != nil {
					return nil, err
				}

				status, _ = persistentVolumeClaim.GetString("status")
				if status != "Bound" {
					getParams := jsonutils.NewDict()
					getParams.Set("cluster", jsonutils.NewString(clusterId))
					getParams.Set("namespace", jsonutils.NewString(namespaceId))

					size := 20
					for i := 0; i < size; i++ {
						if i < size-1 {
							time.Sleep(500 * time.Millisecond)
						}
						persistentVolumeClaim, err = k8s.PersistentVolumeClaims.GetById(s, id, getParams)
						if err != nil {
							return nil, err
						}
						status, _ = persistentVolumeClaim.GetString("status")
						if status == "Bound" {
							break
						}
					}
				}

				if status != "Bound" {
					return nil, errors.New("云硬盘创建失败")
				}

				var volume Volume
				volume.Name = id
				volume.PersistentVolumeClaim.ClaimName = id

				var volumeMount VolumeMount
				volumeMount.Name = id
				volumeMount.MountPath = "/data1"

				r.Template.Spec.Volumes = append(r.Template.Spec.Volumes, volume)
				container.VolumeMounts = append(container.VolumeMounts, volumeMount)
			}

			if dataDisk > 0 {
				createParams := jsonutils.NewDict()
				createParams.Set("cluster", jsonutils.NewString(r.Cluster))
				createParams.Set("namespace", jsonutils.NewString(namespaceId))
				createParams.Set("name", jsonutils.NewString(r.Name+"-data"))
				createParams.Set("size", jsonutils.NewString(strconv.Itoa(sysDisk)+"Gi"))
				createParams.Set("storageClass", jsonutils.NewString(storageClassName))

				persistentVolumeClaim, err := k8s.PersistentVolumeClaims.Create(s, createParams)
				if err != nil {
					return nil, err
				}

				id, err := persistentVolumeClaim.GetString("id")
				if err != nil {
					return nil, err
				}

				status, _ = persistentVolumeClaim.GetString("status")
				if status != "Bound" {
					getParams := jsonutils.NewDict()
					getParams.Set("cluster", jsonutils.NewString(clusterId))
					getParams.Set("namespace", jsonutils.NewString(namespaceId))

					size := 20
					for i := 0; i < size; i++ {
						if i < size-1 {
							time.Sleep(500 * time.Millisecond)
						}
						persistentVolumeClaim, err = k8s.PersistentVolumeClaims.GetById(s, id, getParams)
						if err != nil {
							return nil, err
						}
						status, _ = persistentVolumeClaim.GetString("status")
						if status == "Bound" {
							break
						}
					}
				}

				if status != "Bound" {
					return nil, errors.New("云硬盘创建失败")
				}

				var volume Volume
				volume.Name = id
				volume.PersistentVolumeClaim.ClaimName = id

				var volumeMount VolumeMount
				volumeMount.Name = id
				volumeMount.MountPath = "/data2"

				r.Template.Spec.Volumes = append(r.Template.Spec.Volumes, volume)
				container.VolumeMounts = append(container.VolumeMounts, volumeMount)
			}
		}

		container.SecurityContext.Privileged = false

		//container.Command = []string{
		//	"sleep",
		//	"infinity",
		//}

		container.Command = []string{
			"/bin/sh",
			"-c",
		}

		container.Args = []string{
			"(service ssh start || true) && (systemctl start ssh || true) && (service agent start || true) && (systemctl start agent || true) && (sleep infinity || true)",
		}

		r.Template.Spec.Containers = append(r.Template.Spec.Containers, container)
		r.Template.Spec.RestartPolicy = "Always"

		params1 := jsonutils.Marshal(r)
		deployment, err = k8s.Deployments.Create(s, params1)

		if err != nil {
			return nil, err
		}
		return deployment, nil
	}
	return nil, errors.New("running cluster or namespace not found")
}

func (c *ContainerService) HostDelete(id string) error {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return err
	}

	deployment, err := k8s.Deployments.Get(s, id, nil)
	if err != nil {
		return err
	}

	clusterId, err := deployment.GetString("cluster_id")
	if err != nil {
		return err
	}

	namespace, err := deployment.GetString("namespace")
	if err != nil {
		return err
	}

	body := jsonutils.NewDict()
	body.Set("cluster", jsonutils.NewString(clusterId))
	body.Set("namespace", jsonutils.NewString(namespace))
	body.Set("id", jsonutils.NewString(id))

	_, err = k8s.Deployments.Delete(s, id, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerService) HostGet(id string) (jsonutils.JSONObject, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	deployment, err := k8s.Deployments.Get(s, id, nil)
	if err != nil {
		return nil, err
	}

	ownerName, err := deployment.GetString("name")
	if err != nil {
		return nil, err
	}

	namespace, err := deployment.GetString("namespace")
	if err != nil {
		return nil, err
	}

	clusterId, err := deployment.GetString("cluster_id")
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if len(podsList.Data) == 0 {
		return nil, errors.New("no container found")
	}

	pod := jsonutils.NewDict()
	pod.Set("pod", podsList.Data[0])
	pod.Set("id", jsonutils.NewString(id))
	pod.Set("name", jsonutils.NewString(ownerName))

	return pod, nil
}

//func (c *ContainerService) HostListByRenterId(renterId *int) ([]jsonutils.JSONObject, int, error) {
//	containerList := make([]jsonutils.JSONObject, 0)
//
//	s, err := global.GetCloudClientSession()
//	if err != nil {
//		return nil, 0, err
//	}
//
//	params := jsonutils.NewDict()
//	params.Set("scope", jsonutils.NewString("system"))
//	params.Set("show_fail_reason", jsonutils.NewBool(true))
//	params.Set("details", jsonutils.NewBool(true))
//
//	clusters, err := k8s.KubeClusters.List(s, params)
//	if err != nil {
//		return nil, 0, err
//	}
//	if len(clusters.Data) == 0 {
//		return nil, 0, errors.New("no running cluster found")
//	}
//	for _, cluster := range clusters.Data {
//		status, err := cluster.GetString("status")
//		if err != nil {
//			continue
//		}
//		if status != "running" {
//			continue
//		}
//		clusterId, err := cluster.GetString("id")
//		if err != nil {
//			continue
//		}
//
//		params = jsonutils.NewDict()
//		params.Set("scope", jsonutils.NewString("system"))
//		params.Set("cluster", jsonutils.NewString(clusterId))
//
//		namespace, err := k8s.Namespaces.GetByName(s, "ygpt", params)
//		if err != nil {
//			continue
//		}
//
//		namespaceId, err := namespace.GetString("id")
//		if err != nil {
//			continue
//		}
//
//		params.Set("namespace", jsonutils.NewString(namespaceId))
//		params.Set("details", jsonutils.NewBool(true))
//		params.Set("show_fail_reason", jsonutils.NewBool(true))
//
//		deployments, err := k8s.Deployments.List(s, params)
//		if err != nil {
//			continue
//		}
//
//		podListParams := jsonutils.NewDict()
//		podListParams.Set("cluster", jsonutils.NewString(clusterId))
//		podListParams.Set("namespace", jsonutils.NewString(namespaceId))
//		podListParams.Set("details", jsonutils.NewBool(true))
//		podListParams.Set("scope", jsonutils.NewString("system"))
//		podListParams.Set("show_fail_reason", jsonutils.NewBool(true))
//		podListParams.Set("owner_kind", jsonutils.NewString("Deployment"))
//
//		for _, deployment := range deployments.Data {
//			ownerName, err := deployment.GetString("name")
//			if err != nil {
//				continue
//			}
//
//			id, err := deployment.GetString("id")
//			if err != nil {
//				continue
//			}
//
//			podListParams.Set("owner_name", jsonutils.NewString(ownerName))
//
//			podsList, err := k8s.Pods.List(s, podListParams)
//			if err != nil {
//				continue
//			}
//			if len(podsList.Data) == 0 {
//				continue
//			}
//
//			pod := jsonutils.NewDict()
//			pod.Set("pod", podsList.Data[0])
//			pod.Set("id", jsonutils.NewString(id))
//
//			containerList = append(containerList, pod)
//		}
//	}
//
//	return containerList, len(containerList), nil
//}

func (c *ContainerService) HostListByIds(ids []string) ([]jsonutils.JSONObject, int, error) {
	containerList := make([]jsonutils.JSONObject, 0)

	for _, id := range ids {
		container, err := c.HostGet(id)
		if err != nil {
			continue
		}

		containerList = append(containerList, container)
	}

	return containerList, len(containerList), nil
}

func (c *ContainerService) HostList() ([]jsonutils.JSONObject, int, error) {
	containerList := make([]jsonutils.JSONObject, 0)

	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	clusters, err := k8s.KubeClusters.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	if len(clusters.Data) == 0 {
		return nil, 0, nil
	}
	for _, cluster := range clusters.Data {
		status, err := cluster.GetString("status")
		if err != nil {
			continue
		}
		if status != "running" {
			continue
		}
		clusterId, err := cluster.GetString("id")
		if err != nil {
			continue
		}

		params = jsonutils.NewDict()
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("cluster", jsonutils.NewString(clusterId))
		params.Set("limit", jsonutils.NewInt(0))

		namespace, err := k8s.Namespaces.GetByName(s, "ygpt", params)
		if err != nil {
			continue
		}

		namespaceId, err := namespace.GetString("id")
		if err != nil {
			continue
		}

		params.Set("namespace", jsonutils.NewString(namespaceId))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("show_fail_reason", jsonutils.NewBool(true))

		deployments, err := k8s.Deployments.List(s, params)
		if err != nil {
			continue
		}

		podListParams := jsonutils.NewDict()
		podListParams.Set("cluster", jsonutils.NewString(clusterId))
		podListParams.Set("namespace", jsonutils.NewString(namespaceId))
		podListParams.Set("details", jsonutils.NewBool(true))
		podListParams.Set("scope", jsonutils.NewString("system"))
		podListParams.Set("show_fail_reason", jsonutils.NewBool(true))
		podListParams.Set("owner_kind", jsonutils.NewString("Deployment"))

		for _, deployment := range deployments.Data {
			ownerName, err := deployment.GetString("name")
			if err != nil {
				continue
			}

			id, err := deployment.GetString("id")
			if err != nil {
				continue
			}

			podListParams.Set("owner_name", jsonutils.NewString(ownerName))

			podsList, err := k8s.Pods.List(s, podListParams)
			if err != nil {
				continue
			}
			if len(podsList.Data) == 0 {
				continue
			}

			pod := jsonutils.NewDict()
			pod.Set("pod", podsList.Data[0])
			pod.Set("id", jsonutils.NewString(id))
			pod.Set("name", jsonutils.NewString(ownerName))

			containerList = append(containerList, pod)
		}
	}

	return containerList, len(containerList), nil
}

func (c *ContainerService) HostListByRenterID(renterId *int, pageIndex *int, pageSize *int) ([]jsonutils.JSONObject, int64, error) {
	containerList := make([]jsonutils.JSONObject, 0)

	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	clusters, err := k8s.KubeClusters.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	if len(clusters.Data) == 0 {
		return nil, 0, nil
	}
	for _, cluster := range clusters.Data {
		status, err := cluster.GetString("status")
		if err != nil {
			continue
		}
		if status != "running" {
			continue
		}
		clusterId, err := cluster.GetString("id")
		if err != nil {
			continue
		}

		params = jsonutils.NewDict()
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("cluster", jsonutils.NewString(clusterId))

		namespace, err := k8s.Namespaces.GetByName(s, "ygpt", params)
		if err != nil {
			continue
		}

		namespaceId, err := namespace.GetString("id")
		if err != nil {
			continue
		}

		params.Set("namespace", jsonutils.NewString(namespaceId))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("show_fail_reason", jsonutils.NewBool(true))

		if renterId != nil {
			params.Set("filter", jsonutils.NewString("name.contains('"+strconv.Itoa(*renterId)+"')"))
		}

		deployments, err := k8s.Deployments.List(s, params)
		if err != nil {
			continue
		}

		podListParams := jsonutils.NewDict()
		podListParams.Set("cluster", jsonutils.NewString(clusterId))
		podListParams.Set("namespace", jsonutils.NewString(namespaceId))
		podListParams.Set("details", jsonutils.NewBool(true))
		podListParams.Set("scope", jsonutils.NewString("system"))
		podListParams.Set("show_fail_reason", jsonutils.NewBool(true))
		podListParams.Set("owner_kind", jsonutils.NewString("Deployment"))

		for _, deployment := range deployments.Data {
			ownerName, err := deployment.GetString("name")
			if err != nil {
				continue
			}

			id, err := deployment.GetString("id")
			if err != nil {
				continue
			}

			podListParams.Set("owner_name", jsonutils.NewString(ownerName))

			podsList, err := k8s.Pods.List(s, podListParams)
			if err != nil {
				continue
			}
			if len(podsList.Data) == 0 {
				continue
			}

			pod := jsonutils.NewDict()
			pod.Set("pod", podsList.Data[0])
			pod.Set("id", jsonutils.NewString(id))
			pod.Set("name", jsonutils.NewString(ownerName))

			containerList = append(containerList, pod)
		}
	}

	return containerList, int64(len(containerList)), nil
}
