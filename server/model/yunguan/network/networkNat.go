package network

import (
	"gorm.io/gorm"
	"ygpt/server/model/common/request"
)

type NetworkNat struct {
	gorm.Model
	Name        string `gorm:"column:name;type:text;comment:'NAT 名称'" json:"name"`
	Description string `gorm:"column:description;type:text;comment:'说明'" json:"description"`
	UseType     uint   `gorm:"column:useType;type:tinyint;comment:'用途类别(1.vpc 2.子网 3.空闲)'" json:"useType"`
	BelongName  string `gorm:"column:belong_name;type:text;comment:'所属VPC或者子网名称'" json:"belongName"`
	NatRule     string `gorm:"column:nat_rule;type:text;comment:'nat规则json'" json:"natRule"`
	InnerIpID   int    `gorm:"column:inner_ip_id;type:int;comment:'内网ip地址id'" json:"innerIpID"`
	OuterIpID   string `gorm:"column:outer_ip_id;type:text;comment:'外网ip地址'" json:"outerIpID"`
	DeviceID    int    `gorm:"column:device_id;type:bigint;comment:'所在设备id'" json:"deviceID"`
}

func (o *NetworkNat) TableName() string {
	return "yun_network_nat"
}

type NetworkNatReq struct {
	gorm.Model
	request.PageInfo
	Name string `json:"name"`
}
