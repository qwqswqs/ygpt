package system

import (
	"ygpt/server/global"
)

type SystemTool struct {
	global.GVA_MODEL
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:工具名称"`
	Address     string `json:"address" gorm:"column:address;type:text;comment:所在地址"`
	Description string `json:"description" gorm:"column:description;type:text;comment:作用说明"`
	Working     int    `json:"working" gorm:"column:working;type:int;comment:所在服务器ID"`
	TaskID      string `json:"taskID" gorm:"column:task_id;type:text;comment:进程号"`
}

func (SystemTool) TableName() string {
	return "yun_system_tool"
}
