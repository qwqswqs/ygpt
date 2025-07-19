package cloudpods

import "ygpt/server/service"

type CloudpodsApi struct {
	VirtualMachineApi
	BareMetalApi
	ContainerApi
	ImageApi
	IpApi
	VpcApi
	AlarmApi
	JumpApi
	ResInfoApi
}

var (
	ImageService          = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ImageService
	ContainerService      = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ContainerService
	VirtualMachineService = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.VirtualMachineService
	BareMetalService      = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareMetalService
	ResInfoService        = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ResInfoService
)
