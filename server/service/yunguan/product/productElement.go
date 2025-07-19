package product

import (
	"errors"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/utils"
)

type ProductElemService struct {
}

func (i *ProductElemService) AddElemInfoService(info product.ProductElementInfo) (product.ProductElementInfo, error) {
	err := global.GVA_DB.Create(&info).Error
	return info, err
}

// 修改
func (i *ProductElemService) UpdateElemInfoService(info product.ProductElementInfo) (err error) {
	err = global.GVA_DB.Where("id", info.ID).First(&product.ProductElem{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该元素使用表")
	}

	if info.Status == 0 {
		err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Select("status").Where("id", info.ID).Updates(&product.ProductElem{Status: utils.Pointer(int8(0))}).Error
		if err != nil {
			global.GVA_LOG.Debug(err.Error())
			return errors.New("不存在该工具")
		}
	}
	err = global.GVA_DB.Debug().Model(&product.ProductElem{}).Where("id", info.ID).Updates(&info).Error
	return err
}

// 同步状态修改
func (i *ProductElemService) SyncElemInfoService(info product.ProductElementInfo) (err error) {
	err = global.GVA_DB.Where("type", info.Type).Updates(&info).Error
	return err
}

// QueryElemInfoService 查询
func (i *ProductElemService) QueryElemInfoService(info product.GetElemInfoReq) (list []product.GetElemInfoRes, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GVA_DB.Model(&product.ProductElem{})
	var InfoList []product.GetElemInfoRes
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != 0 {
		db = db.Where("type = ?", info.Type)
	}
	if info.Source != 0 {
		db = db.Where("source = ?", info.Source)
	}
	if info.OpenType != 0 {
		db = db.Where("open_type = ?", info.OpenType)
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	if info.AuditStatus != nil {
		db = db.Where("audit_status = ?", info.AuditStatus)
	}
	if info.UserID != 0 {
		db = db.Where("upload_by = ?", info.UserID)
	}

	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&InfoList).Error
	return InfoList, total, err
}

// QueryElemInfoService 查询
func (i *ProductElemService) QueryAllShareElemInfoService() (list []product.GetElemInfoRes, total int64, err error) {

	db := global.GVA_DB.Model(&product.ProductElementInfo{})
	var InfoList []product.GetElemInfoRes
	db = db.Where("open_type > 1")
	db = db.Where("status != 2")

	err = db.Find(&InfoList).Error
	err = db.Count(&total).Error
	return InfoList, total, err
}

// QueryElemInfoService 查询
func (i *ProductElemService) QueryShareElemInfoService(info product.GetElemInfoReq) (list []product.GetElemInfoRes, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GVA_DB.Model(&product.ProductElementInfo{})
	var InfoList []product.GetElemInfoRes
	db = db.Where("open_type > 1")
	db = db.Where("status != 2")
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != 0 {
		db = db.Where("type = ?", info.Type)
	}
	if info.Source != 0 {
		db = db.Where("source = ?", info.Source)
	}
	if info.UserID != 0 {
		db = db.Where("upload_by = ?", info.UserID)
	}

	err = db.Limit(limit).Offset(offset).Find(&InfoList).Error
	err = db.Count(&total).Error
	return InfoList, total, err
}

// QueryInfoById 查询
func (i *ProductElemService) QueryInfoById(id int) (elem product.ProductElementInfo, err error) {
	db := global.GVA_DB.Model(&product.ProductElementInfo{})
	var InfoList product.ProductElementInfo
	if err != nil {
		return
	}
	err = db.Where("id", id).Find(&InfoList).Error
	return InfoList, err
}

// DeleteElemInfoService 删除
func (i *ProductElemService) DeleteElemInfoService(id uint) (err error) {
	err = global.GVA_DB.Model(&product.ProductElementInfo{}).Debug().Where("id", id).Delete(&product.ProductElementInfo{}).Error
	return err
}
