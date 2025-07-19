package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net"
	"strings"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/service/compute/udp/protocol"
	"ygpt/server/service/compute/udp/server"
)

type InstanceService struct {
}

func (i *InstanceService) Download(privateIp string, downloadDetails []res.DownloadDetail) error {
	downloadMsg := &protocol.DownloadMsg{}
	downloadMsg.Command = protocol.CMD_Download
	downloadMsg.ID = protocol.Counter.GetID()
	// 将 downloadDetails 序列化为 JSON 字节切片
	jsonBytes, er := json.Marshal(downloadDetails)
	if er != nil {
		global.GVA_LOG.Error("error marshaling downloadDetails", zap.Error(er))
		return er
	}
	var err error
	err = json.Unmarshal(jsonBytes, &downloadMsg.Payload.DownloadDetails)
	if er != nil {
		global.GVA_LOG.Error("error unmarshaling downloadDetails", zap.Error(er))
		return er
	}
	//通过索引遍历
	for i := 0; i < len(downloadMsg.Payload.DownloadDetails); i++ {
		// 为每个元素设置 UUID
		downloadMsg.Payload.DownloadDetails[i].UUID = uuid.New().String()
		//所有要下载的文件保存一份下载记录到ProductElementDownload表中
		var elementDownload product.ProductElementDownload
		elementDownload.UUID = downloadMsg.Payload.DownloadDetails[i].UUID
		elementDownload.FileName = downloadMsg.Payload.DownloadDetails[i].FileName
		elementDownload.ProductType = downloadMsg.Payload.DownloadDetails[i].ProductType
		elementDownload.DownloadSource = downloadMsg.Payload.DownloadDetails[i].DownloadSource
		elementDownload.ProductID = downloadMsg.Payload.DownloadDetails[i].ProductID
		elementDownload.ProductName = downloadMsg.Payload.DownloadDetails[i].ProductName
		elementDownload.PrivateIP = privateIp
		elementDownload.Status = 0
		err = global.GVA_DB.Debug().Where("product_id = ? and product_name = ? and download_source = ? and file_name = ? and private_ip = ?",
			downloadMsg.Payload.DownloadDetails[i].ProductID,
			downloadMsg.Payload.DownloadDetails[i].ProductName,
			downloadMsg.Payload.DownloadDetails[i].DownloadSource,
			downloadMsg.Payload.DownloadDetails[i].FileName,
			privateIp).First(&elementDownload).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = global.GVA_DB.Create(&elementDownload).Error
			if err != nil {
				global.GVA_LOG.Error("error save elementDownload", zap.Error(err))
				return err
			}
		} else {
			fmt.Println("hello")
			//更新
			//存在该记录的话，继续使用旧的uuid
			downloadMsg.Payload.DownloadDetails[i].UUID = elementDownload.UUID
			err = global.GVA_DB.Model(&elementDownload).Updates(&product.ProductElementDownload{
				Status:         0,
				ProductType:    downloadMsg.Payload.DownloadDetails[i].ProductType,
				ProductName:    downloadMsg.Payload.DownloadDetails[i].ProductName,
				DownloadSource: downloadMsg.Payload.DownloadDetails[i].DownloadSource,
				FileName:       downloadMsg.Payload.DownloadDetails[i].FileName,
			}).Error
			if err != nil {
				global.GVA_LOG.Error("error update elementDownload", zap.Error(err))
				return err
			}
		}
	}
	ip := net.ParseIP(strings.TrimSpace(privateIp))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("privateIp", privateIp))
		return fmt.Errorf("invalid IP address: %s", privateIp)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, downloadMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, downloadMsg.ID, nil)
	}()

	err = server.WsServer.SendMessage(sessionId, downloadMsg)
	fmt.Println("send message", downloadMsg)
	if err != nil {
		global.GVA_LOG.Error("error send message", zap.Error(err))
		errWithDesc := fmt.Errorf("%s: 实例代理程序离线", err.Error())
		return errWithDesc
	}
	deadLine := time.After(time.Second * 5)
	//var result *protocol.DownloadRetMsg
	select {
	case <-deadLine:
		global.GVA_LOG.Error("download timeout")
		return fmt.Errorf("download timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		_, ok = r.(*protocol.DownloadRetMsg)
		if !ok {
			global.GVA_LOG.Error("error type of DownloadRetMsg receive from task channel")
			return fmt.Errorf("error type of DownloadRetMsg receive from task channel")
		}
		return nil
	}
}

