package renter

import (
	"ygpt/server/global"
)

type RenterRecord struct {
	global.GVA_MODEL
	RenterID int    `json:"renterID" gorm:"column:renter_id;type:bigint;comment:租户ID"`
	Content  string `json:"content" gorm:"column:content;type:text;comment:内容(登录、下线、创建任务、备份)"`
}

func (RenterRecord) TableName() string {
	return "yun_renter_record"
}
