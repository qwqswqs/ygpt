package storage

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

// 虚拟机云硬盘
type VmEvsApi struct {
}

type getVmEvsListReq struct {
	PageIndex   int64   `json:"pageIndex"`
	PageSize    int64   `json:"pageSize"`
	StorageType string  `json:"storageType"` //nfs、rbd、local、baremetal
	DiskType    *string `json:"diskType"`    //不筛选不传，后端解析为nil；sys为系统盘，data为数据盘
	Name        *string `json:"name"`        //不筛选不传，后端解析为nil
	Unused      *bool   `json:"unused"`      //不筛选不传，后端解析为nil；false未挂载，true已挂载
	SizeDesc    string  `json:"sizeDesc"`
	TimeDesc    string  `json:"timeDesc"`
}
type Mounted struct {
	MountedBy   string `json:"mounted_by"`    //挂载的云主机名称
	MountedById string `json:"mounted_by_id"` //挂载的云主机id
	Status      string `json:"status"`        //挂载的云主机状态
}
type getVmEvsListResp struct {
	Id          string    `json:"id"`          //云硬盘id
	Name        string    `json:"name"`        //云硬盘名称
	Type        string    `json:"type"`        //云硬盘类别，sys系统，data数据盘
	StorageType string    `json:"storageType"` //nfs、rbd、local
	Storage     string    `json:"storage"`     //存储总量
	FromVmNas   string    `json:"from_vm_nas"` //来自哪个nas
	Status      string    `json:"status"`      //状态
	Mounted     []Mounted `json:"Mounted"`     //已挂载的云主机列表
	Age         string    `json:"age"`         //创建时间
}

