package res

import (
	"gorm.io/gorm"
	"time"
	"ygpt/server/model/common/request"
)

// 设备信息表
type ResDevice struct {
	gorm.Model
	Type        int       `json:"type" gorm:"column:type;type:int;comment:'硬件类别(1.服务器 2.网络 3.其他)'"`
	Brand       string    `json:"brand" gorm:"column:brand;type:text;comment:'品牌'"`
	Location    string    `json:"location" gorm:"column:location;comment:'所在机房、区域、机柜、机位'"`
	SpecDesc    string    `json:"specDesc" gorm:"column:spec_desc;type:json;comment:'规格描述json'"`
	Description string    `json:"description" gorm:"column:description;type:text;comment:'说明'"`
	InstallDate time.Time `json:"installDate" gorm:"column:install_date;type:datetime(6);comment:'安装日期'"`
	SNNum       string    `json:"snNum" gorm:"column:sn_num;type:text;comment:'SN编号'"`
	Status      int       `json:"status" gorm:"column:status;type:int;comment:'状态(1.正常 2.异常)'"`
	CPUType     string    `json:"cpuType" gorm:"column:cpu_type;type:text;comment:'CPU类别'"`
	GPUType     string    `json:"gpuType" gorm:"column:gpu_type;type:text;comment:'GPU类别'"`
}

func (o *ResDevice) TableName() string {
	return "yun_res_device"
}

type GetResDeviceReq struct {
	request.PageInfo
	Type   int `json:"type"`
	Status int `json:"status"`
}
