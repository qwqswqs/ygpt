package product

import (
	"time"
)

// productBundling表 结构体  ProductBundling
type ProductBinding struct {
	Id           *int       `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;size:19;"`                              //主键
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:255;"`                    //创建者
	CreateUserId *int       `json:"createUserId" form:"createUserId" gorm:"column:create_user_id;comment:创建者id;size:19;"`      //创建者id
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`                      //创建时间
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新者;size:255;"`                    //更新者
	UpdateUserId *int       `json:"updateUserId" form:"updateUserId" gorm:"column:update_user_id;comment:更新者id;size:19;"`      //更新者id
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;"`                      //更新时间
	IsDeleted    *bool      `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识（0=未删除*，1=已删除）;"`           //删除标识（0=未删除*，1=已删除）
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`                            //备注
	NodeId       *int       `json:"nodeId" form:"nodeId" gorm:"column:node_id;comment:节点id;size:19;"`                          //节点id
	ProductId    *int       `json:"productId" form:"productId" gorm:"column:product_id;comment:算力产品id;size:19;"`               //算力产品id
	ElementId    *int       `json:"elementId" form:"elementId" gorm:"column:element_id;comment:ai产品id(-2.全部数据集 -1.全部模型 0.全部ai产品);"` //ai产品id
	Discount     *int       `json:"discount" form:"discount" gorm:"column:discount;default:null;comment:优惠百分比;size:19;"`                    //优惠百分比
	Reduction    *float64    `json:"reduction" form:"reduction" gorm:"column:reduction;type:decimal(12, 2);default:null;comment:优惠金额;size:19;"`                    //优惠金额
}

func (ProductBinding) TableName() string {
	return "yun_productBinding"
}
