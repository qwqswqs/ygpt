package backup

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/backup"
)

type ImageService struct {
}

// AddImageService 添加
func (d *ImageService) AddImageService(model backup.BackupImage) (backup.BackupImage, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// UpdateImageService 修改模型
func (d *ImageService) UpdateImageService(model backup.BackupImage) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&backup.BackupImage{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.BackupImage{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// GetImageListService 查询
func (d *ImageService) GetImageListService(info backup.GetImageListReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&backup.BackupImage{})
	var modelList []backup.BackupImage
	if info.SourceType != 0 {
		db = db.Where("source_type", info.SourceType)
	}
	if info.SourceType == 2 {
		db = db.Where("source_id", info.SourceID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// GetImageByIdService 查询
func (d *ImageService) GetImageByIdService(id int) (image backup.BackupImage, err error) {
	db := global.GVA_DB.Model(&backup.BackupImage{})
	var InfoList backup.BackupImage
	if err != nil {
		return
	}
	err = db.Where("id", id).Find(&InfoList).Error
	return InfoList, err
}

// GetImageListService 查询全部
func (d *ImageService) GetImageAllService() (list interface{}, total int64, err error) {
	db := global.GVA_DB.Model(&backup.BackupImage{})
	var modelList []backup.BackupImage

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteImageService  删除日志
func (d *ImageService) DeleteImageService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&backup.BackupImage{}).Where("id", id).Delete(&backup.BackupImage{}).Error
	return err
}
