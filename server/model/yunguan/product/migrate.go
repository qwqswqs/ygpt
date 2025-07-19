package product

import "ygpt/server/global"

func Migrate() {
	err := global.GVA_DB.AutoMigrate(&ProductSupply{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductComputing{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductElementInfo{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductElementDownload{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductElementRecord{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductConfig{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductCategory{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductElem{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&ProductFile{})
	if err != nil {
		panic(err)
	}
}
