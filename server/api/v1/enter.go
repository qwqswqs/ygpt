package v1

import (
	cloud "ygpt/server/api/v1/cloudpods"
	"ygpt/server/api/v1/compute"
	"ygpt/server/api/v1/example"
	"ygpt/server/api/v1/system"
	"ygpt/server/api/v1/yunguan"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	YunGuanApiGroup yunguan.ApiGroup
	ComputeApiGroup compute.ApiGroup
	CloudApiGroup   cloud.CloudApiGroup
}
