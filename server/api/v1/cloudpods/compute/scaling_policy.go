package compute

import (
	"github.com/gin-gonic/gin"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type ScalingPolicyApi struct{}

type scalingPolicyReq struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
}

func (p *ScalingPolicyApi) List(c *gin.Context) {
	var r scalingPolicyReq
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
	params.Set("scaling_group", jsonutils.NewString(r.Id))
	if r.Name != "" {
		params.Set("filter", jsonutils.NewString("name.contains('"+r.Name+"')"))
	}

	result, err := compute.ScalingPolicy.List(s, params)
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

type scalingPolicyCreateReq struct {
	Unit         string `json:"unit"`
	ScalingGroup string `json:"scaling_group"`
	TriggerType  string `json:"trigger_type"`
	Action       string `json:"action"`
	CoolingTime  int    `json:"cooling_time"`
	Name         string `json:"name"`
	Number       int    `json:"number"`
	CycleTimer   struct {
		CycleType string    `json:"cycle_type"`
		Hour      int       `json:"hour"`
		Minute    int       `json:"minute"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	} `json:"cycleTimer"`
}

func (p *ScalingPolicyApi) Create(c *gin.Context) {
	var r scalingPolicyCreateReq
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

	params := jsonutils.Marshal(r)

	_, err = compute.ScalingPolicy.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

type scalingPolicyDeleteReq struct {
	Id string `json:"id"`
}

func (p *ScalingPolicyApi) Delete(c *gin.Context) {
	var r scalingPolicyDeleteReq
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

	_, err = compute.ScalingPolicy.Delete(s, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

type scalingPolicyDisableReq struct {
	Id string `json:"id"`
}

func (p *ScalingPolicyApi) Disable(c *gin.Context) {
	var r scalingPolicyDisableReq
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

	_, err = compute.ScalingPolicy.PerformAction(s, r.Id, "disable", nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("禁用成功", c)
}

type scalingPolicyEnableReq struct {
	Id string `json:"id"`
}

func (p *ScalingPolicyApi) Enable(c *gin.Context) {
	var r scalingPolicyEnableReq
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

	_, err = compute.ScalingPolicy.PerformAction(s, r.Id, "enable", nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("启用成功", c)
}
