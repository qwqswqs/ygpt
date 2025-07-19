package cloudpods

import (
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	"yunion.io/x/onecloud/pkg/mcclient/modules/monitor"
	"yunion.io/x/onecloud/pkg/mcclient/modules/webconsole"
	options "yunion.io/x/onecloud/pkg/mcclient/options/monitor"
)

type VirtualMachineApi struct {
}

type hostBaseReq struct {
	ID string `json:"id"`
}

type hostListReq struct {
	IDs []string `json:"ids"`
}

type hostBandwidthLimitReq struct {
	ID        string `json:"id"`
	Bandwidth int    `json:"bandwidth"`
}

type vmCreateReq struct {
	Name string `json:"name"`
	//Tag    string `json:"tag"`
	Cpu             int    `json:"cpu"`
	Memory          int    `json:"memory"`
	Storage         int    `json:"storage"`
	SysStorageType  string `json:"sys_storage_type"`
	DataStorageType string `json:"data_storage_type"`
	SysDisk         int    `json:"sys_storage"`
	DataDisk        int    `json:"data_storage"`
	BandWidth       int    `json:"bandwidth"`
}

// 创建云主机
func (v *VirtualMachineApi) HostCreate(c *gin.Context) {
	//var r vmCreateReq
	//err := c.ShouldBindJSON(&r)
	//if err != nil {
	//	response.FailWithMessage(err, c)
	//	return
	//}
	//
	//hostRet, err := VirtualMachineService.HostCreate(r.Name, r.Cpu, r.Memory, r.SysStorageType, r.SysDisk, r.DataDisk)
	//if err != nil {
	//	response.FailWithMessage("创建云主机失败:"+err.Error(), c)
	//	return
	//}
	//response.OkWithDetailed(jsonutils.Marshal(hostRet).PrettyString(), "创建成功", c)
	response.FailWithMessage("该接口已经隐藏!", c)
}

