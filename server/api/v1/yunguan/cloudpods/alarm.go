package cloudpods

import (
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type AlarmApi struct {
}

// 查询云主机
func (p *AlarmApi) GetBaremetalList(c *gin.Context) {
	s, _ := global.GetCloudClientSession()
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	result, _ := modules.Baremetalagents.List(s, params)
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}
