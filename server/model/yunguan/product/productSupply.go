package product

import (
	"time"
)

type APP_MODEL struct {
	Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:主键" json:"id"  form:"id"`
	CreateBy     string    `gorm:"column:create_by;type:varchar(255);default:null;comment:创建者" json:"createBy" form:"createBy"`
	CreateUserId int       `gorm:"column:create_user_id;type:bigint;default:null;comment:创建者id" json:"createUserId" form:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime(3);autoCreateTime:milli;comment:创建时间" json:"createTime" form:"createTime"`
	UpdateBy     string    `gorm:"column:update_by;type:varchar(255);default:null;comment:更新者" json:"updateBy" form:"updateBy"`
	UpdateUserId int       `gorm:"column:update_user_id;type:bigint;default:null;comment:更新者id" json:"updateUserId" form:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime(3);autoUpdateTime:milli;comment:更新时间" json:"updateTime" form:"updateTime"`
	IsDeleted    int       `gorm:"column:is_deleted;type:tinyint;default:0;comment:删除标识（0=未删除*，1=已删除）" json:"isDeleted" form:"isDeleted"`
	Remark       string    `gorm:"column:remark;type:varchar(255);default:null;comment:备注" json:"remark" form:"remark"`
}

type ProductSupply struct {
	APP_MODEL
	Name           string    `json:"name" form:"name" gorm:"column:name;comment:需求标题"`
	Description    string    `json:"description" form:"description" gorm:"column:description;type:text;comment:描述"`
	Price          float64   `json:"price" form:"price" gorm:"column:price;type:double;comment:价格"`
	Inventory      int32     `json:"inventory" form:"inventory" gorm:"column:inventory;type:int;comment:库存"`
	HeadImage      string    `json:"headImage" form:"headImage" gorm:"column:head_image;type:text;comment:首图"`
	Type           int       `json:"type" form:"type"  gorm:"column:type;type:int;comment:1=云上算力产品 2=元素产品 3=供给信息 4=需求信息"`
	SubType        int       `json:"subtype"  form:"subtype"  gorm:"column:subtype;type:int;comment:1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他"`
	LeaseStartTime time.Time `json:"startTime" form:"startTime" gorm:"column:start_time;type:datetime(3);default:null;comment:上架开始时间"`
	LeaseEndTime   time.Time `json:"endTime" form:"endTime" gorm:"column:end_time;type:datetime(3);default:null;comment:上架结束时间"`
	ValidTime      time.Time `json:"validTime" form:"validTime" gorm:"column:valid_time;type:datetime(3);default:null;comment:有效期"`
	Status         int8      `json:"status" form:"status" gorm:"column:status;type:tinyint;default:0;comment:'审核状态（-2 = 未同步 -1=驳回，0=未审核*，1=已通过 3=已同步)';"`
	IsDownload     int8      `json:"isDownload"  form:"isDownload"  gorm:"column:is_download;type:tinyint;default:0;comment:0=可下载 1不可下载"`
	Source         int       `json:"source"  form:"source" gorm:"column:source;type:int;comment:来源 1=平台(自营) 2=第三方 3=运营商 4=智算中心 5=用户 6=节点"`
	UserType       int       `json:"userType"  form:"userType"  gorm:"column:user_type;type:int;comment:个人用户 企业用户"`
	UserId         int64     `json:"userId" form:"userId" gorm:"column:user_id;type:bigint;comment:供应商id"`
	PriceUnit      int64     `json:"unit"  form:"unit" gorm:"column:price_unit;type:int;default:0;comment:0=元 1=日 2=月 3=年 "`
	SpecJson       string    `json:"specJson" form:"specJson" gorm:"column:spec_json;type:text;comment:规格json"`
	RegionId       int64     `json:"regionId"  form:"regionId" gorm:"column:region_id;type:bigint;comment:城市id"`
	User           string    `json:"user"  form:"user" gorm:"column:user;type:varchar(255);comment:联系人"`
	UserPhone      string    `json:"userPhone"  form:"userPhone"  gorm:"column:user_phone;type:varchar(255);comment:联系电话"`
	UserEmail      string    `json:"userEmail"  form:"userEmail" gorm:"column:user_email;type:varchar(255);comment:联系邮箱"`
	ActiveId       int       `json:"activeId"  form:"activeId" gorm:"column:active_id;type:bigint;default:0;comment:活动id"`
	Click          int64     `json:"click" form:"click"   gorm:"column:click;type:bigint;default:0;comment:点击量"`
	//因为涉及到云管的部分资源 给个唯一的关联码
	ReCode     int    `json:"reCode"  form:"reCode"  gorm:"column:re_code;type:int;comment:关联码"`
	ReNodeCode string `json:"reNodeCode"  form:"reNodeCode"  gorm:"column:re_node_code;type:varchar(255);comment:关联节点编码"`
	TypeId     int64  `json:"typeId" form:"typeId"  gorm:"column:type_id;type:bigint;comment:类别id,关联到sys_type"`
	//20250103 新增产品标识
	IsCustom   int64    `json:"isCustom"  form:"isCustom"  gorm:"column:is_custom;type:bigint;default:0;comment:定制产品标识  1=标识定制云算力产品,0标识标配云算力产品"`
	Parameters *float64 `json:"parameters" form:"parameters" gorm:"column:parameters;type:decimal(10,1);default:null"` //2024.2.8 模型产品的参数大小
	Feedback   string   `json:"feedback" form:"feedback" gorm:"column:feedback;type:text;default:null"`
}

func (e *ProductSupply) TableName() string {
	return "yun_product_supply"
}

type GetProductSupplyReq struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Type      int `json:"type"`
}
