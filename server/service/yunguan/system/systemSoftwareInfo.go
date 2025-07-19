package system

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemSoftwareService struct {
}

// AddSystemSoftwareService 添加系统配置
func (m *SystemSoftwareService) AddSystemSoftwareService(model system.SystemSoftware) (system.SystemSoftware, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// UpdateSystemSoftwareService 修改系统配置
func (a *SystemSoftwareService) UpdateSystemSoftwareService(model system.SystemSoftware) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemSoftware{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemSoftware{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// GetSystemSoftwareListService 查询系统配置
func (m *SystemSoftwareService) GetSystemSoftwareListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemSoftware{})
	var modelList []system.SystemSoftware
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteSystemSoftwareService 删除系统配置
func (m *SystemSoftwareService) DeleteSystemSoftwareService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemSoftware{}).Where("id", id).Delete(&system.SystemSoftware{}).Error
	return err
}
