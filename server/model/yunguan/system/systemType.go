package system

import (
	"database/sql/driver"
	"encoding/json"
	"ygpt/server/global"
	"yunion.io/x/pkg/errors"
)

type TypeList []string

// Scan 为自定义的类型实现数据库的读写规则
func (ids *TypeList) Scan(i interface{}) error {
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
func (ids TypeList) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type SystemType struct {
	global.GVA_MODEL
	Num         int      `json:"num" gorm:"type:int;comment:'类别编号'"`
	Name        string   `json:"name" gorm:"type:text;comment:'名称'"`
	Type        TypeList `json:"type" gorm:"type:json;comment:'类别'"`
	Description string   `json:"description" gorm:"type:text;comment:'说明'"`
}

func (SystemType) TableName() string {
	return "yun_system_type"
}
