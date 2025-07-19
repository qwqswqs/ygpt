package product

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
)

type TagList []string

// Scan 为自定义的类型实现数据库的读写规则
func (ids *TagList) Scan(i interface{}) error {
	switch v := i.(type) {
	case []byte:
		// Convert the []byte to a string
		str := string(v)

		// Unmarshal the JSON string into the `DataIds` slice
		err := json.Unmarshal([]byte(str), ids)
		if err != nil {
			return err
		}

	case string:
		// Unmarshal the JSON string directly into the `DataIds` slice
		err := json.Unmarshal([]byte(v), ids)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("gorm.Scan: unsupported type %T", i)
	}

	return nil
}

// Value 为自定义的类型实现数据库的读写规则
func (ids TagList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type ProductElementInfo struct {
	global.GVA_MODEL
	ElementInfoID string    `json:"elementInfoID" gorm:"column:element_info_id;unique;NOT NULL;comment:'元素ID'"`
	Name          string    `json:"name" gorm:"column:name;NOT NULL;comment:'元素名称'"`
	Description   string    `json:"description" gorm:"column:description;comment:'元素描述'"`
	Usage         TagList   `json:"usage" gorm:"column:usage;type:json;comment:'描述元素的应用场景和用途'"`
	UseTime       int       `json:"useTime" gorm:"column:use_time;default:0;comment:'使用次数';"`
	UploadTime    time.Time `json:"uploadTime" gorm:"column:upload_time;comment:'记录元素上传的时间(或者是上次更新时间)'"`
	UploadBy      uint      `json:"uploadBy" gorm:"column:upload_by;comment:'上传用户ID'"`
	Type          int       `json:"type" gorm:"column:type;comment:'元素的类型，支持三种类型：1模型、2数据集、3算法'"`
	Status        int       `json:"status" gorm:"column:status;default:1;comment:'使用状态：1正常、2禁用'"`
	Source        int       `json:"source" gorm:"column:source;default:1;comment:'元素来源：1节点、2租户上传'"`
	MD5           string    `json:"MD5" gorm:"column:md5;comment:'文件校验码''"`
	FileSize      int       `json:"fileSize" gorm:"column:file_size;comment:'文件大小（字节）'"`
	Price         float64   `json:"price" gorm:"column:price;type:decimal(12,2);comment:'文件价格'"`
	IsSynced      int8      `json:"isSynced" gorm:"column:is_synced;default:1;comment:'是否同步（1-未同步，2-已同步）'"`
	OpenType      int       `json:"openType" gorm:"column:open_type;default:1;comment:'资源开放类别（不开放1，节点内可见可用2，开放到算力节点节点内用3，开放到算力节点且可外部用4）'"`
	FilePath      string    `json:"filePath" gorm:"column:file_path;comment:'文件存储URL'"`
	FileName      string    `json:"fileName" gorm:"column:file_name;comment:'文件名'"`
	Endpoint      string    `json:"endpoint" gorm:"column:endpoint;comment:'minio路径'"`
	AccessKey     string    `json:"accessKey" gorm:"column:access_key;comment:'minio访问key'"`
	SecretKey     string    `json:"secretKey" gorm:"column:secret_key;comment:'minio加密key'"`
	BucketName    string    `json:"bucketName" gorm:"column:bucket_name;comment:'minio存储桶'"`
	StartTime     time.Time `json:"startTime" gorm:"column:start_time;type:datetime(6);default:null;comment:'上架开始时间'"`
	EndTime       time.Time `json:"endTime" gorm:"column:end_time;type:datetime(6);default:null;comment:'下架开始时间'"`
	AuditStatus   int       `json:"auditStatus" gorm:"column:audit_status;type:tinyint;default:0;comment:审核状态(-1:驳回，0：等待审核，1通过)"`
	AuditOpinion  string    `json:"auditOpinion" gorm:"column:audit_opinion;type:text;comment:审核意见（主要说明驳回原因）"`
}

func (e *ProductElementInfo) TableName() string {
	return "yun_product_element"
}

type AddElementInfoReq struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Usage       TagList   `json:"usage"`
	Type        int       `json:"type"`
	Price       float64   `json:"price"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}
type GetElemInfoReq struct {
	PageIndex   int       `json:"pageIndex"`
	PageSize    int       `json:"pageSize"`
	Type        int       `json:"type"`   //元素类型
	Source      int       `json:"source"` //数据来源
	UserID      uint      `json:"userID"` //用户ID
	Name        string    `json:"name"`   //元素名称
	OpenType    int       `json:"openType"`
	Status      int       `json:"status"`
	AuditStatus *int      `json:"auditStatus"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}
