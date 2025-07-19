package k8s

import "ygpt/server/global"

func Migrate() {
	err := global.GVA_DB.AutoMigrate(&K8sModel{})
	if err != nil {
		panic(err)
	}
}
