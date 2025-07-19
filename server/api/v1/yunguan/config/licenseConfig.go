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

type LicenseConfigApi struct{}

// 修改系统运行工具
func (d *LicenseConfigApi) UpdateLicenseConfig(c *gin.Context) {
	var model config.LicenseConfigData
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&config.LicenseConfig{}).Where("id", model.ID).Updates(&config.LicenseConfig{
		StartTime: model.StartTime,
		EndTime:   model.EndTime,
		DeviceNum: config.ToBigEndianBytes(uint32(model.DeviceNum)),
		RenterNum: config.ToBigEndianBytes(uint32(model.RenterNum)),
		AdminNum:  config.ToBigEndianBytes(uint32(model.AdminNum)),
	}).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
} // 修改系统运行工具
// 修改系统运行工具
func (d *LicenseConfigApi) UpdateLicensePwd(c *gin.Context) {
	var model config.LicenseConfig
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&config.LicenseConfig{}).Where("id", model.ID).Updates(&config.LicenseConfig{
		LicensePwd: utils.MD5V([]byte(model.LicensePwd)),
	}).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
} // 修改系统运行工具
func (d *LicenseConfigApi) LoginLicenseConfig(c *gin.Context) {
	var model struct {
		LicenseUser string `json:"licenseUser"`
		LicensePwd  string `json:"licensePwd"`
	}
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var licenseData config.LicenseConfig
	err = global.GVA_DB.Model(&config.LicenseConfig{}).Find(&licenseData).Error
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if (model.LicenseUser + "-admin") == (licenseData.LicenseUser + "-admin") {
		if utils.CheckMd5([]byte(model.LicensePwd), licenseData.LicensePwd) {
			response.OkWithMessage("登录成功", c)
			return
		} else {
			response.FailWithMessage("登录失败，密码错误", c)
			return
		}
	}
	response.FailWithMessage("登录失败，用户名错误", c)
	return
}

// 获取配置信息
func (b *LicenseConfigApi) GetLicenseConfig(c *gin.Context) {
	var model config.LicenseConfig
	err := global.GVA_DB.Model(&config.LicenseConfig{}).Find(&model).Error
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	modelData := config.LicenseConfigData{
		ID:          model.ID,
		DeviceNum:   int64(config.FromBigEndianBytes(model.DeviceNum)),
		RenterNum:   int64(config.FromBigEndianBytes(model.RenterNum)),
		AdminNum:    int64(config.FromBigEndianBytes(model.AdminNum)),
		LicensePwd:  model.LicensePwd,
		LicenseUser: model.LicenseUser,
		StartTime:   model.StartTime,
		EndTime:     model.EndTime,
	}
	response.OkWithDetailed(modelData, "获取成功", c)
}
