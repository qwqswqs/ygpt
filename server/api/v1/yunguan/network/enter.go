package network

import "ygpt/server/service"

type NetworkApi struct {
	NetworkVpcApi
	NetworkNatApi
	NetworkIpApi
	NetworkSubApi
}

var (
	IpService  = service.ServiceGroupApp.YunGuanServiceGroup.NetworkServiceGroup.IpService
	subService = service.ServiceGroupApp.YunGuanServiceGroup.NetworkServiceGroup.SubService
	NatService = service.ServiceGroupApp.YunGuanServiceGroup.NetworkServiceGroup.NatService
	VpcService = service.ServiceGroupApp.YunGuanServiceGroup.NetworkServiceGroup.VpcService
)
