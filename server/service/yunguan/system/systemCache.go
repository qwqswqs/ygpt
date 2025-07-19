package system

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/system"
)

type SystemCacheService struct {
}

// 添加日志
func (m *SystemCacheService) AddSystemCacheService(model system.SystemCache) (system.SystemCache, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (a *SystemCacheService) UpdateSystemCacheService(model system.SystemCache) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemCache{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemCache{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *SystemCacheService) GetSystemCacheListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemCache{})
	var modelList []system.SystemCache
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除日志
func (m *SystemCacheService) DeleteSystemCacheService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemCache{}).Where("id", id).Delete(&system.SystemCache{}).Error
	return err
}
