package yunguan

import (
	"ygpt/server/api/v1/yunguan/backup"
	"ygpt/server/api/v1/yunguan/cloudpods"
	"ygpt/server/api/v1/yunguan/config"
	"ygpt/server/api/v1/yunguan/network"
	"ygpt/server/api/v1/yunguan/product"
	"ygpt/server/api/v1/yunguan/renter"
	"ygpt/server/api/v1/yunguan/res"
	"ygpt/server/api/v1/yunguan/system"
)

type ApiGroup struct {
	system.SystemApi
	res.ResApi
	backup.BackupApi
	network.NetworkApi
	product.ProductApi
	renter.RenterGroupApi
	cloudpods.CloudpodsApi
	config.ConfigApi
}
