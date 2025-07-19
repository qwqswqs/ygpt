package service

import (
	"ygpt/server/service/compute"
	"ygpt/server/service/example"
	"ygpt/server/service/system"
	"ygpt/server/service/yunguan"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	YunGuanServiceGroup yunguan.YunGuanServiceGroup
	ComputeServiceGroup compute.ComputeServiceGroup
}
