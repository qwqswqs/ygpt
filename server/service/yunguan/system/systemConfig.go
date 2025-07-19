package system

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/system"
)

type SystemConfigService struct {
}

// 添加系统配置
func (m *SystemConfigService) AddSystemConfigService(model system.SystemConfig) (system.SystemConfig, error) {
	err := global.GVA_DB.Where("key", model.Key).First(&system.SystemConfig{}).Error
	if err == nil {
		return model, errors.New("添加失败,关键字重复")
	}
	err = global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改系统配置
func (a *SystemConfigService) UpdateSystemConfigService(model system.SystemConfig) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemConfig{}).Error
	if err != nil {
		return errors.New("修改失败,不存在该配置")
	}
	err = global.GVA_DB.Where("key", model.Key).First(&system.SystemConfig{}).Error
	if err == nil {
		return errors.New("修改失败,关键字重复")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemConfig{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询系统配置
func (m *SystemConfigService) GetSystemConfigListService(info system.GetSystemConfigReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemConfig{})
	if info.Type != 0 {
		db = db.Where("type = ?", info.Type)
	}
	var modelList []system.SystemConfig
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询系统配置
func (m *SystemConfigService) GetSystemConfigByKeyService(key string) (config system.SystemConfig, err error) {
	db := global.GVA_DB.Model(&system.SystemConfig{})
	db = db.Where("key", key)
	err = db.Find(&config).Error
	return config, err
}

// 删除系统配置
func (m *SystemConfigService) DeleteSystemConfigService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemConfig{}).Where("id", id).Delete(&system.SystemConfig{}).Error
	return err
}
