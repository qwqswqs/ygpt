package system

import "time"

type ProductSysType struct {
	Id           *int       `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;size:19;"`                              //主键
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:255;"`                    //创建者
	CreateUserId *int       `json:"createUserId" form:"createUserId" gorm:"column:create_user_id;comment:创建者id;size:19;"`      //创建者id
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`                      //创建时间
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新者;size:255;"`                    //更新者
	UpdateUserId *int       `json:"updateUserId" form:"updateUserId" gorm:"column:update_user_id;comment:更新者id;size:19;"`      //更新者id
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;"`                      //更新时间
	IsDeleted    *bool      `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;default:0;comment:删除标识（0=未删除*，1=已删除）;"` //删除标识（0=未删除*，1=已删除）
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`                            //备注
	Code         string     `json:"code" form:"code" gorm:"column:code;comment:编号;size:50;"`                                   //编号
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:名称;size:255;"`                                  //名称
	Category     *int8      `json:"category" form:"category" gorm:"column:category;default:0;comment:类别;"`                     //类别
	ParentId     *int       `json:"parent_id" form:"parent_id" gorm:"column:parent_id;type:bigint;comment:父id"`
	Description  string     `json:"description" form:"description" gorm:"column:description;comment:说明;"` //说明
}

// TableName sysType表 SysType自定义表名 sys_type
func (ProductSysType) TableName() string {
	return "yun_product_type"
}

type SysTypeListResponse struct {
	Root   ProductSysType   `json:"root" form:"root"`
	Second []ProductSysType `json:"second" form:"second"`
	Third  []ProductSysType `json:"third" form:"third"`
}
