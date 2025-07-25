// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compute

import (
	"yunion.io/x/cloudmux/pkg/apis/compute"

	"yunion.io/x/onecloud/pkg/apis"
	"yunion.io/x/onecloud/pkg/cloudcommon/types"
)

type TAddressType string

const (
	AddressTypeIPv4 = TAddressType("ipv4")
	AddressTypeIPv6 = TAddressType("ipv6")
)

const (
	NETWORK_TYPE_VPC     = compute.NETWORK_TYPE_VPC
	NETWORK_TYPE_CLASSIC = compute.NETWORK_TYPE_CLASSIC
)

type WireResourceInput struct {
	// 二层网络(ID或Name)的资源
	WireId string `json:"wire_id"`
	// swagger:ignore
	// Deprecated
	// fitler by wire id
	Wire string `json:"wire" yunion-deprecated-by:"wire_id"`
}

type WireFilterListBase struct {
	WireResourceInput

	// 以二层网络名称排序
	OrderByWire string `json:"order_by_wire"`
}

type WireFilterListInput struct {
	VpcFilterListInput
	ZonalFilterListBase

	WireFilterListBase
}

type NetworkResourceInput struct {
	// IP子网（ID或Name）
	NetworkId string `json:"network_id"`
	// swagger:ignore
	// Deprecated
	// filter by networkId
	Network string `json:"network" yunion-deprecated-by:"network_id"`
}

type NetworkFilterListBase struct {
	NetworkResourceInput

	// 以IP子网的名称排序
	OrderByNetwork string `json:"order_by_network"`
}

type NetworkFilterListInput struct {
	WireFilterListInput
	NetworkFilterListBase
}

type NetworkIpMacListInput struct {
	apis.StandaloneAnonResourceListInput

	NetworkId string   `json:"network_id"`
	MacAddr   []string `json:"mac_addr"`
	IpAddr    []string `json:"ip_addr"`
}

type NetworkListInput struct {
	apis.SharableVirtualResourceListInput
	apis.ExternalizedResourceBaseListInput
	SchedtagResourceInput
	WireFilterListInput

	HostResourceInput
	StorageResourceInput

	UsableResourceListInput

	// filter by route table which associate with it
	RouteTableId string `json:"route_table_id"`

	// description: Exact matching ip address in network.
	// example: 10.168.222.1
	Ip []string `json:"ip"`

	// description: Fuzzy matching ip address in network.
	// example: 10.168.222.1
	IpMatch []string `json:"ip_match"`

	IfnameHint []string `json:"ifname_hint"`
	// 起始IP地址
	GuestIpStart []string `json:"guest_ip_start"`
	// 接收IP地址
	GuestIpEnd []string `json:"guest_ip_end"`
	// 掩码
	GuestIpMask []int8 `json:"guest_ip_mask"`
	// 网关地址
	GuestGateway string `json:"guest_gateway"`
	// DNS
	GuestDns []string `json:"guest_dns"`
	// allow multiple dhcp, seperated by ","
	GuestDhcp []string `json:"guest_dhcp"`
	// NTP
	GuestNtp []string `json:"guest_ntp"`

	GuestDomain []string `json:"guest_domain"`

	GuestIp6Start []string `json:"guest_ip6_start"`
	GuestIp6End   []string `json:"guest_ip6_end"`
	GuestIp6Mask  []int8   `json:"guest_ip6_mask"`
	GuestGateway6 []string `json:"guest_gateway6"`
	GuestDns6     []string `json:"guest_dns6"`

	GuestDomain6 []string `json:"guest_domain6"`
	// vlanId 1~4096
	VlanId []int `json:"vlan_id"`
	// 服务器类型
	// example: server
	ServerType []string `json:"server_type"`
	// 分配策略
	AllocPolicy []string `json:"alloc_policy"`
	// 是否加入自动分配地址池
	IsAutoAlloc *bool `json:"is_auto_alloc"`
	// 是否为基础网络（underlay）
	IsClassic *bool `json:"is_classic"`

	// filter by Host schedtag
	HostSchedtagId string `json:"host_schedtag_id"`

	// filter by BGP types
	BgpType []string `json:"bgp_type"`

	HostType string `json:"host_type"`

	// 按起始ip地址排序
	// pattern:asc|desc
	OrderByIpStart string `json:"order_by_ip_start"`
	// 按终止ip地址排序
	// pattern:asc|desc
	OrderByIpEnd string `json:"order_by_ip_end"`
}

