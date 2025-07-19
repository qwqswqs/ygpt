package system

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemToolService struct {
}

// 添加日志
func (m *SystemToolService) AddSystemToolService(model system.SystemTool) (system.SystemTool, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (a *SystemToolService) UpdateSystemToolService(model system.SystemTool) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemTool{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemTool{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *SystemToolService) GetSystemToolListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemTool{})
	var modelList []system.SystemTool
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除日志
func (m *SystemToolService) DeleteSystemToolService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemTool{}).Where("id", id).Delete(&system.SystemTool{}).Error
	return err
}