// 虚拟机云硬盘
func (e *VmEvsApi) GetVmEvsList(c *gin.Context) {
	var r getVmEvsListReq
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
	jointFilters := make([]string, 0)
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("summary_stats", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	filters = append(filters, "disk_type.notin(volume)")
	if r.DiskType != nil {
		params.Set("disk_type", jsonutils.NewString(*r.DiskType))
	}

	if r.Unused != nil {
		params.Set("unused", jsonutils.NewBool(*r.Unused))
	}

	if r.Name != nil {
		filters = append(filters, "name.contains('"+*r.Name+"')")
	}
	if r.SizeDesc != "" {
		params.Set("order_by", jsonutils.NewString("size"))
		params.Set("order", jsonutils.NewString(r.SizeDesc))
	}
	if r.TimeDesc != "" {
		params.Set("order_by", jsonutils.NewString("created_at"))
		params.Set("order", jsonutils.NewString(r.TimeDesc))
	}

	switch r.StorageType {
	case "nfs":
		jointFilters = append(jointFilters, "storages.id(storage_id).storage_type.equals('nfs')")
	case "rbd":
		jointFilters = append(jointFilters, "storages.id(storage_id).storage_type.equals('rbd')")
	case "local":
		jointFilters = append(jointFilters, "storages.id(storage_id).storage_type.equals('local')")
	case "baremetal":
		jointFilters = append(jointFilters, "storages.id(storage_id).storage_type.equals('baremetal')")
	default:
		jointFilters = append(jointFilters, "storages.id(storage_id).storage_type.notin('local','baremetal')")
	}

	params.Set("joint_filter", jsonutils.NewStringArray(jointFilters))
	params.Set("filter", jsonutils.NewStringArray(filters))

	result, err := compute.Disks.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	vmEvsList := make([]getVmEvsListResp, 0)

	for _, data := range result.Data {
		var vmEvs getVmEvsListResp
		// 获取 ID
		vmEvs.Id, _ = data.GetString("id")

		// 获取 Name
		vmEvs.Name, _ = data.GetString("name")

		diskSize, _ := data.Int("disk_size")
		// 获取 Storage
		vmEvs.Storage = strconv.FormatInt(diskSize/1024, 10) + "G"

		vmEvs.Type, _ = data.GetString("disk_type")

		vmEvs.StorageType, _ = data.GetString("storage_type")

		// 获取 storage
		vmEvs.FromVmNas, _ = data.GetString("storage")

		// 获取 Status
		vmEvs.Status, _ = data.GetString("status")

		guests, err := data.GetArray("guests")
		if err == nil {
			if len(guests) > 0 {
				for _, guest := range guests {
					// 获取 MountedBy
					var mount Mounted
					mount.MountedBy, _ = guest.GetString("name")
					mount.MountedById, _ = guest.GetString("id")
					mount.Status, _ = guest.GetString("status")
					vmEvs.Mounted = append(vmEvs.Mounted, mount)
				}
			}
		}

		// 获取创建时间
		vmEvs.Age, _ = data.GetString("created_at")

		vmEvsList = append(vmEvsList, vmEvs)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     vmEvsList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

// 获取回收站硬盘列表
func (e *VmEvsApi) GetRecycleVmEvsList(c *gin.Context) {
	var r getVmEvsListReq
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
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("pending_delete", jsonutils.NewBool(true))
	params.Set("summary_stats", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != nil {
		filters = append(filters, "name.contains('"+*r.Name+"')")
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Disks.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	vmEvsList := make([]getVmEvsListResp, 0)

	for _, data := range result.Data {
		var vmEvs getVmEvsListResp
		// 获取 ID
		vmEvs.Id, _ = data.GetString("id")

		// 获取 Name
		vmEvs.Name, _ = data.GetString("name")

		diskSize, _ := data.Int("disk_size")
		// 获取 Storage
		vmEvs.Storage = strconv.FormatInt(diskSize/1024, 10) + "G"

		vmEvs.Type, _ = data.GetString("disk_type")

		vmEvs.StorageType, _ = data.GetString("storage_type")

		// 获取 storage
		vmEvs.FromVmNas, _ = data.GetString("storage")

		// 获取 Status
		vmEvs.Status, _ = data.GetString("status")

		guests, err := data.GetArray("guests")
		if err == nil {
			if len(guests) > 0 {
				for _, guest := range guests {
					// 获取 MountedBy
					var mount Mounted
					mount.MountedBy, _ = guest.GetString("name")
					mount.MountedById, _ = guest.GetString("id")
					mount.Status, _ = guest.GetString("status")
					vmEvs.Mounted = append(vmEvs.Mounted, mount)
				}
			}
		}

		// 获取创建时间
		vmEvs.Age, _ = data.GetString("created_at")

		vmEvsList = append(vmEvsList, vmEvs)
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     vmEvsList,
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

// 清除回收站硬盘
func (e *VmEvsApi) ClearRecycleVmEvs(c *gin.Context) {
	var r struct {
		ID string `json:"id"`
	}
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
	params.Set("id", jsonutils.NewString(r.ID))
	params.Set("override_pending_delete", jsonutils.NewBool(true))

	_, err = compute.Disks.DeleteWithParam(s, r.ID, params, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("获取成功", c)
}

// 恢复回收站硬盘
func (e *VmEvsApi) RecoverRecycleVmEvs(c *gin.Context) {
	var r struct {
		ID string `json:"id"`
	}
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
	params.Set("id", jsonutils.NewString(r.ID))

	_, err = compute.Disks.PerformAction(s, r.ID, "cancel-delete", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("获取成功", c)
}

type deleteVmEvsReq struct {
	Id string `json:"id"`
}
type deleteVmEvsByIdsReq struct {
	Ids []string `json:"ids"`
}

func (e *VmEvsApi) DeleteVmEvs(c *gin.Context) {
	var r deleteVmEvsReq
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

	_, err = compute.Disks.Delete(s, r.Id, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (e *VmEvsApi) DeleteVmEvsByIds(c *gin.Context) {
	var r deleteVmEvsByIdsReq
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
	for _, id := range r.Ids {
		params := jsonutils.NewDict()

		_, err = compute.Disks.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}
	response.OkWithMessage("删除成功", c)
}

type createVmEvsReq struct {
	Backend    string `json:"backend"`    //存储类型，local、nfs、ceph
	MediumType string `json:"mediumType"` //参考nas的mediumType
	Name       string `json:"name"`       //名称
	StorageId  string `json:"storage_id"` //nas存储的id
	Size       int64  `json:"size"`       //硬盘大小，字节数
}

func (e *VmEvsApi) CreateVmEvs(c *gin.Context) {
	var r createVmEvsReq
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
	params.Set("backend", jsonutils.NewString(r.Backend))
	params.Set("hypervisor", jsonutils.NewString("kvm"))
	params.Set("medium", jsonutils.NewString(r.MediumType))
	params.Set("name", jsonutils.NewString(r.Name))
	//params.Set("prefer_region_id", jsonutils.NewString("default"))
	//params.Set("prefer_zone", jsonutils.NewString("default"))
	//params.Set("project_domain", jsonutils.NewString("default"))
	//params.Set("project_id", jsonutils.NewString("default"))
	params.Set("size", jsonutils.NewInt(r.Size))
	params.Set("storage_id", jsonutils.NewString(r.StorageId))

	_, err = compute.Disks.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

type detachVmEvsReq struct {
	Id   string `json:"id"`
	VmId string `json:"vm_id"`
}

func (e *VmEvsApi) DetachVmEvs(c *gin.Context) {
	var r detachVmEvsReq
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
	params.Set("disk_id", jsonutils.NewString(r.Id))
	params.Set("keep_disk", jsonutils.NewBool(true))
	_, err = compute.Servers.PerformAction(s, r.VmId, "detachdisk", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("卸载成功", c)
}

type attachVmEvsReq struct {
	Id   string `json:"id"`
	VmId string `json:"vmId"`
}

func (e *VmEvsApi) AttachVmEvs(c *gin.Context) {
	var r attachVmEvsReq
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
	params.Set("disk_id", jsonutils.NewString(r.Id))
	_, err = compute.Servers.PerformAction(s, r.VmId, "attachdisk", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("挂载成功", c)
}

type getAttachableVmsReq struct {
	Id string `json:"id"`
}

func (e *VmEvsApi) GetAttachableVms(c *gin.Context) {
	var r getAttachableVmsReq
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
	params.Set("attachable_servers_for_disk", jsonutils.NewString(r.Id))
	params.Set("filter", jsonutils.NewString("status.in(ready, running)"))

	result, err := compute.Servers.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:  jsonutils.Marshal(result.Data).PrettyString(),
		Total: int64(result.Total),
	}, "获取成功", c)
}

type resizeVmEvsReq struct {
	Id   string `json:"id"`
	Size int64  `json:"size"`
}

func (e *VmEvsApi) ResizeVmEvs(c *gin.Context) {
	var r resizeVmEvsReq
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
	params.Set("size", jsonutils.NewInt(r.Size))

	_, err = compute.Disks.PerformAction(s, r.Id, "resize", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("扩容成功", c)
}
