package config

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AlertWays []int64

// Scan 为自定义的类型实现数据库的读写规则
func (ids *AlertWays) Scan(i interface{}) error {
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
func (ids AlertWays) Value() (driver.Value, error) {
	// Marshal the `DataIds` slice into a JSON string
	jsonData, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type AlertConfig struct {
	gorm.Model
	ConfigName      string    `gorm:"comment:'配置名称'" json:"configName"`
	NoticeText      string    `gorm:"comment:'通知模版'" json:"noticeText"`
	MonitorType     int64     `gorm:"comment:'报警类型 1.CPU使用率；2.内存使用率；3.磁盘使用率；4.磁盘读速率；5.磁盘写速率；6.网络入带宽；7.网络出带宽；8.网络收包速率；9.网络发包速率；10.GPU使用率；11.GPU显存使用率；12.GPU温度。'" json:"monitorType"`
	CommonAlertData int64     `gorm:"comment:'一般报警阈值：默认单位%，对于4.磁盘读速率为bps；5.磁盘写速率为bps；6.网络入带宽为bps；7.网络出带宽为bps；8.网络收包速率为发包数/秒；9.网络发包速率为发包数/秒；12.GPU温度为摄氏度。'" json:"commonAlertData"`
	SeriesAlertData int64     `gorm:"comment:'严重报警阈值'" json:"seriesAlertData"`
	AlertUser       AlertWays `gorm:"comment:'提醒用户类型 1租户2管理员3算力圈租户';type:json" json:"alertUser"`
	AlertWay        AlertWays `gorm:"comment:'通知方式 1站内通知2短信3邮箱';type:json" json:"alertWay"`

	AlertConfigType int64  `gorm:"comment:'配置类型 1.租户默认配置；2.租户自定义；3.基础资源配置'" json:"alertConfigType"`
	RenterID        int64  `gorm:"comment:'租户id'" json:"renterID"`
	MonitorTime     string `gorm:"comment:'监控时长'" json:"monitorTime"`
	MonitorInterval string `gorm:"comment:'监控时长间隔'" json:"monitorInterval"`
}

func (o *AlertConfig) TableName() string {
	return "yun_alert_config"
}
