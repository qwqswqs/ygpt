package system

import (
	"errors"
	"regexp"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/yunguan/cloudpods"
	"yunion.io/x/jsonutils"
)

type SystemTicketService struct {
}

// 添加日志
func (m *SystemTicketService) AddSystemTicketService(model system.SystemTicket) (system.SystemTicket, error) {
	err := global.GVA_DB.Create(&model).Error
	return model, err
}

// 修改模型
func (a *SystemTicketService) UpdateSystemTicketService(model system.SystemTicket) (err error) {
	err = global.GVA_DB.Where("id", model.ID).First(&system.SystemTicket{}).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("不存在该工具")
	}
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id", model.ID).Updates(&model).Error
	return err
}

// 查询日志
func (m *SystemTicketService) GetSystemTicketListService(info system.GetTicketReq) (list interface{}, total, totalNum, unAssignNum, unHandleNum, handleNum int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SystemTicket{})
	var modelList []system.SystemTicket
	global.GVA_DB.Model(&system.SystemTicket{}).Where("type", info.TicketType).Count(&totalNum)
	global.GVA_DB.Model(&system.SystemTicket{}).Where("status = ? and type = ?", 1, info.TicketType).Count(&unAssignNum)
	global.GVA_DB.Model(&system.SystemTicket{}).Where("status = ? and type = ?", 2, info.TicketType).Count(&unHandleNum)
	global.GVA_DB.Model(&system.SystemTicket{}).Where("status = ? and type = ?", 3, info.TicketType).Count(&handleNum)
	if info.TicketType != 0 {
		db = db.Where("ticket_type", info.TicketType)
	}
	if info.Status != 0 {
		db = db.Where("status", info.Status)
	}
	if info.TimeDesc != "" {
		db = db.Order("created_at " + info.TimeDesc)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&modelList).Error
	return modelList, total, totalNum, unAssignNum, unHandleNum, handleNum, err
}

// 删除日志
func (m *SystemTicketService) DeleteSystemTicketService(id uint) (err error) {
	// 执行更新操作
	err = global.GVA_DB.Model(&system.SystemTicket{}).Where("id", id).Delete(&system.SystemTicket{}).Error
	return err
}

type RenterResInfo struct {
	ResourceID string `json:"resourceID"` //资源编号
	Name       string `json:"name"`       //资源名称
	Status     string `json:"status"`     //资源状态
	PrivateIp  string `json:"privateIp"`  //内网ip
	Image      string `json:"image"`      //镜像名称
	Username   string `json:"username"`   //用户名
	Password   string `json:"password"`   //密码
}

// DetectUniqueName 检测字符串是否符合特定的命名方式
func DetectUniqueName(name string) bool {
	// 定义正则表达式
	pattern := `^[a-zA-Z0-9]{8}-[a-zA-Z0-9]*-[a-zA-Z]$`
	re := regexp.MustCompile(pattern)

	// 使用正则表达式进行匹配
	return re.MatchString(name)
}