func (i *InstanceService) Install(sshHost string, softWare []string) error {
	installMsg := &protocol.InstallMsg{}
	installMsg.Command = protocol.CMD_Install
	installMsg.ID = protocol.Counter.GetID()

	installMsg.Payload.SoftWare = softWare

	ip := net.ParseIP(strings.TrimSpace(sshHost))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("sshHost", sshHost))
		return fmt.Errorf("invalid IP address: %s", sshHost)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, installMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, installMsg.ID, nil)
	}()

	server.WsServer.SendMessage(sessionId, installMsg)

	deadLine := time.After(time.Second * 5)
	//var result *protocol.DownloadRetMsg
	select {
	case <-deadLine:
		global.GVA_LOG.Error("install timeout")
		return fmt.Errorf("install timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		_, ok = r.(*protocol.InstallRetMsg)
		if !ok {
			global.GVA_LOG.Error("error type of InstallRetMsg receive from task channel")
			return fmt.Errorf("error type of InstallRetMsg receive from task channel")
		}
		return nil
	}
}

func (i *InstanceService) GetInfo(privateIp string) (string, []string, error) {
	getInfoMsg := &protocol.GetInfoReq{}
	getInfoMsg.Command = protocol.CMD_GET_INFO
	getInfoMsg.ID = protocol.Counter.GetID()

	ip := net.ParseIP(strings.TrimSpace(privateIp))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("sshHost", privateIp))
		return "", nil, fmt.Errorf("invalid IP address: %s", privateIp)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, getInfoMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, getInfoMsg.ID, nil)
	}()

	server.WsServer.SendMessage(sessionId, getInfoMsg)

	deadLine := time.After(time.Second * 5)
	var result *protocol.GetInfoRet
	select {
	case <-deadLine:
		global.GVA_LOG.Error("getInfo timeout")
		return "", nil, fmt.Errorf("getInfo timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return "", nil, err
		}
		result, ok = r.(*protocol.GetInfoRet)
		if !ok {
			global.GVA_LOG.Error("error type of InstallRetMsg receive from task channel")
			return "", nil, fmt.Errorf("error type of InstallRetMsg receive from task channel")
		}

		return result.Payload.Database, result.Payload.PyLibrary, nil
	}
}

func (i *InstanceService) AddTask(privateIp string, mode, modelFileUUID string, dataSetFileUUID string, taskName string, frameName string) error {
	addTaskMsg := &protocol.AddTaskMsg{}
	addTaskMsg.Command = protocol.CMD_ADD_TASK
	addTaskMsg.ID = protocol.Counter.GetID()
	addTaskMsg.Payload.Mode = mode
	addTaskMsg.Payload.ModelFileUUID = modelFileUUID
	//从ProductElementDownload表获取模型名称
	var elementDownload product.ProductElementDownload
	if err := global.GVA_DB.Where("uuid = ?", modelFileUUID).First(&elementDownload).Error; err != nil {
		global.GVA_LOG.Error("error get elementDownload", zap.Error(err))
		return err
	}
	addTaskMsg.Payload.ModelName = elementDownload.ProductName
	addTaskMsg.Payload.DataSetFileUUID = dataSetFileUUID
	addTaskMsg.Payload.TaskName = taskName
	addTaskMsg.Payload.FrameName = frameName

	ip := net.ParseIP(strings.TrimSpace(privateIp))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("privateIp", privateIp))
		return fmt.Errorf("invalid IP address: %s", privateIp)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, addTaskMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, addTaskMsg.ID, nil)
	}()

	server.WsServer.SendMessage(sessionId, addTaskMsg)

	deadLine := time.After(time.Second * 5)
	var _ *protocol.AddTaskRetMsg
	select {
	case <-deadLine:
		global.GVA_LOG.Error("AddTask timeout")
		return fmt.Errorf("AddTask timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		_, ok = r.(*protocol.AddTaskRetMsg)
		if !ok {
			global.GVA_LOG.Error("error type of AddTaskRetMsg receive from task channel")
			return fmt.Errorf("error type of AddTaskRetMsg receive from task channel")
		}

		return nil
	}
}

