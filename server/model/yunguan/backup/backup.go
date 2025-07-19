package backup

import (
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type Backup struct {
	global.GVA_MODEL
	Name       string    `json:"name" gorm:"column:name;type:text;comment:'名称'"`
	Size       int       `json:"size" gorm:"column:size;type:bigint;comment:'大小'"`
	RenterID   int       `json:"renterID" gorm:"column:renter_id;type:bigint;comment:'关联租户 id'"`
	UploadTime time.Time `json:"uploadTime" gorm:"column:upload_time;comment:'上传时间'"`
}

func (e *Backup) TableName() string {
	return "yun_backup_back"
}

type BackReq struct {
	global.GVA_MODEL
	request.PageInfo
	Name     string `json:"name"`
	RenterID int    `json:"renterID"`
}
