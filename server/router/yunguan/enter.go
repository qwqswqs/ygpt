package yunguan

import (
	"ygpt/server/router/yunguan/backup"
	"ygpt/server/router/yunguan/cloudpods"
	"ygpt/server/router/yunguan/config"
	"ygpt/server/router/yunguan/network"
	"ygpt/server/router/yunguan/product"
	"ygpt/server/router/yunguan/renter"
	"ygpt/server/router/yunguan/res"
	"ygpt/server/router/yunguan/system"
)

type RouterGroup struct {
	system.SystemRouter
	res.ResRouter
	backup.BackupRouter
	network.NetworkRouter
	product.ProductRouter
	renter.RenterGroupRouter
	cloudpods.CloudpodsRouter
	config.ConfigRouter
}
