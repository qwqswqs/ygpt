package cloudpods

import (
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
)

type BareMetalApi struct {
}

// 创建裸金属
func (b *BareMetalApi) HostCreate(c *gin.Context) {
}

// 关机
func (b *BareMetalApi) HostShutDown(c *gin.Context) {
}

// 开机
func (b *BareMetalApi) HostStart(c *gin.Context) {
}

// 重启
func (b *BareMetalApi) HostRestart(c *gin.Context) {
}

// 挂起
func (b *BareMetalApi) HostSuspend(c *gin.Context) {
}

// 恢复
func (b *BareMetalApi) HostResume(c *gin.Context) {
}

// 裸金属远程连接
func (b *BareMetalApi) GetWebConsoleSSH(c *gin.Context) {
}

// 裸金属监控
func (b *BareMetalApi) HostMonitor(c *gin.Context) {
}

// 获取裸金属列表
func (b *BareMetalApi) HostGet(c *gin.Context) {
}

// 创建硬盘快照
func (b *BareMetalApi) HostCreateDiskSnap(c *gin.Context) {
}

type bareMetalListReq struct {
	Type string `json:"type"`
}

// 获取物理机列表或裸金属列表
func (b *BareMetalApi) HostList(c *gin.Context) {
	var r bareMetalListReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	hosts, total, err := BareMetalService.HostListByIds(r.Type, nil)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     jsonutils.Marshal(hosts).PrettyString(),
		Total:    int64(total),
		Page:     0,
		PageSize: 0,
	}, "获取成功", c)
}
