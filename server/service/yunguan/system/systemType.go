package system

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemTypeService struct {
}

// 添加日志
func (m *SystemTypeService) AddSystemTypeService(model system.SystemType) (system.SystemType, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (a *SystemTypeService) UpdateSystemTypeService(model system.SystemType) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemType{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemType{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *SystemTypeService) GetSystemTypeListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemType{})
	var modelList []system.SystemType
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除日志
func (m *SystemTypeService) DeleteSystemTypeService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemType{}).Where("id", id).Delete(&system.SystemType{}).Error
	return err
}
func (sysTypeService *SystemTypeService) DeleteSysTypeByIds(ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SystemType{}, "id in ?", ids).Error
	return err
}
