package backup

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type BackupImage struct {
	global.GVA_MODEL
	System      string `json:"system" gorm:"column:system;type:text;comment:'系统'"`
	Frame       string `json:"frame" gorm:"column:frame;type:text;comment:'框架'"`
	Language    string `json:"language" gorm:"column:language;type:text;comment:'语言'"`
	Name        string `json:"name" gorm:"column:name;type:text;comment:'镜像名称'"`
	IsPublic    uint   `json:"isPublic" gorm:"column:is_public;type:tinyint;comment:'是否公开(1.是 2.否)'"`
	Size        uint32 `json:"size" gorm:"column:size;type:bigint;comment:'大小'"`
	Description string `json:"description" gorm:"column:description;type:text;comment:'说明'"`
	GenerateWay uint   `json:"generateWay" gorm:"column:generate_way;type:tinyint;comment:'生成方式(1.上传 2.实例)'"`
	FileName    string `json:"fileName" gorm:"column:file_name;type:text;comment:'文件名'"`
	SourceType  int    `json:"sourceType" gorm:" column:source_type;type:tinyint;comment:'来源类别(1.用户 2.租户 3.实例)'"`
	SourceID    int    `json:"sourceID" gorm:"column:source_id;type:bigint;comment:'来源ID(用户,租户ID或实例ID)'"`
}

func (e *BackupImage) TableName() string {
	return "yun_backup_image"
}

type GetUserImageList struct {
	request.PageInfo
}
type GetImageListReq struct {
	request.PageInfo
	SourceType int `json:"sourceType"`
	SourceID   int `json:"sourceID"`
}
