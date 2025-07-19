package storage

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

// 快照
type StoSnapshotApi struct {
}
type getStoSnapshotReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	ServerID  string `json:"serverID"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	OsArch    string `json:"osArch"`
	SizeDesc  string `json:"sizeDesc"`
}

// 获取主机快照列表
func (d *StoSnapshotApi) GetStoSnapshotList(c *gin.Context) {
	var r getStoSnapshotReq
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
	params.Set("server_id", jsonutils.NewString(r.ServerID))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	if r.Status != "" {
		filters = append(filters, "status.in('"+r.Status+"')")
	}
	if r.OsArch != "" {
		params.Set("os_arch", jsonutils.NewString(r.OsArch))
	}
	if r.SizeDesc != "" {
		params.Set("order_by", jsonutils.NewString("size"))
		params.Set("order", jsonutils.NewString(r.SizeDesc))
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.InstanceSnapshots.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "获取成功", c)
}

// 主机回滚
type snapshotResetReq struct {
	HostId     string `json:"hostId"`
	SnapshotId string `json:"snapshotId"`
}

func (d *StoSnapshotApi) SnapshotReset(c *gin.Context) {
	var r snapshotResetReq
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
	param.Set("instance_snapshot", jsonutils.NewString(r.SnapshotId))
	param.Set("auto_start", jsonutils.NewBool(true))
	result, err := compute.Servers.PerformAction(s, r.HostId, "instance-snapshot-reset", param)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

// 删除主机快照
type deleteSnapshotReq struct {
	SnapshotId string `json:"snapshotId"`
}

// 删除主机快照
type deleteSnapshotByIdsReq struct {
	Ids []string `json:"ids"`
}

func (d *StoSnapshotApi) DeleteSnapshot(c *gin.Context) {
	var r deleteSnapshotReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	param := jsonutils.NewDict()
	param.Set("id", jsonutils.NewString(r.SnapshotId))
	result, err := compute.InstanceSnapshots.Delete(s, r.SnapshotId, param)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}
func (d *StoSnapshotApi) DeleteSnapshotByIds(c *gin.Context) {
	var r deleteSnapshotByIdsReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	for _, id := range r.Ids {

		param := jsonutils.NewDict()
		param.Set("id", jsonutils.NewString(id))
		_, err := compute.InstanceSnapshots.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage("删除失败", c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}
