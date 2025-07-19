package system

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemOperateService struct {
}

// 添加日志
func (m *SystemOperateService) AddSystemOperateService(model system.SystemOperate) (system.SystemOperate, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 查询日志
func (m *SystemOperateService) GetSystemOperateListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemOperate{})
	var modelList []system.SystemOperate
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除日志
func (m *SystemOperateService) DeleteSystemOperateService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemOperate{}).Where("id", id).Delete(&system.SystemOperate{}).Error
	return err
}
