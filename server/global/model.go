package global

import (
	"time"

	"gorm.io/gorm"
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
type MinioData struct {
	Endpoint   string `json:"endpoint" `
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	BucketName string `json:"bucketName"`
}
type APP_MODEL struct {
	Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:主键" json:"id"  form:"id"`
	CreateBy     string    `gorm:"column:create_by;type:varchar(255);default:null;comment:创建者" json:"createBy" form:"createBy"`
	CreateUserId int       `gorm:"column:create_user_id;type:bigint;default:null;comment:创建者id" json:"createUserId" form:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime(3);autoCreateTime:milli;comment:创建时间" json:"createTime" form:"createTime"`
	UpdateBy     string    `gorm:"column:update_by;type:varchar(255);default:null;comment:更新者" json:"updateBy" form:"updateBy"`
	UpdateUserId int       `gorm:"column:update_user_id;type:bigint;default:null;comment:更新者id" json:"updateUserId" form:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime(3);autoUpdateTime:milli;comment:更新时间" json:"updateTime" form:"updateTime"`
	IsDeleted    int       `gorm:"column:is_deleted;type:tinyint;default:0;comment:删除标识（0=未删除*，1=已删除）" json:"isDeleted" form:"isDeleted"`
	Remark       string    `gorm:"column:remark;type:varchar(255);default:null;comment:备注" json:"remark" form:"remark"`
}
