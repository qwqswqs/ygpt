package res

import (
	"errors"
	"go.uber.org/zap"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/service/compute/udp"
)

type ResInfoService struct {
}

// 添加
func (m *ResInfoService) AddResInfoService(model res.ResInfo) (res.ResInfo, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *ResInfoService) UpdateResInfoService(model res.ResInfo) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&res.ResInfo{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResInfo{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *ResInfoService) GetResInfoListService(info res.GetResInfoReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&res.ResInfo{})
	var modelList []res.ResInfo
	if info.ImageID != 0 {
		db = db.Where("image_id = ?", info.ImageID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, err
}

// 查询租户资源列表
func (m *ResInfoService) GetRenterResService(info res.GetRenterResReq) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&res.ResInfo{})
	var modelList []res.ResInfo

	//var renterResList []res.GetResInfoRet

	if info.RenterID != 0 {
		db = db.Where("renter_id = ?", info.RenterID)
	}
	if info.Type != 0 {
		db = db.Where("resource_type = ?", info.Type)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error

	//var imageIDList []int
	//for _, v := range modelList {
	//	imageIDList = append(imageIDList, v.ImageID)
	//}
	//
	//var backupImageList []backup.BackupImage
	//err = db.Model(&backup.BackupImage{}).Where("id IN ?", imageIDList).Find(&backupImageList).Error
	//if err != nil {
	//	return
	//}

	//for _, v := range modelList {
	//	var renterRes res.GetResInfoRet
	//	renterRes.ID = v.ID
	//	renterRes.ResourceType = v.ResourceType
	//	renterRes.StartTime = v.StartTime
	//	renterRes.EndTime = v.EndTime
	//	renterRes.Status = v.Status
	//	renterRes.SshHost = v.SshHost
	//	for _, v1 := range backupImageList {
	//		if v.ImageID == int(v1.ID) {
	//			renterRes.Language = v1.Language
	//			renterRes.System = v1.System
	//			renterRes.Frame = v1.Frame
	//			renterRes.ImageName = v1.Name
	//			break
	//		}
	//	}
	//	renterResList = append(renterResList, renterRes)
	//}

	return modelList, total, err
}

// 查看操作
func (m *ResInfoService) GetResInfoDetailService(id uint) (info res.LibraryInfo, err error) {
	var model res.ResInfo
	err = global.GVA_DB.Where("id = ?", id).First(&model).Error
	if err != nil {
		global.GVA_LOG.Error("获取实例信息失败", zap.Error(err))
		return
	}

	udpService := udp.UdpServiceGroup{}
	info.DataBase, info.PyLibrary, err = udpService.GetInfo(model.SshHost)
	if err != nil {
		global.GVA_LOG.Error("获取实例信息失败", zap.Error(err))
		return
	}
	return
}

// DeleteResInfoService  删除日志
func (m *ResInfoService) DeleteResInfoService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&res.ResInfo{}).Where("id", id).Delete(&res.ResInfo{}).Error
	return err
}
