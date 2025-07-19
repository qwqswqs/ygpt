package system

import (
	"time"
	"ygpt/server/global"
)

type SystemSoftware struct {
	global.GVA_MODEL
	//Name 软件名称
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:软件名称"`
	//Type  类别
	Type string `json:"type" gorm:"column:type;type:text;comment:类别"`
	//Version 版本
	Version string `json:"version" gorm:"column:version;type:text;comment:版本"`
	//Url 地址
	Url string `json:"url" gorm:"column:url;type:text;comment:路径"`
	//Size 大小
	Size string `json:"size" gorm:"column:size;type:text;comment:大小"`
	//LastUpdate 最后更新一次时间
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update;default:null;comment:最后一次更新时间"`
	ReSoftwareID int       `json:"reSoftwareID" gorm:"column:re_software_id;type:bigint;comment:关联平台软件id"`
}

func (SystemSoftware) TableName() string {
	return "yun_system_software_info"
}
