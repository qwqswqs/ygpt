package system

import (
	"time"
	"ygpt/server/global"
)

type SystemAlert struct {
	global.GVA_MODEL
	ResourceType int64     `gorm:"comment:'资源类型 1裸金属2虚拟机3容器4宿主机'" json:"resourceType"`
	AlertType    int64     `gorm:"comment:'报警类型 1.CPU使用率;2.内存使用率；3.磁盘使用率；4.磁盘读速率；5.磁盘写速率；6.网络入带宽；7.网络出带宽；8.网络收包速率；9.网络发包速率；10.GPU使用率；11.GPU显存使用率；12.GPU温度。'" json:"alertType"`
	AlertObject  string    `gorm:"comment:'报警对象'" json:"alertObject"`
	RenterID     int64     `gorm:"comment:'租户ID'" json:"renterID"`
	ResourceID   string    `gorm:"comment:'实例ID'" json:"resourceID"`
	ResourceName string    `gorm:"comment:'实例名称'" json:"resourceName"`
	AlertTime    time.Time `gorm:"comment:'报警时间'" json:"alertTime"`
	AlertData    int64     `gorm:"comment:'报警值'" json:"alertData"`
	IsCheck      bool      `gorm:"comment:'是否查看';default:false;" json:"isCheck"`
	AlertLevel   int64     `gorm:"comment:'报警级别'" json:"alertLevel"`
	IsRenterData bool      `gorm:"comment:'租户可查看'" json:"isRenterData"`
}

func (SystemAlert) TableName() string {
	return "yun_system_alert"
}