type GetElemInfoRes struct {
	global.APP_MODEL
	Name        string `json:"name"`
	Description string `json:"description"`
	//Usage        TagList   `json:"usage"`
	//UseTime      int       `json:"useTime"`
	//UploadTime   time.Time `json:"uploadTime"`
	//UploadBy     uint      `json:"uploadBy"`
	Type   int `json:"type"`
	Status int `json:"status"`
	Source int `json:"source"`
	//FilePath     string    `json:"filePath"`
	//FileSize     int       `json:"fileSize"`
	Price    float64 `json:"price"`
	IsSynced int8    `json:"isSynced"`
	OpenType int     `json:"openType"`
	//FileName     string    `json:"fileName"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	AuditStatus  int       `json:"auditStatus"`
	AuditOpinion string    `json:"auditOpinion"`
}
type ProductElem struct {
	global.APP_MODEL
	Name           string    `json:"name" form:"name" gorm:"column:name;comment:需求标题"`
	Description    string    `json:"description" form:"description" gorm:"column:description;type:text;comment:描述"`
	Price          float64   `json:"price" form:"price" gorm:"column:price;type:double;comment:价格"`
	Inventory      int32     `json:"inventory" form:"inventory" gorm:"column:inventory;type:int;comment:库存"`
	Type           int       `json:"type" form:"type"  gorm:"column:type;type:int;comment:1=云上算力产品 2=AI产品 3=供给信息 4=需求信息"`
	SubType        int       `json:"subtype"  form:"subtype"  gorm:"column:subtype;type:int;comment:1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他"`
	LeaseStartTime time.Time `json:"startTime" form:"startTime" gorm:"column:start_time;type:datetime(3);default:null;comment:上架开始时间"`
	LeaseEndTime   time.Time `json:"endTime" form:"endTime" gorm:"column:end_time;type:datetime(3);default:null;comment:上架结束时间"`
	ValidTime      time.Time `json:"validTime" form:"validTime" gorm:"column:valid_time;type:datetime(3);default:null;comment:有效期"`
	Status         *int8     `json:"status" form:"status" gorm:"column:status;type:tinyint;default:0;comment:审核状态（-1=驳回，0=待审核*，1=未上架 2=上架 3=置顶）"`
	IsDownload     int8      `json:"isDownload"  form:"isDownload"  gorm:"column:is_download;type:tinyint;default:0;comment:0=可下载 1不可下载"`
	Source         int       `json:"source"  form:"source" gorm:"column:source;type:int;comment:来源 1=平台(自营) 2=第三方 3=运营商 4=智算中心  5=用户 6=节点"`
	UserType       int       `json:"userType"  form:"userType"  gorm:"column:user_type;type:int;comment:个人用户 企业用户"`
	UserId         int64     `json:"userId" form:"userId" gorm:"column:user_id;type:bigint;comment:供应商id"`
	PriceUnit      int64     `json:"unit"  form:"unit" gorm:"column:price_unit;type:int;default:0;comment:0=元 1=日 2=月 3=年 "`
	RegionId       int64     `json:"regionId"  form:"regionId" gorm:"column:region_id;type:bigint;comment:城市id"`
	ActiveId       int       `json:"activeId"  form:"activeId" gorm:"column:active_id;type:bigint;default:0;comment:活动id"`
	TypeId         int64     `json:"typeId" form:"typeId"  gorm:"column:type_id;type:bigint;comment:类别id,关联到sys_type"`
	Parameters     *float64  `json:"parameters" form:"parameters" gorm:"column:parameters;type:decimal(10,1);default:null"` //2024.2.8 模型产品的参数大小
	IsSynced       int8      `json:"isSynced" gorm:"column:is_synced;default:0;comment:'是否同步（0-未同步，1-已同步）'"`
	OpenType       int       `json:"openType" gorm:"column:open_type;default:1;comment:'资源开放类别（不开放1，节点内可见可用2，开放到算力节点节点内用3，开放到算力节点且可外部用4）'"`
	AuditStatus    int       `json:"auditStatus" gorm:"column:audit_status;type:tinyint;default:0;comment:审核状态(-1:驳回，0：等待审核，1通过)"`
	AuditOpinion   string    `json:"auditOpinion" gorm:"column:audit_opinion;type:text;comment:审核意见（主要说明驳回原因）"`
}

func (ProductElem) TableName() string {
	return "yun_product_elem"
}

