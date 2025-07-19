package renter

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/yunguan/cloudpods"
	"yunion.io/x/jsonutils"
)

type RenterResService struct {
}

// 添加
func (m *RenterResService) AddRenterResService(model renter.RenterRes) (renter.RenterRes, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (m *RenterResService) UpdateRenterResService(model renter.RenterRes) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&renter.RenterRes{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("id", model.ID).Updates(&model).Error
	return err
}

func (m *RenterResService) QueryRenterResInfo(info renter.GetRenterResReqByResourceID) (renter.RenterRes, error) {
	var model renter.RenterRes
	err := global.GVA_DB.Where("resource_id = ?", info.ResourceID).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录未找到，返回空的 renter.RenterRes 和 nil 错误
			return renter.RenterRes{}, nil
		}
		// 其他错误，记录日志并返回错误
		global.GVA_LOG.Error(err.Error())
		return renter.RenterRes{}, err
	}

	return model, nil
}

func (m *RenterResService) QueryRenterResCountByTicketId(ticketId int) (int64, error) {
	var count int64
	err := global.GVA_DB.Model(&renter.RenterRes{}).Where("ticket_id = ?", ticketId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetZoneListService 查询
//
//	func (m *RenterResService) GetRenterResListService(info renter.GetRenterResReq) (list interface{}, total int64, err error) {
//		var resInfos []renter.RenterResInfo
//
//		limit := info.PageSize
//		offset := info.PageSize * (info.Page - 1)
//		db := global.GVA_DB.Model(&renter.RenterRes{})
//		var modelList []renter.RenterRes
//
//		if info.RenterID != 0 {
//			db = db.Where("renter_id = ?", info.RenterID)
//		}
//		if info.Type != 0 {
//			db = db.Where("type = ?", info.Type)
//		}
//
//		err = db.Count(&total).Error
//		if err != nil {
//			return
//		}
//		err = db.Limit(limit).Offset(offset).Find(&modelList).Error
//		if err != nil {
//			return
//		}
//
//		var resourceIds []string
//		for _, item := range modelList {
//			if item.ResourceID != "" {
//				resourceIds = append(resourceIds, item.ResourceID)
//			}
//		}
//
//		vmService := &cloudpods.VirtualMachineService{}
//		containerService := &cloudpods.ContainerService{}
//		bareMetalService := &cloudpods.BareMetalService{}
//		instances := make([]jsonutils.JSONObject, 0)
//
//		if len(resourceIds) == 0 {
//			return resInfos, 0, nil
//		}
//
//		switch info.Type {
//		case 1:
//			instances, _, err = bareMetalService.HostListByIds("bareMetal", resourceIds)
//		case 2:
//			instances, _, err = vmService.HostListByIds(resourceIds)
//		case 3:
//			instances, _, err = containerService.HostListByIds(resourceIds)
//		case 4:
//			return nil, 0, errors.New("参数错误")
//		default:
//			return nil, 0, errors.New("参数错误")
//		}
//		if err != nil {
//			return nil, 0, err
//		}
//
//		for _, item := range modelList {
//			var resInfo renter.RenterResInfo
//			resInfo.ResourceID = item.ResourceID
//
//			var description map[string]interface{}
//			if err := json.Unmarshal([]byte(item.Description), &description); err == nil {
//				if allocTimeStr, ok := description["allocateTime"].(string); ok {
//					resInfo.CreateTime, _ = time.Parse(time.RFC3339, allocTimeStr)
//				}
//
//				if billingType, ok := description["billingType"].(int); ok {
//					if duration, ok := description["duration"].(int); ok {
//						switch billingType {
//						case 1: // 计时
//							resInfo.EndTime = resInfo.CreateTime.Add(time.Duration(duration) * time.Hour)
//						case 2: // 包日
//							resInfo.EndTime = resInfo.CreateTime.Add(time.Duration(duration) * 24 * time.Hour)
//						case 8: // 包月
//							resInfo.EndTime = resInfo.CreateTime.AddDate(0, duration, 0)
//						case 16: // 包年
//							resInfo.EndTime = resInfo.CreateTime.AddDate(duration, 0, 0)
//						default:
//							// 处理未知的计费类型
//						}
//					}
//				}
//			}
//
//			if !resInfo.CreateTime.IsZero() {
//				resInfo.RunTime = calculateRunTime(resInfo.CreateTime)
//			}
//
//			findFlag := false
//			switch item.Type {
//			case 1:
//				for _, bareMetal := range instances {
//					if bareMetalId, _ := bareMetal.GetString("id"); bareMetalId == item.ResourceID {
//						resInfo.Username, resInfo.Password, _ = bareMetalService.HostGetLoginInfo(item.ResourceID)
//						resInfo.Name, _ = bareMetal.GetString("name")
//						resInfo.PrivateIp, _ = bareMetal.GetString("ips")
//						resInfo.Status, _ = bareMetal.GetString("status")
//						if disksInfo, err := bareMetal.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
//							resInfo.Image, _ = disksInfo[0].GetString("image")
//						}
//						findFlag = true
//					}
//				}
//			case 2:
//				for _, vm := range instances {
//					if vmId, _ := vm.GetString("id"); vmId == item.ResourceID {
//						resInfo.Username, resInfo.Password, _ = vmService.HostGetLoginInfo(item.ResourceID)
//						resInfo.Name, _ = vm.GetString("name")
//						resInfo.PrivateIp, _ = vm.GetString("ips")
//						resInfo.Status, _ = vm.GetString("status")
//						if disksInfo, err := vm.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
//							resInfo.Image, _ = disksInfo[0].GetString("image")
//						}
//						findFlag = true
//					}
//				}
//			case 3:
//				for _, deployment := range instances {
//					if podId, _ := deployment.GetString("id"); podId == item.ResourceID {
//						pod, _ := deployment.Get("pod")
//						resInfo.PrivateIp, _ = pod.GetString("podIP")
//						resInfo.Status, _ = pod.GetString("podPhase")
//						if containers, err := pod.GetArray("containers"); len(containers) != 0 && err == nil {
//							resInfo.Name, _ = containers[0].GetString("name")
//							resInfo.Image, _ = containers[0].GetString("image")
//							if parts := strings.Split(resInfo.Image, "/"); len(parts) == 2 {
//								resInfo.Image = parts[1]
//							}
//						}
//						findFlag = true
//					}
//				}
//			}
//
//			// 如果找到匹配的实例，则添加到结果列表中
//			if findFlag {
//				resInfos = append(resInfos, resInfo)
//			}
//		}
//
//		total = int64(len(resInfos))
//		return resInfos, total, nil
//	}
func (m *RenterResService) GetRenterResListService(info renter.GetRenterResReq) (list interface{}, total int64, err error) {
	resInfos := make([]renter.RenterResInfo, 0)
	vmService := &cloudpods.VirtualMachineService{}
	containerService := &cloudpods.ContainerService{}
	bareMetalService := &cloudpods.BareMetalService{}
	instances := make([]jsonutils.JSONObject, 0)

	switch info.Type {
	case 1:
		instances, total, err = bareMetalService.HostListByRenterID("bareMetal", &info.RenterID, &info.Page, &info.PageSize)
		if err != nil {
			return nil, 0, err
		}
		for _, bareMetal := range instances {
			bId, err := bareMetal.GetString("id")
			if err != nil {
				continue
			}
			var resInfo renter.RenterResInfo
			resInfo.ResourceID = bId
			resInfo.Username, resInfo.Password, _ = bareMetalService.HostGetLoginInfo(bId)
			resInfo.Name, _ = bareMetal.GetString("name")
			resInfo.PrivateIp, _ = bareMetal.GetString("ips")
			resInfo.Status, _ = bareMetal.GetString("status")
			if disksInfo, err := bareMetal.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
				resInfo.Image, _ = disksInfo[0].GetString("image")
			}
			resInfos = append(resInfos, resInfo)
		}
	case 2:
		instances, total, err = vmService.HostListByRenterID(&info.RenterID, &info.Page, &info.PageSize)
		if err != nil {
			return nil, 0, err
		}
		for _, vm := range instances {
			vId, err := vm.GetString("id")
			if err != nil {
				continue
			}
			var resInfo renter.RenterResInfo
			resInfo.ResourceID = vId
			resInfo.Username, resInfo.Password, _ = vmService.HostGetLoginInfo(vId)
			resInfo.Name, _ = vm.GetString("name")
			resInfo.PrivateIp, _ = vm.GetString("ips")
			resInfo.Status, _ = vm.GetString("status")
			if disksInfo, err := vm.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
				resInfo.Image, _ = disksInfo[0].GetString("image")
			}
			resInfos = append(resInfos, resInfo)
		}
	case 3:
		instances, total, err = containerService.HostListByRenterID(&info.RenterID, &info.Page, &info.PageSize)
		if err != nil {
			return nil, 0, err
		}
		for _, deployment := range instances {
			var resInfo renter.RenterResInfo
			pod, _ := deployment.Get("pod")
			resInfo.PrivateIp, _ = pod.GetString("podIP")
			resInfo.Status, _ = pod.GetString("podPhase")
			resInfo.ResourceID, _ = deployment.GetString("id")
			if containers, err := pod.GetArray("containers"); len(containers) != 0 && err == nil {
				resInfo.Name, _ = containers[0].GetString("name")
				resInfo.Image, _ = containers[0].GetString("image")
				//if parts := strings.Split(resInfo.Image, "/"); len(parts) == 2 {
				//	resInfo.Image = parts[1]
				//}
			}
			fmt.Println("test:" + resInfo.ResourceID + "," + resInfo.Name)
			resInfos = append(resInfos, resInfo)
		}
	default:
		return nil, 0, errors.New("参数错误")
	}

	return resInfos, total, nil
}

