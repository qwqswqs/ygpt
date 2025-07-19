package common

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	modules "yunion.io/x/onecloud/pkg/mcclient/modules/logger"
)

type CloudCommonLogApi struct{}
type getLogReq struct {
	request.PageInfo
	Id          string `json:"id"`
	IsContainer bool   `json:"isContainer"`
	ObjType     string `json:"objType"`
	OwnerKind   string `json:"ownerKind"`
	OwnerName   string `json:"ownerName"`
	NameSpace   string `json:"namespace"`
	ClusterId   string `json:"clusterId"`
	NextMarker  int64  `json:"nextMarker"`
	Name        string `json:"name"`
	Severity    string `json:"severity"`  //风险级别
	IsSuccess   *bool  `json:"isSuccess"` //结果
}

func (d *CloudCommonLogApi) GetLogsList(c *gin.Context) {
	var r getLogReq
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
	params.Set("obj_id", jsonutils.NewString(r.Id))
	params.Set("limit", jsonutils.NewInt(20))
	if r.ObjType != "" {
		params.Set("obj_type", jsonutils.NewString(r.ObjType))
	}

	if r.NextMarker > 0 {
		params.Set("paging_marker", jsonutils.NewInt(r.NextMarker))
	}

	if r.IsContainer {
		params.Set("owner_kind", jsonutils.NewString(r.OwnerKind))
		params.Set("owner_name", jsonutils.NewString(r.OwnerName))
		params.Set("namespace", jsonutils.NewString(r.NameSpace))
		params.Set("cluster", jsonutils.NewString(r.ClusterId))
	}
	if r.IsSuccess != nil {
		params.Set("success", jsonutils.NewBool(*r.IsSuccess))
	}
	if r.Name != "" {
		params.Set("obj_name", jsonutils.NewString(r.Name))
	}
	if r.Severity != "" {
		params.Set("severity", jsonutils.NewString(r.Severity))
	}

	logs, err := modules.Actions.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(logs).PrettyString(), "获取成功", c)
}
