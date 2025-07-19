package network

import "ygpt/server/global"

func Migrate() {
	err := global.GVA_DB.AutoMigrate(&NetworkIp{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&NetworkVpc{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&NetworkNat{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&NetworkSub{})
	if err != nil {
		panic(err)
	}
}