func (v *VirtualMachineApi) HostListByIds(c *gin.Context) {
	var r hostListReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	hostRet, total, err := VirtualMachineService.HostListByIds(r.IDs)
	if err != nil {
		response.FailWithMessage("获取云主机信息失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     jsonutils.Marshal(hostRet).PrettyString(),
		Total:    int64(total),
		Page:     0,
		PageSize: 0,
	}, "获取成功", c)
}
func (v *VirtualMachineApi) HostBandwidthLimit(c *gin.Context) {
	var r hostBandwidthLimitReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = VirtualMachineService.HostBandwidthLimit(r.ID, r.Bandwidth)
	if err != nil {
		response.FailWithMessage("限制带宽失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("限制带宽成功", c)
}

// 删除云主机
func (v *VirtualMachineApi) HostDelete(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = VirtualMachineService.HostDelete(r.ID)
	if err != nil {
		response.FailWithMessage("删除云主机失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// 关机
func (v *VirtualMachineApi) HostShutDown(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := compute.Servers.PerformAction(s, r.ID, "stop", params)
	if err != nil {
		response.FailWithMessage("关机失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "关机成功", c)
}

// 开机
func (v *VirtualMachineApi) HostStart(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := compute.Servers.PerformAction(s, r.ID, "start", params)
	if err != nil {
		response.FailWithMessage("开机失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "开机成功", c)
}

// 重启
func (v *VirtualMachineApi) HostRestart(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := compute.Servers.PerformAction(s, r.ID, "restart", params)
	if err != nil {
		response.FailWithMessage("重启失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "重启成功", c)
}

// 挂起
func (v *VirtualMachineApi) HostSuspend(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := compute.Servers.PerformAction(s, r.ID, "suspend", params)
	if err != nil {
		response.FailWithMessage("挂起失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "挂起成功", c)
}

// 恢复
func (v *VirtualMachineApi) HostResume(c *gin.Context) {
	var r hostBaseReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := compute.Servers.PerformAction(s, r.ID, "resume", params)
	if err != nil {
		response.FailWithMessage("恢复失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "恢复成功", c)
}

type webConsoleReq struct {
	ID          string `json:"id"`          //必传
	Port        int64  `json:"port"`        //可不传，默认22
	Username    string `json:"username"`    //可不传
	Password    string `json:"password"`    //可不传
	ConsoleType string `json:"consoleType"` //必传，ssh或vnc
}

// 云主机远程连接
func (v *VirtualMachineApi) GetWebConsoleSSH(c *gin.Context) {
	var r webConsoleReq

	var result jsonutils.JSONObject

	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("云主机远程连接失败:"+err.Error(), c)
		return
	}

	switch r.ConsoleType {
	case "ssh":
		vm, err := compute.Servers.Get(s, r.ID, nil)
		if err != nil {
			response.FailWithMessage("云主机远程连接失败:"+err.Error(), c)
			return
		}

		ip, err := vm.GetString("ips")
		if err != nil {
			response.FailWithMessage("云主机远程连接失败:"+err.Error(), c)
			return
		}

		if r.Port == 0 {
			r.Port = 22
		}

		params := jsonutils.NewDict()
		params.Set("id", jsonutils.NewString(r.ID))
		params.Set("ip", jsonutils.NewString(ip))
		params.Set("port", jsonutils.NewInt(r.Port))
		params.Set("type", jsonutils.NewString("server"))
		if r.Username != "" && r.Password != "" {
			params.Set("username", jsonutils.NewString(r.Username))
			params.Set("password", jsonutils.NewString(r.Password))
		}

		body := jsonutils.NewDict()
		body.Set("webconsole", params)

		result, err = webconsole.WebConsole.DoSshConnect(s, r.ID, body)
		if err != nil {
			response.FailWithMessage("云主机远程连接失败:"+err.Error(), c)
			return
		}
	case "vnc":
		result, err = webconsole.WebConsole.DoServerConnect(s, r.ID, nil)
		if err != nil {
			response.FailWithMessage("云主机远程连接失败:"+err.Error(), c)
			return
		}
	default:
		response.FailWithMessage("云主机远程连接失败:未知连接类型", c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "云主机远程连接成功", c)
}

// 云主机监控
func (v *VirtualMachineApi) HostMonitor(c *gin.Context) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}

	// 创建查询选项
	opts := &options.MetricQueryOptions{
		Interval: "30m",
	}

	// 获取查询输入
	input, err := opts.GetQueryInput()
	if err != nil {
		response.FailWithMessage("监控失败：无法获取输入", c)
	}

	// 执行查询
	resp, err := monitor.UnifiedMonitorManager.PerformQuery(s, input)
	if err != nil {
		response.FailWithMessage("监控失败", c)
	}

	response.OkWithDetailed(jsonutils.Marshal(resp).PrettyString(), "获取成功", c)
}

// 获取云主机列表
func (v *VirtualMachineApi) HostGet(c *gin.Context) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("filter", jsonutils.NewString("hypervisor.notin(baremetal,container)"))
	result, _ := compute.Servers.List(s, params)

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

type createDiskReq struct {
	Disk         string `json:"disk"`
	GenerateName string `json:"generateName"`
}

// 创建硬盘快照
func (v *VirtualMachineApi) HostCreateDiskSnap(c *gin.Context) {
	var r createDiskReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("disk", jsonutils.NewString(r.Disk))
	params.Set("generate_name", jsonutils.NewString(r.GenerateName))
	result, err := compute.Snapshots.Create(s, params)
	if err != nil {
		response.FailWithMessage("创建失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 创建主机快照
type createServerSnapReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (v *VirtualMachineApi) HostCreateServerSnap(c *gin.Context) {
	var r createServerSnapReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("generate_name", jsonutils.NewString(r.Name))
	result, err := compute.Servers.PerformAction(s, r.ID, "instance-snapshot", params)

	if err != nil {
		response.FailWithMessage("创建失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 回滚主机快照
type resetServerSnapReq struct {
	ServerID string `json:"serverID"`
	SnapID   string `json:"snapID"`
}

func (v *VirtualMachineApi) HostResetServerSnap(c *gin.Context) {
	var r resetServerSnapReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("instance_snapshot", jsonutils.NewString(r.SnapID))
	params.Set("auto_start", jsonutils.NewBool(false))
	result, err := compute.Servers.PerformAction(s, r.ServerID, "instance-snapshot-reset", params)

	if err != nil {
		response.FailWithMessage("创建失败", c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}