type NetworkResourceInfoBase struct {
	// IP子网名称
	Network string `json:"network"`
}

type NetworkResourceInfo struct {
	NetworkResourceInfoBase

	// 二层网络ID
	WireId string `json:"wire_id"`

	WireResourceInfo
}

type NetworkIpMacCreateInput struct {
	apis.StandaloneAnonResourceCreateInput

	NetworkId string `json:"network_id"`
	MacAddr   string `json:"mac_addr"`
	IpAddr    string `json:"ip_addr"`
}

type NetworkCreateInput struct {
	apis.SharableVirtualResourceCreateInput

	// description: ip range of guest, if not set, you shoud set guest_ip_start,guest_ip_end and guest_ip_mask params
	// example: 10.168.222.1/24
	GuestIpPrefix string `json:"guest_ip_prefix"`

	// description: ip range of guest ip start, if set guest_ip_prefix, this parameter will be ignored
	// example: 10.168.222.1
	GuestIpStart string `json:"guest_ip_start"`

	// description: ip range of guest ip end, if set guest_ip_prefix, this parameter will be ignored
	// example: 10.168.222.100
	GuestIpEnd string `json:"guest_ip_end"`

	// description: ip range of guest ip mask, if set guest_ip_prefix, this parameter will be ignored
	// example: 24
	// maximum: 30
	// minimum: 12
	GuestIpMask int8 `json:"guest_ip_mask"`

	IfnameHint string `json:"ifname_hint"`

	// description: guest gateway
	// example: 192.168.222.1
	GuestGateway string `json:"guest_gateway"`

	// description: ipv6 range of guest, if not set, you shoud set guest_ip6_start,guest_ip6_end and guest_ip6_mask params
	// example: 3ffe:3200:2001:2300::/64
	GuestIp6Prefix string `json:"guest_ip6_prefix"`
	// description: ipv6 range of guest ip start, if set guest_ip6_prefix, this parameter will be ignored
	// example: 3ffe:3200:2001:2300::1
	GuestIp6Start string `json:"guest_ip6_start"`
	// description: ipv6 range of guest ip end, if set guest_ip6_prefix, this parameter will be ignored
	// example: 3ffe:3200:2001:2300:ffff:ffff:ffff:ffff
	GuestIp6End string `json:"guest_ip6_end"`
	// description: ipv6 range of guest ip mask, if set guest_ip6_prefix, this parameter will be ignored
	// example: 64
	// maximum: 126
	// minimum: 48
	GuestIp6Mask uint8 `json:"guest_ip6_mask"`

	// description: guest gateway of IPv6
	// example: 3ffe:3200:2001:2300::1
	GuestGateway6 string `json:"guest_gateway6"`

	// description: guest dns
	// example: 114.114.114.114,8.8.8.8
	GuestDns string `json:"guest_dns"`

	// description: guest dhcp
	// example: 192.168.222.1,192.168.222.4
	GuestDHCP string `json:"guest_dhcp"`

	// description: guest ntp
	// example: cn.pool.ntp.org,0.cn.pool.ntp.org
	GuestNtp string `json:"guest_ntp"`

	// swagger:ignore
	WireId string `json:"wire_id"`

	// description: wire id or name
	Wire string `json:"wire"`

	// description: zone id or name
	Zone string `json:"zone"`

	// description: vpc id or name
	Vpc string `json:"vpc"`

	// description: server type
	// enum: ["guest","baremetal","pxe","ipmi","hostlocal"]
	// default: guest
	ServerType TNetworkType `json:"server_type"`

	// 是否加入自动分配地址池
	IsAutoAlloc *bool `json:"is_auto_alloc"`

	// VlanId
	VlanId *int `json:"vlan_id"`

	// deprecated
	Vlan *int `json:"vlan" yunion-deprecated-by:"vlan_id"`

	// 线路类型
	BgpType string `json:"bgp_type"`

	// 是否申请ip时自动挂载公网ip，仅对aws生效
	AssignPublicIp bool `json:"assign_public_ip"`
}