// 获取所有租户资源列表
//func (m *RenterResService) GetAllRenterResListService() (list interface{}, err error) {
//	db := global.GVA_DB.Model(&renter.RenterRes{})
//	var modelList []renter.RenterRes
//
//	err = db.Find(&modelList).Error
//	if err != nil {
//		return
//	}
//
//	return modelList, nil
//}

func calculateRunTime(createTime time.Time) string {
	duration := time.Since(createTime)
	years := duration / (365 * 24 * time.Hour)
	remainingDays := (duration % (365 * 24 * time.Hour)) / (24 * time.Hour)
	remainingHours := (duration % (24 * time.Hour)) / time.Hour

	if years > 0 {
		return fmt.Sprintf("%d年%d天", years, remainingDays)
	} else if remainingDays > 0 {
		return fmt.Sprintf("%d天%d小时", remainingDays, remainingHours)
	}
	return fmt.Sprintf("%d小时", remainingHours)
}

// DeleteRenterService  删除日志
func (m *RenterResService) DeleteRenterResService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("id", id).Delete(&renter.RenterRes{}).Error
	return err
}

// 绑定资源
func (m *RenterResService) BindResInfoService(model renter.RenterRes) (err error) {
	var resInfo res.ResInfo
	err = global.GVA_DB.Model(&res.ResInfo{}).Where("id = ?", model.ResourceID).First(&resInfo).Error
	if err != nil {
		return errors.New("不存在该资源")
	}

	if resInfo.RenterID != -1 {
		return errors.New("该资源已被分配")
	}

	switch resInfo.Status {
	case 1:
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("id = ?", model.ID).First(&model).Error
		resInfo.Status = 2
		resInfo.RenterID = model.RenterID
		err = global.GVA_DB.Model(&res.ResInfo{}).Where("id = ?", resInfo.ID).Updates(resInfo).Error
		if err != nil {
			return err
		}

		err = global.GVA_DB.Where("id = ?", model.ID).First(&model).Error
		if err != nil {
			return errors.New("不存在该资源分配订单")
		}
		//算力平台跳转节点租户端的接口路由，由云管前端的地址 + "/#/jump"路由组成
		model.Url = global.GVA_CONFIG.System.ServeIP + "/#/jump"
		model.SshHost = resInfo.SshHost
		model.SshPort = resInfo.SshPort
		model.SshUser = resInfo.SshUser
		model.SshPwd = resInfo.SshPwd
		model.ResourceID = resInfo.ResourceID
		model.Status = 2
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("id = ?", model.ID).Updates(model).Error
		if err != nil {
			return err
		}
		err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id = ?", model.TicketId).Update("status", 2).Error
		return err
	case 2:
		return errors.New("该资源已使用")
	case 3:
		return errors.New("该资源已损坏")
	case 4:
		return errors.New("该资源已注销")
	default:
		return errors.New("该资源未知状态，无法使用")
	}
}

func (m *RenterResService) ReleaseResInfoService(model renter.RenterRes) (err error) {
	err = global.GVA_DB.Where("id = ?", model.ID).First(&model).Error
	if err != nil {
		return err
	}

	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("id = ?", model.ID).Update("status", 4).Error
	if err != nil {
		return err
	}

	var resInfo res.ResInfo
	err = global.GVA_DB.Model(&res.ResInfo{}).Where("id = ?", model.ResourceID).First(&resInfo).Error

	resInfo.RenterID = 0
	resInfo.Status = 1

	err = global.GVA_DB.Model(&res.ResInfo{}).Where("id = ?", model.ResourceID).Updates(resInfo).Error

	return err
}

// 获取订单相关的租户资源列表
func (m *RenterResService) GetResInfoByTicketService(model renter.RenterRes) (list interface{}, err error) {
	var data []renter.RenterRes
	err = global.GVA_DB.Model(&renter.RenterRes{}).Where("ticket_id = ?", model.TicketId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, err
}
