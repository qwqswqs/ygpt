package baseRes

import "ygpt/server/service"

type BaseResApiGroup struct {
	BaseHostApi
	BaseK8SApi
	BaseDevicesApi
}

var (
	bareHostService = service.ServiceGroupApp.YunGuanServiceGroup.CloudpodsServiceGroup.BareHostService
)
