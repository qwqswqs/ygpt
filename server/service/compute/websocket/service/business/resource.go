package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"ygpt/server/global"
	sysmodel "ygpt/server/model/system"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service"
	"ygpt/server/service/compute/websocket/protocol"
	"yunion.io/x/jsonutils"
)

type ResourceService struct{}

type specJson struct {
	Name       string   `json:"name"`
	Value      []string `json:"value"`
	ValueIndex int      `json:"valueIndex"`
}
type params struct {
	ProductId       int
	ImageName       string
	Cpu             int
	Memory          int
	GpuType         string
	GpuMemory       int
	GpuNumber       int
	Storage         int
	SysStorageType  string
	DataStorageType string
	SysDisk         int
	DataDisk        int
	BandWidth       int
}

func (r *ResourceService) HandleOrderTicketRequest(message *protocol.OrderTicketRequest) error {
	var existTicket system.SystemTicket
	if err := global.GVA_DB.Where("order_id = ?", message.Payload.OrderID).First(&existTicket).Error; err == nil {
		//如果工单已经存在，则直接返回
		retMsg := protocol.OrderTicketResponse{}
		retMsg.ID = message.ID
		retMsg.Command = protocol.OrderTicketRetCmd
		retMsg.RetCode = protocol.SuccessCodeCmd
		retMsg.RetMsg = protocol.SuccessMsgCmd
		return global.GVA_WS.SendMessageDebug(retMsg)
	}
	var ticket system.SystemTicket
	ticket.Type = 2
	//此处创建的工单只关注报价与上传合同发票
	ticket.TicketType = 4
	ticket.IsQuotation = message.Payload.IsQuote
	ticket.IsThird = message.Payload.IsThird
	var renterInfo renter.Renter
	if err := global.GVA_DB.Where("username = ?", message.Payload.BuyerPhone).First(&renterInfo).Error; err != nil {
		global.GVA_LOG.Error("Error finding renter info:", zap.Error(err))
		return err
	}
	ticket.RenterID = renterInfo.RenterID
	ticket.Description = message.Payload.OrderDetailJSON
	ticket.CreatePerson = 0
	ticket.Status = 2
	//如果该订单是标准产品，则不需要报价
	ticket.IsQuotation = message.Payload.IsQuote
	//如果本节点是自营节点，那么所有资源分配工单都不需要报价或处理
	if ticket.IsThird == 0 {
		ticket.Status = 2
	}
	ticket.OrderId = message.Payload.OrderID
	err := global.GVA_DB.Create(&ticket).Error
	if err != nil {
		global.GVA_LOG.Error("Error creating ticket:", zap.Error(err))
		return err
	}
	retMsg := protocol.OrderTicketResponse{}
	retMsg.ID = message.ID
	retMsg.Command = protocol.OrderTicketRetCmd
	retMsg.RetCode = protocol.SuccessCodeCmd
	retMsg.RetMsg = protocol.SuccessMsgCmd
	return global.GVA_WS.SendMessageDebug(retMsg)
}

