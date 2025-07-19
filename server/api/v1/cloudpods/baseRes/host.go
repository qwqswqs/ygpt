package baseRes

import (
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"yunion.io/x/onecloud/pkg/mcclient/modules/webconsole"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type BaseHostApi struct {
}
type getBaseHostReq struct {
	PageIndex   int64  `json:"pageIndex"`
	PageSize    int64  `json:"pageSize"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     *bool  `json:"enabled"`
	OsArch      string `json:"osArch"`
}

// 获取宿主机列表
func (d *BaseHostApi) GetBaseHostList(c *gin.Context) {
	var r getBaseHostReq
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
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("baremetal", jsonutils.NewBool(false))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Enabled != nil {
		params.Set("enabled", jsonutils.NewBool(*r.Enabled))
	}
	if r.OsArch != "" {
		params.Set("os_arch", jsonutils.NewString(r.OsArch))
	}
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Hosts.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 获取物理机列表
func (d *BaseHostApi) GetBaseBareHostList(c *gin.Context) {
	var r getBaseHostReq
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
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Description != "" {
		filters = append(filters, "description.contains('"+r.Description+"')")
	}
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("detail", jsonutils.NewBool(true))
	params.Set("baremetal", jsonutils.NewBool(true))
	params.Set("host_type", jsonutils.NewString("baremetal"))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	if r.Enabled != nil {
		params.Set("enabled", jsonutils.NewBool(*r.Enabled))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Hosts.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

type getBaseHostIdReq struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// 获取物理机IPMI账号密码
func (d *BaseHostApi) GetImpiInfo(c *gin.Context) {
	var r getBaseHostIdReq
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
	result, err := compute.Hosts.GetIpmiInfo(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 获取物理机登录账号密码
func (d *BaseHostApi) GetLoginInfo(c *gin.Context) {
	var r getBaseHostIdReq
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
	var result jsonutils.JSONObject
	params := jsonutils.NewDict()
	if r.Status == "host" {
		result, err = compute.Hosts.GetLoginInfo(s, r.ID, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	} else {
		result, err = compute.Servers.GetLoginInfo(s, r.ID, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 物理机修改状态
func (d *BaseHostApi) UpdateBareHostStatus(c *gin.Context) {
	var r getBaseHostIdReq
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
	result, err := compute.Hosts.PerformAction(s, r.ID, r.Status, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "成功", c)
}

// 物理机修改备注
func (d *BaseHostApi) UpdateBareHostDescription(c *gin.Context) {
	var r getBaseHostIdReq
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
	params.Set("description", jsonutils.NewString(r.Description))
	result, err := compute.Hosts.Update(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "成功", c)
}

type deleteBareHostByIds struct {
	Ids []string `json:"ids"`
}

// 删除物理机
func (d *BaseHostApi) DeleteBareHost(c *gin.Context) {
	var r getBaseHostIdReq
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
	result, err := compute.Hosts.Delete(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

// 批量删除物理机
func (d *BaseHostApi) DeleteBareHostByIds(c *gin.Context) {
	var r deleteBareHostByIds
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
		_, err := compute.Hosts.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

// 创建物理机
type addBareHostReq struct {
	Mac      string `json:"mac"`
	Net      string `json:"net"`
	Ip       string `json:"ip"`
	Password string `json:"password"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// 新建物理机
func (d *BaseHostApi) AddBareHost(c *gin.Context) {
	var r addBareHostReq
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
	params.Set("access_mac", jsonutils.NewString(r.Mac))
	params.Set("access_net", jsonutils.NewString(r.Net))
	params.Set("enable_pxe_boot", jsonutils.NewBool(true))
	params.Set("ipmi_ip_addr", jsonutils.NewString(r.Ip))
	params.Set("ipmi_password", jsonutils.NewString(r.Password))
	params.Set("ipmi_username", jsonutils.NewString(r.Username))
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("no_prepare", jsonutils.NewBool(false))
	params.Set("project_domain", jsonutils.NewString("default"))
	result, err := compute.Hosts.Create(s, params)
	if err != nil {
		var errorResponse ErrorResponse
		// 解析 JSON 字符串
		json.Unmarshal([]byte(err.Error()), &errorResponse)
		message := ""
		if errorResponse.Error.Code == 409 {
			message = "创建失败，检查名称、ip或mac是否已被使用"
		} else if errorResponse.Error.Code == 400 {
			message = "创建失败，IP 地址超出了系统中网络的IP范围"
		} else {
			message = "创建失败"
		}
		response.FailWithMessage(message, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 远程连接物理机
type webConsoleReq struct {
	ID string `json:"id"`
}

// 云主机远程连接
func (p *BaseHostApi) GetWebConsoleSSH(c *gin.Context) {
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
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := webconsole.WebConsole.DoBaremetalConnect(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 监控宿主机
type MonitorDataReq struct {
	ID       string `json:"id"`
	Time     string `json:"time"`
	Interval string `json:"interval"`
}

// 监控宿主机CPU
func (p *BaseHostApi) GetHostMonitorCPUData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorCpuData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"cpu": cpuData,
	}, "获取成功", c)
}

// 监控宿主机内存
func (p *BaseHostApi) GetHostMonitorMemData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorMemData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"mem": cpuData,
	}, "获取成功", c)
}

// 监控宿主机磁盘使用率
func (p *BaseHostApi) GetHostMonitorDiskData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorDiskData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"disk": cpuData,
	}, "获取成功", c)
}

// 监控宿主机磁盘读速率
func (p *BaseHostApi) GetHostMonitorDiskReadData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorDiskRead(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"diskRead": cpuData,
	}, "获取成功", c)
}

// 监控宿主机磁盘写速率
func (p *BaseHostApi) GetHostMonitorDiskWriteData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorDiskWrite(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"diskWrite": cpuData,
	}, "获取成功", c)
}

// 监控宿主机网络入流量
func (p *BaseHostApi) GetHostMonitorBpsRecvData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorBpsRecvData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"bpsRecv": cpuData,
	}, "获取成功", c)
}

// 监控宿主机网络出流量
func (p *BaseHostApi) GetHostMonitorBpsSentData(c *gin.Context) {
	var r MonitorDataReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//cpu
	cpuData, err := bareHostService.MonitorBpsSentData(r.ID, r.Time, r.Interval)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(gin.H{
		"bpsSent": cpuData,
	}, "获取成功", c)
}
