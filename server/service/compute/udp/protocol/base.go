package protocol

import (
	"sync"

	"github.com/CDUwenbojin/websocket"
)

const (
	CMD_LOGIN               websocket.MessageCmd = "Login"
	CMD_LOGIN_RET           websocket.MessageCmd = "LoginRet"
	CMD_KEEP                websocket.MessageCmd = "Keep"
	CMD_KEEP_RET            websocket.MessageCmd = "KeepRet"
	CMD_Monitor             websocket.MessageCmd = "Monitor"
	CMD_Monitor_RET         websocket.MessageCmd = "MonitorRet"
	CMD_Alarm               websocket.MessageCmd = "Alarm"
	CMD_Alarm_RET           websocket.MessageCmd = "AlarmRet"
	CMD_Download            websocket.MessageCmd = "Download"
	CMD_Download_RET        websocket.MessageCmd = "DownloadRet"
	CMD_Download_Finish     websocket.MessageCmd = "DownloadFinish"
	CMD_Download_Finish_RET websocket.MessageCmd = "DownloadFinishRet"
	CMD_Install             websocket.MessageCmd = "Install"
	CMD_Install_RET         websocket.MessageCmd = "InstallRet"
	CMD_Install_Finish      websocket.MessageCmd = "InstallFinish"
	CMD_Install_Finish_RET  websocket.MessageCmd = "InstallFinishRet"
	CMD_GET_INFO            websocket.MessageCmd = "GetInfo"
	CMD_GET_INFO_RET        websocket.MessageCmd = "GetInfoRet"
	CMD_ADD_TASK            websocket.MessageCmd = "AddTask"
	CMD_ADD_TASK_RET        websocket.MessageCmd = "AddTaskRet"
	CMD_LIST_TASK           websocket.MessageCmd = "ListTask"
	CMD_LIST_TASK_RET       websocket.MessageCmd = "ListTaskRet"
	CMD_MANAGE_TASK         websocket.MessageCmd = "ManageTask"
	CMD_MANAGE_TASK_RET     websocket.MessageCmd = "TaskManageRet"
)

//type MessageCmd string

type BaseMsg struct {
	Command websocket.MessageCmd `json:"Command"`
	ID      int64                `json:"ID"`
}

type BaseRetMsg struct {
	BaseMsg
	RetCode int64  `json:"RetCode"`
	RetMsg  string `json:"RetMsg"`
}

var Counter = &counter{}

type counter struct {
	mu sync.Mutex
	ID int64
}

func (c *counter) GetID() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ID++
	if c.ID > 50000 {
		c.ID = 0
	}
	return c.ID
}

// CPUInfo 包含了CPU的相关信息
type CPUInfo struct {
	Model       string  `json:"model"`        // CPU型号
	Cores       int     `json:"cores"`        // CPU核数
	UsedPercent float64 `json:"used_percent"` // CPU使用率
}

// MemoryInfo 包含了内存的相关信息
type MemoryInfo struct {
	Total       string  `json:"total"`        // 内存总量
	Available   string  `json:"available"`    // 可用内存
	Used        string  `json:"used"`         // 已使用内存
	UsedPercent float64 `json:"used_percent"` // 内存使用率
}

// DiskInfo 包含了磁盘的相关信息
type DiskInfo struct {
	Total       string  `json:"total"`        // 磁盘总量
	Used        string  `json:"used"`         // 已使用磁盘空间
	Free        string  `json:"free"`         // 可用磁盘空间
	UsedPercent float64 `json:"used_percent"` // 磁盘使用率
}

// HostInfo 包含了主机的相关信息
type HostInfo struct {
	Hostname   string `json:"hostname"`    // 主机名
	OS         string `json:"os"`          // 操作系统
	Platform   string `json:"platform"`    // 平台信息
	KernelArch string `json:"kernel_arch"` // 内核架构
}

type SystemInfo struct {
	CPU  CPUInfo    `json:"cpu"`  //cpu信息
	Mem  MemoryInfo `json:"mem"`  //内存信息
	Disk DiskInfo   `json:"disk"` //硬盘信息
	Host HostInfo   `json:"host"` //主机信息
	//Net  interface{} `json:"net"`
}

type TaskInfo struct {
	TaskID          string `json:"taskID"`          //任务ID
	Mode            string `json:"mode"`            //训练模式:train  推理模式:inference
	ModelFilePath   string `json:"modelFilePath"`   //模型文件路径
	ModelName       string `json:"modelName"`       //模型名字
	DataSetFilePath string `json:"dataSetFilePath"` //数据集文件路径
	TaskName        string `json:"taskName"`        //任务名称
	FrameName       string `json:"frameName"`       //框架名称
	ResultMsg       string `json:"resultMsg"`       //输出
}
