package protocol

// 认证
type LoginMsg struct {
	BaseMsg
	Payload struct {
		PrivateIps []string `json:"PrivateIps"`
	} `json:"Payload"`
}

// 保活
type KeepMsg struct {
	BaseMsg
	Payload struct {
	} `json:"Payload"`
}

// 监控
type MonitorMsg struct {
	BaseMsg
	Payload struct {
		SystemInfo SystemInfo `json:"systemInfo"`
	} `json:"Payload"`
}

// 报警
type AlarmMsg struct {
	BaseMsg
}

// 下载
type DownloadMsg struct {
	BaseMsg
	Payload struct {
		DownloadDetails []DownloadDetail `json:"DownloadDetails"`
	} `json:"Payload"`
}

// 下载详情
type DownloadDetail struct {
	UUID           string `json:"UUID"`
	DownloadSource int    `json:"DownloadSource"`
	//两种下载方式：1.url下载(目前不用)，传URL;2.minio下载,传文件名
	ProductID   int    `json:"ProductID"`
	ProductName string `json:"ProductName"`
	ProductType int    `json:"ProductType"`
	FileName    string `json:"FileName"`
	//Url         string `json:"Url"`
}

// 下载完成
type DownloadFinishMsg struct {
	BaseMsg
	Payload struct {
		UUID   string `json:"UUID"`
		Status int    `json:"status"` // 0:失败，1：成功
	} `json:"Payload"`
}

// 安装
type InstallMsg struct {
	BaseMsg
	Payload struct {
		SoftWare []string `json:"SoftWare"`
	} `json:"Payload"`
}

type InstallFinishMsg struct {
	BaseMsg
	Payload struct {
	} `json:"Payload"`
}

type GetInfoReq struct {
	BaseMsg
	Payload struct {
	} `json:"Payload"`
}

type AddTaskMsg struct {
	BaseMsg
	Payload struct {
		Mode            string `json:"mode"` //训练还是验证
		ModelFileUUID   string `json:"modelFileUUID"`
		ModelName       string `json:"modelName"`
		DataSetFileUUID string `json:"dataSetFileUUID"`
		TaskName        string `json:"taskName"`
		FrameName       string `json:"frameName"`
	} `json:"Payload"`
}

type ListTaskMsg struct {
	BaseMsg
	Payload struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	} `json:"Payload"`
}

type ManageTask struct {
	BaseMsg
	Payload struct {
		Method int8   `json:"method"` // 1：停止（pause），2：恢复(resume)，3：终止(stop),4：报告(querytask),5：杀死进程(killtask)
		TaskID string `json:"taskID"`
	} `json:"Payload"`
}
