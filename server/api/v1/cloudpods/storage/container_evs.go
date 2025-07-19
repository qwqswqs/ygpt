package storage

import (
	"github.com/gin-gonic/gin"
	"strings"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/utils"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type ContainerEvsApi struct{}

type getContainerEvsListReq struct {
	PageIndex int64   `json:"pageIndex"`
	PageSize  int64   `json:"pageSize"`
	Name      *string `json:"name"`    //不筛选不传，后端解析为nil
	Unused    *bool   `json:"unused"`  //不筛选不传，后端解析为nil；false未挂载，true已挂载
	Cluster   string  `json:"cluster"` //k8s集群ID，必传
}

type getContainerEvsListResp struct {
	Id               string   `json:"id"`                 //云硬盘id
	Name             string   `json:"name"`               //云硬盘名称
	Storage          string   `json:"storage"`            //存储总量
	FromContainerNas string   `json:"from_container_nas"` //来自哪个容器nas
	Status           string   `json:"status"`             //状态
	MountedBys       []string `json:"mounted_bys"`        //挂载的容器名称
	Age              string   `json:"age"`                //创建时间
}

// 容器云硬盘
func (e *ContainerEvsApi) GetContainerEvsList(c *gin.Context) {
	var r getContainerEvsListReq
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

	if r.Cluster == "" {
		response.FailWithMessage("参数错误，未知k8s集群", c)
		return
	}

	namespaceId := utils.GetNamespaceId(r.Cluster)
	if namespaceId == "" {
		response.FailWithMessage("未找到云管平台命名空间", c)
		return
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("namespace", jsonutils.NewString(namespaceId))

	if r.Unused != nil {
		params.Set("unused", jsonutils.NewBool(*r.Unused))
	}

	if r.Name != nil {
		params.Set("filter", jsonutils.NewString("name.contains('"+*r.Name+"')"))
	}

	result, err := k8s.PersistentVolumeClaims.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	containerEvsList := make([]getContainerEvsListResp, 0)

	for _, data := range result.Data {
		var containerEvs getContainerEvsListResp
		// 获取 ID
		if id, ok := k8s.PersistentVolumeClaims.Get_Id(data).(string); ok {
			containerEvs.Id = id
		}

		// 获取 Name
		if name, ok := k8s.PersistentVolumeClaims.Get_Name(data).(string); ok {
			containerEvs.Name = name
		}

		// 获取 Storage
		containerEvs.Storage, _ = data.GetString("capacity", "storage")

		// 获取 StorageClass
		if storageClass, ok := k8s.PersistentVolumeClaims.Get_StorageClass(data).(string); ok {
			containerEvs.FromContainerNas = storageClass
		}

		// 获取 Status
		if status, ok := k8s.PersistentVolumeClaims.Get_Status(data).(string); ok {
			containerEvs.Status = status
		}

		// 获取 MountedBy
		pods, _ := data.GetArray("mountedBy")

		for _, pod := range pods {
			podName, _ := pod.GetString()
			// 找到第一个 '-' 的位置
			firstDashIndex := strings.Index(podName, "-")
			// 找到第二个 '-' 的位置
			secondDashIndex := strings.Index(podName[firstDashIndex+1:], "-") + firstDashIndex + 1

			// 提取所需的部分
			if secondDashIndex > firstDashIndex {
				containerEvs.MountedBys = append(containerEvs.MountedBys, podName[:secondDashIndex])
			} else {
				containerEvs.MountedBys = append(containerEvs.MountedBys, podName)
			}
		}

		// 获取 Age
		if age, ok := k8s.PersistentVolumeClaims.Get_Age(data).(string); ok {
			containerEvs.Age = age
		}

		containerEvsList = append(containerEvsList, containerEvs)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     containerEvsList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

type createContainerEvsReq struct {
	Cluster      string `json:"cluster"`      //k8s集群Id
	Name         string `json:"name"`         //云硬盘名称
	Size         string `json:"size"`         //磁盘大小，比如12Gi，目前只以Gi为单位
	ContainerNas string `json:"containerNas"` //所属容器Nas名称
}

func (e *ContainerEvsApi) CreateContainerEvs(c *gin.Context) {
	var r createContainerEvsReq
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

	if r.Cluster == "" {
		response.FailWithMessage("参数错误，未知k8s集群", c)
		return
	}

	namespaceId := utils.GetNamespaceId(r.Cluster)
	if namespaceId == "" {
		response.FailWithMessage("未找到云管平台命名空间", c)
		return
	}

	params := jsonutils.NewDict()
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("namespace", jsonutils.NewString(namespaceId))
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("size", jsonutils.NewString(r.Size))
	params.Set("storageClass", jsonutils.NewString(r.ContainerNas))

	_, err = k8s.PersistentVolumeClaims.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

type deleteContainerEvsReq struct {
	Cluster string `json:"cluster"` //k8s集群id
	Id      string `json:"id"`      //云硬盘id
}
type deleteContainerEvsByIdsReq struct {
	Cluster string   `json:"cluster"` //k8s集群id
	Ids     []string `json:"ids"`     //云硬盘id
}

func (e *ContainerEvsApi) DeleteContainerEvs(c *gin.Context) {
	var r deleteContainerEvsReq
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

	if r.Cluster == "" {
		response.FailWithMessage("参数错误，未知k8s集群", c)
		return
	}

	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.Id))
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("namespace", jsonutils.NewString("ygpt"))

	_, err = k8s.PersistentVolumeClaims.Delete(s, r.Id, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
func (e *ContainerEvsApi) DeleteContainerEvsByIds(c *gin.Context) {
	var r deleteContainerEvsByIdsReq
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

	if r.Cluster == "" {
		response.FailWithMessage("参数错误，未知k8s集群", c)
		return
	}
	for _, id := range r.Ids {
		params := jsonutils.NewDict()
		params.Set("id", jsonutils.NewString(id))
		params.Set("cluster", jsonutils.NewString(r.Cluster))
		params.Set("namespace", jsonutils.NewString("ygpt"))

		_, err = k8s.PersistentVolumeClaims.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
