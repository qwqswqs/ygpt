package backup

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/backup"
)

type SnapshotService struct {
}

// 添加
func (d *SnapshotService) AddSnapshotService(model backup.BackupSnapshot) (backup.BackupSnapshot, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (d *SnapshotService) UpdateSnapshotService(model backup.BackupSnapshot) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&backup.BackupSnapshot{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.BackupSnapshot{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (d *SnapshotService) GetSnapshotListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&backup.BackupSnapshot{})
	var modelList []backup.BackupSnapshot

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询日志
func (d *SnapshotService) GetUserSnapshotListService(info backup.GetBakSnapReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&backup.BackupSnapshot{})
	var modelList []backup.BackupSnapshot
	if info.RenterID != 0 {
		db = db.Where("renter_id", info.RenterID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteSnapshotService  删除日志
func (d *SnapshotService) DeleteSnapshotService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.BackupSnapshot{}).Where("id", id).Delete(&backup.BackupSnapshot{}).Error
	return err
}
