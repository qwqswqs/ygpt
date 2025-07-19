package product

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
)

type ProductComputingService struct {
}

// 添加
func (m *ProductComputingService) AddProductComputingService(model product.ProductComputing) (product.ProductComputing, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ProductComputingService) UpdateProductComputingService(model product.ProductComputing) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&product.ProductComputing{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	if model.AuditStatus == 0 {
		err = global.GVA_DB.Debug().Model(&product.ProductComputing{}).Select("audit_status").Where("id", model.ID).Updates(&product.ProductComputing{AuditStatus: 0}).Error
		if err != nil {
			global.GVA_LOG.Debug(err.Error())
			return errors.New("不存在该工具")
		}
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductComputing{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 同步所有
func (m *ProductComputingService) SyncedAllProductComputingService() (err error) {
	err = global.GVA_DB.Model(&product.ProductComputing{}).Where("is_synced", 1).Updates(&product.ProductComputing{
		IsSynced: 2,
	}).Error
	return err
}

// GetComputingInfoListService 查询
func (m *ProductComputingService) GetComputingInfoListService(info product.GetProductComputingReq) (list interface{}, isSyncedNum, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&product.ProductComputing{})
	var modelList []product.ProductComputing
	if info.Type != 0 {
		db = db.Where("type", info.Type)
	}
	if info.IsCustom != nil {
		db = db.Where("is_custom = ?", info.IsCustom)
	}
	if info.IsSynced != nil {
		db = db.Where("is_synced = ?", info.IsSynced)
	}
	if info.AuditStatus != nil {
		db = db.Where("audit_status = ?", info.AuditStatus)
	}
	if info.Name != "" {
		db = db.Where("name like ?", "%"+info.Name+"%")
	}
	if info.TimeDesc == "desc" {
		db = db.Order("created_at desc")
	} else {
		db = db.Order("created_at asc")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	err = global.GVA_DB.Model(&product.ProductComputing{}).Where("is_synced", 1).Count(&isSyncedNum).Error
	if err != nil {
		return
	}
	return modelList, isSyncedNum, total, err
}

// DeleteProductComputingService  删除日志
func (m *ProductComputingService) DeleteProductComputingService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&product.ProductComputing{}).Where("id", id).Delete(&product.ProductComputing{}).Error
	return err
}
