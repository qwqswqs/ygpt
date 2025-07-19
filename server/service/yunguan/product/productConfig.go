package product

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
)

type ProductConfigService struct {
}

// 添加
func (m *ProductConfigService) AddProductConfigService(model product.ProductConfig) (product.ProductConfig, error) {
	if !errors.Is(global.GVA_DB.Where("name = ? && type = ?", model.Name, model.Type).First(&model).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return model, errors.New("该产品配置已存在")
	}
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ProductConfigService) UpdateProductConfigService(model product.ProductConfig) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&product.ProductConfig{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductConfig{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// GetComputingInfoListService 查询
func (m *ProductConfigService) GetConfigListService(info product.GetProductConfigListReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&product.ProductConfig{})
	var modelList []product.ProductConfig
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != 0 {
		db = db.Where("JSON_CONTAINS(type, ?)", strconv.Itoa(info.Type))
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// GetComputingInfoListService 查询
func (m *ProductConfigService) GetAllConfigListService() (list interface{}, total int64, err error) {

	db := global.GVA_DB.Model(&product.ProductConfig{})
	var modelList []product.GetAllProductConfigListRes
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteProductConfigService  删除日志
func (m *ProductConfigService) DeleteProductConfigService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductConfig{}).Where("id", id).Delete(&product.ProductConfig{}).Error
	return err
}
