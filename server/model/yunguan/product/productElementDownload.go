package product

import (
	"ygpt/server/global"
)

type ProductElementDownload struct {
	global.GVA_MODEL
	//对于模型数据集下载，只有以下两个字段是必需的
	UUID     string `json:"uuid" form:"uuid" gorm:"column:uuid;comment:;type:varchar(255);"`
	FileName string `json:"fileName" form:"fileName" gorm:"column:file_name;comment:'文件名称';type:varchar(255);"`

	ProductType    int    `json:"productType" form:"productType" gorm:"column:product_type;comment:'AI(元素产品)类型';type:int;"`
	DownloadSource int    `json:"downloadSource" form:"downloadSource" gorm:"column:download_source;comment:'下载源：1.平台2.本地（节点）';type:int;"`
	ProductID      int    `json:"productID" form:"productID" gorm:"column:product_id;comment:'产品ID';type:int;"`
	ProductName    string `json:"productName" form:"productName" gorm:"column:product_name;comment:'产品名称';type:varchar(255);"`
	PrivateIP      string `json:"privateIP" form:"privateIP" gorm:"column:private_ip;comment:'下载该文件的实例IP';type:varchar(255);"`
	Status         int    `json:"status" form:"status" gorm:"column:status;comment:'文件状态:-1：下载失败；0-99：下载中；100.下载成功';type:int;"`
}

func (e *ProductElementDownload) TableName() string {
	return "yun_product_element_download"
}
