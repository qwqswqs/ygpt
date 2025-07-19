package product

import (
	"ygpt/server/global"
)

type ProductCategory struct {
	global.GVA_MODEL
	Key   string `json:"key" gorm:"column:key;type:varchar(255);comment:键"`
	Value string `json:"value" gorm:"column:value;type:varchar(255);comment:值"`
}

func (e *ProductCategory) TableName() string {
	return "yun_product_category"
}
