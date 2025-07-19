package product

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/product"
)

type ProductElementRecordService struct {
}

// 添加
func (m *ProductElementRecordService) AddProductElementRecordService(model product.ProductElementRecord) (product.ProductElementRecord, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ProductElementRecordService) UpdateProductElementRecordService(model product.ProductElementRecord) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&product.ProductElementRecord{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductElementRecord{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// GetComputingInfoListService 查询
func (m *ProductElementRecordService) GetElementRecordListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&product.ProductElementRecord{})
	var modelList []product.ProductElementRecord
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// DeleteProductElementRecordService  删除日志
func (m *ProductElementRecordService) DeleteProductElementRecordService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductElementRecord{}).Where("id", id).Delete(&product.ProductElementRecord{}).Error
	return err
}
