package system

import "ygpt/server/global"

type SystemLog struct {
	global.GVA_MODEL
	Level int    `json:"level" gorm:"column:level;type:tinyint;not null;comment:级别"`
	Info  string `json:"info" gorm:"column:info;type:text;comment:日志信息"`
	Birth string `json:"birth" gorm:"column:birth;type:text;comment:日志源"`
}

func (SystemLog) TableName() string {
	return "yun_system_log"
}
