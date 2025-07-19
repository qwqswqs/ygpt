package product

import (
	"database/sql/driver"
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/model/common/request"

	"github.com/pkg/errors"
)

type ContentList []string
type PriceList []int

// Scan 为自定义的类型实现数据库的读写规则
func (ids *ContentList) Scan(i interface{}) error {
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
func (ids ContentList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// Scan 为自定义的类型实现数据库的读写规则
func (ids *PriceList) Scan(i interface{}) error {
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
func (ids PriceList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type ProductConfig struct {
	global.GVA_MODEL
	Tag   string      `json:"tag" gorm:"column:tag;type:text;comment:'标签'"`
	Name  string      `json:"name" gorm:"column:name;type:text;comment:'规格名称'"`
	Type  PriceList   `json:"type" gorm:"column:type;type:json;comment:'资源类别（1.裸金属、2.云主机、3.容器、4.存储）'"`
	Value ContentList `json:"value" gorm:"column:value;type:json;comment:'规格内容'"`
}

func (e *ProductConfig) TableName() string {
	return "yun_product_config"
}

type GetProductConfigListReq struct {
	request.PageInfo
	Type int    `json:"type"`
	Name string `json:"name"`
}
type GetAllProductConfigListRes struct {
	Tag   string      `json:"tag"`
	Name  string      `json:"name"`
	Type  PriceList   `json:"type"`
	Value ContentList `json:"value"`
}
