package system

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type SystemConfig struct {
	global.GVA_MODEL
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:配置名称"`
	Key         string `json:"key" gorm:"column:key;type:text;comment:配置关键字"`
	Value       string `json:"value" gorm:"column:value;type:json;comment:配置值"`
	Description string `json:"description" gorm:"column:description;type:text;comment:配置描述"`
	Type        int    `json:"type" gorm:"column:type;type:int;comment:'配置类别(1.安全 2.备份 3.系统)'"`
}

func (SystemConfig) TableName() string {
	return "yun_system_config"
}

type GetSystemConfigReq struct {
	request.PageInfo
	Type int `json:"type"`
}
