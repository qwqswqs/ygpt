package res

import "ygpt/server/global"

func Migrate() {
	err := global.GVA_DB.AutoMigrate(&ResInfo{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ResDisk{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ResDevice{})
	if err != nil {
		panic(err)
	}
}
