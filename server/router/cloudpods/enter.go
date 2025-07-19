package yunguan

import (
	"ygpt/server/router/cloudpods/baseRes"
	"ygpt/server/router/cloudpods/common"
	"ygpt/server/router/cloudpods/compute"
	"ygpt/server/router/cloudpods/image"
	"ygpt/server/router/cloudpods/k8s"
	"ygpt/server/router/cloudpods/network"
	"ygpt/server/router/cloudpods/storage"
)

type RouterGroup struct {
	compute.ComputeRouter
	image.ImageRouterGroup
	network.NetRouterGroup
	storage.StorageRouterGroup
	baseRes.BaseResRouterGroup
	k8s.K8sRouterGroup
	common.CloudCommonRouter
}
