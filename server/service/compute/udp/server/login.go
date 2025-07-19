package server

import (
	"ygpt/server/global"
	"ygpt/server/model/yunguan/renter"
)

func LoginQuery(privateIps []string) bool {
	var count int64
	global.GVA_DB.Model(&renter.RenterRes{}).Where("private_ip IN ?", privateIps).Count(&count)
	if count == 1 {
		return true
	} else if count > 1 {
		global.GVA_LOG.Error("private_ip重复")
		return true
	}
	return false
}
