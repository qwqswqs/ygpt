package renter

import (
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type Renter struct {
	global.GVA_MODEL
	RenterID      int    `json:"renterID" gorm:"column:renter_id;type:int;comment:租户平台ID"`
	Username      string `json:"username" gorm:"column:username;type:text;comment:用户名"`
	Email         string `json:"email" gorm:"comment:邮箱"`
	CompanyName   string `json:"companyName" gorm:"column:company_name;type:text;comment:'公司名'"`
	Source        int    `json:"source" gorm:"column:source;type:tinyint;comment:用户来源(1.本地注册 2.算力注册)"`
	Status        int    `json:"status" gorm:"column:status;type:tinyint;comment:状态(1.正常 2.过期)"`
	ReceiveNotice int    `json:"receiveNotice" gorm:"comment:是否接受报警信息;default:7"`
}

func (Renter) TableName() string {
	return "yun_renter"
}

type GetRenterListReq struct {
	request.PageInfo
	Username    string `json:"username"`
	CompanyName string `json:"companyName"`
	Source      int    `json:"source"`
	Status      int    `json:"status"`
	TimeDesc    string `json:"timeDesc"`
}
