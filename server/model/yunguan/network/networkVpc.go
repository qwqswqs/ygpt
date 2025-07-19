package network

import (
	"gorm.io/gorm"
	"ygpt/server/model/common/request"
)

type NetworkVpc struct {
	gorm.Model
	Name           string `gorm:"column:name;type:text;comment:'名称'" json:"name"`
	Zone           string `gorm:"column:zone;type:text;comment:'区域'" json:"zone"`
	IpID           int    `gorm:"column:ip_id;type:bigint;comment:'IP地址id'" json:"ipID"`
	AccessRule     string `gorm:"column:access_rule;type:json;comment:'访问规则json'" json:"accessRule"`
	RouterRule     string `gorm:"column:router_rule;type:json;comment:'路由规则json'" json:"routerRule"`
	SafeRule       string `gorm:"column:safe_rule;type:json;comment:'安全规则json'" json:"safeRule"`
	Manager        string `gorm:"column:manager;type:text;comment:'管理人员'" json:"manager"`
	BelongLocation string `gorm:"column:belong_location;type:text;comment:'归属单位'" json:"belongLocation"`
	Description    string `gorm:"column:description;type:text;comment:'说明'" json:"description"`
}

func (o *NetworkVpc) TableName() string {
	return "yun_network_vpc"
}

type NetworkVpcReq struct {
	gorm.Model
	request.PageInfo
	Name   string `json:"name"`
	Status int    `json:"status"`
}
