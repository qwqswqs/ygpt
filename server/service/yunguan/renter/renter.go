package renter

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/renter"
)

type RenterService struct {
}

// 添加
func (m *RenterService) AddRenterService(model renter.Renter) (renter.Renter, error) {
	if !errors.Is(global.GVA_DB.Where("username = ?", model.Username).First(&model).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return model, errors.New("用户名已存在")
	}
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *RenterService) UpdateRenterService(model renter.Renter) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&renter.Renter{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.Renter{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 修改登录时间
func (m *RenterService) UpdateRenterLoginTimeService(id uint) (err error) {
	return global.GVA_DB.Model(&renter.Renter{}).
		Select("login_time").
		Where("renter_id=?", id).
		Updates(map[string]interface{}{
			"login_time": time.Now(),
		}).Error
}

// GetZoneListService 查询
func (m *RenterService) GetRenterListService(info renter.GetRenterListReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&renter.Renter{})
	var modelList []renter.Renter
	if info.Username != "" {
		db.Where("username LIKE  ?", "%"+info.Username+"%")
	}
	if info.CompanyName != "" {
		db.Where("company_name LIKE  ?", "%"+info.CompanyName+"%")
	}
	if info.Source != 0 {
		db.Where("source = ?", info.Source)
	}
	if info.Status != 0 {
		db.Where("status = ?", info.Status)
	}
	if info.TimeDesc == "desc" {
		db.Order("created_at desc")
	} else {
		db.Order("created_at asc")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

func (m *RenterService) GetRenterResListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	type resModel struct {
		Id           uint      `json:"id"`
		TenantId     int       `json:"tenantId"`
		Username     string    `json:"username"`
		ResourceType int       `json:"resourceType"`
		OrderId      int       `json:"orderId"`
		Status       int       `json:"status"`
		CreationDate time.Time `json:"created_at"`
	}

	var resList []renter.RenterRes
	var renterList []renter.Renter
	var modelList []resModel
	// First, count the total number of records
	err = global.GVA_DB.Model(&renter.RenterRes{}).Count(&total).Error
	if err != nil {
		return modelList, total, err
	}
	err = global.GVA_DB.Model(&renter.RenterRes{}).Order("created_at DESC").Limit(limit).Offset(offset).Scan(&resList).Error
	if err != nil {
		return modelList, total, err
	}

	var renterIdList []int

	for _, v := range resList {
		renterIdList = append(renterIdList, v.RenterID)
	}

	if len(renterIdList) > 0 {
		db := global.GVA_DB.Model(&renter.Renter{}).Order("created_at DESC").Where("renter_id IN ?", renterIdList).Find(&renterList)
		if db.Error != nil {
			global.GVA_LOG.Error("查询租户错误", zap.Error(db.Error))
			return modelList, total, errors.New("查询租户错误")
		}
	}

	for _, v := range resList {
		var res resModel
		res.Id = v.ID
		res.TenantId = v.RenterID
		res.ResourceType = v.Type
		res.OrderId = v.OrderId
		res.Status = v.Status
		for _, v1 := range renterList {
			if v.RenterID == v1.RenterID {
				res.Id = v.ID
				res.TenantId = v.RenterID
				res.Username = v1.Username
				res.ResourceType = v.Type
				res.OrderId = v.OrderId
				res.Status = v.Status
				res.CreationDate = v1.CreatedAt
				break
			}
		}
		modelList = append(modelList, res)
	}

	return modelList, total, err
}

// DeleteRenterService  删除日志
func (m *RenterService) DeleteRenterService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.Renter{}).Where("id", id).Delete(&renter.Renter{}).Error
	return err
}
