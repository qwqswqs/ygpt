package compute

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/service"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	compute "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	devtool "yunion.io/x/onecloud/pkg/mcclient/modules/devtool"
	"yunion.io/x/onecloud/pkg/mcclient/modules/scheduler"
	webconsole "yunion.io/x/onecloud/pkg/mcclient/modules/webconsole"
)

type ComputeServersApi struct{}
type hostDetails struct {
	Id         string `json:"id"`
	ServerType string `json:"serverType"` //vm 云主机;baremetal 裸金属
}

func (p *ComputeServersApi) HostDetails(c *gin.Context) {
	var r hostDetails
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
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("show_emulated", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	switch r.ServerType {
	case "vm":
		params.Set("filter", jsonutils.NewString("hypervisor.notin(baremetal,container)"))
	case "baremetal":
		params.Set("hypervisor", jsonutils.NewString("baremetal"))
	default:
		response.FailWithMessage(errors.New("未知实例类型:"+r.ServerType), c)
		return
	}

	serverDetails, err := compute.Servers.GetById(s, r.Id, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(serverDetails.Interface(), "获取成功", c)
}

type getBareListReq struct {
	PageIndex    int64  `json:"pageIndex"`
	PageSize     int64  `json:"pageSize"`
	Name         string `json:"name"`
	SecGroupID   string `json:"secGroupID"`
	ScalingGroup string `json:"scalingGroup"`
	IpDesc       string `json:"ipDesc"`
	TimeDesc     string `json:"timeDesc"`
	Status       string `json:"status"`
	OsDist       string `json:"osDist"`
}

// 恢复回收主机
func (p *ComputeServersApi) RecoverRecycleHost(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.PerformAction(s, r.ID, "cancel-delete", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "清除成功", c)
}

// 清除回收主机
func (p *ComputeServersApi) ClearRecycleHost(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.DeleteWithParam(s, r.ID, params, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "清除成功", c)
}

// 批量清除回收主机
func (p *ComputeServersApi) ClearRecycleHostByIds(c *gin.Context) {
	var r hostBaseByIdsReq
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
		params.Set("id", jsonutils.NewString(id))
		params.Set("override_pending_delete", jsonutils.NewBool(true))
		_, err := compute.Servers.DeleteWithParam(s, id, params, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}
	response.OkWithMessage("清除成功", c)
}

// 获取回收站列表
func (p *ComputeServersApi) GetRecycleServerList(c *gin.Context) {
	var r getBareListReq
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
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.TimeDesc != "" {
		params.Set("order_by", jsonutils.NewString("auto_delete_at"))
		params.Set("order", jsonutils.NewString(r.TimeDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Servers.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 获取云主机登录账号密码
func (d *ComputeServersApi) GetLoginInfo(c *gin.Context) {
	var r webConsoleReq
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
	result, err := compute.Servers.GetLoginInfo(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 查询裸金属列表
func (p *ComputeServersApi) GetBareServerList(c *gin.Context) {
	var r getBareListReq
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
	params.Set("details", jsonutils.NewBool(true))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("hypervisor", jsonutils.NewString("baremetal"))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.IpDesc != "" {
		params.Set("order_by", jsonutils.NewString("ips"))
		params.Set("order_by_ip", jsonutils.NewString(r.IpDesc))
		params.Set("order", jsonutils.NewString(r.IpDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Servers.List(s, params)
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

// 查询云主机列表
func (p *ComputeServersApi) GetServerList(c *gin.Context) {
	var r getBareListReq
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
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	filters = append(filters, "hypervisor.notin(baremetal,container)")
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.SecGroupID != "" {
		params.Set("secgroup_id", jsonutils.NewString(r.SecGroupID))
		params.Set("brand", jsonutils.NewString("OneCloud"))
	}
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.ScalingGroup != "" {
		params.Set("scaling_group", jsonutils.NewString(r.ScalingGroup))
	}
	if r.OsDist != "" {
		params.Set("os_dist", jsonutils.NewString(r.OsDist))
	}

	if r.IpDesc != "" {
		params.Set("order_by", jsonutils.NewString("ips"))
		params.Set("order_by_ip", jsonutils.NewString(r.IpDesc))
		params.Set("order", jsonutils.NewString(r.IpDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Servers.List(s, params)
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

type webConsoleReq struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}

// 云主机远程连接
func (p *ComputeServersApi) GetWebConsoleSSH(c *gin.Context) {
	var r webConsoleReq
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
	params.Set("ip", jsonutils.NewString(r.IP))
	params.Set("port", jsonutils.NewString("22"))
	params.Set("type", jsonutils.NewString("server"))
	result, err := webconsole.WebConsole.DoSshConnect(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

type createDiskReq struct {
	Disk         string `json:"disk"`
	GenerateName string `json:"generateName"`
}

// 创建硬盘快照
func (p *ComputeServersApi) CreateDiskSnap(c *gin.Context) {
	var r createDiskReq
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
	params.Set("disk", jsonutils.NewString(r.Disk))
	params.Set("generate_name", jsonutils.NewString(r.GenerateName))
	result, err := compute.Snapshots.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

type hostBaseReq struct {
	ID string `json:"ID"`
}
type hostBaseByIdsReq struct {
	Ids []string `json:"ids"`
}

// 删除
func (p *ComputeServersApi) HostDelete(c *gin.Context) {
	var r hostBaseReq
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
	param.Set("disable_delete", jsonutils.NewBool(false))
	param.Set("protected", jsonutils.NewBool(false))
	_, err = compute.Servers.Update(s, r.ID, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	params.Set("delete_snapshots", jsonutils.NewBool(true))
	params.Set("delete_disks", jsonutils.NewBool(true))
	params.Set("delete_eip", jsonutils.NewBool(true))
	result, err := compute.Servers.Delete(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	var model renter.RenterRes
	err = global.GVA_DB.Where("resource_id = ?", r.ID).First(&model).Error
	if err != nil {
		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
		return
	}
	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", r.ID).Updates(&renter.RenterRes{
		EndTime: time.Now(),
	}).Error
	if err != nil {
		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

// 批量删除
func (p *ComputeServersApi) HostDeleteByIds(c *gin.Context) {
	var r hostBaseByIdsReq
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
		param := jsonutils.NewDict()
		param.Set("disable_delete", jsonutils.NewBool(false))
		param.Set("protected", jsonutils.NewBool(false))
		_, err = compute.Servers.Update(s, id, param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		params := jsonutils.NewDict()
		params.Set("id", jsonutils.NewString(id))
		params.Set("delete_snapshots", jsonutils.NewBool(true))
		params.Set("delete_disks", jsonutils.NewBool(true))
		params.Set("delete_eip", jsonutils.NewBool(true))
		_, err := compute.Servers.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		var model renter.RenterRes
		err = global.GVA_DB.Where("resource_id = ?", id).First(&model).Error
		if err != nil {
			continue
		}
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", id).Updates(&renter.RenterRes{
			EndTime: time.Now(),
		}).Error
	}

	response.OkWithMessage("删除成功", c)
}

// 关机
func (p *ComputeServersApi) HostShutDown(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.PerformAction(s, r.ID, "stop", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "关机成功", c)
}

// 开机
func (p *ComputeServersApi) HostStart(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.PerformAction(s, r.ID, "start", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "开机成功", c)
}

// 重启
func (p *ComputeServersApi) HostRestart(c *gin.Context) {
	var r hostBaseReq
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
	params.Set("is_force", jsonutils.NewBool(true))
	result, err := compute.Servers.PerformAction(s, r.ID, "restart", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "重启成功", c)
}

// 挂起
func (p *ComputeServersApi) HostSuspend(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.PerformAction(s, r.ID, "suspend", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "挂起成功", c)
}

// 恢复
func (p *ComputeServersApi) HostResume(c *gin.Context) {
	var r hostBaseReq
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
	result, err := compute.Servers.PerformAction(s, r.ID, "resume", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "恢复成功", c)
}

type diskData struct {
	Backend   string `json:"backend"`
	Disk_type string `json:"disk_type"`
	Image_id  string `json:"image_id"`
	Index     int    `json:"index"`
	//Medium    string `json:"medium"` 不需要介质类型，如果传入不对，会导致错误
	Size int `json:"size"`
}
type netData struct {
	Network string `json:"network"`
}
type device struct {
	Model  string `json:"model"`
	Vendor string `json:"vendor"`
}

type hostCreateReq struct {
	Auto_start       bool       `json:"auto_start"`
	Deploy_telegraf  bool       `json:"deploy_telegraf"`
	Prefer_region    string     `json:"prefer_region"`
	Prefer_zone      string     `json:"prefer_zone"`
	Duration         string     `json:"duration"`
	CDrom            string     `json:"cdrom"`
	Project_id       string     `json:"project_id"`
	Count            int        `json:"count"`
	Hypervisor       string     `json:"hypervisor"`
	Dry_run          bool       `json:"dry_run"`
	Description      string     `json:"description"`
	Disks            []diskData `json:"disks"`
	Generate_name    string     `json:"generate_name"`
	Hostname         string     `json:"hostname"`
	Nets             []netData  `json:"nets"`
	Sku              string     `json:"sku"`
	EipBw            int        `json:"eip_bw"`
	EipChargeType    string     `json:"eip_charge_type"`
	Vcpu_count       int        `json:"vcpu_count"`
	Vmem_size        int        `json:"vmem_size"`
	Machine          string     `json:"machine"`
	OsArch           string     `json:"os_arch"`
	Isolated_devices []device   `json:"isolated_devices"`
}

func GetZoneID() string {

	s, err := global.GetCloudClientSession()
	if err != nil {
		return ""
	}
	params := jsonutils.NewDict()
	result, err := compute.Zones.List(s, params)
	if err != nil {
		return ""
	}
	jsonString := jsonutils.Marshal(result).PrettyString()
	startIndex := strings.Index(jsonString, `"id": "`) + len(`"id": "`)
	if startIndex == -1+len(`"id": "`) {
		fmt.Println("id字段未找到")
		return ""
	}

	// 查找id字段的结束位置（即下一个双引号的位置）
	endIndex := strings.Index(jsonString[startIndex:], `"`)
	if endIndex == -1 {
		fmt.Println("id字段值未正确结束")
		return ""
	}
	endIndex += startIndex

	// 提取id字段的值
	id := jsonString[startIndex:endIndex]
	fmt.Println("ID:", id)
	return id
}
func GetWireID() string {

	s, err := global.GetCloudClientSession()
	if err != nil {
		return ""
	}
	params := jsonutils.NewDict()
	result, err := compute.Wires.List(s, params)
	if err != nil {
		return ""
	}
	jsonString := jsonutils.Marshal(result).PrettyString()
	startIndex := strings.Index(jsonString, `"id": "`) + len(`"id": "`)
	if startIndex == -1+len(`"id": "`) {
		fmt.Println("id字段未找到")
		return ""
	}

	// 查找id字段的结束位置（即下一个双引号的位置）
	endIndex := strings.Index(jsonString[startIndex:], `"`)
	if endIndex == -1 {
		fmt.Println("id字段值未正确结束")
		return ""
	}
	endIndex += startIndex

	// 提取id字段的值
	id := jsonString[startIndex:endIndex]
	fmt.Println("ID:", id)
	return id
}

// 获取可使用的GPU信息
//
//	func (p *ComputeServersApi) GetGpuList(c *gin.Context) {
//		s, err := global.GetCloudClientSession()
//		if err != nil {
//			response.FailWithMessage(err, c)
//			return
//		}
//		podListParams := jsonutils.NewDict()
//		podListParams.Set("show_emulated", jsonutils.NewBool(true))
//		podListParams.Set("resource_type", jsonutils.NewString("shared"))
//		podListParams.Set("project_domain", jsonutils.NewString("default"))
//
//		result, err := compute.Capabilities.List(s, podListParams)
//		if err != nil {
//			response.FailWithMessage(err, c)
//			return
//		}
//		response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
//	}

func (p *ComputeServersApi) GetGpuList(c *gin.Context) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	gpuListParams := jsonutils.NewDict()
	gpuListParams.Set("scope", jsonutils.NewString("system"))
	gpuListParams.Set("show_fail_reason", jsonutils.NewBool(true))
	gpuListParams.Set("with_meta", jsonutils.NewBool(true))
	gpuListParams.Set("details", jsonutils.NewBool(true))
	gpuListParams.Set("show_baremetal_isolated_devices", jsonutils.NewBool(false))
	gpuListParams.Set("filter", jsonutils.NewString("guest_id.isnullorempty()"))

	result, err := compute.IsolatedDevices.List(s, gpuListParams)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	// 创建一个用于存储 GPU 型号和对应主机信息的 map
	gpuInfo := make(map[string]map[string]int)

	for _, data := range result.Data {
		hostEnabled, err := data.Bool("host_enabled")
		if err != nil {
			continue
		}
		if !hostEnabled {
			continue
		}

		model, err := data.GetString("model")
		if err != nil {
			continue
		}

		host, err := data.GetString("host")
		if err != nil {
			continue
		}

		// 初始化 GPU 型号的 map
		if _, exists := gpuInfo[model]; !exists {
			gpuInfo[model] = make(map[string]int)
		}

		// 检查主机是否存在，如果不存在则初始化为 0
		if _, exists := gpuInfo[model][host]; !exists {
			gpuInfo[model][host] = 0
		}

		// 增加主机对应的 GPU 数量
		gpuInfo[model][host]++
	}

	// 将结果转换为 JSON 格式
	response.OkWithDetailed(gpuInfo, "获取成功", c)
}

// 创建云主机
func (p *ComputeServersApi) HostCreate(c *gin.Context) {
	var r hostCreateReq
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
	//r.Disks[0].Image_id=""
	r.Auto_start = true
	r.Deploy_telegraf = true
	r.Prefer_region = "default"
	r.Prefer_zone = GetZoneID()
	r.Project_id = s.GetProjectId()
	r.Count = 1
	r.Hypervisor = "kvm"
	r.Dry_run = true
	r.Machine = "pc"
	r.OsArch = "x86"
	if r.EipBw != 0 {
		r.Hostname = r.Generate_name
		r.EipChargeType = "bandwidth"
	}
	_, err = compute.ServerSkus.GetByName(s, r.Sku, nil)
	if err != nil { // sku不存在，创建sku
		skuParams := jsonutils.NewDict()
		skuParams.Set("cpu_core_count", jsonutils.NewInt(int64(r.Vcpu_count)))
		skuParams.Set("memory_size_mb", jsonutils.NewInt(int64(r.Vmem_size)))
		skuInfo, err := compute.ServerSkus.Create(s, skuParams)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		r.Sku, err = skuInfo.GetString("name")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	params := jsonutils.Marshal(r)
	_, err = compute.Servers.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	_, err = scheduler.SchedManager.DoForecast(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	r.Dry_run = false
	params = jsonutils.Marshal(r)
	result, err := compute.Servers.Create(s, params)

	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

type BaremetalCreateReq struct {
	ProjectId    string `json:"project_id"` //
	Count        int    `json:"count"`      //
	VmemSize     int    `json:"vmem_size"`
	VcpuCount    int    `json:"vcpu_count"`
	GenerateName string `json:"generate_name"`
	Hypervisor   string `json:"hypervisor"` //
	AutoStart    bool   `json:"auto_start"` //
	Vdi          string `json:"vdi"`        //
	Disks        []struct {
		Size       int64  `json:"size"`
		ImageId    string `json:"image_id,omitempty"`
		Mountpoint string `json:"mountpoint,omitempty"`
	} `json:"disks"`
	BaremetalDiskConfigs []BaremetalDiskConfig `json:"baremetal_disk_configs"`
	Nets                 []struct {
		Network string `json:"network"`
		Exit    bool   `json:"exit"`
	} `json:"nets"`
	PreferHost      string   `json:"prefer_host"`
	PreferRegion    string   `json:"prefer_region"` //
	PreferZone      string   `json:"prefer_zone"`   //
	Password        string   `json:"password"`
	ImageId         string   `json:"imageId"`
	IsolatedDevices []device `json:"isolated_devices"`
}

type migrateReq struct {
	VmId        string `json:"vmId"`
	HostId      string `json:"hostId"`
	MigrateType string `json:"migrateType"` //cold 冷迁移、hot 热迁移、down 宕机
}

// 迁移
func (p *ComputeServersApi) Migrate(c *gin.Context) {
	var r migrateReq
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

	switch r.MigrateType {
	case "cold":
		params.Set("auto_start", jsonutils.NewBool(true))
		params.Set("prefer_host", jsonutils.NewString(r.HostId))

		// 冷迁移
		_, err = compute.Servers.PerformAction(s, r.VmId, "migrate", params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	case "hot":
		params.Set("prefer_host", jsonutils.NewString(r.HostId))
		params.Set("quickly_finish", jsonutils.NewBool(true))
		params.Set("skip_cpu_check", jsonutils.NewBool(true))
		params.Set("skip_kernel_check", jsonutils.NewBool(true))

		// 热迁移
		_, err = compute.Servers.PerformAction(s, r.VmId, "live-migrate", params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	case "down":
		params.Set("rescue_mode", jsonutils.NewBool(true))
		params.Set("prefer_host", jsonutils.NewString(r.HostId))

		// 宕机迁移
		_, err = compute.Servers.PerformAction(s, r.VmId, "migrate", params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	default:
		response.FailWithMessage(errors.New("未知迁移类型:"+r.MigrateType), c)
		return
	}

	response.OkWithMessage("迁移中,请稍等!", c)
}

type manageGPUReq struct {
	VmId      string   `json:"vmId"`      //云主机id
	AttachIds []string `json:"attachIds"` //要挂载的gpu的Id
	DetachIds []string `json:"detachIds"` //要卸载的gpu的Id
}

func (p *ComputeServersApi) ManageGpu(c *gin.Context) {
	var r manageGPUReq
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
	params.Set("auto_start", jsonutils.NewBool(true))
	params.Set("add_devices", jsonutils.NewStringArray(r.AttachIds))
	params.Set("del_devices", jsonutils.NewStringArray(r.DetachIds))

	_, err = compute.Servers.PerformAction(s, r.VmId, "set-isolated-device", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("操作成功", c)
}

type getGPUInfoReq struct {
	HostId string `json:"hostId"` //云主机所在宿主机id
}

func (p *ComputeServersApi) GetGpuInfo(c *gin.Context) {
	var r getGPUInfoReq
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
	params.Set("limit", jsonutils.NewInt(0))

	filters := make([]string, 0)
	filters = append(filters, "host_id.in("+r.HostId+")")
	filters = append(filters, "dev_type.notin(USB,NIC,NVME-PT)")

	params.Set("filter", jsonutils.Marshal(filters))

	result, err := compute.IsolatedDevices.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	gpuArray := jsonutils.NewArray()

	for _, data := range result.Data {
		id, _ := data.GetString("id")
		model, _ := data.GetString("model")
		vendorDeviceId, _ := data.GetString("vendor_device_id")
		addr, _ := data.GetString("addr")
		guest, _ := data.GetString("guest")
		guestId, _ := data.GetString("guest_id")
		guestStatus, _ := data.GetString("guest_status")

		gpuInfo := jsonutils.NewDict()
		gpuInfo.Set("id", jsonutils.NewString(id))
		gpuInfo.Set("vendorDeviceId", jsonutils.NewString(vendorDeviceId))
		gpuInfo.Set("addr", jsonutils.NewString(addr))
		gpuInfo.Set("model", jsonutils.NewString(model))
		gpuInfo.Set("vmName", jsonutils.NewString(guest))
		gpuInfo.Set("vmId", jsonutils.NewString(guestId))
		gpuInfo.Set("vmStatus", jsonutils.NewString(guestStatus))

		gpuArray.Add(gpuInfo)
	}

	response.OkWithDetailed(gpuArray.Interface(), "获取成功", c)
}

type migrateForecastReq struct {
	VmId        string `json:"vmId"`
	MigrateType string `json:"migrateType"` //cold 冷迁移、hot 热迁移、down 宕机
	PageIndex   int64  `json:"pageIndex"`
	PageSize    int64  `json:"pageSize"`
}

// 迁移探测
func (p *ComputeServersApi) MigrateForecast(c *gin.Context) {
	var r migrateForecastReq
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

	switch r.MigrateType {
	case "cold":
		params.Set("live_migrate", jsonutils.NewBool(false))
		params.Set("skip_cpu_check", jsonutils.NewBool(true))
		params.Set("skip_kernel_check", jsonutils.NewBool(true))
	case "hot":
		params.Set("live_migrate", jsonutils.NewBool(true))
		params.Set("skip_cpu_check", jsonutils.NewBool(true))
		params.Set("skip_kernel_check", jsonutils.NewBool(true))
	case "down":
		params.Set("live_migrate", jsonutils.NewBool(false))
		params.Set("is_rescue_mode", jsonutils.NewBool(true))
	default:
		response.FailWithMessage(errors.New("未知迁移类型:"+r.MigrateType), c)
		return
	}

	// 迁移探测可迁移的宿主机
	result, err := compute.Servers.PerformAction(s, r.VmId, "migrate-forecast", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	var data jsonutils.JSONObject
	isParse := result.Contains("data")
	if isParse {
		dataArray, err := result.GetArray("data")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		if len(dataArray) == 0 {
			response.FailWithMessage(errors.New("无宿主可选"), c)
			return
		}
		data = dataArray[0]
	} else {
		data = result
	}

	filterStr := "id.notin("
	filteredCandidates, err := data.GetArray("filtered_candidates")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	for _, candidate := range filteredCandidates {
		id, err := candidate.GetString("id")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		filterStr += id
		filterStr += ","
	}

	filterStr = strings.TrimSuffix(filterStr, ",") + ")"

	hostParams := jsonutils.NewDict()
	hostParams.Set("scope", jsonutils.NewString("system"))
	hostParams.Set("show_fail_reason", jsonutils.NewBool(true))
	hostParams.Set("host_type", jsonutils.NewString("hypervisor"))
	hostParams.Set("enabled", jsonutils.NewInt(1))
	hostParams.Set("host_status", jsonutils.NewString("online"))
	hostParams.Set("os_arch", jsonutils.NewString("x86"))
	hostParams.Set("server_id_for_network", jsonutils.NewString(r.VmId))
	hostParams.Set("project_domain", jsonutils.NewString("default"))
	hostParams.Set("filter", jsonutils.NewString(filterStr))
	hostParams.Set("details", jsonutils.NewBool(true))
	hostParams.Set("limit", jsonutils.NewInt(r.PageSize))
	hostParams.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	hostResult, err := compute.Hosts.List(s, hostParams)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	var page int
	if hostResult.Limit != 0 {
		page = (hostResult.Offset / hostResult.Limit) + 1
	}
	response.OkWithDetailed(response.PageResult{
		List:     jsonutils.Marshal(hostResult.Data).Interface(),
		Total:    int64(hostResult.Total),
		Page:     page,
		PageSize: hostResult.Limit,
	}, "获取成功", c)
}

type BaremetalDiskConfig struct {
	Conf    string `json:"conf"`
	Driver  string `json:"driver"`
	Count   int    `json:"count"`
	Range   []int  `json:"range"`
	Adapter int    `json:"adapter"`
	Type    string `json:"type"`
}

func (p *ComputeServersApi) BaremetalCreate(c *gin.Context) {
	var r BaremetalCreateReq

	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	if r.ImageId == "" {
		response.FailWithMessage("镜像参数错误", c)
		return
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	r.PreferRegion = "default"
	r.PreferZone = GetZoneID()
	r.ProjectId = s.GetProjectId()
	r.Count = 1
	r.Hypervisor = "baremetal"
	r.AutoStart = true
	r.Vdi = "vnc"

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("baremetal", jsonutils.NewBool(true))
	params.Set("host_type", jsonutils.NewString("baremetal"))
	params.Set("with_meta", jsonutils.NewBool(true))

	//寻找物理机
	baremetal, err := compute.Hosts.GetById(s, r.PreferHost, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	spec, err := baremetal.Get("spec")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	cpu, err := spec.Int("cpu")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	cpuStr := "cpu:" + strconv.FormatInt(cpu, 10) + "/"

	mem, err := spec.Int("mem")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	memStr := "mem:" + strconv.FormatInt(mem, 10) + "M/"

	r.VmemSize = int(mem)
	r.VcpuCount = int(cpu)

	nicCount, err := baremetal.Int("nic_count")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	nicStr := "nic:" + strconv.FormatInt(nicCount-1, 10)

	disk, err := spec.Get("disk")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	diskMap, err := disk.GetMap()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	diskStr := ""
	for driverName, driverObj := range diskMap {
		diskInfs, err := driverObj.GetArray("adapter0")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		if len(diskInfs) <= 0 {
			continue
		}

		diskInfo := diskInfs[0]
		count, err := diskInfo.Int("count")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		size, err := diskInfo.Int("size")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		diskType, err := diskInfo.GetString("type")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		diskStr = diskStr + "disk:" + driverName + "_adapter0_" + diskType + "_" + strconv.FormatInt(size/1024, 10) + "Gx" + strconv.FormatInt(count, 10) + "/"
	}

	sysInfo, err := baremetal.Get("sys_info")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	manufacture, err := sysInfo.GetString("manufacture")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	// 将空格替换为下划线
	manufacture = strings.ReplaceAll(manufacture, " ", "_")

	manufactureStr := "manufacture:" + manufacture + "/"

	model, err := sysInfo.GetString("model")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 将空格替换为下划线
	model = strings.ReplaceAll(model, " ", "_")

	// 构建 modelStr
	modelStr := "model:" + model + "/"

	skuStr := cpuStr + diskStr + manufactureStr + memStr + modelStr + nicStr

	podListParams := jsonutils.NewDict()
	podListParams.Set("show_emulated", jsonutils.NewBool(true))
	podListParams.Set("resource_type", jsonutils.NewString("shared"))
	//podListParams.Set("project_domain", jsonutils.NewString("default"))
	podListParams.Set("host_type", jsonutils.NewString("baremetal"))

	Capabilities, err := compute.Capabilities.List(s, podListParams)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	if len(Capabilities.Data) > 0 {
		capabilitie := Capabilities.Data[0]
		specs, err := capabilitie.Get("specs")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		hosts, err := specs.Get("hosts")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		if !hosts.Contains(skuStr) {
			response.FailWithMessage("未找到物理机", c)
			return
		}

		isolatedDevices, _ := hosts.GetArray("isolated_devices")

		jsonutils.Update(&r.IsolatedDevices, isolatedDevices)
	}

	//赋值磁盘参数
	diskMap, err = spec.GetMap("disk")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	driver, err := spec.GetString("driver")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	strogeType, err := baremetal.GetString("storage_type")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	for key, value := range diskMap {
		if key == driver {
			diskInfos, err := value.GetArray("adapter0")
			if err != nil {
				response.FailWithMessage(err, c)
				return
			}
			if len(diskInfos) > 0 {
				diskInfo := diskInfos[0]
				size, err := diskInfo.Int("size")
				if err != nil {
					response.FailWithMessage(err, c)
					return
				}

				startIndex, err := diskInfo.Int("start_index")
				if err != nil {
					response.FailWithMessage(err, c)
					return
				}

				r.Disks = []struct {
					Size       int64  `json:"size"`
					ImageId    string `json:"image_id,omitempty"`
					Mountpoint string `json:"mountpoint,omitempty"`
				}{
					{
						Size:    size - 30*1024,
						ImageId: r.ImageId,
					},
					{
						Size:       -1,
						Mountpoint: "/(系统)",
					},
				}

				var baremetalDiskConfig BaremetalDiskConfig

				baremetalDiskConfigs := make([]BaremetalDiskConfig, 0)

				baremetalDiskConfig.Conf = "none"
				baremetalDiskConfig.Type = strogeType
				baremetalDiskConfig.Driver = driver
				baremetalDiskConfig.Count = 1
				baremetalDiskConfig.Adapter = 0
				baremetalDiskConfig.Range = append(baremetalDiskConfig.Range, int(startIndex))

				baremetalDiskConfigs = append(baremetalDiskConfigs, baremetalDiskConfig)

				r.BaremetalDiskConfigs = baremetalDiskConfigs

				break
			}
		}
	}

	if r.BaremetalDiskConfigs == nil || len(r.BaremetalDiskConfigs) == 0 {
		response.FailWithMessage("磁盘参数错误", c)
		return
	}

	_, err = scheduler.SchedManager.DoForecast(s, jsonutils.Marshal(r))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	result, err := compute.Servers.Create(s, jsonutils.Marshal(r))
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 查询套餐列表
func (p *ComputeServersApi) GetResSkuList(c *gin.Context) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	result, err := compute.ServerSkus.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

type setPasswordReq struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 重置密码
func (p *ComputeServersApi) ResetPassword(c *gin.Context) {
	var r setPasswordReq
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
	params.Set("auto_start", jsonutils.NewBool(true))
	params.Set("reset_password", jsonutils.NewBool(true))
	if r.Password != "" {
		params.Set("password", jsonutils.NewString(r.Password))
	}
	params.Set("username", jsonutils.NewString(r.Username))

	result, err := compute.Servers.PerformAction(s, r.ID, "set-password", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "设置成功", c)
}

type updateNameReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error struct {
		Class   string    `json:"class"`
		Code    int       `json:"code"`
		Details string    `json:"details"`
		Time    time.Time `json:"request"`
	} `json:"error"`
}

// 修改名称
func (p *ComputeServersApi) UpdateName(c *gin.Context) {
	var r updateNameReq
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
	params.Set("name", jsonutils.NewString(r.Name))

	result, err := compute.Servers.Update(s, r.ID, params)
	if err != nil {
		var errorResponse ErrorResponse
		// 解析 JSON 字符串
		json.Unmarshal([]byte(err.Error()), &errorResponse)
		message := ""
		if errorResponse.Error.Code == 409 {
			message = "与系统已有主机重名"
		} else if errorResponse.Error.Code == 400 {
			message = "修改失败,请重新输入名称"
		}
		response.FailWithMessage(message, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "修改成功", c)
}

// 监控云主机
type MonitorDataReq struct {
	ID       string `json:"id"`
	Time     string `json:"time"`
	Interval string `json:"interval"`
}

// 监控云主机CPU
func (p *ComputeServersApi) GetHostMonitorCPUData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorCpuData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"cpu": cpuData,
	}, "获取成功", c)
}

// 监控云主机内存
func (p *ComputeServersApi) GetHostMonitorMemData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorMemData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"mem": cpuData,
	}, "获取成功", c)
}

// 监控云主机磁盘使用率
func (p *ComputeServersApi) GetHostMonitorDiskData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorDiskData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"disk": cpuData,
	}, "获取成功", c)
}

// 监控云主机磁盘读速率
func (p *ComputeServersApi) GetHostMonitorDiskReadData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorDiskRead(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"diskRead": cpuData,
	}, "获取成功", c)
}

// 监控云主机磁盘写速率
func (p *ComputeServersApi) GetHostMonitorDiskWriteData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorDiskWrite(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"diskWrite": cpuData,
	}, "获取成功", c)
}

// 监控云主机网络入流量
func (p *ComputeServersApi) GetHostMonitorBpsRecvData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorBpsRecvData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"bpsRecv": cpuData,
	}, "获取成功", c)
}

// 监控云主机网络出流量
func (p *ComputeServersApi) GetHostMonitorBpsSentData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorBpsSentData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"bpsSent": cpuData,
	}, "获取成功", c)
}

// 监控云主机网络收包数
func (p *ComputeServersApi) GetHostMonitorPpsRecvData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorPpsRecvData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"ppsRecv": cpuData,
	}, "获取成功", c)
}

// 监控云主机网络发包数
func (p *ComputeServersApi) GetHostMonitorPpsSentData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := serverService.MonitorPpsSentData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"ppsSent": cpuData,
	}, "获取成功", c)
}

// 监控GPU使用率
func (p *ComputeServersApi) GetHostMonitorGPUData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//gpu
	gpuData, err := serverService.MonitorGpuData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"utilization_gpu": gpuData,
	}, "获取成功", c)
}

