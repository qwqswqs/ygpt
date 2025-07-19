package network

import (
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type NetworkSecgroupApi struct {
}
type getSecgroupReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	Name      string `json:"name"`
	NumDesc   string `json:"numDesc"`
}
type getSecgroupRuleReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	Direction string `json:"direction"`
	Id        string `json:"id"`
}

// 获取VPC列表
func (d *NetworkSecgroupApi) GetSecgroupList(c *gin.Context) {
	var r getSecgroupReq
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
	params.Set("detail", jsonutils.NewBool(true))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		params.Set("filter", jsonutils.NewString("name.contains('"+r.Name+"')"))
	}
	if r.NumDesc != "" {
		params.Set("order_by", jsonutils.NewString("guest_cnt"))
		params.Set("order", jsonutils.NewString(r.NumDesc))
	}
	result, err := compute.SecGroups.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}
func (d *NetworkSecgroupApi) GetSecgroupRuleList(c *gin.Context) {
	var r getSecgroupRuleReq
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
	params.Set("order_by", jsonutils.NewString("priority"))
	params.Set("order", jsonutils.NewString("desc"))
	params.Set("direction", jsonutils.NewString(r.Direction))
	params.Set("secgroup", jsonutils.NewString(r.Id))
	params.Set("detail", jsonutils.NewBool(true))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	result, err := compute.SecGroupRules.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 新增VPC
type addSecgroupReq struct {
	Name string `json:"name"`
}
type addSecgroupRuleReq struct {
	Action      string `json:"action"`
	Direction   string `json:"direction"`
	Priority    int64  `json:"priority"`
	Secgroup    string `json:"secgroup"`
	Type        string `json:"type"`
	Ports       string `json:"ports"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	Cidr        string `json:"cidr"`
}

func (d *NetworkSecgroupApi) AddSecgroup(c *gin.Context) {
	var r addSecgroupReq
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
	param.Set("project_id", jsonutils.NewString(s.GetProjectId()))
	result, err := compute.SecGroups.Create(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}
func (d *NetworkSecgroupApi) AddSecgroupRule(c *gin.Context) {
	var r addSecgroupRuleReq
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
	param.Set("action", jsonutils.NewString(r.Action))
	param.Set("cidr", jsonutils.NewString(r.Cidr))
	param.Set("description", jsonutils.NewString(r.Description))
	param.Set("direction", jsonutils.NewString(r.Direction))
	param.Set("ports", jsonutils.NewString(r.Ports))
	param.Set("priority", jsonutils.NewInt(r.Priority))
	param.Set("protocol", jsonutils.NewString(r.Protocol))
	param.Set("secgroup", jsonutils.NewString(r.Secgroup))
	param.Set("type", jsonutils.NewString(r.Type))
	result, err := compute.SecGroupRules.Create(s, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 删除VPC
type deleteSecgroupReq struct {
	Id string `json:"id"`
}

// 删除VPC
type deleteSecgroupByIdsReq struct {
	Ids []string `json:"ids"`
}

func (d *NetworkSecgroupApi) DeleteSecgroup(c *gin.Context) {
	var r deleteSecgroupReq
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
	result, err := compute.SecGroups.Delete(s, r.Id, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}
func (d *NetworkSecgroupApi) DeleteSecgroupByIds(c *gin.Context) {
	var r deleteSecgroupByIdsReq
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
		_, err := compute.SecGroups.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
func (d *NetworkSecgroupApi) DeleteSecgroupRule(c *gin.Context) {
	var r deleteSecgroupReq
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
	result, err := compute.SecGroupRules.Delete(s, r.Id, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}
func (d *NetworkSecgroupApi) DeleteSecgroupRuleByIds(c *gin.Context) {
	var r deleteSecgroupByIdsReq
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
		_, err := compute.SecGroupRules.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

// 删除VPC
type updateSecgroupRuleReq struct {
	Id          string `json:"id"`
	Action      string `json:"action"`
	Cidr        string `json:"cidr"`
	Description string `json:"description"`
	Ports       string `json:"ports"`
	Priority    int64  `json:"priority"`
	Protocol    string `json:"protocol"`
}

func (d *NetworkSecgroupApi) UpdateSecgroupRule(c *gin.Context) {
	var r updateSecgroupRuleReq
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
	param.Set("action", jsonutils.NewString(r.Action))
	param.Set("cidr", jsonutils.NewString(r.Cidr))
	param.Set("description", jsonutils.NewString(r.Description))
	param.Set("ports", jsonutils.NewString(r.Ports))
	param.Set("priority", jsonutils.NewInt(r.Priority))
	param.Set("protocol", jsonutils.NewString(r.Protocol))
	result, err := compute.SecGroupRules.Update(s, r.Id, param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "修改成功", c)
}

type addServerSecgroupRuleReq struct {
	Ids      []string `json:"id"`
	SecGroup string   `json:"secGroup"`
}

// 绑定虚拟机
func (d *NetworkSecgroupApi) AddSeverSecgroupRule(c *gin.Context) {
	var r addServerSecgroupRuleReq
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
	secId := []string{
		r.SecGroup,
	}
	for _, id := range r.Ids {
		param := jsonutils.NewDict()
		param.Set("id", jsonutils.NewString(id))
		param.Set("secgroup_ids", jsonutils.NewStringArray(secId))
		_, err = compute.Servers.PerformAction(s, id, "add-secgroup", param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}
	response.OkWithData("修改成功", c)
}

// 解绑虚拟机
func (d *NetworkSecgroupApi) RevokeSeverSecgroupRule(c *gin.Context) {
	var r addServerSecgroupRuleReq
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
	secId := []string{
		r.SecGroup,
	}
	for _, id := range r.Ids {
		param := jsonutils.NewDict()
		param.Set("id", jsonutils.NewString(id))
		param.Set("secgroup_ids", jsonutils.NewStringArray(secId))
		_, err = compute.Servers.PerformAction(s, id, "revoke-secgroup", param)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}
	response.OkWithData("修改成功", c)
}
