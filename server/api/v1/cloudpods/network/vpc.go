package network

import (
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type NetworkVpcApi struct {
}
type getVpcReq struct {
	PageIndex  int64  `json:"pageIndex"`
	PageSize   int64  `json:"pageSize"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	Status     string `json:"status"`
	IpNumDesc  string `json:"ipNumDesc"`
	IsExternal string `json:"isExternal"`
}

// 获取VPC列表
func (d *NetworkVpcApi) GetVpcList(c *gin.Context) {
	var r getVpcReq
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
	//params.Set("usable", jsonutils.NewBool(true))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.contains('"+r.Status+"')")
	}
	if r.Id != "" {
		filters = append(filters, "id.contains('"+r.Id+"')")
	}
	if r.IpNumDesc != "" {
		params.Set("order_by", jsonutils.NewString("network_count"))
		params.Set("order", jsonutils.NewString(r.IpNumDesc))
	}
	if r.IsExternal != "" {
		params.Set("external_access_mode", jsonutils.NewString(r.IsExternal))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Vpcs.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 新增VPC
type addVpcReq struct {
	Name       string `json:"name"`
	Cidr       string `json:"cidr"`
	IsExternal bool   `json:"isExternal"`
}

func (d *NetworkVpcApi) AddVpc(c *gin.Context) {
	var r addVpcReq
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
	param.Set("cidr_block", jsonutils.NewString(r.Cidr))

	param.Set("cloudregion_id", jsonutils.NewString("default"))
	param.Set("external_access_mode", jsonutils.NewString("none"))
	param.Set("project_domain", jsonutils.NewString("default"))
	if r.IsExternal {
		param.Set("external_access_mode", jsonutils.NewString("eip-distgw"))
	} else {
		param.Set("external_access_mode", jsonutils.NewString("none"))
	}
	result, err := compute.Vpcs.Create(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 删除VPC
type deleteVpcReq struct {
	Id string `json:"id"`
}
type deleteVpcByIds struct {
	Ids []string `json:"ids"`
}

func (d *NetworkVpcApi) DeleteVpc(c *gin.Context) {
	var r deleteVpcReq
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
	result, err := compute.Vpcs.Delete(s, r.Id, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}
func (d *NetworkVpcApi) DeleteVpcByIds(c *gin.Context) {
	var r deleteVpcByIds
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
		_, err := compute.Vpcs.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
