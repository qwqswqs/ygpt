package cloudpods

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type IpApi struct {
}
type getNetIpReq struct {
	PageIndex  int64  `json:"pageIndex"`
	PageSize   int64  `json:"pageSize"`
	ServerType string `json:"server_type"`
	Name       string `json:"name"`
}

// 查询IP列表
func (p *IpApi) GetNetIPList(c *gin.Context) {
	var r getNetIpReq
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
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.ServerType != "" {
		params.Set("server_type", jsonutils.NewString(r.ServerType))
	}
	if r.Name != "" {
		params.Set("filter", jsonutils.NewString("name.contains('"+r.Name+"')"))
	}
	result, _ := compute.Networks.List(s, params)

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// todo:删除IP子网
// todo:新增IP子网
type addNetIpReq struct {
	guest_gateway  int64 `json:"guest_gateway"`
	guest_ip_end   int64 `json:"guest_ip_end"`
	guest_ip_mask  int64 `json:"guest_ip_mask"`
	guest_ip_start int64 `json:"guest_ip_start"`
	name           int64 `json:"name"`
}

// 新增IP子网
func (p *IpApi) AddNetIP(c *gin.Context) {
	var r getNetIpReq
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
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	result, _ := compute.Networks.List(s, params)
	result, _ = compute.ProjectMappings.List(s, params)

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// todo:修改IP子网
