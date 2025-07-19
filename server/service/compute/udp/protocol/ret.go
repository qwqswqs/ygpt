package protocol

// 认证
type LoginRetMsg struct {
	BaseRetMsg
	Payload struct {
		Minio struct {
			AccessKey       string `json:"accessKey"`
			BucketName      string `json:"bucketName"`
			Endpoint        string `json:"endpoint"`
			SecretAccessKey string `json:"secretAccessKey"`
			UseSSL          int    `json:"useSSL"`
		} `json:"minio"`
	} `json:"Payload"`
}

func (l *LoginRetMsg) SetError(retMsg string, id int64) {
	l.Command = CMD_LOGIN_RET
	l.ID = id
	l.RetCode = -1
	l.RetMsg = retMsg
}

func (l *LoginRetMsg) SetSuccess(id int64) {
	l.Command = CMD_LOGIN_RET
	l.ID = id
	l.RetCode = 0
	l.RetMsg = "success"
}

// 保活
type KeepRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

// 监控
type MonitorRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

// 报警
type AlarmRetMsg struct {
	BaseRetMsg
}

// 下载
type DownloadRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

// 下载完成
type DownloadFinishRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

// 安装
type InstallRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

type InstallFinishRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

type GetInfoRet struct {
	BaseRetMsg
	Payload struct {
		Database  string   `json:"Database"`
		PyLibrary []string `json:"PyLibrary"`
	} `json:"Payload"`
}

type AddTaskRetMsg struct {
	BaseRetMsg
	Payload struct {
	} `json:"Payload"`
}

type ListTaskRetMsg struct {
	BaseRetMsg
	Payload struct {
		TaskInfos []TaskInfo `json:"taskInfos"`
		Total     int        `json:"total"`
	} `json:"Payload"`
}

type ManageTaskRetMsg struct {
	BaseRetMsg
}
