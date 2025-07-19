package renter

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/renter"
)

type RenterTaskService struct {
}

// 添加
func (m *RenterTaskService) AddRenterTaskService(model renter.RenterTask) (renter.RenterTask, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *RenterTaskService) UpdateRenterTaskService(model renter.RenterTask) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&renter.RenterTask{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.RenterTask{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// GetZoneListService 查询
func (m *RenterTaskService) GetRenterTaskListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&renter.RenterTask{})
	var modelList []renter.RenterTask

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteRenterService  删除日志
func (m *RenterTaskService) DeleteRenterTaskService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.RenterTask{}).Where("id", id).Delete(&renter.RenterTask{}).Error
	return err
}
