package baseRes

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type BaseDevicesApi struct {
}
type getBaremetalGpuInfoReq struct {
	HostId string `json:"hostId"`
}
type getBaremetalGpuInfoResq struct {
	Model string `json:"model"`
	Num   int    `json:"num"`
}

// 获取gpu列表
func (d *BaseDevicesApi) GetBaremetalGpuList(c *gin.Context) {
	var r getBaremetalGpuInfoReq
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
	params.Set("show_baremetal_isolated_devices", jsonutils.NewBool(true))
	params.Set("summary_stats", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("host_id", jsonutils.NewString(r.HostId))

	result, err := compute.IsolatedDevices.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	var res getBaremetalGpuInfoResq
	res.Num = result.Total
	if result.Total > 0 {
		res.Model, _ = result.Data[0].GetString("model")
	}

	response.OkWithDetailed(res, "获取成功", c)
}
