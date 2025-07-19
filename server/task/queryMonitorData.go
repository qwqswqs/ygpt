package task

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/renter"
	modelsys "ygpt/server/model/yunguan/system"
	"ygpt/server/service"
	"ygpt/server/service/yunguan/system"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type MonitorData struct {
	TimeStamp int64
	Value     float64
	Object    string
}

type Instance struct {
	RenterRes     renter.RenterRes
	InstancesName string
	InstanceID    string
}

func QueryMonitorData() {
	serverService := service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ServerService
	bareHostService := service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareHostService
	//监控查询条件
	time := "1h"     //查询范围1个小时内
	interval := "5m" //取值间隔5分钟
	// 获取所有云主机和裸金属列表（容器不做监控报警）
	s, err := global.GetCloudClientSession()
	if err != nil {
		return
	}
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("filter", jsonutils.NewString("hypervisor.notin(container)"))
	result, err := compute.Servers.List(s, params)
	if err != nil {
		global.GVA_LOG.Error("查询实例失败", zap.Any("err", err))
		return
	}
	// 查询数据库中的云主机，裸金属
	var renterRes []renter.RenterRes
	if err = global.GVA_DB.Model(&renter.RenterRes{}).Where("type != ? and status = ? and resource_id != ?", 3, 2, "").Find(&renterRes).Error; err != nil {
		global.GVA_LOG.Error("查询数据库中的云主机，裸金属失败", zap.Any("err", err))
		return
	}
	// 过滤云主机，裸金属
	instances := make([]Instance, 0)
	instance := Instance{}
	for _, r := range renterRes {
		for _, v := range result.Data {
			id, err := v.GetString("id")
			if err != nil {
				global.GVA_LOG.Error("获取云主机id失败", zap.Any("err", err))
				return
			}
			if r.ResourceID == id {
				instance.RenterRes = r
				name, err := v.GetString("name")
				if err != nil {
					global.GVA_LOG.Error("获取云主机名称失败", zap.Any("err", err))
					return
				}
				instance.InstancesName = name
				instance.InstanceID = id
				instances = append(instances, instance)
			}
		}
	}

	for _, r := range instances {
		cpuData, err := serverService.MonitorCpuData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		//global.GVA_LOG.Info("CPU监控数据", zap.Any("cpuData", cpuData))
		memData, err := serverService.MonitorMemData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		diskData, err := serverService.MonitorDiskData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		drData, err := serverService.MonitorDiskRead(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		dwData, err := serverService.MonitorDiskWrite(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		brData, err := serverService.MonitorBpsRecvData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		bsData, err := serverService.MonitorBpsSentData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		//prData, err := serverService.MonitorPpsRecvData(r.RenterRes.ResourceID, time, interval)
		//if err != nil {
		//	global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		//}
		//psData, err := serverService.MonitorPpsSentData(r.RenterRes.ResourceID, time, interval)
		//if err != nil {
		//	global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		//}
		gpuData, err := serverService.MonitorGpuData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		gpuMemData, err := serverService.MonitorGpuMemData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		gpuTemData, err := serverService.MonitorGpuTemperatureData(r.RenterRes.ResourceID, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		//global.GVA_LOG.Info("监控数据", zap.Any("cpuData", cpuData), zap.Any("memData", memData), zap.Any("diskData", diskData),
		//	zap.Any("drData", drData), zap.Any("dwData", dwData), zap.Any("brData", brData), zap.Any("bsData", bsData),
		//	zap.Any("gpuData", gpuData), zap.Any("gpuMemData", gpuMemData),
		//	zap.Any("gpuTemData", gpuTemData))

		//
		checkAlertStrategy(1, HandleMonitorData(cpuData, "cpu"), r)
		checkAlertStrategy(2, HandleMonitorData(memData, "mem"), r)
		checkAlertStrategy(3, HandleMonitorData(diskData, "disk"), r)
		checkAlertStrategy(4, HandleMonitorData(drData, "diskRead"), r)
		checkAlertStrategy(5, HandleMonitorData(dwData, "diskWrite"), r)
		checkAlertStrategy(6, HandleMonitorData(brData, "bpsRecv"), r)
		checkAlertStrategy(7, HandleMonitorData(bsData, "bpsSent"), r)
		//checkAlertStrategy(8,HandleMonitorData(prData, "network"),r)
		//checkAlertStrategy(9,HandleMonitorData(psData, "network"),r)
		checkAlertStrategy(10, HandleMonitorData(gpuData, "utilization_gpu"), r)
		checkAlertStrategy(11, HandleMonitorData(gpuMemData, "utilization_memory"), r)
		checkAlertStrategy(12, HandleMonitorData(gpuTemData, "temperature_gpu"), r)
	}
	// 获取宿主机列表
	params = jsonutils.NewDict()
	params.Set("baremetal", jsonutils.NewBool(false))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	result, err = compute.Hosts.List(s, params)
	if err != nil {
		return
	}
	for _, v := range result.Data {
		id, err := v.GetString("id")
		if err != nil {
			global.GVA_LOG.Error("获取宿主机id失败", zap.Any("err", err))
		}
		ins := Instance{}
		ins.InstancesName, err = v.GetString("name")
		if err != nil {
			global.GVA_LOG.Error("获取宿主机名称失败", zap.Any("err", err))
		}
		ins.InstanceID = id
		cpuData, err := bareHostService.MonitorCpuData(id, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		//global.GVA_LOG.Info("CPU监控数据", zap.Any("cpuData", cpuData))
		memData, err := bareHostService.MonitorMemData(id, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		diskData, err := bareHostService.MonitorDiskData(id, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		//drData, err := bareHostService.MonitorDiskRead(id, time, interval)
		//if err != nil {
		//	global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		//}
		//dwData, err := bareHostService.MonitorDiskWrite(id, time, interval)
		//if err != nil {
		//	global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		//}
		brData, err := bareHostService.MonitorBpsRecvData(id, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}
		bsData, err := bareHostService.MonitorBpsSentData(id, time, interval)
		if err != nil {
			global.GVA_LOG.Error("查询监控数据失败", zap.Any("err", err))
		}

		checkAlertStrategy(1, HandleMonitorData(cpuData, "cpu"), ins)
		checkAlertStrategy(2, HandleMonitorData(memData, "hostMem"), ins)
		checkAlertStrategy(3, HandleMonitorData(diskData, "disk"), ins)
		//checkAlertStrategy(4, HandleMonitorData(drData, "diskRead"), ins)
		//checkAlertStrategy(5, HandleMonitorData(dwData, "diskWrite"), ins)
		checkAlertStrategy(6, HandleMonitorData(brData, "bpsRecv"), ins)
		checkAlertStrategy(7, HandleMonitorData(bsData, "bpsSent"), ins)
	}
	//resultMap := jsonutils.Marshal(result).Interface()
	//global.GVA_LOG.Info("宿主机列表", zap.Any("result", resultMap))
}

func HandleMonitorData(data interface{}, monitorName string) []MonitorData {
	monitorData := make([]MonitorData, 0)
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		global.GVA_LOG.Error("监控数据断言失败", zap.Any("err", "data is not map[string]interface{}"))
		return monitorData
	}
	if seriesTotal, ok := dataMap["series_total"].(int64); ok {
		if seriesTotal <= 0 {
			return monitorData
		}
	} else {
		return monitorData
	}
	series, ok := dataMap["series"].([]interface{})
	if !ok {
		global.GVA_LOG.Error("监控数据断言失败", zap.Any("err", "series is not []map[string]interface{}"))
		return monitorData
	}
	//  以下处理逻辑是找出最高的监控指标的值并返回,并且认为一种监控指标只有一个监控对象（不考虑cpu多核心，多磁盘，多网卡等因素）
	for _, seriesItem := range series {
		Item, ok := seriesItem.(map[string]interface{})
		if !ok {
			return monitorData
		}
		points, ok := Item["points"].([]interface{})
		if !ok {
			return monitorData
		}

		//  取最后一个点与策略进行比较
		value := 0.0
		timeStamp := int64(0)
		// 宿主机的内存监控指标比较特殊
		if monitorName == "hostMem" {
			var err error
			pointsLen := len(points)
			if pointsLen < 3 {
				return monitorData
			}
			point, ok := points[pointsLen-1].([]interface{})
			if !ok {
				return monitorData
			}
			value, err = toFloat64(point[0])
			if err != nil {
				return monitorData
			}
			timeStamp, ok = point[2].(int64)
			if !ok {
				return monitorData
			}
		} else {
			var err error
			pointsLen := len(points)
			if pointsLen < 2 {
				return monitorData
			}
			point, ok := points[pointsLen-1].([]interface{})
			if !ok {
				return monitorData
			}
			value, err = toFloat64(point[0])
			if err != nil {
				return monitorData
			}
			timeStamp, ok = point[1].(int64)
			if !ok {
				return monitorData
			}
		}
		lastPointMonitorData := MonitorData{
			Value:     value,
			TimeStamp: timeStamp,
		}
		// 只保留值最大的监控数据
		if len(monitorData) > 0 {
			if monitorData[len(monitorData)-1].Value < lastPointMonitorData.Value {
				monitorData = monitorData[:len(monitorData)-1]
			}
		}
		monitorData = append(monitorData, lastPointMonitorData)
		// 取倒数第二个点与策略进行比较
		//if pointsLen-2 >= 0 {
		//	penultimatePointMonitorData := MonitorData{
		//		Value: points[pointsLen-2][0].(float64),
		//	}
		//	monitorData = append(monitorData, penultimatePointMonitorData)
		//}
	}
	return monitorData
}

func checkAlertStrategy(monitorType int, monitorData []MonitorData, i Instance) {
	for _, monitorDataItem := range monitorData {
		ept := renter.RenterRes{}
		if i.RenterRes == ept {
			err := Alert(monitorDataItem.Value, monitorType, 4, 0, i.InstanceID, i.InstancesName)
			if err != nil {
				//global.GVA_LOG.Error("告警失败", zap.Error(err))
				return
			}
		} else {
			err := Alert(monitorDataItem.Value, monitorType, i.RenterRes.Type, int64(i.RenterRes.RenterID), i.RenterRes.ResourceID, i.InstancesName)
			if err != nil {
				//global.GVA_LOG.Error("告警失败", zap.Error(err))
				return
			}
		}
	}
}

func Alert(alertValue float64, alertType int, resourceType int, renterID int64, ResourceID string, ResourceName string) error {
	// 检查最近最新的报警记录
	//newestAlertData := modelsys.SystemAlert{}
	//if err := global.GVA_DB.Where("resource_type = ? AND renter_id = ? AND resource_id = ? AND alert_type = ?", resourceType, renterID, instanceID, alertType).Last(&newestAlertData).Error; err == nil {
	//	// 如果存在且距离现在小于10分钟，则不发送报警,避免频繁报警
	//	if newestAlertData.AlertTime.Unix() > 0 && newestAlertData.AlertTime.Unix()+600 > time.Now().Unix() {
	//		global.GVA_LOG.Info("报警已存在，距离现在小于10分钟，不发送报警", zap.Any("newestAlertData", newestAlertData))
	//		return nil
	//	}
	//}
	alertService := system.SystemAlertService{}
	alertData := modelsys.SystemAlert{}
	alertData.AlertData = int64(alertValue)
	alertData.ResourceType = int64(resourceType)
	alertData.RenterID = renterID
	alertData.ResourceID = ResourceID
	alertData.ResourceName = ResourceName
	alertData.ResourceType = int64(resourceType)
	alertData.AlertType = int64(alertType)
	return alertService.BeforeAddAlert(alertData)
}

func toFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case json.Number:
		return v.Float64()
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