type SNetworkNics struct {
	// 云主机网卡数量
	Vnics int `json:"vnics"`
	// IPv4地址数量
	Vnics4 int `json:"vnics4"`
	// IPv6地址数量
	Vnics6 int `json:"vnics6"`
	// 裸金属网卡数量
	BmVnics int `json:"bm_vnics"`
	// 负载均衡网卡数量
	LbVnics int `json:"lb_vnics"`
	// 浮动Ip网卡数量
	EipVnics int `json:"eip_vnics"`
	// VIP数量
	GroupVnics int `json:"group_vnics"`

	// 预留IPv4数量
	ReserveVnics4 int `json:"reserve_vnics4"`
	// 预留IPv6数量
	ReserveVnics6 int `json:"reserve_vnics6"`

	// RDS网卡数量
	RdsVnics int `json:"rds_vnics"`
	// NAT网关网卡数量
	NatVnics      int `json:"nat_vnics"`
	BmReusedVnics int `json:"bm_reused_vnics"`

	// 弹性网卡数量
	NetworkinterfaceVnics int `json:"networkinterface_vnics"`

	// 已使用端口数量
	PortsUsed int `json:"ports_used"`

	Total int `json:"total"`

	// 已使用IPv4端口数量
	Ports6Used int `json:"ports6_used"`

	Total6 int `json:"total6"`
}

func (self *SNetworkNics) SumTotal() {
	self.Total = self.Vnics4 +
		self.BmVnics +
		self.LbVnics +
		self.LbVnics +
		self.EipVnics +
		self.GroupVnics +
		self.ReserveVnics4 +
		self.RdsVnics +
		self.NetworkinterfaceVnics +
		self.NatVnics -
		self.BmReusedVnics
	self.Total6 = self.Vnics6 +
		self.ReserveVnics6
	self.PortsUsed = self.Total
	self.Ports6Used = self.Total6
}

type NetworkDetails struct {
	apis.SharableVirtualResourceDetails
	WireResourceInfo

	SNetwork
	SNetworkNics

	// 是否是内网
	Exit bool `json:"exit"`
	// 端口数量
	Ports int `json:"ports"`

	Dns string `json:"dns"`

	// 路由信息
	Routes    []types.SRoute             `json:"routes"`
	Schedtags []SchedtagShortDescDetails `json:"schedtags"`

	IsClassic bool `json:"is_classic"`

	AdditionalWires []SSimpleWire `json:"additional_wires"`
}

type NetworkIpMacDetails struct {
	apis.StandaloneAnonResourceDetails

	NetworkId string `json:"network_id"`
	IpAddr    string `json:"ip_addr"`
	MacAddr   string `json:"mac_addr"`
}

type NetworkReserveIpInput struct {
	apis.Meta

	// description: reserved ip list
	// required: true
	// example: [10.168.222.131, 10.168.222.134]
	Ips []string `json:"ips"`

	// description: the comment
	// example: reserve ip for test
	Notes  string `json:"notes"`
	Status string `json:"status"`
	// description: The reserved cycle
	// required: false
	Duration string `json:"duration"`
}

type NetworkReleaseReservedIpInput struct {
	apis.Meta

	// description: IP to be released
	// required: true
	// example: 10.168.222.121
	Ip string `json:"ip"`
}

