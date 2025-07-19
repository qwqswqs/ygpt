package product

import (
	"database/sql/driver"
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/model/common/request"

	"github.com/pkg/errors"
)

// Scan 为自定义的类型实现数据库的读写规则
func (ids *SpecList) Scan(i interface{}) error {
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
func (ids SpecList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type SpecList []SpecData
type SpecData struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
	// 	Type  PriceList   `json:"type"`
	Value ContentList `json:"value"`
}

type StoreObj struct {
	Tag   string   `json:"tag"`
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type StoreList []StoreObj

// Scan 为自定义的类型实现数据库的读写规则
func (ids *StoreList) Scan(i interface{}) error {
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
func (ids StoreList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type PriceType struct {
	Base struct {
		Day   float64 `json:"day"`
		Hour  float64 `json:"hour"`
		Year  float64 `json:"year"`
		Month float64 `json:"month"`
	} `json:"base"`
	SysDisk struct {
		Day   []float64 `json:"day"`
		Hour  []float64 `json:"hour"`
		Year  []float64 `json:"year"`
		Month []float64 `json:"month"`
	} `json:"sysDisk"`
	DataDisk struct {
		Day   []float64 `json:"day"`
		Hour  []float64 `json:"hour"`
		Year  []float64 `json:"year"`
		Month []float64 `json:"month"`
	} `json:"dataDisk"`
	Bandwidth struct {
		Day   []float64 `json:"day"`
		Hour  []float64 `json:"hour"`
		Year  []float64 `json:"year"`
		Month []float64 `json:"month"`
	} `json:"bandwidth"`
}

// Scan 为自定义的类型实现数据库的读写规则
func (ids *PriceType) Scan(i interface{}) error {
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
func (ids PriceType) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type Bind struct {
	Model []struct {
		ProductId   int     `json:"aiProductId"`
		ProductName string  `json:"productName"`
		Reduction   float64 `json:"reduction"`
		Discount    int     `json:"discount"`
	} `json:"model"`
	Dataset []struct {
		ProductId   int     `json:"aiProductId"`
		ProductName string  `json:"productName"`
		Reduction   float64 `json:"reduction"`
		Discount    int     `json:"discount"`
	} `json:"dataset"`
}

func (ids *Bind) Scan(i interface{}) error {
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
func (ids Bind) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type ProductComputing struct {
	global.GVA_MODEL
	Name         string      `json:"name" gorm:"column:name;NOT NULL;type:text;comment:'产品名称'"`
	Type         int         `json:"type" gorm:"column:type;type:tinyint;comment:'资源的类别，可以是 1.裸金属、2.云主机、3.容器、4.存储'"`
	Description  string      `json:"description" gorm:"column:description;type:text;comment:'云算力产品描述'"`
	Specs        SpecList    `json:"specs" gorm:"column:specs;type:json;comment:'产品规格描述json(描述包含的资源规格和价格范围)'"`
	Status       int         `json:"status" gorm:"column:status;type:tinyint;comment:'状态(1.已保存，2.已上架，3.已下架)'"`
	IsSynced     int8        `json:"isSynced" gorm:"column:is_synced;default:1;comment:'是否同步（1-未同步，2-已同步）'"`
	Billing      ContentList `json:"billing" gorm:"column:billing;type:json;comment:'计费方式'"`
	Price        PriceType   `json:"price" gorm:"column:price;type:json;comment:'价格'"`
	IsEnv        int         `json:"isEnv" gorm:"type:int;comment:'是否可以自定义环境'"`
	IsDemand     int         `json:"isDemand" gorm:"type:int;comment:'是否可以自定义需求'"`
	IsAi         int         `json:"isAi" gorm:"type:int;comment:'是否可以自定义AI'"`
	IsCustom     int         `json:"isCustom" gorm:"column:is_custom;type:tinyint;comment:是否为定制产品"`
	Env          ContentList `json:"env" gorm:"column:env;type:json;comment:自定义环境"`
	Os           ContentList `json:"os" gorm:"column:os;type:json;comment:操作系统"`
	Bind         Bind        `json:"bind" gorm:"column:bind;type:json;default:null;comment:绑定AI产品列表"`
	AuditStatus  int         `json:"auditStatus" gorm:"column:audit_status;default:0;type:tinyint;comment:审核状态(-1:驳回，0：等待审核，1通过)"`
	AuditOpinion string      `json:"auditOpinion" gorm:"column:audit_opinion;type:text;comment:审核意见（主要说明驳回原因）"`
	Storage      StoreList   `json:"storage" gorm:"column:storage;type:json;default:null;comment:存储"`
	ResourceId   string      `json:"resourceId" gorm:"column:resource_id;type:varchar(255);comment:物理机资源ID"`
}

func (e *ProductComputing) TableName() string {
	return "yun_product_computing"
}

type GetProductComputingReq struct {
	request.PageInfo
	Type        int    `json:"type"`
	Name        string `json:"name"`
	AuditStatus *int   `json:"auditStatus"`
	IsCustom    *int   `json:"isCustom"`
	IsSynced    *int   `json:"isSynced"`
	TimeDesc    string `json:"timeDesc"`
}
