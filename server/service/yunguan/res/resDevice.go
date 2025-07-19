package res

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/res"
)

type ResDeviceService struct {
}

// 添加
func (m *ResDeviceService) AddResDeviceService(model res.ResDevice) (res.ResDevice, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ResDeviceService) UpdateResDeviceService(model res.ResDevice) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&res.ResDevice{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResDevice{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *ResDeviceService) GetResDeviceListService(info res.GetResDeviceReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&res.ResDevice{})
	var modelList []res.ResDevice

	if info.Type != 0 {
		db.Where("type = ?", info.Type)
	}
	if info.Status != 0 {
		db.Where("status = ?", info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询服务器
func (m *ResDeviceService) GetAllResDeviceListService() (list interface{}, total int64, err error) {
	db := global.GVA_DB.Model(&res.ResDevice{})
	var modelList []res.ResDevice

	db.Where("type = 1")

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteResDeviceService  删除日志
func (m *ResDeviceService) DeleteResDeviceService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResDevice{}).Where("id", id).Delete(&res.ResDevice{}).Error
	return err
}
