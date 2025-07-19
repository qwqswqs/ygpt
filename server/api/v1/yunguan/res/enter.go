package res

import "ygpt/server/service"

type ResApi struct {
	DiskApi
	ResDeviceApi
	ResInfoApi
}

var (
	DiskService      = service.ServiceGroupApp.YunGuanServiceGroup.ResServiceGroup.DiskService
	resDeviceService = service.ServiceGroupApp.YunGuanServiceGroup.ResServiceGroup.ResDeviceService
	resInfoService   = service.ServiceGroupApp.YunGuanServiceGroup.ResServiceGroup.ResInfoService
)