func (i *InstanceService) ListTask(privateIp string, page, pageSize int) ([]protocol.TaskInfo, int, error) {
	listTaskProgressMsg := &protocol.ListTaskMsg{}
	listTaskProgressMsg.Command = protocol.CMD_LIST_TASK
	listTaskProgressMsg.ID = protocol.Counter.GetID()

	listTaskProgressMsg.Payload.Page = page
	listTaskProgressMsg.Payload.PageSize = pageSize

	ip := net.ParseIP(strings.TrimSpace(privateIp))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("privateIp", privateIp))
		return nil, 0, fmt.Errorf("invalid IP address: %s", privateIp)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, listTaskProgressMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, listTaskProgressMsg.ID, nil)
	}()

	server.WsServer.SendMessage(sessionId, listTaskProgressMsg)

	deadLine := time.After(time.Second * 5)
	var result *protocol.ListTaskRetMsg
	select {
	case <-deadLine:
		global.GVA_LOG.Error("ListTaskProgress timeout")
		return nil, 0, fmt.Errorf("ListTaskProgress timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return nil, 0, err
		}
		result, ok = r.(*protocol.ListTaskRetMsg)
		if !ok {
			global.GVA_LOG.Error("error type of ListTaskRetMsg receive from task channel")
			return nil, 0, fmt.Errorf("error type of ListTaskRetMsg receive from task channel")
		}

		return result.Payload.TaskInfos, result.Payload.Total, nil
	}
}

func (i *InstanceService) ManageTask(privateIp string, taskId string, method int) error {
	manageTaskMsg := &protocol.ManageTask{}
	manageTaskMsg.Command = "TaskManage"
	manageTaskMsg.ID = protocol.Counter.GetID()
	manageTaskMsg.Payload.Method = int8(method)
	manageTaskMsg.Payload.TaskID = taskId

	ip := net.ParseIP(strings.TrimSpace(privateIp))
	if ip == nil {
		global.GVA_LOG.Error("Invalid IP address", zap.String("privateIp", privateIp))
		return fmt.Errorf("invalid IP address: %s", privateIp)
	}

	_, sessionId := server.WsServer.GetClientInfo(ip.String(), "")

	resultCh := make(chan any, 1)

	global.HandleMsgTransfer("store", sessionId, manageTaskMsg.ID, resultCh)

	defer func() {
		global.HandleMsgTransfer("delete", sessionId, manageTaskMsg.ID, nil)
	}()
	global.GVA_LOG.Info("manageTaskMsg", zap.Any("manageTaskMsg", manageTaskMsg))
	err := server.WsServer.SendMessage(sessionId, manageTaskMsg)
	if err != nil {
		global.GVA_LOG.Error("manageTaskMsg send error", zap.Error(err))
		return err
	}
	deadLine := time.After(time.Second * 5)
	//var result *protocol.ListTaskRetMsg
	select {
	case <-deadLine:
		global.GVA_LOG.Error("manageTaskMsg timeout")
		return fmt.Errorf("manageTaskMsg timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			return err
		}
		_, ok = r.(*protocol.ManageTaskRetMsg)
		if !ok {
			global.GVA_LOG.Error("error type of manageTaskMsg receive from task channel")
			return fmt.Errorf("error type of manageTaskMsg receive from task channel")
		}

		return nil
	}
}
