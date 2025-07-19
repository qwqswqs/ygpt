package storage

import (
	"github.com/gin-gonic/gin"
	"strconv"
	compute2 "ygpt/server/api/v1/cloudpods/compute"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

// 虚拟机存储池
type VmNasApi struct {
}
type getVmNasListReq struct {
	PageIndex   int64  `json:"pageIndex"`
	PageSize    int64  `json:"pageSize"`
	StorageType string `json:"storageType"` //nfs、rbd
	Name        string `json:"name"`
	Enable      string `json:"enable"`
}
type storageConf struct {
	NfsHost      string `json:"nfsHost"`      //nfs服务地址，nfs时有效
	NfsSharedDir string `json:"nfsSharedDir"` //共享目录，nfs时有效
	MonHost      string `json:"monHost"`      //ceph集群监控器地址，rbd时有效
	Pool         string `json:"pool"`         //存储池，rbd时有效
}

type usedInfo struct {
	Capacity            string `json:"capacity"`            //物理容量
	ActualCapacityUsed  string `json:"actualCapacityUsed"`  //物理容量已使用
	VirtualCapacity     string `json:"virtualCapacity"`     //虚拟容量
	VirtualUsedCapacity string `json:"virtualUsedCapacity"` //虚拟容量已使用
}

type getVmNasListResp struct {
	Id          string      `json:"id"`          //Nas的id
	Name        string      `json:"name"`        //Nas名称
	Status      string      `json:"status"`      //状态
	Enabled     bool        `json:"enabled"`     //启用状态
	StorageConf storageConf `json:"storageConf"` //存储配置
	UsedInfo    usedInfo    `json:"usedInfo"`    //用量统计
	MediumType  string      `json:"mediumType"`
	StorageType string      `json:"storageType"` //nfs、rbd、local、baremetal
	Age         string      `json:"age"`         //创建时间
}

// 获取存储池列表
func (d *VmNasApi) GetVmNasList(c *gin.Context) {
	var r getVmNasListReq
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

	filters := make([]string, 0)
	params := jsonutils.NewDict()
	params.Set("is_baremetal", jsonutils.NewBool(false))
	params.Set("summary_stats", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))

	switch r.StorageType {
	case "nfs":
		filters = append(filters, "storage_type.equals('nfs')")
	case "rbd":
		filters = append(filters, "storage_type.equals('rbd')")
	case "local":
		filters = append(filters, "storage_type.equals('local')")
	case "baremetal":
		params.Set("is_baremetal", jsonutils.NewBool(true))
	default:
		filters = append(filters, "storage_type.notin('local')")
	}

	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}

	params.Set("filter", jsonutils.NewStringArray(filters))
	if r.Enable == "true" {
		params.Set("enabled", jsonutils.NewBool(true))
	} else if r.Enable == "false" {
		params.Set("enabled", jsonutils.NewBool(false))
	}

	if r.PageSize == 0 {
		params.Set("limit", jsonutils.NewInt(r.PageSize))
	} else {
		params.Set("limit", jsonutils.NewInt(r.PageSize))
		params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	}
	result, err := compute.Storages.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	vmNasList := make([]getVmNasListResp, 0)

	for _, data := range result.Data {
		var vmNas getVmNasListResp
		// 获取 ID
		vmNas.Id, _ = data.GetString("id")

		// 获取 Name
		vmNas.Name, _ = data.GetString("name")

		// 获取 Status
		vmNas.Status, _ = data.GetString("status")

		// 获取 Enabled
		vmNas.Enabled, _ = data.Bool("enabled")

		vmNas.MediumType, _ = data.GetString("medium_type")

		vmNas.StorageType, _ = data.GetString("storage_type")

		// 获取 存储配置
		vmNas.StorageConf.Pool, _ = data.GetString("storage_conf", "pool")
		vmNas.StorageConf.MonHost, _ = data.GetString("storage_conf", "mon_host")
		vmNas.StorageConf.NfsHost, _ = data.GetString("storage_conf", "nfs_host")
		vmNas.StorageConf.NfsSharedDir, _ = data.GetString("storage_conf", "nfs_shared_dir")

		// 获取用量统计
		capacity, _ := data.Int("capacity")
		actualCapacityUsed, _ := data.Int("actual_capacity_used")
		virtualCapacity, _ := data.Int("virtual_capacity")
		virtualUsedCapacity, _ := data.Int("used_capacity")

		vmNas.UsedInfo.Capacity = strconv.FormatInt(capacity/1024, 10) + "G"
		vmNas.UsedInfo.ActualCapacityUsed = strconv.FormatInt(actualCapacityUsed/1024, 10) + "G"
		vmNas.UsedInfo.VirtualCapacity = strconv.FormatInt(virtualCapacity/1024, 10) + "G"
		vmNas.UsedInfo.VirtualUsedCapacity = strconv.FormatInt(virtualUsedCapacity/1024, 10) + "G"

		// 获取 Age
		vmNas.Age, _ = data.GetString("created_at")

		vmNasList = append(vmNasList, vmNas)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     vmNasList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

type deleteVmNasReq struct {
	Id string `json:"id"`
}
type deleteVmNasByIdsReq struct {
	Ids []string `json:"ids"`
}

// 删除存储
func (d *VmNasApi) DeleteVmNas(c *gin.Context) {
	var r deleteVmNasReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.Id))
	_, err = compute.Storages.Delete(s, r.Id, params)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// 批量删除存储