type ProductFile struct {
	Id           int        `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;size:19;"`                                 //主键
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:255;"`                       //创建者
	CreateUserId *int       `json:"createUserId" form:"createUserId" gorm:"column:create_user_id;comment:创建者id;size:19;"`         //创建者id
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间;"`          //创建时间
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新者;size:255;"`                       //更新者
	UpdateUserId *int       `json:"updateUserId" form:"updateUserId" gorm:"column:update_user_id;comment:创建者id;size:19;"`         //创建者id
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:更新时间;"`                         //更新时间
	IsDeleted    int        `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识（0=未删除*，1=已删除）;"`              //删除标识（0=未删除*，1=已删除）
	Remark       string     `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255;"`                               //备注
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:文件名;size:255;"`                                    //文件名
	Url          string     `json:"url" form:"url" gorm:"column:url;comment:文件地址;size:255;"`                                      //文件地址
	Tag          string     `json:"tag" form:"tag" gorm:"column:tag;comment:tag;size:255;"`                                       //tag
	Key          string     `json:"key" form:"key" gorm:"column:key;comment:文件地址;size:255;"`                                      //key
	Description  string     `json:"description" form:"description" gorm:"column:description;comment:描述;size:1000;"`               //描述
	Size         *int64     `json:"size" form:"size" gorm:"column:size;comment:文件大小;type:bigint;"`                                //文件大小
	Unit         string     `json:"unit" form:"unit" gorm:"column:unit;comment:单位;size:255;"`                                     //单位
	Downloads    *int       `json:"downloads" form:"downloads" gorm:"column:downloads;type:bigint;comment:下载量;size:10;default:0"` //下载量
	Visits       *int       `json:"visits" form:"visits"  gorm:"column:visits;type:bigint;comment:访问量;default:0"`                 //查看量
	Favorite     *int       `json:"favorite"  form:"favorite"  gorm:"column:favorite;type:bigint;comment:收藏量;default:0"`          //收藏量
	Type         *int       `json:"type" form:"type" gorm:"column:type;comment:类型（1=产品图片，2=产品说明书，3=模型，4=数据集）default:0;"`          //类型（1=产品图片，2=产品说明书，3=模型，4=数据集）
	TypeId       *int       `json:"typeId" form:"typeId" gorm:"column:type_id;comment:关联id;size:19;default:0"`                    //关联id
	ProductId    *int       `json:"productId" form:"productId" gorm:"column:product_id;comment:产品id;size:19;"`                    //产品id
}

func (ProductFile) TableName() string {
	return "yun_product_file"
}

type ProductCreateV2Req struct {
	//Name 标题
	Name string `json:"name"  form:"name"`
	//Source 0=平台 1=用户 2=节点
	Source int `json:"source" form:"source"`
	//Description 描述
	Description string `json:"description"  form:"description"`
	//Price 价格
	Price float64 `json:"price" form:"price"`
	//Inventory  库存
	Inventory int32 `json:"inventory" form:"inventory"`
	//HeadImage 图片地址
	HeadImage string `json:"headImage" form:"headImage"`
	//Type 1=云上算力产品 2=元素产品 3=供给信息 4=需求信息
	Type int `json:"type"   form:"type"`
	//SubType 1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他
	SubType int `json:"subType"  form:"subType"`
	//StartTime 开始时间
	LeaseStartTime time.Time `json:"startTime"  form:"startTime"`
	//EndTime 结束时间
	LeaseEndTime time.Time `json:"endTime"  form:"endTime"`
	//IsDownload 0=可下载 1不可下载"
	IsDownload int8 `json:"isDownload"  form:"isDownload" `
	// 状态
	Status int `json:"status" form:"status"`

	//SpecJson 规格的JSON参数
	SpecJson string `json:"specJson" form:"specJson" `
	//RegionId 城市id
	RegionId int `json:"regionId"  form:"regionId"`

	//User 联系人
	User string `json:"user"  form:"user"`
	//UserPhone 联系电话
	UserPhone string `json:"userPhone"  form:"userPhone"`
	//UserEmail 联系电话
	UserEmail string `json:"userEmail"  form:"userEmail"`
	//ValidTime 有效时间
	ValidTime time.Time `json:"validTime" form:"validTime"`
	//TypeId 关联id
	TypeId int64 `json:"typeId" form:"typeId" `
	//uploadType
	UploadType *int `json:"uploadType"  form:"uploadType"`

	//Url  文件名称
	Url *string `json:"url"  form:"url"`
	// 模型参数
	Parameter *float64  `json:"parameter" form:"parameter"`
	Files     []elFiles `json:"files" form:"files"`
	FileName  *string   `json:"fileName" form:"fileName"`
	Size      int64     `json:"size" form:"size"`
}

type elFiles struct {
	Name       string `json:"name" form:"name"`
	Percentage int    `json:"percentage" form:"percentage"`
	//Raw        string `json:"raw" form:"raw"`
	//Response   string `json:"response" form:"response"`
	Size   int    `json:"size" form:"size"`
	Status string `json:"status" form:"status"`
	//Uid    int64 `json:"uid" form:"uid"`
}
type ProductSearch3 struct {
	request.PageInfo
	//Type 1=云上算力产品，2=元素产品，3=供给信息，4=需求信息
	Type *int `json:"type"  form:"type"`
	//SubType 1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他
	SubType *int `json:"subType"  form:"subType"`
	//UserId 用户id
	UserId *int `json:"userId"  form:"userId"`
	//TypeId   种类id
	TypeId *int `json:"typeId" form:"typeId"`
	//Status  审核状态（-1=驳回，0=待审核*，1=通过）
	Status   *int `json:"status" form:"status"`
	IsSynced *int `json:"isSynced" form:"isSynced"`

	//Source 来源 0=平台 1=用户 2=节点
	Source *int `json:"source"  form:"source"`

	//ReNodeCode 关联节点编码
	ReNodeCode *string `json:"reNodeCode" form:"reNodeCode" `
	PriceDesc  string  `json:"priceDesc" form:"priceDesc"`
	EleType    *int    `json:"eleType" form:"eleType"`
}
