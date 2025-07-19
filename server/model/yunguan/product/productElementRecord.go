package product

import (
	"time"
	"ygpt/server/global"
)

type ProductElementRecord struct {
	global.GVA_MODEL
	ElementID int       `json:"elementID" gorm:"column:element_id;type:bigint;comment:'关联元素ID'"`
	RenterID  int       `json:"renterID" gorm:"column:renter_id;type:bigint;comment:'租户ID'"`
	UseTime   time.Time `json:"useTime" gorm:"column:use_time;type:datetime(6);comment:'使用时间戳'"`
	EndTime   time.Time `json:"endTime" gorm:"column:end_time;type:datetime(6);comment:'结束时间'"`
}

func (e *ProductElementRecord) TableName() string {
	return "yun_product_element_record"
}
