package res

import (
	"time"
	"ygpt/server/model/common/request"

	"gorm.io/gorm"
)

// 资源信息表
type ResInfo struct {
	gorm.Model
	ResourceID   string    `json:"resourceID" gorm:"column:resource_id;type:text;comment:'cloudpods资源ID'"`
	ResourceType int       `json:"resourceType" gorm:"column:resource_type;type:tinyint;comment:'资源类别（1存储、2云主机、3容器、4裸金属）'"`
	StartTime    time.Time `json:"startTime" gorm:"type:datetime(6);comment:'启用时间'"`
	EndTime      time.Time `json:"endTime" gorm:"column:end_time;comment:'结束时间'"`
	ServerID     int       `json:"serverID" gorm:"type:bigint;comment:'所属服务器ID'"`
	SpecDesc     string    `json:"specDesc" gorm:"type:json;comment:'规格json描述'"`
	RenterID     int       `json:"renterID" gorm:"type:bigint;comment:'所属租户ID(0表示没租户用)'"`
	Status       int       `json:"status" gorm:"type:tinyint;comment:'当前状态(1.空闲，2.使用，3.损坏，4.注销)'"`
	IPID         int       `json:"ipID" gorm:"comment:'所属子网id(0表示不属于任何子网)'"`
	ImageID      int       `json:"imageID" gorm:"comment:'所用镜像id'"`
	Url          string    `json:"url" gorm:"type:text;comment:'访问地址'"`
	SshHost      string    `json:"sshHost" gorm:"column:ssh_host;type:varchar(255);comment:'ssh地址'"`
	SshPort      int       `json:"sshPort" gorm:"column:ssh_port;type:int;comment:'ssh端口'"`
	SshUser      string    `json:"sshUser" gorm:"column:ssh_user;type:varchar(255);comment:'ssh用户名'"`
	SshPwd       string    `json:"sshPwd" gorm:"column:ssh_pwd;type:text;comment:'ssh密码(空表示和算力圈密码一致)'"`
}

func (o *ResInfo) TableName() string {
	return "yun_res_info"
}

type GetResInfoReq struct {
	request.PageInfo
	ImageID int `json:"imageID"`
}

type GetResInfoRet struct {
	ID           uint      `json:"id"`
	ResourceType int       `json:"resourceType"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	RunTime      int64     `json:"runTime"`
	ImageName    string    `json:"imageName"`
	SshHost      string    `json:"sshHost"`
	Status       int       `json:"status"`
	System       string    `json:"system"`
	Frame        string    `json:"frame"`
	Language     string    `json:"language"`
}

type LibraryInfo struct {
	DataBase  string   `json:"dataBase"`
	PyLibrary []string `json:"pyLibrary"`
}
type GetRenterResReq struct {
	request.PageInfo
	Type     int `json:"type"`
	RenterID int `json:"renterID"`
}

type MonitorResReq struct {
	PrivateIp string `json:"privateIp"` //实例私有ip地址
}

type InstallReq struct {
	PrivateIp string   `json:"privateIp"`
	SoftWare  []string `json:"softWare"`
}

type DownloadReq struct {
	PrivateIp       string           `json:"privateIp"`
	DownloadDetails []DownloadDetail `json:"DownloadDetails"`
}

type DownloadDetail struct {
	DownloadSource int `json:"DownloadSource"`
	//两种下载方式：1.url下载，传URL;2.minio下载,传文件名
	ProductID   int    `json:"ProductID"`
	ProductName string `json:"ProductName"`
	ProductType int    `json:"ProductType"`
	FileName    string `json:"FileName"`
	Url         string `json:"Url"`
}

type AddTaskReq struct {
	PrivateIp     string `json:"privateIp"` //实例私有ip地址
	Mode          string `json:"mode"`      //训练模式:train  推理模式:inference
	ModelFileUUID string `json:"modelFileUUID"`
	//ModelName       string `json:"modelName"`
	DataSetFileUUID string `json:"dataSetFileUUID"`
	TaskName        string `json:"taskName"`  //任务名称
	FrameName       string `json:"frameName"` //框架名称
}

type ListTaskReq struct {
	request.PageInfo
	PrivateIp string `json:"privateIp"` //实例私有ip地址
}
type ManageTaskReq struct {
	request.PageInfo
	PrivateIp string `json:"privateIp"` //实例私有ip地址\
	TaskId    string `json:"id"`
	Method    int    `json:"manageMethod"`
}
