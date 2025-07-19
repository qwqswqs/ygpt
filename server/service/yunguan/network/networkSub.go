package network

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/network"
)

type SubService struct {
}

// 添加
func (m *SubService) AddSubService(model network.NetworkSub) (network.NetworkSub, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *SubService) UpdateSubService(model network.NetworkSub) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&network.NetworkSub{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkSub{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *SubService) GetSubListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&network.NetworkSub{})
	var modelList []network.NetworkSub

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询日志
func (m *SubService) GetAllSubListService() (list interface{}, total int64, err error) {
	db := global.GVA_DB.Model(&network.NetworkSub{})
	var modelList []network.NetworkSub

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteSubService  删除日志
func (m *SubService) DeleteSubService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkSub{}).Where("id", id).Delete(&network.NetworkSub{}).Error
	return err
}
