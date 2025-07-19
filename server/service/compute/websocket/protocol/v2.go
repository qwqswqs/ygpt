package protocol

import (
	"time"
)

// 基础格式===================================================================

const (
	LoginCmd                  = "login"
	LoginRetCmd               = "loginRet"
	KeepCmd                   = "keep"
	KeepRetCmd                = "keepRet"
	LogoutCmd                 = "logout"
	LogoutRetCmd              = "logoutRet"
	MonitorCmd                = "monitor"
	MonitorRetCmd             = "monitorRet"
	TicketCmd                 = "ticket"
	TicketRetCmd              = "ticketRet"
	OrderTicketCmd            = "orderTicket"
	OrderTicketRetCmd         = "orderTicketRet"
	QuoteCmd                  = "quote"
	QuoteRetCmd               = "quoteRet"
	AllocateCmd               = "allocate"
	AllocateRetCmd            = "allocateRet"
	ReleaseCmd                = "release"
	ReleaseRetCmd             = "releaseRet"
	YgModifyCmd               = "ygModify"
	YgModifyRetCmd            = "ygModifyRet"
	TenantInfoQueryCmd        = "tenantInfoQuery"
	TenantInfoQueryRetCmd     = "tenantInfoQueryRet"
	TenantControlCmd          = "tenantControl"
	TenantControlRetCmd       = "tenantControlRet"
	SoftVersionCmd            = "softVersion"
	SoftVersionRetCmd         = "softVersionRet"
	SoftDownloadCmd           = "softDownload"
	SoftDownloadRetCmd        = "softDownloadRet"
	ProductStatusChangeCmd    = "productStatusChange"
	ProductStatusChangeRetCmd = "productStatusChangeRet"
	PushCmd                   = "push"
	PushRetCmd                = "pushRet"
	OrderPushCmd              = "orderPush"
	OrderPushRetCmd           = "orderPushRet"
	SupplyInfoReleaseCmd      = "supplyInfoRelease"
	SupplyInfoReleaseRetCmd   = "supplyInfoReleaseRet"
	ManageSupplyInfoCmd       = "manageSupplyInfo"
	ManageSupplyInfoRetCmd    = "manageSupplyInfoRet"
	DemandInfoReleaseCmd      = "demandInfoRelease"
	DemandInfoReleaseRetCmd   = "demandInfoReleaseRet"
	ManageDemandInfoCmd       = "manageDemandInfo"
	ManageDemandInfoRetCmd    = "manageDemandInfoRet"
	ElementUploadCmd          = "elementUpload"    //2024 11 26
	ElementUploadRetCmd       = "elementUploadRet" //2024 11 26
	ElementManageCmd          = "elementManage"    //2024 11 26
	ElementManageRetCmd       = "elementManageRet" //2024 11 26
	CopyAiCmd                 = "copyAiProductToCloudResource"
	CopyAiRetCmd              = "copyAiProductToCloudResourceRet"
	TokenReqCmd               = "token"
	TokenRetCmd               = "tokenRet"
	NodeStaticsCmd            = "nodeStatics"
	NodeStaticsRetCmd         = "nodeStaticsRet"
	TenantFinanceCmd          = "tenantFinance"
	TenantFinanceRetCmd       = "tenantFinanceRet"

	// RegionCmd 获取城市列表(从info_region获取数据)
	RegionCmd    = "region"
	RegionRetCmd = "regionRet"
	//三级标题

	ThirdTagCmd             = "thirdTag"
	ThirdTagRetCmd          = "thirdTagRet"
	SuccessCodeCmd          = 0
	SuccessMsgCmd           = "操作成功"
	UpdateCodeCmd           = 1
	UpdateMsgCmd            = "更新成功"
	ErrCodeCmd              = -400
	ErrMsgCmd               = "操作失败"
	NoLoginErrorCode        = 101
	NoLoginErrorMsg         = "节点未登录"
	NoAuthErrorCode         = 102
	NoAuthErrorMsg          = "节点授权码错误"
	NoLoginErrorCode2       = 103
	NoLoginErrorMsg2        = "节点未登录，请先登录"
	LoginBeOccupiedCode     = 104
	LoginBeOccupiedMsg      = "该节点在异地登录，当前设备上的节点已登出算力平台"
	ExistCode               = 200
	ExistMsgCmd             = "产品已存在"
	WrongMsgCode            = 201
	WrongMsgMsg             = "参数格式错误"
	NoParamErrorCode        = 202
	NoParamErrorMsg         = "缺少参数，缺乏要求的信息"
	WrongProductIDErrorCode = 203
	WrongProductIDErrorMsg  = "修改失败，产品编号不存在"
	WrongDeleteErrorCode    = 204
	WrongDeleteErrorMsg     = "删除失败，产品不存在"
	TicketCheckCode         = 300
	TicketCheckMsg          = "成功生成资源分配工单，请等待节点回应"
	AllocateErrorCode       = 301
	AllocateErrorMsg        = "资源分配失败，订单不存在"
	RepeatAllocateErrorCode = 302
	RepeatAllocateErrorMsg  = "资源分配失败，订单要求的资源已分配"
	SpaceCreateErrorCode    = 303
	SpaceCreateErrorMsg     = "创建空间失败，要求的软件不存在"
	WrongRecycleErrorCode   = 304
	WrongRecycleErrorMsg    = "回收失败，资源不存在"
	AllocateSuccessCode     = 305
	AllocateSuccessMsg      = "资源分配成功"
	AllocateFailedCode      = 306
	AllocateFailedMsg       = "资源分配失败"
)

