package server

import (
	"fmt"
	"github.com/CDUwenbojin/websocket"
	"go.uber.org/zap"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/service/compute/udp/protocol"
	"yunion.io/x/log"
)

func (u *WebSocketServer) HandleLoginMsg(sessionId websocket.SessionID, msg *protocol.LoginMsg) error {
	global.GVA_LOG.Debug("HandleLoginMsg")

	//isExisted := LoginQuery(msg.Payload.PrivateIps)

	loginRetMsg := &protocol.LoginRetMsg{}
	loginRetMsg.Payload.Minio.SecretAccessKey = global.GVA_CONFIG.Minio.SecretKey
	loginRetMsg.Payload.Minio.AccessKey = global.GVA_CONFIG.Minio.AccessKey
	loginRetMsg.Payload.Minio.Endpoint = global.GVA_CONFIG.Minio.EndPoint
	loginRetMsg.Payload.Minio.BucketName = global.GVA_CONFIG.Minio.BucketName
	loginRetMsg.Payload.Minio.UseSSL = global.GVA_CONFIG.Minio.UseSSL
	//if !isExisted {
	//	loginRetMsg.SetError("Login failed", msg.ID)
	//} else {
	loginRetMsg.SetSuccess(msg.ID)
	u.SetClientInfo(msg.Payload.PrivateIps[0], sessionId)
	log.Infoln("ip", msg.Payload.PrivateIps[0])
	global.HandleMsgTransfer("store", sessionId, msg.ID, nil)
	//}
	return u.SendMessage(sessionId, loginRetMsg)
}
func (u *WebSocketServer) HandleMonitorMsg(sessionId websocket.SessionID, msg *protocol.MonitorMsg) error {
	global.GVA_LOG.Debug("HandleMonitorMsg")
	// 处理监控信息
	//monitorService := MonitorService{}
	//err := monitorService.HandleMonitorInfo(sessionId, &msg.Payload.SystemInfo)
	//if err != nil {
	//	global.GVA_LOG.Error("HandleMonitorMsg", zap.Error(err))
	//}
	//monitorRetMsg := &protocol.MonitorRetMsg{}
	//monitorRetMsg.Command = "MonitorRet"
	//monitorRetMsg.ID = msg.ID
	//if err != nil {
	//	monitorRetMsg.RetCode = -1
	//	monitorRetMsg.RetMsg = err.Error()
	//}
	//return u.SendMessage(sessionId, monitorRetMsg)
	return nil
}

func (u *WebSocketServer) HandleAlarmMsg(sessionId websocket.SessionID, msg *protocol.AlarmMsg) error {
	global.GVA_LOG.Debug("HandleAlarmMsg")

	alarmRetMsg := &protocol.AlarmRetMsg{}

	alarmRetMsg.Command = "AlarmRet"
	alarmRetMsg.ID = msg.ID

	return u.SendMessage(sessionId, alarmRetMsg)
}

func (u *WebSocketServer) HandleDownloadRetMsg(sessionId websocket.SessionID, msg *protocol.DownloadRetMsg) error {
	global.GVA_LOG.Debug("HandleDownloadRetMsg")

	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}

func (u *WebSocketServer) HandleDownloadFinishMsg(sessionId websocket.SessionID, msg *protocol.DownloadFinishMsg) error {
	global.GVA_LOG.Debug("HandleDownloadFinishMsg")
	global.GVA_LOG.Info("HandleDownloadFinishMsg", zap.Any("msg", msg))
	//处理下载完成；下载完成或失败更新状态
	var elementDownload product.ProductElementDownload
	if err := global.GVA_DB.Where("uuid = ?", msg.Payload.UUID).First(&elementDownload).Error; err != nil {
		global.GVA_LOG.Error("error find elementDownload", zap.Error(err))
		return err
	}
	elementDownload.Status = msg.Payload.Status
	if err := global.GVA_DB.Save(&elementDownload).Error; err != nil {
		global.GVA_LOG.Error("error save elementDownload", zap.Error(err))
		return err
	}
	downloadFinishRetMsg := &protocol.DownloadFinishRetMsg{}
	downloadFinishRetMsg.Command = "DownloadFinishRet"
	downloadFinishRetMsg.ID = msg.ID

	return u.SendMessage(sessionId, downloadFinishRetMsg)
}
func (u *WebSocketServer) HandleInstallRetMsg(sessionId websocket.SessionID, msg *protocol.InstallRetMsg) error {
	global.GVA_LOG.Debug("HandleInstallRetMsg")

	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}

func (u *WebSocketServer) HandleInstallFinishMsg(sessionId websocket.SessionID, msg *protocol.LoginMsg) error {
	global.GVA_LOG.Debug("HandleInstallFinishMsg")
	return nil
}

func (u *WebSocketServer) HandleGetInfoRetMsg(sessionId websocket.SessionID, msg *protocol.GetInfoRet) error {
	global.GVA_LOG.Debug("HandleGetInfoRetMsg")

	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}

func (u *WebSocketServer) HandleAddTaskRetMsg(sessionId websocket.SessionID, msg *protocol.AddTaskRetMsg) error {
	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}

func (u *WebSocketServer) HandleListTaskRetMsg(sessionId websocket.SessionID, msg *protocol.ListTaskRetMsg) error {
	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}
func (u *WebSocketServer) HandleManageTaskRetMsg(sessionId websocket.SessionID, msg *protocol.ManageTaskRetMsg) error {
	global.GVA_LOG.Info("HandleManageTaskRetMsg", zap.Any("msg", msg))
	msgChan, ok := global.HandleMsgTransfer("get", sessionId, msg.ID, nil)
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if msg.RetCode < 0 {
		msgChan <- fmt.Errorf(msg.RetMsg)
		return nil
	}
	msgChan <- msg
	return nil
}
func (u *WebSocketServer) HandleKeepMsg(sessionId websocket.SessionID, msg *protocol.KeepMsg) error {
	global.GVA_LOG.Debug("HandleKeepMsg")

	retMsg := &protocol.KeepRetMsg{}
	retMsg.Command = protocol.CMD_KEEP_RET
	retMsg.ID = msg.ID
	retMsg.RetCode = 0
	retMsg.RetMsg = "Keep alive success"

	return u.SendMessage(sessionId, retMsg)
}
