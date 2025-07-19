package network

import (
	"gorm.io/gorm"
	"time"
)

type NetworkSub struct {
	gorm.Model
	Name        string    `gorm:"column:name;type:text;comment:'地址段名称'" json:"name"`
	InnerIpID   int       `gorm:"column:inner_ip_id;type:int;comment:'内网ip地址id'" json:"innerIpID"`
	GatewayIp   string    `gorm:"column:gateway_ip;type:text;comment:'网关地址'" json:"gatewayIp"`
	Dns         string    `gorm:"column:dns;type:text;comment:'DNS地址'" json:"dns"`
	Description string    `gorm:"column:description;type:text;comment:'说明'" json:"description"`
	Header      string    `gorm:"column:header;type:text;comment:'负责人'" json:"header"`
	Status      int       `gorm:"column:status;type:int;comment:'当前状态(1.空闲 2.运行)'" json:"status"`
	RenterID    int       `gorm:"column:renter_id;type:int;comment:'租户id'" json:"renterID"`
	RentTime    time.Time `gorm:"column:rent_time;type:datetime(6);comment:'租用时间'" json:"rentTime"`
	BelongVpc   int       `gorm:"column:belong_vpc;type:int;comment:'所属VPC(独立的为0)'" json:"belongVpc"`
}

func (o *NetworkSub) TableName() string {
	return "yun_network_sub"
}