type NetworkPurgeInput struct {
	apis.Meta
}

type NetworkMergeInput struct {
	apis.Meta

	// description: network id or name to be merged
	// required: true
	// example: test-network
	Target string `json:"target"`
}

type NetworkSplitInput struct {
	apis.Meta

	// description: The middle - separated IP must belong to the network
	// required: true
	// example: 10.168.222.181
	SplitIp string `json:"split_ip"`

	// description: another network name after split
	// required: false
	Name string `json:"name"`
}

type NetworkTryCreateNetworkInput struct {
	apis.Meta

	Ip          string `json:"ip"`
	Mask        int    `json:"mask"`
	IsOnPremise bool   `json:"is_on_premise"`

	ServerType TNetworkType `json:"server_type"`
}

type NetworkSyncInput struct {
	apis.Meta
}

type NetworkIpMacUpdateInput struct {
	apis.StandaloneAnonResourceBaseUpdateInput

	MacAddr string `json:"mac_addr"`
	IpAddr  string `json:"ip_addr"`
}

type NetworkUpdateInput struct {
	apis.SharableVirtualResourceBaseUpdateInput

	// 起始IP地址
	GuestIpStart string `json:"guest_ip_start"`
	// 结束IP地址
	GuestIpEnd string `json:"guest_ip_end"`
	// 掩码
	GuestIpMask int8 `json:"guest_ip_mask"`
	// 网关地址
	GuestGateway *string `json:"guest_gateway"`
	// DNS
	GuestDns *string `json:"guest_dns"`
	// allow multiple dhcp, seperated by ","
	GuestDhcp *string `json:"guest_dhcp"`
	// NTP
	GuestNtp *string `json:"guest_ntp"`

	// 起始IP6地址
	GuestIp6Start *string `json:"guest_ip6_start"`
	// 结束IP6地址
	GuestIp6End *string `json:"guest_ip6_end"`
	// IP6子网掩码长度
	GuestIp6Mask *uint8 `json:"guest_ip6_mask"`
	// IP6网关地址
	GuestGateway6 *string `json:"guest_gateway6"`

	GuestDomain *string `json:"guest_domain"`

	VlanId *int `json:"vlan_id"`

	// 分配策略
	AllocPolicy string `json:"alloc_policy"`

	// 是否加入自动分配地址池
	IsAutoAlloc *bool `json:"is_auto_alloc"`

	// 更新网络类型
	ServerType TNetworkType `json:"server_type"`
}

type GetNetworkAddressesInput struct {
	// 获取资源的范围，例如 project|domain|system
	Scope string `json:"scope"`
}

type GetNetworkAddressesOutput struct {
	// IP子网地址记录
	Addresses []SNetworkUsedAddress `json:"addresses"`

	// IPv6子网地址记录
	Addresses6 []SNetworkUsedAddress `json:"addresses6"`
}

type GetNetworkAvailableAddressesInput struct {
}

type GetNetworkAvailableAddressesOutput struct {
	// IP子网地址记录
	Addresses []string `json:"addresses"`

	// IPv6子网地址记录
	Addresses6 []string `json:"addresses6"`
}

type NetworkSetBgpTypeInput struct {
	apis.Meta

	// description: new BgpType name
	// required: true
	// example: ChinaTelecom, BGP, etc.
	BgpType string `json:"bgp_type"`
}

type NetworkIpMacBatchCreateInput struct {
	NetworkId string            `json:"network_id"`
	IpMac     map[string]string `json:"ip_mac"`
}

type NetworkSwitchWireInput struct {
	apis.Meta

	// description: new wire Id or name
	// required: true
	// example: bcast0
	WireId string `json:"wire_id"`
}

type NetworSyncAdditionalWiresInput struct {
	apis.Meta

	WireIds []string `json:"wire_ids"`
}
