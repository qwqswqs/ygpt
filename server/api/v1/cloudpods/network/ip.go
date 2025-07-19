package network

import (
	"encoding/json"
	compute2 "ygpt/server/api/v1/cloudpods/compute"
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type NetIpApi struct {
}
type getNetIpReq struct {
	PageIndex   int64    `json:"pageIndex"`
	PageSize    int64    `json:"pageSize"`
	ServerType  []string `json:"server_type"`
	VpcId       string   `json:"vpcId"`
	Name        string   `json:"name"`
	GatewayDesc string   `json:"gatewayDesc"`
}

// 查询IP列表
func (p *NetIpApi) GetNetIPList(c *gin.Context) {
	var r getNetIpReq
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
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if len(r.ServerType) != 0 {
		params.Set("server_type", jsonutils.NewStringArray(r.ServerType))
	}
	if r.VpcId != "" {
		params.Set("vpc", jsonutils.NewString(r.VpcId))
	}
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.GatewayDesc != "" {
		params.Set("order_by", jsonutils.NewString("guest_gateway"))
		params.Set("order", jsonutils.NewString(r.GatewayDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Networks.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

type addNetIpReq struct {
	Name       string `json:"name"`
	ServerType string `json:"serverType"`
	IpPrefix   string `json:"ipPrefix"`
	IpMask     int64  `json:"ipMask"`
	IpGateway  string `json:"ipGateway"`
	IpStart    string `json:"ipStart"`
	IpEnd      string `json:"ipEnd"`
	Type       string `json:"type"`
	VpcId      string `json:"vpcId"`
	ID         string `json:"id"`
	WireID     string `json:"wireId"`
}

// 新增IP子网
func (p *NetIpApi) AddNetIP(c *gin.Context) {

	var r addNetIpReq
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
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("alloc_policy", jsonutils.NewString("none"))
	params.Set("project_id", jsonutils.NewString(s.GetProjectId()))

	if r.Type == "default" {
		params.Set("guest_gateway", jsonutils.NewString(r.IpGateway))
		params.Set("guest_ip_end", jsonutils.NewString(r.IpEnd))
		params.Set("guest_ip_start", jsonutils.NewString(r.IpStart))
		params.Set("guest_ip_mask", jsonutils.NewInt(r.IpMask))
		params.Set("server_type", jsonutils.NewString(r.ServerType))
		params.Set("wire_id", jsonutils.NewString(compute2.GetWireID()))
	} else {
		params.Set("guest_ip_prefix", jsonutils.NewString(r.IpPrefix))
		params.Set("zone", jsonutils.NewString(compute2.GetZoneID()))
		params.Set("vpc", jsonutils.NewString(r.VpcId))
	}

	result, err := compute.Networks.Create(s, params)
	if err != nil {
		var errorResponse ErrorResponse

		// 解析 JSON 字符串
		json.Unmarshal([]byte(err.Error()), &errorResponse)
		message := ""
		if errorResponse.Error.Code == 409 {
			message = "子网重名"
		} else if errorResponse.Error.Code == 400 {
			message = "ip与vpc 中的现有ip冲突"
		}
		response.FailWithMessage(message, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 编辑IP子网
func (p *NetIpApi) UpdateNetIP(c *gin.Context) {

	var r addNetIpReq
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
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("alloc_policy", jsonutils.NewString("none"))
	params.Set("platform_type", jsonutils.NewString("idc"))

	params.Set("guest_gateway", jsonutils.NewString(r.IpGateway))
	params.Set("guest_ip_end", jsonutils.NewString(r.IpEnd))
	params.Set("guest_ip_start", jsonutils.NewString(r.IpStart))
	params.Set("guest_ip_mask", jsonutils.NewInt(r.IpMask))
	params.Set("server_type", jsonutils.NewString(r.ServerType))
	params.Set("wire_id", jsonutils.NewString(r.WireID))

	result, err := compute.Networks.Update(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "修改成功", c)
}

type deleteNetIpReq struct {
	ID string `json:"id"`
}
type deleteNetIpByIdsReq struct {
	Ids []string `json:"ids"`
}

// 删除IP
func (p *NetIpApi) DeleteNetIP(c *gin.Context) {
	var r deleteNetIpReq
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
	result, err := compute.Networks.Delete(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

// 批量删除IP
func (p *NetIpApi) DeleteNetIpByIds(c *gin.Context) {
	var r deleteNetIpByIdsReq
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
		_, err := compute.Networks.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
