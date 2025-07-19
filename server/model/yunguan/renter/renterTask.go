package renter

import (
	"time"
	"ygpt/server/global"
)

type RenterTask struct {
	global.GVA_MODEL
	DatasetID   int       `json:"datasetID" gorm:"column:dataset_id;type:int;comment:数据集ID"`
	ModelID     int       `json:"modelID" gorm:"column:model_id;type:int;comment:模型ID"`
	AlgoID      int       `json:"algoID" gorm:"column:algo_id;type:int;comment:算法ID"`
	Name        string    `json:"name" gorm:"column:name;type:text;comment:任务名称"`
	Description string    `json:"description" gorm:"column:description;type:text;comment:任务描述"`
	InstanceID  int       `json:"instanceID" gorm:"column:instance_id;type:int;comment:任务所在实例资源id"`
	StartTime   time.Time `json:"startTime" gorm:"column:start_time;default:null;comment:开始时间"`
	EndTime     time.Time `json:"endTime" gorm:"column:end_time;default:null;comment:结束时间"`
	RenterID    int       `json:"renterID" gorm:"column:renter_id;type:int;comment:租户ID"`
	Type        int       `json:"type" gorm:"column:type;type:int;comment:任务类型(1.推理、2.训练、3.其他)"`
}

func (RenterTask) TableName() string {
	return "yun_renter_task"
}
