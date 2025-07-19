package product

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
)

type ProductSupplyService struct {
}

// 添加
func (m *ProductSupplyService) AddProductSupplyService(model product.ProductSupply) (product.ProductSupply, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ProductSupplyService) UpdateProductSupplyService(model product.ProductSupply) (err error) {
	err = global.GVA_DB.Where("id", model.Id).First(&product.ProductSupply{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	if model.Status == 0 {
		err = global.GVA_DB.Model(&product.ProductSupply{}).Select("status").Where("id", model.Id).Updates(&product.ProductSupply{Status: 0}).Error
		if err != nil {
			global.GVA_LOG.Debug(err.Error())
			return errors.New("不存在该工具")
		}
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductSupply{}).Where("id", model.Id).Updates(&model).Error
	return err
}

// GetSupplyInfoListService 查询
func (m *ProductSupplyService) GetSupplyInfoListService(info product.GetProductSupplyReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GVA_DB.Model(&product.ProductSupply{})
	var modelList []product.ProductSupply

	if info.Type != 0 {
		db = db.Where("type = ?", info.Type)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteProductSupplyService  删除日志
func (m *ProductSupplyService) DeleteProductSupplyService(id int64) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductSupply{}).Where("id", id).Delete(&product.ProductSupply{}).Error
	return err
}
