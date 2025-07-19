package network

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/network"
)

type IpService struct {
}

// 添加
func (m *IpService) AddIpService(model network.NetworkIp) (network.NetworkIp, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *IpService) UpdateIpService(model network.NetworkIp) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&network.NetworkIp{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkIp{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *IpService) GetIpListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&network.NetworkIp{})
	var modelList []network.NetworkIp

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询日志
func (m *IpService) GetAllIpListService() (list interface{}, total int64, err error) {
	db := global.GVA_DB.Model(&network.NetworkIp{})
	var modelList []network.NetworkIp
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteIpService  删除日志
func (m *IpService) DeleteIpService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkIp{}).Where("id", id).Delete(&network.NetworkIp{}).Error
	return err
}
