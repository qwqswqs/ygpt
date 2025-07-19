package system

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/yunguan/config"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils/email"
	"ygpt/server/utils/sms"
)

type SystemAlertService struct {
}

// 添加系统配置
func (m *SystemAlertService) AddSystemAlertService(model system.SystemAlert) (system.SystemAlert, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 生成警告时做的准备
func (m *SystemAlertService) BeforeAddAlert(model system.SystemAlert) error {
	var data config.AlertConfig
	db := global.GVA_DB.Model(&config.AlertConfig{})
	if model.ResourceType == 4 {
		// 获取基础资源的报警策略
		err := global.GVA_DB.Model(&config.AlertConfig{}).Where("monitor_type = ? && alert_config_type = ?", model.AlertType, 3).First(&data).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 没有找到符合条件的记录
				return err
			} else {
				// 处理其他类型的错误
				fmt.Printf("查询失败: %v\n", err)
				return err
			}
		}
	} else {
		// 查找自定义策略
		err := db.Where("monitor_type = ? && renter_id = ?", model.AlertType, model.RenterID).First(&data).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 没有找到符合条件的记录
				//查找默认的配置
				err := global.GVA_DB.Model(&config.AlertConfig{}).Where("monitor_type = ? && alert_config_type = ?", model.AlertType, 1).First(&data).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// 没有找到符合条件的记录
						return err
					} else {
						// 处理其他类型的错误
						fmt.Printf("查询失败: %v\n", err)
						return err
					}
				}
			} else {
				// 处理其他类型的错误
				fmt.Printf("查询失败: %v\n", err)
				return err
			}
		}
	}
	var alertData system.SystemAlert
	alertData.AlertTime = time.Now()
	alertData.ResourceType = model.ResourceType
	alertData.ResourceID = model.ResourceID
	alertData.AlertType = model.AlertType
	alertData.AlertData = model.AlertData
	alertData.RenterID = model.RenterID
	alertData.ResourceName = model.ResourceName
	if model.AlertData >= data.SeriesAlertData {
		alertData.AlertLevel = 2
	} else if model.AlertData >= data.CommonAlertData {
		alertData.AlertLevel = 1
	} else {
		return errors.New("数据正常")
	}
	var renterData renter.Renter
	err := global.GVA_DB.Model(&renter.Renter{}).Where("id = ?", model.RenterID).Find(&renterData).Error
	if err != nil {
		return err
	}
	//for _, user := range data.AlertUser {
	//	if user == 1 {
	alertData.IsRenterData = true
	for _, i := range data.AlertWay {
		if i == 1 {
			//global.GVA_LOG.Info("租户报警", zap.Any("alertData", alertData))
			err = global.GVA_DB.Create(&alertData).Error
			if err != nil {
				return err
			}
		}
		if i == 2 {
			err = sms.SendSmsCommon(global.GVA_CONFIG.Sms.TemplateAlertId, data.NoticeText, renterData.Username)
			if err != nil {
				return err
			}
		}
		if i == 3 {
			err = email.SendEmail2(renterData.Email, data.NoticeText)
			if err != nil {
				return err
			}
		}
	}
	//}
	//if user == 2 {
	//	alertData.IsRenterData = false
	//	err = global.GVA_DB.Create(&alertData).Error
	//	if err != nil {
	//		return err
	//	}
	//}
	//if user == 3 {
	//	//TODO:算力圈租户查看信息
	//}
	//}

	return err
}

// 查询系统配置
func (m *SystemAlertService) GetSystemAlertListService(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemAlert{})
	var modelList []system.SystemAlert
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 删除系统配置
func (m *SystemAlertService) DeleteSystemAlertService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemAlert{}).Where("id", id).Delete(&system.SystemAlert{}).Error
	return err
}
