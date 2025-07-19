package backup

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type BackupSnapshot struct {
	global.GVA_MODEL
	Name         string `json:"name" gorm:"column:name;type:text;comment:'快照名称'"`
	InstanceID   int    `json:"instanceID" gorm:"column:instance_id;type:bigint;comment:'实例资源ID'"`
	Size         int    `json:"size" gorm:"column:size;type:bigint;comment:'快照大小'"`
	RenterID     uint   `json:"renterID" gorm:"column:renter_id;type:bigint;comment:'租户 id'"`
	Description  string `json:"description" gorm:"column:description;type:text;comment:'说明'"`
	SpecDesc     string `json:"specDesc" gorm:"column:spec_desc;type:json;comment:'描述json'"`
	SnapshotType int    `json:"snapshotType" gorm:"column:snapshot_type;type:tinyint;comment:'快照类型(1.磁盘 2.内存 3.实例)'"`
}

func (e *BackupSnapshot) TableName() string {
	return "yun_backup_snapshot"
}

type GetBakSnapReq struct {
	request.PageInfo
	RenterID int `json:"renterID"`
}
