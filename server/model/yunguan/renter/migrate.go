package renter

import "ygpt/server/global"

// Migrate 自动迁移数据库
func Migrate() {
	err := global.GVA_DB.AutoMigrate(&Renter{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&RenterTask{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&RenterRecord{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&RenterRes{})
	if err != nil {
		panic(err)
	}
}
