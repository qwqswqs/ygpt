package renter

import "ygpt/server/service"

type RenterGroupApi struct {
	RenterApi
	RenterTaskApi
	RenterResApi
	RenterRecordApi
}

var (
	renterService       = service.ServiceGroupApp.YunGuanServiceGroup.RenterServiceGroup.RenterService
	renterResService    = service.ServiceGroupApp.YunGuanServiceGroup.RenterServiceGroup.RenterResService
	renterRecordService = service.ServiceGroupApp.YunGuanServiceGroup.RenterServiceGroup.RenterRecordService
	renterTaskService   = service.ServiceGroupApp.YunGuanServiceGroup.RenterServiceGroup.RenterTaskService
)
