package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils"
)

type SystemConfigApi struct {
}

// 添加系统运行工具
func (d *SystemConfigApi) AddSystemConfig(c *gin.Context) {
	var r system.SystemConfig
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	serviceReturn, err := systemConfigService.AddSystemConfigService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改系统运行工具
func (d *SystemConfigApi) UpdateSystemConfig(c *gin.Context) {
	var model system.SystemConfig
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemConfigService.UpdateSystemConfigService(model)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 获取系统运行工具列表
func (b *SystemConfigApi) GetSystemConfigList(c *gin.Context) {
	var pageInfo system.GetSystemConfigReq
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
	list, total, err := systemConfigService.GetSystemConfigListService(pageInfo)
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

// 删除系统运行工具
func (b *SystemConfigApi) DeleteSystemConfig(c *gin.Context) {
	var model system.SystemConfig
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemConfigService.DeleteSystemConfigService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
