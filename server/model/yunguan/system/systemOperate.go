package system

import (
	"time"
	"ygpt/server/global"
)

type SystemOperate struct {
	global.GVA_MODEL
	Type        int       `json:"type" gorm:"type:int;comment:'操作类型 1.一键巡查 2.病毒扫描 3.漏洞扫描'"`
	StartTime   time.Time `json:"startTime" gorm:"type:datetime;comment:'开始时间'"`
	EndTime     time.Time `json:"endTime" gorm:"type:datetime;comment:'结束时间'"`
	ProblemNum  int       `json:"problemNum" gorm:"type:int;comment:'发现问题数量'"`
	ProblemText string    `json:"problemText" gorm:"type:text;comment:'问题记录文本'"`
	Operator    int       `json:"operator" gorm:"type:int;comment:'操作人ID'"`
}

func (SystemOperate) TableName() string {
	return "yun_system_operate"
}
