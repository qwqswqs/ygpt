package backup

import "ygpt/server/global"

// Migrate 自动迁移数据库
func Migrate() {
	err := global.GVA_DB.AutoMigrate(&Backup{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&BackupImage{})
	if err != nil {
		panic(err)
	}
	err = global.GVA_DB.AutoMigrate(&BackupSnapshot{})
	if err != nil {
		panic(err)
	}
}
