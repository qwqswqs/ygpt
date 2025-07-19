package system

import "ygpt/server/global"

// Migrate 自动迁移数据库
func Migrate() {
	err := global.GVA_DB.AutoMigrate(&SystemLog{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemTool{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemConfig{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemCache{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemTicket{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemAlert{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemSoftware{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemType{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&SystemOperate{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductRegion{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductSysType{})
	if err != nil {
		panic(err)
	}
}