type MessageCmd string

type BaseMsg struct {
	Command MessageCmd `json:"command"`
	ID      int        `json:"ID"`
}

type BaseRetMsg struct {
	BaseMsg
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

// Login 云管登录请求结构体===================================================================
type Login struct {
	BaseMsg
	Payload struct {
		CloudType       int    `json:"cloudType"` // 1表示公有云, 2表示私有云, 3表示混合云
		NodeCode        string `json:"nodeCode"`  //节点编码
		AuthCode        string `json:"authCode"`  // 授权码
		InstanceEnvJson string `json:"instanceEnvJson"`
	} `json:"payload"`
}

// LoginResponse 云管登录响应结构体
type LoginResponse struct {
	BaseRetMsg
	Payload struct {
		NewNodeCode  string `json:"newNodeCode"`  // 新授权码
		LoginCode    string `json:"loginCode"`    // 登录码
		MsgKey       string `json:"msgKey"`       // 消息密钥
		MinioHost    string `json:"minioHost"`    // IP地址
		MinioPort    string `json:"minioPort"`    // 端口号
		AccessKey    string `json:"accessKey"`    // 用户名
		SecretKey    string `json:"secretKey"`    // 密码
		UseSSL       int    `json:"useSSL"`       // 是否开启SSL
		UploadBucket string `json:"uploadBucket"` // 上传桶名称
	} `json:"payload"`
}

// Keep 云管保活请求结构体===================================================================
type Keep struct {
	BaseMsg
	Payload struct {
		NodeCode string `json:"nodeCode"` // 授权码
	} `json:"payload"`
}

// KeepResponse 云管保活响应结构体
type KeepResponse struct {
	BaseRetMsg
	Payload struct {
		NodeCode string `json:"nodeCode"` // 授权码
	} `json:"payload"`
}

// Logout 云管登出请求结构体===================================================================
type Logout struct {
	BaseMsg
	Payload struct {
		NodeCode string `json:"nodeCode"` // 授权码
	} `json:"payload"`
}

// LogoutResponse 云管登出响应结构体
type LogoutResponse struct {
	BaseRetMsg
}

// Monitor 运行监控请求结构体===================================================================
type Monitor struct {
	BaseMsg
	NodeCode string `json:"nodeCode"` // 授权码
}

// MonitorResponse 运行监控响应结构体
type MonitorResponse struct {
	BaseRetMsg
	Payload struct {
		LoginCode         string  `json:"loginCode"`         // 登录码
		VirtualMachineNum int     `json:"virtualMachineNum"` // 云主机数量
		PhysicsMachineNum int     `json:"physicsMachineNum"` // 物理机数量
		CPUPoint          float64 `json:"CPUPoint"`          // CPU使用率
		MemPoint          float64 `json:"memPoint"`          // 内存使用率
		DiskPoint         float64 `json:"diskPoint"`         // 磁盘使用率
		ImageNum          int     `json:"imageNum"`          // 镜像数量
	} `json:"payload"`
}

// Alert 告警处理请求结构体===================================================================
type Alert struct {
	BaseMsg
	Payload struct {
		NodeCode           string `json:"nodeCode"`           // 授权码
		WarningID          string `json:"warningID"`          // 订单ID
		Level              uint8  `json:"level"`              // 告警级别
		WarningDetailsJSON string `json:"warningDetailsJSON"` // 具体告警信息
	} `json:"payload"`
}

// AlertResponse 告警处理响应结构体
type AlertResponse struct {
	BaseRetMsg
	Payload struct {
		Result string `json:"result"` // 处理结果
		Method string `json:"method"` // 处理方法
	} `json:"payload"`
}

// TicketRequest 工单生成请求结构体===================================================================
type TicketRequest struct {
	BaseMsg
	Payload struct {
		TicketsID         int    `json:"ticketsID"`         // 工单编号
		TicketsDetailJSON string `json:"ticketsDetailJSON"` // 工单详细描述
	} `json:"payload"`
}

// TicketResponse 工单生成响应结构体
type TicketResponse struct {
	BaseRetMsg
}

// TicketRet 异步工单处理响应结构体
type TicketRet struct {
	BaseRetMsg
	Payload struct {
		TicketsID int    `json:"ticketsID"` // 工单编号
		Result    string `json:"result"`    // 具体处理结果
	} `json:"payload"`
}

// OrderTicketRequest 订单工单生成请求结构体===================================================================
type OrderTicketRequest struct {
	BaseMsg
	Payload struct {
		OrderID         int    `json:"ticketsID"` // 工单编号
		IsQuote         int    `json:"isQuote"`
		IsThird         int    `json:"isThird"`
		BuyerPhone      string `json:"buyerPhone"`        //凭此可以在云管查到对应租户
		OrderDetailJSON string `json:"ticketsDetailJSON"` // 工单详细描述
	} `json:"payload"`
}

type OrderTicketResponse struct {
	BaseRetMsg
}

// QuoteRequest 报价请求结构体===================================================================
type QuoteRequest struct {
	BaseMsg
	Payload struct {
		YgTicketID int     `json:"ticketID"`
		OrderID    int     `json:"orderID"`    // 订单编号
		QuotePrice float64 `json:"quotePrice"` // 报价
		//此处使用了解析前端请求时用的结构体，字段名称顺序一致，方便使用
		Contract Contract `json:"contract"`
		Invoice  Invoice  `json:"invoice"`
	} `json:"payload"`
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
type QuoteResponse struct {
	BaseRetMsg
}

// PackageUpload 套餐上报请求结构体===================================================================
type PackageUpload struct {
	BaseMsg
	NodeCode string `json:"nodeCode"` // 授权码
	Payload  struct {
		Packages []Package `json:"Packages"`
	} `json:"payload"`
}

type Package struct {
	PackageID string           `json:"packageID"` // 子套餐ID
	Payload   []PackagePayload `json:"payload"`
}

type PackagePayload struct {
	SourceType string        `json:"source_type"` // "container/base/vm/ele"
	Spec       []PackageSpec `json:"spec"`
}

type PackageSpec struct {
	Name  string `json:"name"`  // 规格名称
	Unit  string `json:"unit"`  // 单位
	Data  int    `json:"data"`  // 数量
	Price int    `json:"price"` // 单价
}

// PackageUploadResponse 套餐上报响应结构体
type PackageUploadResponse struct {
	BaseRetMsg
	Payload struct {
		ResourceAuditStatus string `json:"resourceAuditStatus"` // 资源审核状态
	} `json:"payload"`
}

// SupplyInfoReleaseRequest 供给信息发布请求结构体===================================================================
type SupplyInfoReleaseRequest struct {
	BaseMsg
	NodeCode           string `json:"nodeCode"` // 云管授权码
	ProductDetailsJson string `json:"productDetailsJSON"`
}

// SupplyInfoReleaseResponse 供给信息发布响应结构体
type SupplyInfoReleaseResponse struct {
	BaseRetMsg
	SupplyInfoID int `json:"supplyInfoID"`
}

// SupplyInfoManageRequest 供给信息管理请求结构体===================================================================
type SupplyInfoManageRequest struct {
	BaseMsg
	Payload struct {
		NodeCode     string `json:"nodeCode"`
		Method       int8   `json:"method"`       // 1:查询, 2:上架, 3:下架
		SupplyInfoID int    `json:"supplyDemand"` // 可选项，查询时不需要指定
	} `json:"payload"`
}

// SupplyInfoManageResponse 供给信息管理响应结构体
type SupplyInfoManageResponse struct {
	BaseRetMsg
	Method             int8   `json:"method"`       // 1.查询，2.上架，3.下架
	SupplyInfoID       int    `json:"supplyDemand"` // 可选项，查询时不需要指定
	ProductDetailsJson string `json:"productDetailsJSON"`
}

type SupplyInfoPayload struct {
	Status         int8   `json:"status"`
	SupplyInfoID   int    `json:"supplyInfoID"`
	SupplyInfoName string `json:"supplyInfoName"`
	TapTimes       int    `json:"tapTimes"`
}

// DemandInfoReleaseRequest 需求信息发布请求结构体===================================================================
type DemandInfoReleaseRequest struct {
	BaseMsg
	Payload struct {
		NodeCode           string `json:"nodeCode"` // 授权码
		ProductDetailsJson string `json:"productDetailsJSON"`
	} `json:"payload"`
}

// DemandInfoReleaseResponse 需求信息发布响应结构体
type DemandInfoReleaseResponse struct {
	BaseRetMsg
	DemandInfoID int `json:"demandInfoID"`
}

// DemandInfoManageRequest 需求信息管理请求结构体===================================================================
type DemandInfoManageRequest struct {
	BaseMsg
	Payload struct {
		NodeCode     string `json:"nodeCode"`
		Method       int8   `json:"method"`       // 1:查询, 2:上架, 3:下架
		DemandInfoID int    `json:"demandInfoID"` // 可选项，查询时不需要指定
	} `json:"payload"`
}

// DemandInfoManageResponse 需求信息管理响应结构体
type DemandInfoManageResponse struct {
	BaseRetMsg
	Method             int8   `json:"method"`       // 1.查询，2.上架，3.下架
	DemandInfoID       int    `json:"demandInfoID"` // 可选项，查询时不需要指定
	ProductDetailsJson string `json:"productDetailsJSON"`
}

type DemandInfoPayload struct {
	DemandInfoID   int    `json:"demandInfoID"`
	DemandInfoName string `json:"demandInfoName"`
	TapTimes       int    `json:"tapTimes"`
}

// ElementRelease 元素产品发布请求结构体===================================================================
type ElementRelease struct {
	BaseMsg
	NodeCode     string      `json:"nodeCode"` // 授权码
	ProductEle   interface{} `json:"productEle"`
	ProductFiles interface{} `json:"productFiles"`
}

// ElementReleaseResponse 元素产品发布响应结构体
type ElementReleaseResponse struct {
	BaseRetMsg
	ElementID int `json:"elementID"`
}

// ElementManageRequest 元素产品管理请求结构体===================================================================
type ElementManageRequest struct {
	BaseMsg
	NodeCode  string `json:"nodeCode"`  // 授权码
	Method    int8   `json:"method"`    // 1.查询，2.上架，3.下架
	ElementID int    `json:"elementID"` // 可选，查询所有元素产品是不需要指定的
}

// ElementManageResponse 元素产品管理应答结构体
type ElementManageResponse struct {
	BaseRetMsg
	Method             int8   `json:"method"` // 1.查询，2.上架，3.下架
	ProductDetailsJson string `json:"productDetailsJSON"`
	ElementID          int    `json:"elementID"` // 可选，查询所有元素产品是不需要指定的
}

type ElementManagePayload struct {
	ElementID     int    `json:"elementID"`
	ElementName   string `json:"elementName"`
	OnSale        int    `json:"onSale"`
	UpLoaded      int    `json:"upLoaded"`
	PurchaseTimes int    `json:"purchaseTimes"`
	TapTimes      int    `json:"tapTimes"`
}

// CopyAiProductToCloudResource 复制AI产品到云资源===================================================================
type CopyAiProductToCloudResource struct {
	BaseMsg
	Payload struct {
		CloudResourceID string           `json:"CloudResourceID"`
		DownloadDetails []DownloadDetail `json:"DownloadDetails"`
	} `json:"Payload"`
}

// 下载详情
type DownloadDetail struct {
	//下载源：1.平台2.本地（节点）
	DownloadSource int `json:"DownloadSource"`
	//两种下载方式：1.url下载(目前不用)，传URL;2.minio下载,传文件名
	ProductID   int    `json:"ProductID"`
	ProductName string `json:"ProductName"`
	ProductType int    `json:"ProductType"`
	FileName    string `json:"FileName"`
	Url         string `json:"Url"`
}
type CopyAiProductToCloudResourceResponse struct {
	BaseRetMsg
}

// SoftVersionRequest 请求软件版本结构体===================================================================
type SoftVersionRequest struct {
	BaseMsg
	NodeCode string `json:"nodeCode"`
}

// SoftVersionResponse 请求软件版本应答结构体
type SoftVersionResponse struct {
	BaseRetMsg
	SoftWaresJson []string `json:"softWares"`
}

// SoftDownloadRequest 请求软件详情结构体===================================================================
type SoftDownloadRequest struct {
	BaseMsg
	NodeCode   string `json:"nodeCode"`
	SoftWareID int    `json:"softWareID"`
}

// SoftDownloadResponse 请求软件详情应答结构体
type SoftDownloadResponse struct {
	BaseRetMsg
	Url     string `json:"url"`
	Version string `json:"version"`
	Size    string `json:"size"`
}

// PackageAllocateRequest 资源分配请求结构体===================================================================
// 请求节点分配一台实例，这台实例可能是一个订单中的若干实例之一
type PackageAllocateRequest struct {
	BaseMsg
	ReUserID            int    `json:"reUserID"` //该订单所属的平台userID
	YgAccount           string `json:"ygAccount"`
	YgPassword          string `json:"ygPassword"`
	Phone               string `json:"phone"`
	OrderID             int    `json:"orderID"`
	OrderComputingID    int    `json:"orderComputingID"`
	ProductID           int    `json:"productID"` //该产品在云管节点上的ID，由于裸金属资源的特殊性，分配过程需要这个ID
	OrderType           int    `json:"orderType"` // 1: 创建新资源，2: 扩容，3: 创建混合资源
	ResourceType        int    `json:"resourceType"`
	IsCustom            int    `json:"isCustom"`
	AllocateDetailsJson string `json:"allocateDetailsJSON"`
	OrderDetailJSON     string `json:"orderDetailJSON"` // 订单详细描述：这个实例所属订单的详细描述，这个订单中的所有实例都会包含相同的描述
	AiList              string `json:"aiList"`
	UseTime             int64  `json:"useTime"` // 该实例的可用时长，单位小时
}

type OrderDetailJSON struct {
	Code        string    `json:"code"`
	CreatedTime time.Time `json:"createdTime"`
	BillingType int       `json:"billingType"`
	Duration    int       `json:"duration"`
	ProductId   int       `json:"productId"`
	ProductName string    `json:"productName"`
	ProductNum  int       `json:"productNum"`
	SpecJson    string    `json:"specJson"`
	Env         string    `json:"env"`
	Os          string    `json:"os"`
	AiList      string    `json:"aiList"`
	Storage     string    `json:"storage"`
}

// PackageAllocateResponse 资源分配应答结构体
type PackageAllocateResponse struct {
	BaseRetMsg
	NodeCode            string    `json:"nodeCode"`
	Time                time.Time `json:"time"`
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime"`
	CompositeResourceID string    `json:"composite_resourceID"`
	OrderID             int       `json:"orderID"`
	OrderComputingID    int       `json:"orderComputingID"`
	Url                 string    `json:"url"`
	SshAddress          string    `json:"sshAddress"`
	SshPort             int       `json:"sshPort"`
	SshAccount          string    `json:"sshAccount"`
	SshSecret           string    `json:"sshSecret"`
}

// ReleaseRequest 资源释放请求结构体===================================================================
type ReleaseRequest struct {
	BaseMsg
	CompositeResourceID string `json:"composite_resourceID"`
	Reason              string `json:"reason"` // 释放原因
}

// ReleaseResponse 资源释放应答结构体
type ReleaseResponse struct {
	BaseRetMsg
	Time                time.Time `json:"time"`                 // 释放时间
	CompositeResourceID string    `json:"composite_resourceID"` // 必须唯一
}

// TenantInfoRequest 租户信息查询请求结构体===================================================================
type TenantInfoRequest struct {
	BaseMsg
	TenantID int `json:"tenantID"` // 用户ID
}

// TenantInfoResponse 租户信息查询应答结构体
type TenantInfoResponse struct {
	BaseRetMsg
	ResourcesInfoJson string `json:"resourcesInfoJson"`
	OperationInfoJson string `json:"operationInfoJson"`
}

// TenantAccountControlRequest 租户账户控制请求结构体===================================================================
type TenantAccountControlRequest struct {
	BaseMsg
	TenantID    int    `json:"tenantID"`    // 用户ID
	Demand      int    `json:"demand"`      // 控制方法：1.冻结 2.解冻3.修改密码
	NewPassword string `json:"newPassword"` // 可选
}

type TenantAccountControlResponse struct {
	BaseRetMsg
}

// YgModifyRequest 云管系统配置修改请求结构体===================================================================
type YgModifyRequest struct {
	BaseMsg
	//1.允许私有云接入2.允许私有云自动升级3.修改管理员信息4.查询运管统计信息5.查询云管运行信息
	Method      int         `json:"method"`
	Access      int         `json:"access"`
	AutoUpgrade int         `json:"autoUpgrade"`
	AdminModify AdminModify `json:"adminModify"`
	Target      string      `json:"target"`
}

type AdminModify struct {
	NickName string `json:"nickName"`
	Contact  string `json:"contact"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type YgModifyResponse struct {
	BaseRetMsg
	NodeCode   string `json:"nodeCode"`
	Method     int    `json:"method"`
	Statistics `json:"statistics"`
	Run        `json:"run"`
}

type Statistics struct {
	Vm        int `json:"vm"`
	BareMetal int `json:"bareMetal"`
	Container int `json:"container"`
}

type Run struct {
	Status string `json:"status"`
}

// ProductStatusChangeRequest 算力平台产品状态变动通知请求===================================================================
// 适用算力，AI，供需状态变动通知
type ProductStatusChangeRequest struct {
	BaseMsg
	YgProductID int `json:"ygProductID"`
	ProductType int `json:"productType"` //1.算力2.元素3.供需
	//1.上架2.下架3.审核不通过 4通过 5.删除
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

// ProductStatusChangeResponse 云管平台产品状态变动通知应答结构体
type ProductStatusChangeResponse struct {
	BaseRetMsg
	YgProductID int `json:"ygProductID"`
}

// PushRequest 云管平台可用资源规格推送请求结构体===================================================================
type PushRequest struct {
	BaseMsg
	NodeCode           string `json:"nodeCode"`
	ProductDetailsJson string `json:"productDetailsJson"`
}

// PushResponse 云管平台可用资源规格推送应答结构体
type PushResponse struct {
	BaseRetMsg
	ProductID int `json:"productID"`
}

// OrderPushRequest 云管平台资源订购情况推送请求结构体===================================================================
type OrderPushRequest struct {
	BaseMsg
	NodeCode           string `json:"nodeCode"`
	Contact            string `json:"contact"`
	Information        string `json:"information"`
	ProductDetailsJson string `json:"productDetailsJson"`
}

type OrderPushPayload struct {
	SourceType string          `json:"source_type"` // "container/base/vm"
	Spec       []OrderPushSpec `json:"spec"`
}

type OrderPushSpec struct {
	Name string `json:"name"` // "cpu/mem/disk/gpu"
	Unit string `json:"unit"` // 单位
	Data int    `json:"data"` // 数量
}

// OrderPushResponse 云管平台资源订购情况推送应答结构体
type OrderPushResponse struct {
	BaseRetMsg
	ProductID int `json:"productID"`
}

// GetTypeReq 获取元素类型===================================================================
type GetTypeReq struct {
	Command string `json:"command"`
	ID      int64  `json:"ID"`
	Payload struct {
	}
}

// DataSetModelRes 获取元素类型响应
type DataSetModelRes struct {
	BaseRetMsg
	Payload struct {
		Res []Res
	} `json:"payload"`
}

type Res struct {
	Root   ProductDictionary   `json:"root"`
	Second []ProductDictionary `json:"second"`
	Third  []ProductDictionary `json:"third"`
}
type ProductDictionary struct {
	//global.APP_MODEL
	Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:主键" json:"id"  form:"id"`
	CreateBy     string    `gorm:"column:create_by;type:varchar(255);default:null;comment:创建者" json:"createBy" form:"createBy"`
	CreateUserId int64     `gorm:"column:create_user_id;type:bigint;default:null;comment:创建者id" json:"createUserId" form:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime(3);autoCreateTime:milli;comment:创建时间" json:"createTime" form:"createTime"`
	UpdateBy     string    `gorm:"column:update_by;type:varchar(255);default:null;comment:更新者" json:"updateBy" form:"updateBy"`
	UpdateUserId int64     `gorm:"column:update_user_id;type:bigint;default:null;comment:更新者id" json:"updateUserId" form:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime(3);autoUpdateTime:milli;comment:更新时间" json:"updateTime" form:"updateTime"`
	IsDeleted    int64     `gorm:"column:is_deleted;type:tinyint;default:0;comment:删除标识（0=未删除*，1=已删除）" json:"isDeleted" form:"isDeleted"`
	Remark       string    `gorm:"column:remark;type:varchar(255);default:null;comment:备注" json:"remark" form:"remark"`
	//ID       int64   `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:主键" json:"id"  form:"id"`
	Type     int64  `json:"type" form:"type" gorm:"column:type;type:int;comment:1=机柜 2=带宽 3=算力 4=硬件"`
	Title    string `json:"title" form:"title" gorm:"column:title;type:varchar(255);comment:标题名称"`
	ParentId int64  `json:"parentId" form:"parentId" gorm:"column:parent_id;type:int;comment:父id"`
}

// RegionReq 请求地域信息===================================================================
type RegionReq struct {
	BaseMsg
	Payload struct {
	}
}

// RegionRes 地域信息响应
type RegionRes struct {
	BaseRetMsg
	AuthCode string `json:"authCode"`
	Payload  struct {
		Regions []InfoRegion
	} `json:"payload"`
}

type InfoRegion struct {
	Id           *int       `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;size:19;"`                              //主键
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:255;"`                    //创建者
	CreateUserId *int       `json:"createUserId" form:"createUserId" gorm:"column:create_user_id;comment:创建者id;size:19;"`      //创建者id
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间;"`       //创建时间
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新者;size:255;"`                    //更新者
	UpdateUserId *int       `json:"updateUserId" form:"updateUserId" gorm:"column:update_user_id;comment:创建者id;size:19;"`      //创建者id
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间;"`       //更新时间
	IsDeleted    int64      `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;default:0;comment:删除标识（0=未删除*，1=已删除）;"` //删除标识（0=未删除*，1=已删除）
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`                            //备注
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:名称;size:191;"`                                  //名称
	Code         string     `json:"code" form:"code" gorm:"column:code;comment:编码;size:191;"`                                  //编码
	Type         *int       `json:"type" form:"type" gorm:"column:type;comment:类型（1=国，2=省，3=市，4=区县，5=方向n区）;"`                  //类型（1=国，2=省，3=市，4=区县，5=方向n区）
	TypeId       *int       `json:"typeId" form:"typeId" gorm:"column:type_id;comment:关联id;size:19;"`                          //关联id
}

// TokenReq 平台用户访问节点指定资源请求token===================================================================
type TokenReq struct {
	BaseMsg
	OrderComputingID int64  `json:"orderComputingID"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

// TokenRes 平台用户访问节点请求token响应
type TokenRes struct {
	BaseRetMsg
	OrderComputingID int64 `json:"orderComputingID"`
	Payload          struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"payload"`
}

// UploadNodeStaticsRequest 上传节点状态请求===================================================================
type UploadNodeStaticsRequest struct {
	BaseMsg
	NodeCode    string `json:"nodeCode"`
	StaticsJson string `json:"staticsJson"`
}

// UploadNodeStaticsResponse 上传节点状态响应
type UploadNodeStaticsResponse struct {
	BaseRetMsg
}

// 获取租户财务明细

type GetTenantFinanceDetailRequest struct {
	BaseMsg
	Payload ReqInfo `json:"payload"`
}

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   // 关键字
}
type ReqInfo struct {
	PageInfo   PageInfo
	NodeCode   *string `json:"nodeCode" form:"nodeCode"`
	BuyerPhone *string `json:"buyerPhone" form:"buyerPhone"`
	//SupplyUser 供应者
	UserId *int `json:"userId"  form:"userId"`
	////BuyUser 购买者
	//BuyUser string `json:"buyUser" form:"buyUser"`
	//Status 1=未结算 2=人工结算 3=自动结算
	BuyUserId *int `json:"buyUserId" form:"buyUserId"`
	Status    *int `json:"status" form:"status"`
	//ProductType 产品类型 和product表subType保持一致
	ProductType *int `json:"productType" form:"productType"`
	// PayType 1=线上 2=线下
	PayType *int `json:"payType" form:"payType"`
	//StartTime  开始时间
	StartTime *string `json:"startTime" form:"startTime"`
	//EndTime  结束时间
	EndTime *string `json:"endTime" form:"endTime"`
	//SortWay 排序方式  asc 升序 desc 降序
	SortWay *string `json:"sortWay"  form:"sortWay"`
}

type GetTenantFinanceDetailResponse struct {
	BaseRetMsg
	Payload struct {
		List     interface{} `json:"list"`
		Total    int         `json:"total"`
		Page     int         `json:"page"`
		PageSize int         `json:"pageSize"`
	}
}
