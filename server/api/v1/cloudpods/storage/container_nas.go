package storage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/utils"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type ContainerNasApi struct {
}

type getContainerNasListReq struct {
	PageIndex int64   `json:"pageIndex"`
	PageSize  int64   `json:"pageSize"`
	Name      *string `json:"name"`    //不筛选不传，后端解析为nil
	Cluster   string  `json:"cluster"` //k8s集群ID，必传
}
type getContainerNasListResp struct {
	Id     string `json:"id"`      //Nas的id
	Name   string `json:"name"`    //Nas名称
	Status string `json:"status"`  //状态
	Pool   string `json:"pool"`    //存储池名称
	FsType string `json:"fs_type"` //文件系统
	Age    string `json:"age"`     //创建时间
}

// 容器存储池
func (e *ContainerNasApi) GetContainerNasList(c *gin.Context) {
	var r getContainerNasListReq
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
		response.FailWithMessage(errors.New("参数错误，未知k8s集群"), c)
		return
	}

	namespaceId := utils.GetNamespaceId(r.Cluster)
	if namespaceId == "" {
		response.FailWithMessage(err, c)
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

	if r.Name != nil {
		params.Set("filter", jsonutils.NewString("name.contains('"+*r.Name+"')"))
	}

	result, err := k8s.Storageclass.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	containerNasList := make([]getContainerNasListResp, 0)

	for _, data := range result.Data {
		var containerNas getContainerNasListResp
		// 获取 ID
		if id, ok := k8s.Storageclass.Get_Id(data).(string); ok {
			containerNas.Id = id
		}

		// 获取 Name
		if name, ok := k8s.Storageclass.Get_Name(data).(string); ok {
			containerNas.Name = name
		}

		// 获取 Status
		containerNas.Status, _ = data.GetString("status")

		// 获取 pool
		containerNas.Pool, _ = data.GetString("parameters", "pool")

		// 获取 FsType
		containerNas.FsType, _ = data.GetString("parameters", "csi.storage.k8s.io/fstype")

		// 获取 Age
		if age, ok := k8s.PersistentVolumeClaims.Get_Age(data).(string); ok {
			containerNas.Age = age
		}

		containerNasList = append(containerNasList, containerNas)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     containerNasList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

type createContainerNasReq struct {
	Name          string `json:"name"`          //名称
	Pool          string `json:"pool"`          //ceph存储池名称
	CephClusterId string `json:"cephClusterId"` //ceph集群ID
	CsiFsType     string `json:"csiFsType"`     //文件系统，传ext4或xfs
	Cluster       string `json:"cluster"`       //k8s集群ID，必传
	SecretId      string `json:"secretId"`      //选择的保密字典id
}

func (e *ContainerNasApi) CreateContainerNas(c *gin.Context) {
	var r createContainerNasReq
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

	cephCSIRBD := jsonutils.NewDict()
	cephCSIRBD.Set("clusterId", jsonutils.NewString(r.CephClusterId))
	cephCSIRBD.Set("csiFsType", jsonutils.NewString(r.CsiFsType))
	cephCSIRBD.Set("imageFeatures", jsonutils.NewString("layering"))
	cephCSIRBD.Set("pool", jsonutils.NewString(r.Pool))
	cephCSIRBD.Set("secretName", jsonutils.NewString(r.SecretId))
	cephCSIRBD.Set("secretNamespace", jsonutils.NewString(namespaceId))
	params := jsonutils.NewDict()
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("provisioner", jsonutils.NewString("rbd.csi.ceph.com"))
	params.Set("cephCSIRBD", cephCSIRBD)

	_, err = k8s.Storageclass.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

type deleteContainerNasReq struct {
	Id      string `json:"id"`      //存储池id
	Cluster string `json:"cluster"` //k8s集群ID，必传
}
type deleteContainerNasByIdsReq struct {
	Ids     []string `json:"ids"`     //存储池id
	Cluster string   `json:"cluster"` //k8s集群ID，必传
}

func (e *ContainerNasApi) DeleteContainerNas(c *gin.Context) {
	var r deleteContainerNasReq
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
	params.Set("cluster", jsonutils.NewString(r.Cluster))

	_, err = k8s.Storageclass.Delete(s, r.Id, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (e *ContainerNasApi) DeleteContainerNasByIds(c *gin.Context) {
	var r deleteContainerNasByIdsReq
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
	params.Set("cluster", jsonutils.NewString(r.Cluster))

	for _, id := range r.Ids {
		_, err = k8s.Storageclass.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}
	response.OkWithMessage("删除成功", c)
}

type getCephClusterIdReq struct {
	Cluster string `json:"cluster"` //k8s集群ID，必传
}

type getCephClusterIdResp struct {
	CephClusterId string `json:"cephClusterId"` //ceph集群ID
}

func (e *ContainerNasApi) GetCephClusterId(c *gin.Context) {
	var r getCephClusterIdReq
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

	cephCsiParams := jsonutils.NewDict()
	cephCsiParams.Set("type", jsonutils.NewString("cephCSI"))
	cephConf, err := k8s.KubeClusters.GetSpecific(s, r.Cluster, "component-setting", cephCsiParams)
	if err != nil {
		response.FailWithMessage("该k8s未部署ceph-csi，请于基础资源-k8s集群处进行部署", c)
		return
	}

	configs, err := cephConf.GetArray("cephCSI", "config")
	if err != nil {
		response.FailWithMessage("该k8s未部署ceph-csi，请于基础资源-k8s集群处进行部署", c)
		return
	}
	if len(configs) == 0 {
		response.FailWithMessage("该k8s未部署ceph-csi，请于基础资源-k8s集群处进行部署", c)
		return
	}

	config := configs[0]

	var resp getCephClusterIdResp

	resp.CephClusterId, err = config.GetString("clusterId")
	if err != nil {
		response.FailWithMessage("该k8s未部署ceph-csi，请于基础资源-k8s集群处进行部署", c)
		return
	}

	response.OkWithDetailed(resp, "获取成功", c)
}

type getCephSecretListReq struct {
	Cluster string `json:"cluster"` //k8s集群ID，必传
}

type getCephSecretListResp struct {
	Id     string `json:"id"`     //保密字典id
	Name   string `json:"name"`   //保密字典名称
	Status string `json:"status"` //状态
	Age    string `json:"age"`    //创建时间
}

// 获取ceph的账号密码
func (e *ContainerNasApi) GetCephSecretList(c *gin.Context) {
	var r getCephSecretListReq
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
	params.Set("limit", jsonutils.NewInt(0))
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("namespace", jsonutils.NewString(namespaceId))
	params.Set("filter", jsonutils.NewString("type.contains('yunion.io/ceph-csi')"))

	result, err := k8s.Secrets.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	cephSecretList := make([]getCephSecretListResp, 0)

	for _, data := range result.Data {
		var cephSecret getCephSecretListResp
		// 获取 ID
		if id, ok := k8s.Secrets.Get_Id(data).(string); ok {
			cephSecret.Id = id
		}

		// 获取 Name
		if name, ok := k8s.Secrets.Get_Name(data).(string); ok {
			cephSecret.Name = name
		}

		// 获取 Status
		cephSecret.Status, _ = data.GetString("status")

		// 获取 Age
		if age, ok := k8s.Secrets.Get_Age(data).(string); ok {
			cephSecret.Age = age
		}

		cephSecretList = append(cephSecretList, cephSecret)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     cephSecretList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

type createCephSecretReq struct {
	Cluster string `json:"cluster"` //k8s集群ID，必传
	UserId  string `json:"userId"`  //用户名
	UserKey string `json:"userKey"` //密钥
}

func (e *ContainerNasApi) CreateCephSecret(c *gin.Context) {
	var r createCephSecretReq
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

	_, err = k8s.Secrets.GetByName(s, r.UserId, params)
	if err == nil {
		response.FailWithMessage("已存在该用户，请删除后重新添加", c)
		return
	}

	cephCSI := jsonutils.NewDict()
	cephCSI.Set("userId", jsonutils.NewString(r.UserId))
	cephCSI.Set("userKey", jsonutils.NewString(r.UserKey))
	params.Set("name", jsonutils.NewString(r.UserId))
	params.Set("type", jsonutils.NewString("yunion.io/ceph-csi"))
	params.Set("cephCSI", cephCSI)

	_, err = k8s.Secrets.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

type delCephSecretReq struct {
	Cluster string `json:"cluster"` //k8s集群ID，必传
	Id      string `json:"id"`      //保密字典id
}

func (e *ContainerNasApi) DeleteCephSecret(c *gin.Context) {
	var r delCephSecretReq

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
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("namespace", jsonutils.NewString("ygpt"))
	params.Set("id", jsonutils.NewString(r.Id))

	_, err = k8s.Secrets.Delete(s, r.Id, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

type getCephPoolsReq struct {
	CephClusterId string `json:"cephClusterId"` //ceph集群ID
	SecretId      string `json:"secretId"`      //选择的保密字典id
	Cluster       string `json:"cluster"`       //k8s集群ID，必传
}
type getCephPoolsResp struct {
	Pools []string `json:"pools"` //存储池名数组
}

func (e *ContainerNasApi) GetCephPools(c *gin.Context) {
	var r getCephPoolsReq

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

	cephCSIRBD := jsonutils.NewDict()
	cephCSIRBD.Set("clusterId", jsonutils.NewString(r.CephClusterId))
	cephCSIRBD.Set("secretName", jsonutils.NewString(r.SecretId))
	cephCSIRBD.Set("secretNamespace", jsonutils.NewString(namespaceId))
	params := jsonutils.NewDict()
	params.Set("cluster", jsonutils.NewString(r.Cluster))
	params.Set("name", jsonutils.NewString("rbd-climc"))
	params.Set("provisioner", jsonutils.NewString("rbd.csi.ceph.com"))
	params.Set("cephCSIRBD", cephCSIRBD)

	result, err := k8s.Storageclass.PerformClassAction(s, "connection-test", params)
	if err != nil {
		response.FailWithMessage("连接Ceph集群失败，请检查认证密钥是否正确。", c)
		return
	}

	var resp getCephPoolsResp

	resp.Pools, err = jsonutils.GetStringArray(result, "cephCSIRBD", "pools")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(resp, "获取成功", c)
}
