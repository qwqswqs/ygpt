package network

import (
	"gorm.io/gorm"
)

type NetworkIp struct {
	gorm.Model
	Name       string `gorm:"column:name;type:text;comment:'地址段名称'" json:"name"`
	StartIp    string `gorm:"column:start_ip;type:text;comment:'开始地址'" json:"startIp"`
	EndIp      string `gorm:"column:end_ip;type:text;comment:'结束地址'" json:"endIp"`
	GatewayIp  string `gorm:"column:gateway_ip;type:text;comment:'网关地址'" json:"gatewayIp"`
	MaskBits   int    `gorm:"column:mask_bits;type:int;comment:'掩码位数'" json:"maskBits"`
	IpNum      int    `gorm:"column:ip_num;type:int;comment:'地址个数'" json:"ipNum"`
	BelongType int    `gorm:"column:belong_type;type:int;comment:'归属类别(1.vpc内 2.子网内 3.vpc外部 4.子网外)'" json:"belongType"`
	BelongID   int    `gorm:"column:belong_id;type:int;comment:'归属ID'" json:"belongID"`
}

func (o *NetworkIp) TableName() string {
	return "yun_network_ip"
}
