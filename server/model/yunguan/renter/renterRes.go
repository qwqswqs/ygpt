package renter

import (
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type RenterRes struct {
	global.GVA_MODEL
	RenterID         int    `json:"renterID" gorm:"column:renter_id;type:bigint;comment:实际使用了Renter表的renter_id(等同sys_user的ID)"`
	Type             int    `json:"type" gorm:"column:type;type:tinyint;comment:'资源类别:1.裸金属2.云主机3.容器'"`
	UseTime          int64  `json:"useTime" gorm:"column:use_time;type:bigint;comment:'使用时长'"`
	Description      string `json:"description" gorm:"column:description;type:text;comment:'选购资源描述json'"`
	PrivateIp        string `json:"privateIp" gorm:"column:private_ip;type:varchar(255);comment:'IPv4内网地址'"`
	Url              string `json:"url" gorm:"type:text;comment:'访问地址'"`
	SshHost          string `json:"sshHost" gorm:"column:ssh_host;type:varchar(255);comment:'ssh地址'"`
	SshPort          int    `json:"sshPort" gorm:"column:ssh_port;type:int;comment:'ssh端口'"`
	SshUser          string `json:"sshUser" gorm:"column:ssh_user;type:varchar(255);comment:'ssh用户名'"`
	SshPwd           string `json:"sshPwd" gorm:"column:ssh_pwd;type:text;comment:'ssh密码(空表示和算力圈密码一致)'"`
	ResourceID       string `json:"resourceID" gorm:"column:resource_id;type:text;comment:'cloudpods资源ID'"`
	OrderId          int    `json:"orderId" gorm:"column:order_id;type:bigint;comment:'算力平台订单ID'"`
	OrderComputingId int    `json:"orderComputingId" gorm:"column:order_computing_id;type:bigint;comment:'算力平台算力订单ID'"`
	TicketId         int    `json:"ticketId" gorm:"column:ticket_id;type:bigint;comment:'工单ID'"`
	Status           int    `json:"status" gorm:"column:status;type:tinyint;comment:状态(1.待处理 2.执行中 3.执行完 4.已回收)"`
	AiList           string `json:"aiList" form:"aiList" gorm:"column:ai_list;type:json;comment:捆绑ai产品列表"`

	DataSource int       `json:"dataSource" gorm:"comment:数据来源(1.本地 2.算力圈);default:1"`
	StartTime  time.Time `json:"startTime" gorm:"comment:实例开始时间;default:null"`
	EndTime    time.Time `json:"endTime" gorm:"comment:实例释放时间;default:null"`
}

func (RenterRes) TableName() string {
	return "yun_renter_res"
}

type GetRenterResReq struct {
	request.PageInfo
	RenterID int `json:"renterID"`
	Type     int `json:"type"`
}

type GetRenterResReqByResourceID struct {
	ResourceID string `json:"resourceID"`
}

type GetRenterResCountByTicketId struct {
	TicketId int `json:"ticketId"`
}

type RenterResInfo struct {
	ResourceID string `json:"resourceID"` //资源编号
	Name       string `json:"name"`       //资源名称
	Status     string `json:"status"`     //资源状态
	PrivateIp  string `json:"privateIp"`  //内网ip
	Image      string `json:"image"`      //镜像名称
	Username   string `json:"username"`   //用户名
	Password   string `json:"password"`   //密码
}