func (d *VmNasApi) DeleteVmNasByIds(c *gin.Context) {
	var r deleteVmNasByIdsReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	for _, id := range r.Ids {
		params := jsonutils.NewDict()
		params.Set("id", jsonutils.NewString(id))
		_, err = compute.Storages.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage("删除失败", c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

type nfsNas struct {
	Host string `json:"host"`
	Dir  string `json:"dir"`
}
type cephNas struct {
	RbdKey     string `json:"rbd_key"`
	RbdMonHost string `json:"rbd_mon_host"`
	RbdPool    string `json:"rbd_pool"`
}
type addVmNasReq struct {
	Name        string  `json:"name"`
	StorageType string  `json:"storageType"` //nfs、ceph
	Nfs         nfsNas  `json:"nfs"`
	Ceph        cephNas `json:"ceph"`
}

// 添加存储
func (d *VmNasApi) AddVmNas(c *gin.Context) {
	var r addVmNasReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	params := jsonutils.NewDict()
	params.Set("area", jsonutils.NewString("default"))
	params.Set("capacity", jsonutils.NewInt(0))
	params.Set("medium_type", jsonutils.NewString("rotate"))
	params.Set("zone", jsonutils.NewString(compute2.GetZoneID()))
	params.Set("name", jsonutils.NewString(r.Name))

	switch r.StorageType {
	case "nfs":
		params.Set("nfs_host", jsonutils.NewString(r.Nfs.Host))
		params.Set("nfs_shared_dir", jsonutils.NewString(r.Nfs.Dir))
		params.Set("storage_type", jsonutils.NewString("nfs"))
	case "rbd":
		params.Set("rbd_mon_host", jsonutils.NewString(r.Ceph.RbdMonHost))
		params.Set("rbd_key", jsonutils.NewString(r.Ceph.RbdKey))
		params.Set("rbd_pool", jsonutils.NewString(r.Ceph.RbdPool))
		params.Set("storage_type", jsonutils.NewString("rbd"))
	default:
		response.FailWithMessage("未知存储类型", c)
		return
	}

	_, err = compute.Storages.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

type disableVmNasReq struct {
	Id string `json:"id"`
}

// 禁用存储
func (d *VmNasApi) DisableVmNas(c *gin.Context) {
	var r disableVmNasReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("禁用失败", c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.Id))
	_, err = compute.Storages.PerformAction(s, r.Id, "disable", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("禁用成功", c)
}

type enableVmNasReq struct {
	Id string `json:"id"`
}

// 启用存储
func (d *VmNasApi) EnableVmNas(c *gin.Context) {
	var r enableVmNasReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("启用失败", c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.Id))
	_, err = compute.Storages.PerformAction(s, r.Id, "enable", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("启用成功", c)
}

type vmNasConnHostReq struct {
	ID      string   `json:"id"`
	HostIds []string `json:"hostIds"`
	Dir     string   `json:"dir"`
}

// NAS存储关联宿主机
func (d *VmNasApi) VmNasConnHost(c *gin.Context) {
	var r vmNasConnHostReq
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
	params.Set("dir", jsonutils.NewString(r.Dir))
	params.Set("mount_point", jsonutils.NewString(r.Dir))
	//params.Set("host", jsonutils.NewStringArray(r.HostIds))

	for _, hostId := range r.HostIds {
		_, err = compute.Hoststorages.Attach(s, hostId, r.ID, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("关联成功", c)
}

type vmNasDisConnHostReq struct {
	Id     string `json:"id"`
	HostId string `json:"hostId"`
}

// NAS存储取消关联宿主机
func (d *VmNasApi) VmNasDisConnHost(c *gin.Context) {
	var r vmNasDisConnHostReq
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

	_, err = compute.Storages.Put(s, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	_, err = compute.Hoststorages.Detach(s, r.HostId, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("取消关联成功", c)
}

type vmNasConnHostListReq struct {
	Id string `json:"id"`
}

// NAS存储关联的宿主机列表
func (d *VmNasApi) VmNasConnHostList(c *gin.Context) {
	var r vmNasConnHostListReq
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

	param := jsonutils.NewDict()
	param.Set("scope", jsonutils.NewString("system"))
	param.Set("show_fail_reason", jsonutils.NewBool(true))
	param.Set("storage", jsonutils.NewString(r.Id))
	param.Set("details", jsonutils.NewBool(true))

	hosts, err := compute.Hosts.List(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  jsonutils.Marshal(hosts.Data).PrettyString(),
		Total: int64(hosts.Total),
	}, "获取成功", c)
}

type vmNasNoConnHostListReq struct {
	Id string `json:"id"`
}

func (d *VmNasApi) VmNasNoConnHostList(c *gin.Context) {
	var r vmNasNoConnHostListReq
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

	param := jsonutils.NewDict()
	param.Set("hypervisor", jsonutils.NewString("kvm"))
	param.Set("storage_not_attached", jsonutils.NewBool(true))
	param.Set("storage", jsonutils.NewString(r.Id))

	hosts, err := compute.Hosts.List(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  jsonutils.Marshal(hosts.Data).PrettyString(),
		Total: int64(hosts.Total),
	}, "获取成功", c)
}
