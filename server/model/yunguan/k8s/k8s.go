package k8s

import "gorm.io/gorm"

type K8sModel struct {
	gorm.Model
	Name     string `gorm:"column:name;type:text;comment:'集群名称'" json:"name"`
	K8sId    string `gorm:"column:k8s_id;type:text;comment:'集群ID'" json:"k8sId"`
	ApiIP    string `gorm:"column:api_ip;type:text;comment:'集群IP'" json:"apiIP"`
	Endpoint string `gorm:"column:endpoint;type:text;comment:'集群访问端口'" json:"endpoint"`
}

func (o *K8sModel) TableName() string {
	return "yun_k8s"
}
