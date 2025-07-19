package system

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemLogService struct {
}

// 添加日志
func (m *SystemLogService) AddSystemLogService(model system.SystemLog) (system.SystemLog, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 查询日志
func (m *SystemLogService) GetSystemLogListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemLog{})
	var modelList []system.SystemLog
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除日志
func (m *SystemLogService) DeleteSystemLogService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemLog{}).Where("id", id).Delete(&system.SystemLog{}).Error
	return err
}
