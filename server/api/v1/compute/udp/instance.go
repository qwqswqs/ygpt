package udp

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/res"
)

type InstanceApi struct {
}

func (i *InstanceApi) Monitor(c *gin.Context) {
	var r res.MonitorResReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	//monitorInfo, err := InstanceService.Monitor(r.PrivateIp)
	//if err != nil {
	//	global.GVA_LOG.Error("监控失败!", zap.Error(err))
	//	response.FailWithMessage("监控失败", c)
	//	return
	//}
	//response.OkWithDetailed(monitorInfo, "获取成功", c)
}

func (i *InstanceApi) Download(c *gin.Context) {
	var r res.DownloadReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = InstanceService.Download(r.PrivateIp, r.DownloadDetails)
	if err != nil {
		global.GVA_LOG.Error("下载失败!", zap.Error(err))
		response.FailWithMessage("下载失败", c)
		return
	}
	response.OkWithMessage("正在下载中，请稍等", c)
}

func (i *InstanceApi) Install(c *gin.Context) {
	var r res.InstallReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = InstanceService.Install(r.PrivateIp, r.SoftWare)
	if err != nil {
		global.GVA_LOG.Error("安装失败!", zap.Error(err))
		response.FailWithMessage("安装失败", c)
		return
	}
	response.OkWithMessage("正在安装中，请稍等", c)
}

func (i *InstanceApi) AddTask(c *gin.Context) {
	var r res.AddTaskReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = InstanceService.AddTask(r.PrivateIp, r.Mode, r.ModelFileUUID, r.DataSetFileUUID, r.TaskName, r.FrameName)
	if err != nil {
		global.GVA_LOG.Error("添加任务失败!", zap.Error(err))
		response.FailWithMessage("添加任务失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("添加任务成功", c)
}

func (i *InstanceApi) ListTask(c *gin.Context) {
	var r res.ListTaskReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	ret, total, err := InstanceService.ListTask(r.PrivateIp, r.Page, r.PageSize)
	if err != nil {
		global.GVA_LOG.Error("获取任务列表失败!", zap.Error(err))
		response.FailWithMessage("获取任务列表失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{List: ret, Total: int64(total)}, "获取任务列表成功", c)
}

func (i *InstanceApi) ManageTask(c *gin.Context) {
	var r res.ManageTaskReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = InstanceService.ManageTask(r.PrivateIp, r.TaskId, r.Method)
	if err != nil {
		global.GVA_LOG.Error("管理任务失败!", zap.Error(err))
		response.FailWithMessage("管理任务失败", c)
		return
	}
	response.OkWithMessage("管理任务成功", c)
}