func (m *SystemTicketService) GetNoAssignRes(resType int) (list interface{}, total int64, err error) {
	resInfos := make([]renter.RenterResInfo, 0)
	vmService := &cloudpods.VirtualMachineService{}
	containerService := &cloudpods.ContainerService{}
	bareMetalService := &cloudpods.BareMetalService{}
	instances := make([]jsonutils.JSONObject, 0)

	switch resType {
	case 1:
		instances, _, err = bareMetalService.HostList("bareMetal")
		if err != nil {
			return nil, 0, err
		}
		for _, bareMetal := range instances {
			bId, err := bareMetal.GetString("id")
			if err != nil {
				continue
			}
			var resInfo renter.RenterResInfo
			resInfo.Username, resInfo.Password, _ = bareMetalService.HostGetLoginInfo(bId)
			resInfo.Name, err = bareMetal.GetString("name")
			if err != nil {
				continue
			}
			isNoAssign := DetectUniqueName(resInfo.Name)
			if isNoAssign { //满足命名，表明该资源已分配给租户
				continue
			}
			resInfo.ResourceID = bId
			resInfo.PrivateIp, _ = bareMetal.GetString("ips")
			resInfo.Status, _ = bareMetal.GetString("status")
			if disksInfo, err := bareMetal.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
				resInfo.Image, _ = disksInfo[0].GetString("image")
			}
			resInfos = append(resInfos, resInfo)
		}
	case 2:
		instances, _, err = vmService.HostList()
		if err != nil {
			return nil, 0, err
		}
		for _, vm := range instances {
			vId, err := vm.GetString("id")
			if err != nil {
				continue
			}
			var resInfo renter.RenterResInfo
			resInfo.Username, resInfo.Password, _ = vmService.HostGetLoginInfo(vId)
			resInfo.Name, err = vm.GetString("name")
			if err != nil {
				continue
			}
			isNoAssign := DetectUniqueName(resInfo.Name)
			if isNoAssign { //满足命名，表明该资源已分配给租户
				continue
			}
			resInfo.ResourceID = vId
			resInfo.PrivateIp, _ = vm.GetString("ips")
			resInfo.Status, _ = vm.GetString("status")
			if disksInfo, err := vm.GetArray("disks_info"); len(disksInfo) != 0 && err == nil {
				resInfo.Image, _ = disksInfo[0].GetString("image")
			}
			resInfos = append(resInfos, resInfo)
		}
	case 3:
		instances, _, err = containerService.HostList()
		if err != nil {
			return nil, 0, err
		}
		for _, deployment := range instances {
			var resInfo renter.RenterResInfo
			resInfo.Name, err = deployment.GetString("name")
			if err != nil {
				continue
			}
			isNoAssign := DetectUniqueName(resInfo.Name)
			if isNoAssign { //满足命名，表明该资源已分配给租户
				continue
			}

			resInfo.ResourceID, err = deployment.GetString("id")
			if err != nil {
				continue
			}

			pod, _ := deployment.Get("pod")
			resInfo.PrivateIp, _ = pod.GetString("podIP")
			resInfo.Status, _ = pod.GetString("podPhase")
			if containers, err := pod.GetArray("containers"); len(containers) != 0 && err == nil {
				resInfo.Image, _ = containers[0].GetString("image")
				//if parts := strings.Split(resInfo.Image, "/"); len(parts) == 2 {
				//	resInfo.Image = parts[1]
				//}
			}
			resInfos = append(resInfos, resInfo)
		}
	default:
		return nil, 0, errors.New("参数错误")
	}

	total = int64(len(resInfos))
	return resInfos, total, nil
}

func (m *SystemTicketService) GetAllRes(resType int) (list interface{}, total int64, err error) {
	resInfos := make([]renter.RenterResInfo, 0)
	vmService := &cloudpods.VirtualMachineService{}
	containerService := &cloudpods.ContainerService{}
	bareMetalService := &cloudpods.BareMetalService{}
	instances := make([]jsonutils.JSONObject, 0)

	switch resType {
	case 1:
		instances, _, err = bareMetalService.HostList("bareMetal")
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
		instances, _, err = vmService.HostList()
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
		instances, _, err = containerService.HostList()
		if err != nil {
			return nil, 0, err
		}
		for _, deployment := range instances {
			var resInfo renter.RenterResInfo
			resInfo.ResourceID, err = deployment.GetString("id")
			if err != nil {
				continue
			}
			resInfo.Name, _ = deployment.GetString("name")
			pod, _ := deployment.Get("pod")
			resInfo.PrivateIp, _ = pod.GetString("podIP")
			resInfo.Status, _ = pod.GetString("podPhase")
			if containers, err := pod.GetArray("containers"); len(containers) != 0 && err == nil {
				resInfo.Image, _ = containers[0].GetString("image")
				//if parts := strings.Split(resInfo.Image, "/"); len(parts) == 2 {
				//	resInfo.Image = parts[1]
				//}
			}
			resInfos = append(resInfos, resInfo)
		}
	default:
		return nil, 0, errors.New("参数错误")
	}

	total = int64(len(resInfos))
	return resInfos, total, nil
}
