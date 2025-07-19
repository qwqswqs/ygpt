package compute

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type ScalingGroupApi struct{}

func (p *ScalingGroupApi) List(c *gin.Context) {
	var r scalingConfigListReq
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
	params.Set("details", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	if r.Name != "" {
		params.Set("filter", jsonutils.NewString("name.contains('"+r.Name+"')"))
	}

	result, err := compute.ScalingGroup.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	var page int
	if result.Limit != 0 {
		page = (result.Offset / result.Limit) + 1
	}

	response.OkWithDetailed(response.PageResult{
		List:     jsonutils.Marshal(result.Data).Interface(),
		Total:    int64(result.Total),
		Page:     page,
		PageSize: result.Limit,
	}, "获取成功", c)
}

type scalingGroupCreateReq struct {
	Name            string   `json:"name"`
	ScalingConfigId string   `json:"scalingConfigId"`
	Vpc             string   `json:"vpc"`
	MaxNumber       int      `json:"maxNumber"`
	BeginNumber     int      `json:"beginNumber"`
	MinNumber       int      `json:"minNumber"`
	ShrinkPrinciple string   `json:"shrinkPrinciple"`
	Networks        []string `json:"networks"`
}

type createReq struct {
	Cloudregion          string   `json:"cloudregion"`
	Project              string   `json:"project"`
	GenerateName         string   `json:"generate_name"`
	Brand                string   `json:"brand"`
	GuestTemplateId      string   `json:"guest_template_id"`
	Vpc                  string   `json:"vpc"`
	MaxInstanceNumber    int      `json:"max_instance_number"`
	DesireInstanceNumber int      `json:"desire_instance_number"`
	MinInstanceNumber    int      `json:"min_instance_number"`
	ShrinkPrinciple      string   `json:"shrink_principle"`
	HealthCheckMode      string   `json:"health_check_mode"`
	HealthCheckCycle     int      `json:"health_check_cycle"`
	HealthCheckGov       int      `json:"health_check_gov"`
	Networks             []string `json:"networks"`
	Hypervisor           string   `json:"hypervisor"`
}

func (p *ScalingGroupApi) Create(c *gin.Context) {
	var r scalingGroupCreateReq
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

	var req createReq
	req.Networks = r.Networks
	req.Vpc = r.Vpc
	req.MinInstanceNumber = r.MinNumber
	req.MaxInstanceNumber = r.MaxNumber
	req.DesireInstanceNumber = r.BeginNumber
	req.GuestTemplateId = r.ScalingConfigId
	req.ShrinkPrinciple = r.ShrinkPrinciple
	req.GenerateName = r.Name
	req.Project = s.GetProjectId()
	req.Hypervisor = "kvm"
	req.HealthCheckMode = "normal"
	req.HealthCheckGov = 180
	req.HealthCheckCycle = 300
	req.Brand = "OneCloud"
	req.Cloudregion = "default"

	params := jsonutils.Marshal(req)

	_, err = compute.ScalingGroup.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

type scalingGroupDeleteReq struct {
	Id string `json:"id"`
}
type scalingGroupDeleteByIdsReq struct {
	Ids []string `json:"ids"`
}

func (p *ScalingGroupApi) Delete(c *gin.Context) {
	var r scalingGroupDeleteReq
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

	_, err = compute.ScalingGroup.Delete(s, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
func (p *ScalingGroupApi) DeleteByIds(c *gin.Context) {
	var r scalingGroupDeleteByIdsReq
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
		_, err = compute.ScalingGroup.Delete(s, id, nil)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

type scalingGroupDisableReq struct {
	Id string `json:"id"`
}

func (p *ScalingGroupApi) Disable(c *gin.Context) {
	var r scalingGroupDisableReq
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

	_, err = compute.ScalingGroup.PerformAction(s, r.Id, "disable", nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("禁用成功", c)
}

type scalingGroupEnableReq struct {
	Id string `json:"id"`
}

func (p *ScalingGroupApi) Enable(c *gin.Context) {
	var r scalingGroupEnableReq
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

	_, err = compute.ScalingGroup.PerformAction(s, r.Id, "enable", nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("启用成功", c)
}

type detachReq struct {
	VmIds          []string `json:"vmIds"`          //要移除的云主机id
	ScalingGroupId string   `json:"scalingGroupId"` //弹性伸缩组id
}

func (p *ScalingGroupApi) Detach(c *gin.Context) {
	var r detachReq
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
	params.Set("delete_server", jsonutils.NewBool(true))
	params.Set("scaling_group", jsonutils.NewString(r.ScalingGroupId))
	_ = compute.Servers.BatchPerformAction(s, r.VmIds, "detach-scaling-group", params)

	response.OkWithMessage("移除成功", c)
}
