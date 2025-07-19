package cloudpods

import (
	"ygpt/server/api/v1/cloudpods/baseRes"
	"ygpt/server/api/v1/cloudpods/common"
	"ygpt/server/api/v1/cloudpods/compute"
	"ygpt/server/api/v1/cloudpods/image"
	"ygpt/server/api/v1/cloudpods/k8s"
	"ygpt/server/api/v1/cloudpods/network"
	"ygpt/server/api/v1/cloudpods/storage"
)

type CloudApiGroup struct {
	compute.ComputeApiGroup
	image.ImageApiGroup
	network.NetApiGroup
	storage.StorageApiGroup
	baseRes.BaseResApiGroup
	k8s.K8sAPiGroup
	common.CloudCommonApiGroup
}
