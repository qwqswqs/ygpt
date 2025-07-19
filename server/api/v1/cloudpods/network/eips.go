package network

import (
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type NetworkEipApi struct {
}
type getEipReq struct {
	PageIndex     int64  `json:"pageIndex"`
	PageSize      int64  `json:"pageSize"`
	Name          string `json:"name"`
	BandWidthDesc string `json:"bandWidthDesc"`
	IpDesc        string `json:"ipDesc"`
}

// 获取弹性公网IP列表
func (d *NetworkEipApi) GetEipList(c *gin.Context) {
	var r getEipReq
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
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		params.Set("filter", jsonutils.NewString("name.contains('"+r.Name+"')"))
	}
	if r.IpDesc != "" {
		params.Set("order_by", jsonutils.NewString("ips"))
		params.Set("order", jsonutils.NewString(r.IpDesc))
	}
	if r.BandWidthDesc != "" {
		params.Set("order_by", jsonutils.NewString("bandwidth"))
		params.Set("order", jsonutils.NewString(r.BandWidthDesc))
	}
	result, err := compute.Elasticips.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 新增VPC
type addEipReq struct {
	Name      string `json:"name"`
	Network   string `json:"network"`
	Bandwidth int64  `json:"bandwidth"`
}

func (d *NetworkEipApi) AddEip(c *gin.Context) {
	var r addEipReq
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
	param.Set("name", jsonutils.NewString(r.Name))
	param.Set("network", jsonutils.NewString(r.Network))
	param.Set("bandwidth", jsonutils.NewInt(r.Bandwidth))
	param.Set("charge_type", jsonutils.NewString("bandwidth"))
	param.Set("vpc", jsonutils.NewString("default"))

	param.Set("tenant", jsonutils.NewString(s.GetTenantId()))
	result, err := compute.Elasticips.Create(s, param)
	if err != nil {
		var errorResponse ErrorResponse

		// 解析 JSON 字符串
		json.Unmarshal([]byte(err.Error()), &errorResponse)
		message := ""
		if errorResponse.Error.Code == 409 {
			message = "子网重名"
		} else {
			message = "创建失败"
		}
		response.FailWithMessage(message, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 删除VPC
type deleteEipReq struct {
	Id string `json:"id"`
}

// 删除VPC
type deleteEipByIdsReq struct {
	Ids []string `json:"ids"`
}

func (d *NetworkEipApi) DeleteEip(c *gin.Context) {
	var r deleteEipReq
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
	param.Set("id", jsonutils.NewString(r.Id))
	result, err := compute.Elasticips.Delete(s, r.Id, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}
func (d *NetworkEipApi) DeleteEipByIds(c *gin.Context) {
	var r deleteEipByIdsReq
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
		param.Set("id", jsonutils.NewString(id))
		_, err := compute.Elasticips.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

// 获取可绑定主机列表
type getBindServerReq struct {
	Id string `json:"id"`
}

func (d *NetworkEipApi) GetBindServerList(c *gin.Context) {
	var r getBindServerReq
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
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("without_eip", jsonutils.NewBool(true))
	params.Set("eip_associable", jsonutils.NewBool(true))
	params.Set("filter", jsonutils.NewString("status.in(ready, running)"))
	params.Set("usable_server_for_eip", jsonutils.NewString(r.Id))
	result, err := compute.Servers.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

// 绑定
type bindEipServerReq struct {
	Id        string `json:"id"`
	ServerId  string `json:"serverId"`
	BandWidth int64  `json:"bandwidth"`
}

func (d *NetworkEipApi) BindEipServer(c *gin.Context) {
	var r bindEipServerReq
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
	params.Set("id", jsonutils.NewString(r.Id))
	params.Set("instance_id", jsonutils.NewString(r.ServerId))
	params.Set("instance_type", jsonutils.NewString("server"))
	result, err := compute.Elasticips.PerformAction(s, r.Id, "associate", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "绑定成功", c)
}

// 解绑
func (d *NetworkEipApi) UnbindEipServer(c *gin.Context) {
	var r bindEipServerReq
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
	result, err := compute.Elasticips.PerformAction(s, r.Id, "dissociate", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "解绑成功", c)
}

func (d *NetworkEipApi) ModEipBandWidth(c *gin.Context) {
	var r bindEipServerReq
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
	params.Set("id", jsonutils.NewString(r.Id))
	params.Set("bandwidth", jsonutils.NewInt(r.BandWidth))
	result, err := compute.Elasticips.PerformAction(s, r.Id, "change-bandwidth", params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "修改成功", c)
}
