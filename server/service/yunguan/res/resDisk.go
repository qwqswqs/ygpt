package res

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/res"
)

type DiskService struct{}

// 添加
func (d *DiskService) AddDiskService(model res.ResDisk) (res.ResDisk, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (d *DiskService) UpdateDiskService(model res.ResDisk) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&res.ResDisk{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResDisk{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (d *DiskService) GetDiskListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&res.ResDisk{})
	var modelList []res.ResDisk

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteDiskService  删除日志
func (d *DiskService) DeleteDiskService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResDisk{}).Where("id", id).Delete(&res.ResDisk{}).Error
	return err
}
