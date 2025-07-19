package system

import (
	"time"
	"ygpt/server/global"
)

type SystemCache struct {
	global.GVA_MODEL
	Address      string    `json:"address" gorm:"column:address;type:varchar(255);not null;comment:地址"`
	Size         int       `json:"size" gorm:"column:size;type:bigint;comment:大小"`
	UploadTime   time.Time `json:"uploadTime" gorm:"column:upload_time;type:time;comment:上传时间"`
	LastDownload time.Time `json:"lastDownload" gorm:"column:last_download;default:null;comment:最近下载时间"`
	Type         string    `json:"type" gorm:"column:type;type:text;comment:类别"`
	UserType     int       `json:"userType" gorm:"column:user_type;type:text;comment:上传用户类别"`
	UserID       int       `json:"userID" gorm:"column:user_id;type:text;comment:上传用户ID"`
}

func (SystemCache) TableName() string {
	return "yun_system_cache"
}
