package compute

import "ygpt/server/service"

type ComputeApiGroup struct {
	ComputeServersApi
	ScalingConfigApi
	ScalingGroupApi
	ScalingPolicyApi
}

var (
	serverService = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.ServerService
)
