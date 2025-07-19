package product

import (
	"errors"
	"gorm.io/gorm"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
)

type ProductCategoryService struct {
}

// 添加
func (m *ProductCategoryService) AddProductCategoryService(model product.ProductCategory) (product.ProductCategory, error) {
	if !errors.Is(global.GVA_DB.Where("`key` = ? AND value = ?", model.Key, model.Value).First(&model).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return model, errors.New("该产品配置已存在")
	}
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// GetComputingInfoListService 查询
func (m *ProductCategoryService) GetCategoryListService() (list interface{}, total int64, err error) {
	db := global.GVA_DB.Model(&product.ProductCategory{})
	var modelList []product.ProductCategory
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&modelList).Error
	return modelList, total, err
}

// DeleteProductCategoryService  删除日志
func (m *ProductCategoryService) DeleteProductCategoryService(model product.ProductCategory) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductCategory{}).Where("ID", model.ID).Delete(&product.ProductCategory{}).Error
	return err
}
