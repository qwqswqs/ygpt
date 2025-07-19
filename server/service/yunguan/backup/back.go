package backup

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/backup"
)

type BackService struct {
}

// 添加
func (d *BackService) AddBackService(model backup.Backup) (backup.Backup, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (d *BackService) UpdateBackService(model backup.Backup) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&backup.Backup{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.Backup{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (d *BackService) GetBackListService(info backup.BackReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&backup.Backup{})
	var modelList []backup.Backup
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.RenterID != 0 {
		db = db.Where("renter_id = ?", info.RenterID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteBackService  删除日志
func (d *BackService) DeleteBackService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.Backup{}).Where("id", id).Delete(&backup.Backup{}).Error
	return err
}
