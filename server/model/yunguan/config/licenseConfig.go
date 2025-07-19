package config

import "gorm.io/gorm"

type LicenseConfig struct {
	gorm.Model
	StartTime   int64  `gorm:"comment:'开始时间'" json:"startTime"`
	EndTime     int64  `gorm:"comment:'结束时间'" json:"endTime"`
	DeviceNum   []byte `gorm:"comment:'设备数量';type:BLOB" json:"deviceNum"`
	RenterNum   []byte `gorm:"comment:'租户数量';type:BLOB" json:"renterNum"`
	AdminNum    []byte `gorm:"comment:'管理员数量';type:BLOB" json:"adminNum"`
	LicenseUser string `gorm:"comment:'License用户名'" json:"licenseUser"`
	LicensePwd  string `gorm:"comment:'License密码'" json:"licensePwd"`
}

func (o *LicenseConfig) TableName() string {
	return "yun_license_config"
}

type LicenseConfigData struct {
	ID          uint   `json:"id"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	DeviceNum   int64  `json:"deviceNum"`
	RenterNum   int64  `json:"renterNum"`
	AdminNum    int64  `json:"adminNum"`
	LicenseUser string `json:"licenseUser"`
	LicensePwd  string `json:"licensePwd"`
}
