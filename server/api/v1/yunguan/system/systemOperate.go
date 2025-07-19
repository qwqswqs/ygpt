package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils"
)

type SystemOperateApi struct {
}

// 获取日志列表
func (b *SystemOperateApi) GetSystemOperateList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := systemOperateService.GetSystemOperateListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 删除模型
func (b *SystemOperateApi) DeleteSystemOperate(c *gin.Context) {
	var model system.SystemOperate
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemOperateService.DeleteSystemOperateService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