func (r *ResourceService) AllocateResource(message *protocol.PackageAllocateRequest) error {
	//构建应答
	response := protocol.PackageAllocateResponse{}
	response.ID = message.ID
	response.Command = protocol.AllocateRetCmd
	response.NodeCode = global.GVA_CONFIG.System.NodeCode
	//查数据库
	var renterInfo renter.Renter
	var user sysmodel.SysUser
	//用户是否存在
	err := global.GVA_DB.Where("phone = ?", message.Phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//注册新用户和租户身份
			err = r.registerNewRenter(&user, &renterInfo, message)
			if err != nil {
				global.GVA_LOG.Error("Error creating user:", zap.Error(err))
				return err
			}
		}
	} else {
		//存在，则查询租户信息
		err = global.GVA_DB.Where("renter_id = ?", user.ID).First(&renterInfo).Error
		if err != nil {
			global.GVA_LOG.Error("Error finding renter info:", zap.Error(err))
			return err
		}
	}
	//查租户资源表资源是否已经分配
	var renterRes renter.RenterRes
	err = global.GVA_DB.Where("renter_id = ? and order_id = ? and order_computing_id = ?",
		renterInfo.RenterID,
		message.OrderID,
		message.OrderComputingID).First(&renterRes).Error
	if err == nil {
		//资源已经分配
		response.RetCode = 302
		response.RetMsg = "订单处理中或已经分配"
		return global.GVA_WS.SendMessageDebug(response)
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//如果是其他错误则返回
			return err
		} else {
			//解析
			var AllocateDetailsMap map[string]interface{}
			err = json.Unmarshal([]byte(message.AllocateDetailsJson), &AllocateDetailsMap)
			if err != nil {
				global.GVA_LOG.Error("Error decoding JSON:", zap.Error(err))
				return err
			}
			var newRenterRes renter.RenterRes
			//创建尚未分配资源的租户资源记录
			newRenterRes.RenterID = renterInfo.RenterID
			newRenterRes.Type = message.ResourceType
			newRenterRes.Description = message.OrderDetailJSON
			newRenterRes.OrderId = message.OrderID
			newRenterRes.OrderComputingId = message.OrderComputingID
			//newRenterRes.TicketId = int(ticket.ID)
			newRenterRes.Status = 1
			newRenterRes.DataSource = 2
			newRenterRes.AiList = "[]"
			newRenterRes.UseTime = message.UseTime
			if err = global.GVA_DB.Create(&newRenterRes).Error; err != nil {
				global.GVA_LOG.Error("Error creating ticket:", zap.Error(err))
				return err
			} else {
				//从报文中解析出资源分配的具体参数
				err, param := r.unmarshalParams(AllocateDetailsMap, message.ProductID)
				if err != nil {
					global.GVA_LOG.Error("Error unmarshalling JSON:", zap.Error(err))
					return err
				}
				if message.IsCustom != 1 {
					//标准产品可以直接分配资源，定制产品先跳过，需要运维人员在工单管理来分配
					//开启协程,分配资源
					go func() {
						r.cloudPodsResourceAllocate(param, message, &newRenterRes)
					}()
				} else {
					//创建定制产品资源分配的工单
					var ticket system.SystemTicket
					ticket.Type = 2
					ticket.TicketType = 2
					//因为这个工单只关注资源分配，在此工单中不需要报价或处理
					ticket.IsQuotation = 0
					ticket.IsThird = 0
					ticket.RenterID = renterInfo.RenterID
					ticket.Description = message.OrderDetailJSON
					ticket.OrderId = message.OrderID
					ticket.CreatePerson = 0
					ticket.Status = 1
					if err := global.GVA_DB.Create(&ticket).Error; err != nil {
						global.GVA_LOG.Error("Error creating ticket:", zap.Error(err))
						return err
					}
					// 将工单ID更新到资源表中
					err := global.GVA_DB.Model(&newRenterRes).Updates(&renter.RenterRes{
						TicketId: int(ticket.ID),
					}).Error
					if err != nil {
						global.GVA_LOG.Error("Error updating ticket:", zap.Error(err))
						return err
					}
				}
				response.RetCode = protocol.SuccessCodeCmd
				response.RetMsg = protocol.SuccessMsgCmd + ":" + "资源分配中，等待节点回应"
				return global.GVA_WS.SendMessageDebug(response)
			}
		}
	}
}

