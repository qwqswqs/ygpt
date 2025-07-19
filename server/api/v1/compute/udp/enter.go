package udp

import "ygpt/server/service"

type UdpApi struct {
	InstanceApi
}

var (
	InstanceService = service.ServiceGroupApp.ComputeServiceGroup.UdpServiceGroup.InstanceServiceGroup.InstanceService
)
