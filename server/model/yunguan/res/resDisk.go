package res

import (
	"gorm.io/gorm"
	"time"
	"ygpt/server/model/common/request"
)

type ResDisk struct {
	gorm.Model
	RegisterTime  time.Time `json:"registerTime" gorm:"column:register_time;type:datetime;comment:'登记时间'"`
	Status        int       `json:"status" gorm:"column:status;type:tinyint;comment:'状态（使用，空闲，损坏）'"`
	Capacity      uint32    `json:"capacity" gorm:"column:capacity;type:bigint;comment:'容量大小'"`
	StorageType   int       `json:"storageType" gorm:"column:storage_type;type:tinyint;comment:'存储类别'"`
	Brand         string    `json:"brand" gorm:"column:brand;type:text;comment:'品牌'"`
	Specification string    `json:"specification" gorm:"column:specification;type:json;comment:'规格json'"`
	SerialNumber  string    `json:"serialNumber" gorm:"column:serial_number;type:text;comment:'编号'"`
	DiskFormat    string    `json:"diskFormat" gorm:"column:disk_format;type:text;comment:'磁盘格式'"`
	ServerID      int       `json:"serverID" gorm:"column:server_id;type:text;comment:'所属服务器ID'"`
}

func (o *ResDisk) TableName() string {
	return "yun_res_disk"
}

type ResDiskReq struct {
	gorm.Model
	request.PageInfo
	Status int `json:"status"`
}