// 监控GPU显存使用率
func (p *ComputeServersApi) GetHostMonitorGPUMemData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//gpu
	gpuMemData, err := serverService.MonitorGpuMemData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"utilization_memory": gpuMemData,
	}, "获取成功", c)
}

// 监控GPU温度
func (p *ComputeServersApi) GetHostMonitorGpuTemperature(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	gpuTemData, err := serverService.MonitorGpuTemperatureData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"temperature_gpu": gpuTemData,
	}, "获取成功", c)
}

type ForecastAgentReq struct {
	Id string `json:"id"` // 实例id
}

func (p *ComputeServersApi) ForecastAgent(c *gin.Context) {
	var r ForecastAgentReq
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
	params.Set("server_id", jsonutils.NewString(r.Id))
	params.Set("details", jsonutils.NewBool(false))
	params.Set("limit", jsonutils.NewInt(1))
	result, err := devtool.DevToolScriptApplyRecords.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result.Data).Interface(), "获取成功", c)
}

type ForecastAgentCanInstallReq struct {
	Id string `json:"id"` // 实例id
}

func (p *ComputeServersApi) ForecastAgentCanInstall(c *gin.Context) {
	var r ForecastAgentCanInstallReq
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

	result, err := compute.Servers.GetSpecific(s, r.Id, "sshable", nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	methods, err := result.GetArray("method_tried")
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	agentCanInstall := true

	for _, method := range methods {
		sshable, err := method.Bool("sshable")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		if !sshable {
			agentCanInstall = false
			break
		}
	}

	response.OkWithDetailed(agentCanInstall, "获取成功", c)
}

type SetAgentCanInstallReq struct {
	Id       string `json:"id"`       //实例id
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
	Port     int64  `json:"port"`     //ssh端口
}

func (p *ComputeServersApi) SetAgentCanInstall(c *gin.Context) {
	var r SetAgentCanInstallReq
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
	params.Set("user", jsonutils.NewString(r.Username))
	params.Set("password", jsonutils.NewString(r.Password))
	params.Set("port", jsonutils.NewInt(r.Port))

	_, err = compute.Servers.PerformAction(s, r.Id, "make-sshable", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("设置成功,请稍等片刻后安装agent", c)
}

type InstallAgentReq struct {
	Id string `json:"id"` // 实例id
}

func (p *ComputeServersApi) InstallAgent(c *gin.Context) {
	var r InstallAgentReq
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
	params.Set("server_id", jsonutils.NewString(r.Id))
	params.Set("auto_choose_proxy_endpoint", jsonutils.NewBool(true))

	_, err = devtool.DevToolScripts.PerformAction(s, "monitor agent", "apply", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("安装中,请稍等!", c)
}
