package config

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/config"
	"ygpt/server/utils"
)

type AlertConfigApi struct{}

// 修改系统运行工具
func (d *AlertConfigApi) AddAlertConfig(c *gin.Context) {
	var model config.AlertConfig
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if model.AlertConfigType == 2 {
		model.RenterID = int64(utils.GetUserID(c))
	}
	err = global.GVA_DB.Model(&config.AlertConfig{}).Create(&model).Error
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"createAt": time.Now().Unix()}, "创建成功", c)
} // 修改系统运行工具
// 修改系统运行工具
func (d *AlertConfigApi) UpdateAlertConfig(c *gin.Context) {
	var model config.AlertConfig
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&config.AlertConfig{}).Where("id", model.ID).Updates(&model).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
} // 修改系统运行工具
// 修改系统运行工具
func (d *AlertConfigApi) QueryAlertConfig(c *gin.Context) {
	var model struct {
		PageIndex       int `json:"pageIndex"`
		PageSize        int `json:"pageSize"`
		MonitorType     int `json:"monitorType"`
		AlertConfigType int `json:"alertConfigType"`
	}
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var dataList []config.AlertConfig
	var total int64
	db := global.GVA_DB.Model(&config.AlertConfig{})
	if model.MonitorType != 0 {
		db.Where("monitor_type=?", model.MonitorType)
	}
	db.Where("alert_config_type=?", model.AlertConfigType)
	if model.AlertConfigType == 2 {
		db.Where("renter_id=?", utils.GetUserID(c))
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("created_at DESC").Limit(model.PageSize).Offset(model.PageSize * (model.PageIndex - 1)).Find(&dataList).Error
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     dataList,
		Total:    total,
		Page:     model.PageIndex,
		PageSize: model.PageSize,
	}, "获取成功", c)
} // 修改系统运行工具
func (d *AlertConfigApi) DeleteAlertConfig(c *gin.Context) {
	var model struct {
		ID int `json:"id"`
	}
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&config.AlertConfig{}).Where("id", model.ID).Delete(&config.AlertConfig{}).Error
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