func (r *ResourceService) registerNewRenter(user *sysmodel.SysUser, renterInfo *renter.Renter, message *protocol.PackageAllocateRequest) error {
	var authorities []sysmodel.SysAuthority
	authorities = append(authorities, sysmodel.SysAuthority{
		AuthorityId: 777,
	})
	newUser := &sysmodel.SysUser{
		UUID:        uuid.UUID{},
		Username:    message.Phone,
		NickName:    message.YgAccount,
		Password:    message.YgPassword,
		HeaderImg:   "",
		AuthorityId: 777,
		Authorities: authorities,
		Enable:      1,
		Phone:       message.Phone,
		Email:       "",
		Source:      2,
	}
	if !errors.Is(global.GVA_DB.Where("username = ?", newUser.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		err := errors.New("用户名已注册")
		return err
	}
	// 否则 附加uuid 注册
	newUser.UUID = uuid.Must(uuid.NewV4())
	err := global.GVA_DB.Create(&newUser).Error
	if err != nil {
		global.GVA_LOG.Error("注册用户失败!", zap.Error(err))
		return err
	} else {
		//同时注册相应的租户
		renterInfo.RenterID = int(newUser.ID)
		renterInfo.Username = message.Phone
		renterInfo.Source = 2
		renterInfo.Status = 1
		err = global.GVA_DB.Create(&renterInfo).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceService) cloudPodsResourceAllocate(param params, message *protocol.PackageAllocateRequest, newRenterRes *renter.RenterRes) {
	//云主机资源分配在ygpt/serve/service/yunguan/cloudpods/virtual_machine，方法名为hostcreate
	var resultJson jsonutils.JSONObject
	var resourceId string
	var err error
	switch message.ResourceType {
	//裸金属
	case 1:
		global.GVA_LOG.Info("裸金属资源分配接口",
			zap.Int("product", param.ProductId),
			zap.String("imageName", param.ImageName),
			zap.Int("renterID", newRenterRes.RenterID),
		)
		if param.ProductId != 0 && param.ImageName != "" {
			resultJson, err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareMetalService.HostCreate(param.ProductId, param.ImageName, newRenterRes.RenterID)
		} else {
			global.GVA_LOG.Error("裸金属资源分配接口参数错误", zap.Error(errors.New("参数错误")))
		}
		resourceId, err = r.allocateResultCheck(message.ResourceType, err, resultJson, param.BandWidth, *newRenterRes)
	//云主机
	case 2:
		global.GVA_LOG.Info("云主机资源分配接口",
			zap.String("imageName", param.ImageName),
			zap.Int("cpu", param.Cpu),
			zap.Int("memory", param.Memory),
			zap.String("storageType", param.SysStorageType),
			zap.Int("sysDisk", param.SysDisk),
			zap.Int("dataDisk", param.DataDisk),
			zap.Int("GPU", param.GpuNumber))
		if param.Cpu != 0 && param.Memory != 0 && param.Storage != 0 {
			resultJson, err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService.HostCreate(param.ImageName, param.Cpu, param.Memory, param.SysStorageType, param.SysDisk, param.DataDisk, param.GpuNumber, newRenterRes.RenterID)
		} else {
			global.GVA_LOG.Error("云主机资源分配接口参数错误", zap.Error(errors.New("参数错误")))
		}
		resourceId, err = r.allocateResultCheck(message.ResourceType, err, resultJson, param.BandWidth, *newRenterRes)
	//容器
	case 3:
		global.GVA_LOG.Info("容器资源分配接口",
			zap.String("imageName", param.ImageName),
			zap.Int("cpu", param.Cpu),
			zap.Int("memory", param.Memory),
			zap.Int("GPU", param.GpuNumber))
		if param.Cpu != 0 && param.Memory != 0 {
			resultJson, err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ContainerService.HostCreate(param.ImageName, param.Cpu, param.Memory, param.SysDisk, param.DataDisk, param.GpuNumber, newRenterRes.RenterID)
		} else {
			global.GVA_LOG.Error("容器资源分配接口参数错误", zap.Error(errors.New("参数错误")))
		}
		resourceId, err = r.allocateResultCheck(message.ResourceType, err, resultJson, param.BandWidth, *newRenterRes)
	//存储，该类型暂未实现
	case 4:
	default:
		global.GVA_LOG.Error("分配资源错误", zap.Error(errors.New("未知资源类型")))
	}
	//分配失败，创建工单--------------------------------------------------------------------
	if err != nil {
		global.GVA_LOG.Error("资源分配失败", zap.Error(err))
		var ticket system.SystemTicket
		ticket.Type = 2
		ticket.TicketType = 1
		ticket.Description = message.OrderDetailJSON
		ticket.IsQuotation = 0
		ticket.IsThird = 0
		ticket.RenterID = newRenterRes.RenterID
		ticket.ProblemDescription = err.Error()
		ticket.Status = 1
		ticket.OrderId = message.OrderID
		ticket.CreatePerson = 0
		if err = global.GVA_DB.Where("id = ?", newRenterRes.TicketId).Create(&ticket).Error; err != nil {
			global.GVA_LOG.Error("创建工单失败", zap.Error(err))
		} else {
			//更新租户资源状态
			//更新资源ID到租户资源记录，ssh信息要等分配完成
			var updatesRenterRes renter.RenterRes
			updatesRenterRes.TicketId = int(ticket.ID)
			updatesRenterRes.Status = 1
			if err = global.GVA_DB.Model(&updatesRenterRes).Where("id = ?", newRenterRes.ID).Updates(&updatesRenterRes).Error; err != nil {
				global.GVA_LOG.Error("更新资源ID错误:", zap.Error(err))
			}
		}
	}

	//分配结束，发送异步通知-----------------------------------------------------------------
	if resourceId != "" {
		//然后调用AllocateResponse，发送应答
		err = r.AllocateResponse(err, message.OrderComputingID, resourceId, false)
		if err != nil {
			global.GVA_LOG.Error("资源分配结果：", zap.Error(err))
		}
	}
}

func (r *ResourceService) allocateResultCheck(resType int, err error, resultJson jsonutils.JSONObject, bandWidth int, newRenterRes renter.RenterRes) (string, error) {
	if err != nil {
		global.GVA_LOG.Error("资源分配接口返回错误", zap.Error(err))
		return "", err
	} else {
		if resultJson != nil {
			//更新资源ID到租户资源记录，ssh信息要等分配完成
			var updatesRenterRes renter.RenterRes
			idValue, err := resultJson.GetString("id")
			if err != nil {
				global.GVA_LOG.Error("获取资源ID错误", zap.Error(err))
				return "", fmt.Errorf("获取资源ID错误" + err.Error())
			} else {
				updatesRenterRes.ResourceID = idValue
				updatesRenterRes.Url = global.GVA_CONFIG.System.ServeIP + "/#/jump"
			}
			if err = global.GVA_DB.Model(&updatesRenterRes).Where("id = ?", newRenterRes.ID).Updates(&updatesRenterRes).Error; err != nil {
				global.GVA_LOG.Error("更新资源ID错误:", zap.Error(err))
				return "", fmt.Errorf("更新资源ID错误" + err.Error())
			} else {
				switch resType {
				case 1:
					err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareMetalService.HostForecast(resultJson)
				case 2:
					// 探测云主机状态，err为空代表云主机已经创建成功，renter_res记录已经更新
					err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService.HostForecast(resultJson)
				case 3:
					//容器会立即创建，不需要检查，视为创建成功了
					now := time.Now()
					res := &renter.RenterRes{}
					if global.GVA_DB.Where("id = ?", newRenterRes.ID).First(&res).Error != nil {
						global.GVA_LOG.Error("查询租户资源记录错误:", zap.Error(err))
						return "", fmt.Errorf("查询租户资源记录错误" + err.Error())
					}
					end := time.Now().Add(time.Duration(res.UseTime) * time.Hour)
					err = global.GVA_DB.Model(&newRenterRes).Where("id = ?", newRenterRes.ID).Updates(&renter.RenterRes{StartTime: now, EndTime: end}).Error
				default:
					global.GVA_LOG.Error("分配资源错误", zap.Error(errors.New("未知资源类型")))
				}
				if err == nil {
					switch resType {
					case 1:
						//裸金属不能带宽限制
						err = nil
					case 2:
						err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService.HostBandwidthLimit(idValue, bandWidth)
					case 3:
						//TODO : 容器暂未实现此功能
						err = nil
					default:
						global.GVA_LOG.Error("分配资源错误", zap.Error(errors.New("未知资源类型")))
						return "", fmt.Errorf("分配资源错误未知资源类型")
					}
					if err != nil {
						global.GVA_LOG.Error("资源分配带宽限制失败：", zap.Error(err))
					}
					return idValue, nil
				} else {
					global.GVA_LOG.Error("资源分配状态检测失败：", zap.Error(err))
					return "", err
				}
			}
		} else {
			global.GVA_LOG.Error("资源分配：", zap.Error(errors.New("资源分配失败,资源分配进程未完成")))
			return "", errors.New("资源分配失败,资源分配进程未完成")
		}
	}
}

func (r *ResourceService) unmarshalParams(allocateDetailsMap map[string]interface{}, productId int) (error, params) {
	var param params
	// 初始化参数
	param.ProductId = productId
	param.ImageName = ""
	param.Memory = 0
	param.Cpu = 0
	param.GpuType = ""
	param.GpuMemory = 0
	param.GpuNumber = 0
	param.Storage = 0
	param.SysStorageType = ""
	param.DataStorageType = ""
	param.SysDisk = 0
	param.DataDisk = 0
	param.BandWidth = 0
	value, ok := allocateDetailsMap["os"].(string)
	if ok {
		param.ImageName = value
	}
	//specJson是一个数组
	var specJsons []specJson
	err := json.Unmarshal([]byte(allocateDetailsMap["specJson"].(string)), &specJsons)
	if err != nil {
		global.GVA_LOG.Error("Invalid specJson format")
		return err, param
	} else {
		for _, spec := range specJsons {
			switch spec.Name {
			case "CPU":
				cupWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(cupWithUnit, "核") {
					//规格中cpu为数字+核
					param.Cpu, err = strconv.Atoi(strings.Split(cupWithUnit, "核")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid cpu format")
						continue
					}
				} else {
					global.GVA_LOG.Error("Invalid cpu format")
					continue
				}
			case "内存":
				memoryWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(memoryWithUnit, "GB") {
					//规格中memory为数字+GB
					param.Memory, err = strconv.Atoi(strings.Split(memoryWithUnit, "GB")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid memory format")
						continue
					}
				} else {
					global.GVA_LOG.Error("Invalid memory format")
					continue
				}
			case "系统盘":
				systemDiskWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(systemDiskWithUnit, "GB") {
					//规格中storage为数字+GB
					systemDisk, err := strconv.Atoi(strings.Split(systemDiskWithUnit, "GB")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid storage format")
						continue
					}
					param.Storage += systemDisk
					param.SysDisk = systemDisk
				} else {
					global.GVA_LOG.Error("Invalid storage format")
					continue
				}
			case "数据盘":
				dataDiskWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(dataDiskWithUnit, "GB") {
					// 规格中storage为数字+GB
					dataDisk, err := strconv.Atoi(strings.Split(dataDiskWithUnit, "GB")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid storage format")
						continue
					}
					// 将数据盘加到存储中
					param.Storage += dataDisk
					param.DataDisk = dataDisk
				} else if strings.Contains(dataDiskWithUnit, "TB") {
					// 规格中storage为数字+TB
					dataDisk, err := strconv.Atoi(strings.Split(dataDiskWithUnit, "TB")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid storage format")
						continue
					}
					param.Storage += dataDisk * 1024
					param.DataDisk = dataDisk * 1024
				} else {
					global.GVA_LOG.Error("Invalid storage format")
					continue
				}
			case "带宽":
				bandwidthWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(bandwidthWithUnit, "Mbps") {
					//规格中bandwidth为数字+Mbps
					param.BandWidth, err = strconv.Atoi(strings.Split(bandwidthWithUnit, "Mbps")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid bandwidth format")
						continue
					}
				} else if strings.Contains(bandwidthWithUnit, "M") {
					//规格中bandwidth为数字+M
					param.BandWidth, err = strconv.Atoi(strings.Split(bandwidthWithUnit, "M")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid bandwidth format")
						continue
					}
				} else {
					global.GVA_LOG.Error("Invalid bandwidth format")
					continue
				}
			case "GPU数量":
				gpuNumberWithUnit := spec.Value[spec.ValueIndex]
				if strings.Contains(gpuNumberWithUnit, "卡") {
					//规格中gpuNumber为数字+卡
					param.GpuNumber, err = strconv.Atoi(strings.Split(gpuNumberWithUnit, "卡")[0])
					if err != nil {
						global.GVA_LOG.Error("Invalid gpuNumber format")
						continue
					}
				} else {
					global.GVA_LOG.Error("Invalid gpuNumber format")
					continue
				}
			}
		}
	}
	var storageSlice []map[string]any
	var storageObject map[string]any
	err = json.Unmarshal([]byte(allocateDetailsMap["storage"].(string)), &storageSlice)
	if err != nil {
		global.GVA_LOG.Error("Invalid storage format")
		if err = json.Unmarshal([]byte(allocateDetailsMap["storage"].(string)), &storageObject); err != nil {
			global.GVA_LOG.Error("Invalid storage format")
			return err, param
		} else {
			storageSlice = append(storageSlice, storageObject)
		}
	}
	//storageJson是一个数组,但是业务上只会有一个元素[{"tag": "nas", "name": "文件存储"}]
	for _, storage := range storageSlice {
		//cloudPods支持系统盘和数据盘使用不同的存储方式（块存储和文件存储），但业务上不允许所以这里只判断一个
		if storage["tag"].(string) == "nas" {
			param.DataStorageType = "nfs"
			param.SysStorageType = "nfs"
		} else if storage["tag"].(string) == "block" {
			param.SysStorageType = "local"
			param.DataStorageType = "local"
		} else {
			global.GVA_LOG.Error("Invalid storage format")
			return err, param
		}
	}
	return nil, param
}

func (r *ResourceService) AllocateResponse(errMsg error, orderComputingId int, resourceID string, isCustom bool) error {
	//errMsg为空表示分配资源没有出错，非空代表出错 ******************************************8888
	//查数据库
	var renterRes renter.RenterRes
	var resMsg protocol.PackageAllocateResponse
	if errMsg != nil && orderComputingId != 0 {
		//分配失败的处理
		resMsg.ID = -1
		resMsg.Command = protocol.AllocateRetCmd
		resMsg.RetCode = protocol.AllocateFailedCode
		resMsg.RetMsg = protocol.AllocateFailedMsg
		resMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
		resMsg.OrderComputingID = orderComputingId
		return global.GVA_WS.SendMessageDebug(resMsg)
	}
	//异步应答，将报文ID设置-1,防止干扰平台的重发
	resMsg.ID = -1
	resMsg.Command = protocol.AllocateRetCmd
	resMsg.RetCode = protocol.AllocateSuccessCode
	resMsg.RetMsg = protocol.AllocateSuccessMsg
	resMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	//分配完成的时间
	resMsg.Time = time.Now()
	if err := global.GVA_DB.Where("resource_id = ?", resourceID).First(&renterRes).Error; err != nil {
		return err
	} else {
		if isCustom {
			// 定制产品由前端手动分配，分配是否成功未知
			go r.CustomAllocateResponse(renterRes)
		} else {
			resMsg.StartTime = renterRes.StartTime
			resMsg.EndTime = renterRes.EndTime
			resMsg.CompositeResourceID = renterRes.ResourceID
			resMsg.OrderID = renterRes.OrderId
			resMsg.OrderComputingID = renterRes.OrderComputingId
			resMsg.Url = renterRes.Url
			resMsg.SshAddress = renterRes.SshHost
			resMsg.SshPort = renterRes.SshPort
			resMsg.SshAccount = renterRes.SshUser
			resMsg.SshSecret = renterRes.SshPwd
			return global.GVA_WS.SendMessageDebug(resMsg)
		}
	}
	return nil
}

func (r *ResourceService) CustomAllocateResponse(renterRes renter.RenterRes) {
	var resMsg protocol.PackageAllocateResponse
	//定制产品由前端手动分配，分配是否成功未知
	//检查分配状态
	var CreatErr error
	CreatErr = nil
	if renterRes.Type == 1 {
		CreatErr = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareMetalService.HostForecastById(renterRes.ResourceID)
	} else if renterRes.Type == 2 {
		CreatErr = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService.HostForecastById(renterRes.ResourceID)
	}
	if CreatErr == nil {
		resMsg.Time = time.Now()
		if err := global.GVA_DB.Model(&renterRes).Where("resource_id = ?", renterRes.ResourceID).First(&renterRes).Error; err != nil {
			global.GVA_LOG.Error("查询失败")
			return
		}

		resMsg.StartTime = renterRes.StartTime
		resMsg.EndTime = renterRes.EndTime
		//异步应答，将报文ID设置-1,防止干扰平台的重发
		resMsg.ID = -1
		resMsg.Command = protocol.AllocateRetCmd
		resMsg.RetCode = protocol.AllocateSuccessCode
		resMsg.RetMsg = protocol.AllocateSuccessMsg
		resMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
		resMsg.StartTime = renterRes.StartTime
		resMsg.EndTime = renterRes.EndTime
		resMsg.CompositeResourceID = renterRes.ResourceID
		resMsg.OrderID = renterRes.OrderId
		resMsg.OrderComputingID = renterRes.OrderComputingId
		resMsg.Url = renterRes.Url
		resMsg.SshAddress = renterRes.SshHost
		resMsg.SshPort = renterRes.SshPort
		resMsg.SshAccount = renterRes.SshUser
		resMsg.SshSecret = renterRes.SshPwd
	} else {
		//分配失败的处理
		resMsg.ID = -1
		resMsg.Command = protocol.AllocateRetCmd
		resMsg.RetCode = protocol.AllocateFailedCode
		resMsg.RetMsg = protocol.AllocateFailedMsg
		resMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
		resMsg.OrderComputingID = renterRes.OrderComputingId
	}
	err := global.GVA_WS.SendMessageDebug(resMsg)
	if err != nil {
		global.GVA_LOG.Error("发送资源分配应答失败")
	}
}

func (r *ResourceService) ReleaseResource(message *protocol.ReleaseRequest) error {
	response := protocol.ReleaseResponse{}
	response.ID = message.ID
	response.Command = protocol.ReleaseRetCmd
	//释放资源,更新租户资源表
	var renterRes renter.RenterRes
	renterRes.Status = 4
	renterRes.EndTime = time.Now()
	err := global.GVA_DB.Where("resource_id = ?", message.CompositeResourceID).Updates(&renterRes).Error
	if err != nil {
		return err
	}
	global.GVA_LOG.Info("回收资源", zap.String("id:", message.CompositeResourceID))
	err = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService.HostDelete(message.CompositeResourceID)
	if err != nil {
		return err
	}
	response.RetCode = 0
	response.RetMsg = "释放成功"
	response.Time = renterRes.EndTime
	response.CompositeResourceID = message.CompositeResourceID
	return global.GVA_WS.SendMessageDebug(response)
}
