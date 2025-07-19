package yunguan

import (
	"ygpt/server/service/yunguan/backup"
	"ygpt/server/service/yunguan/cloudpods"
	"ygpt/server/service/yunguan/network"
	"ygpt/server/service/yunguan/product"
	"ygpt/server/service/yunguan/renter"
	"ygpt/server/service/yunguan/res"
	"ygpt/server/service/yunguan/system"
)

type YunGuanServiceGroup struct {
	system.SystemServiceGroup
	res.ResServiceGroup
	backup.BackupServiceGroup
	network.NetworkServiceGroup
	product.ProductServiceGroup
	renter.RenterServiceGroup
	cloudpods.CloudpodsServiceGroup
}
