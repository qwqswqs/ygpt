package compute

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type ScalingConfigApi struct{}

type scalingConfigListReq struct {
	Name      string `json:"name"`
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
}

// 伸缩配置列表
func (p *ScalingConfigApi) List(c *gin.Context) {
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

	result, err := compute.GuestTemplate.List(s, params)
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

type Content struct {
	Auto_start       bool       `json:"auto_start"`
	Deploy_telegraf  bool       `json:"deploy_telegraf"`
	Prefer_region    string     `json:"prefer_region"`
	Prefer_zone      string     `json:"prefer_zone"`
	Duration         string     `json:"duration"`
	CDrom            string     `json:"cdrom"`
	Project_id       string     `json:"project_id"`
	Count            int        `json:"__count__"`
	Hypervisor       string     `json:"hypervisor"`
	Description      string     `json:"description"`
	Disks            []diskData `json:"disks"`
	Generate_name    string     `json:"generate_name"`
	Nets             []netData  `json:"nets"`
	Sku              string     `json:"sku"`
	Vcpu_count       int        `json:"vcpu_count"`
	Vmem_size        int        `json:"vmem_size"`
	OsArch           string     `json:"os_arch"`
	Isolated_devices []device   `json:"isolated_devices"`

	Bios    string `json:"bios"`
	Vdi     string `json:"vdi"`
	Vga     string `json:"vga"`
	Machine string `json:"machine"`
}

type ScalingConfigCreateReq struct {
	Name    string  `json:"name"`
	Content Content `json:"content"`
}

// 创建伸缩配置
func (p *ScalingConfigApi) Create(c *gin.Context) {
	var r ScalingConfigCreateReq
	err := c.ShouldBindJSON(&r.Content)
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
	r.Content.Auto_start = true
	r.Content.Deploy_telegraf = true
	r.Content.Count = 1
	r.Content.Hypervisor = "kvm"
	r.Content.OsArch = "x86"
	r.Content.Machine = "pc"
	r.Name = r.Content.Generate_name

	_, err = compute.ServerSkus.GetByName(s, r.Content.Sku, nil)
	if err != nil { // sku不存在，创建sku
		skuParams := jsonutils.NewDict()
		skuParams.Set("cpu_core_count", jsonutils.NewInt(int64(r.Content.Vcpu_count)))
		skuParams.Set("memory_size_mb", jsonutils.NewInt(int64(r.Content.Vmem_size)))
		skuInfo, err := compute.ServerSkus.Create(s, skuParams)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}

		r.Content.Sku, err = skuInfo.GetString("name")
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	params := jsonutils.Marshal(r)

	_, err = compute.GuestTemplate.Create(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

type scalingConfigDeleteReq struct {
	Id string `json:"id"`
}
type scalingConfigDeleteByIdsReq struct {
	Ids []string `json:"ids"`
}

// 删除伸缩配置
func (p *ScalingConfigApi) Delete(c *gin.Context) {
	var r scalingConfigDeleteReq
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

	_, err = compute.GuestTemplate.Delete(s, r.Id, nil)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// 删除伸缩配置
func (p *ScalingConfigApi) DeleteByIds(c *gin.Context) {
	var r scalingConfigDeleteByIdsReq
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
		_, err = compute.GuestTemplate.Delete(s, id, nil)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
