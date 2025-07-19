package network

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/network"
)

type NatService struct {
}

// 添加
func (m *NatService) AddNatService(model network.NetworkNat) (network.NetworkNat, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *NatService) UpdateNatService(model network.NetworkNat) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&network.NetworkNat{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkNat{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *NatService) GetNatListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&network.NetworkNat{})
	var modelList []network.NetworkNat

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteNatService  删除日志
func (m *NatService) DeleteNatService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&network.NetworkNat{}).Where("id", id).Delete(&network.NetworkNat{}).Error
	return err
}
