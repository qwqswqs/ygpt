package system

import (
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type SystemTicket struct {
	global.GVA_MODEL
	Type               int       `json:"type" gorm:"column:type;type:int;comment:类别(1.系统 2.平台)"`
	TicketType         int       `json:"ticketType" gorm:"column:ticket_type;type:int;comment:工单类别(1.标准产品 2.定制产品 3.投诉与建议 4.合同与发票(报价))"`
	IsQuotation        int       `json:"isQuotation" gorm:"column:is_quotation;type:tinyint;comment:是否需要报价(0.否 1.是)"`
	IsThird            int       `json:"isThird" gorm:"column:is_third;type:tinyint;comment:是否为第三方(0.否 1.是)"`
	RenterID           int       `json:"renterID" gorm:"column:renter_id;type:bigint;comment:所属租户ID"`
	Description        string    `json:"description" gorm:"column:description;type:text;comment:工单描述json"`
	ProblemDescription string    `json:"problemDescription" gorm:"column:problem_description;type:text;comment:故障问题描述json"`
	CreatePerson       int       `json:"createPerson" gorm:"column:create_person;type:int;comment:创建用户ID(0表示系统)"`
	AssignPerson       int       `json:"assignPerson" gorm:"column:assign_person;type:int;comment:分配用户ID"`
	HandlePerson       int       `json:"handlePerson" gorm:"column:handle_person;type:int;comment:处理用户ID"`
	HandelCondition    int       `json:"handelCondition" gorm:"column:handel_condition;type:int;comment:完成情况"`
	HandelTime         time.Time `json:"handelTime" gorm:"column:handel_time;default:null;comment:处理完成时间"`
	AssignTime         time.Time `json:"assignTime" gorm:"column:assign_time;default:null;comment:分配时间"`
	Status             int       `json:"status" gorm:"column:status;type:int;comment:状态(1.待处理 2.已处理)"`
	OrderId            int       `json:"orderId" gorm:"column:order_id;type:bigint;default:null;comment:平台订单ID"`
	Quotation          *float64  `json:"quotation" gorm:"column:quotation;type:decimal(12, 2);default:null;comment:报价"`
	InvoiceJson        string    `json:"invoiceJson" gorm:"column:invoice_json;type:json;default:null;comment:发票json"`
	ContractJson       string    `json:"contractJson" gorm:"column:contract_json;type:json;default:null;comment:合同json"`
}

func (SystemTicket) TableName() string {
	return "yun_system_ticket"
}

type GetTicketReq struct {
	request.PageInfo
	TicketType int    `json:"ticketType"`
	Status     int    `json:"status"`
	TimeDesc   string `json:"timeDesc"`
}

type UpdateAndUploadReq struct {
	global.GVA_MODEL
	Type            int       `json:"type"`
	TicketType      int       `json:"ticketType"`
	IsQuotation     int       `json:"isQuotation"`
	IsThird         int       `json:"isThird"`
	RenterID        int       `json:"renterID"`
	Description     string    `json:"description"`
	CreatePerson    int       `json:"createPerson"`
	AssignPerson    int       `json:"assignPerson" `
	HandlePerson    int       `json:"handlePerson"`
	HandelCondition int       `json:"handelCondition"`
	HandelTime      time.Time `json:"handelTime"`
	AssignTime      time.Time `json:"assignTime" `
	Status          int       `json:"status"`
	OrderId         int       `json:"orderId"`
	Quotation       *float64  `json:"quotation"`
	Contract        Contract  `json:"contractJson"`
	Invoice         Invoice   `json:"invoiceJson"`
}
type Contract struct {
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	SignedTime   time.Time `json:"signedTime"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	File         string    `json:"file"`
	Amount       float64   `json:"amount"`
	CustomerName string    `json:"customerName"`
	SupplierName string    `json:"supplierName"`
}
type Invoice struct {
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	SignedTime time.Time `json:"signedTime"`
	StartTime  time.Time `json:"startTime"`
	//EndTime      string  `json:"endTime"`
	File         string  `json:"file"`
	Amount       float64 `json:"amount"`
	CustomerName string  `json:"customerName"`
	SupplierName string  `json:"supplierName"`
	TaxRate      float64 `json:"taxRate"`
	TaxAmount    float64 `json:"taxAmount"`
}
