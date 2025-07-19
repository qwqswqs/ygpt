package product

import (
	"time"
)

// infoRegion表 结构体  InfoRegion
type InfoRegion struct {
	Id           *int       `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;size:19;"`                              //主键
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:255;"`                    //创建者
	CreateUserId *int       `json:"createUserId" form:"createUserId" gorm:"column:create_user_id;comment:创建者id;size:19;"`      //创建者id
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间;"`       //创建时间
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新者;size:255;"`                    //更新者
	UpdateUserId *int       `json:"updateUserId" form:"updateUserId" gorm:"column:update_user_id;comment:创建者id;size:19;"`      //创建者id
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间;"`       //更新时间
	IsDeleted    int        `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;default:0;comment:删除标识（0=未删除*，1=已删除）;"` //删除标识（0=未删除*，1=已删除）
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`                            //备注
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:名称;size:191;"`                                  //名称
	Code         string     `json:"code" form:"code" gorm:"column:code;comment:编码;size:191;"`                                  //编码
	Type         *int       `json:"type" form:"type" gorm:"column:type;comment:类型（1=国，2=省，3=市，4=区县，5=方向n区）;"`                  //类型（1=国，2=省，3=市，4=区县，5=方向n区）
	ParentId     *int       `json:"parentId"  form:"parentId"  gorm:"column:parent_id;type:bigint;comment:父id"`
	//TypeId       *int       `json:"typeId" form:"typeId" gorm:"column:type_id;comment:关联id;size:19;"`                          //关联id
}

// TableName infoRegion表 InfoRegion自定义表名 info_region
func (InfoRegion) TableName() string {
	return "yun_product_region"
}
