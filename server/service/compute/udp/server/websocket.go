package server

import (
	"sync"
	"ygpt/server/global"
	"ygpt/server/service/compute/udp/protocol"

	"github.com/CDUwenbojin/websocket"
	gmap "github.com/doraemonkeys/sync-gmap"
	"go.uber.org/zap"
)

var WsServer WebSocketServer

type WebSocketServer struct {
	wsSrv             *websocket.Server
	counter           Counter
	messageHandlers   websocket.MessageHandlerMap
	loginIpMap        gmap.SyncMap[string, websocket.SessionID]
	loginSessionIDMap gmap.SyncMap[websocket.SessionID, string]
}

type Counter struct {
	mu sync.Mutex
	ID int64
}

func (c *Counter) GetID() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ID++
	if c.ID > 50000 {
		c.ID = 0
	}
	return c.ID
}

func (w *WebSocketServer) GetClientInfo(ip string, sessionId websocket.SessionID) (string, websocket.SessionID) {
	if sessionId != "" {
		ip, _ = w.loginSessionIDMap.Load(sessionId)
		return ip, sessionId
	} else if ip != "" {
		sessionId, _ = w.loginIpMap.Load(ip)
		return ip, sessionId
	} else {
		return "", ""
	}
}

func (w *WebSocketServer) SetClientInfo(ip string, sessionId websocket.SessionID) {
	w.loginIpMap.Store(ip, sessionId)
	w.loginSessionIDMap.Store(sessionId, ip)
}

func (w *WebSocketServer) RemoveClientInfo(sessionId websocket.SessionID) {
	ip, _ := w.loginSessionIDMap.Load(sessionId)
	w.loginSessionIDMap.Delete(sessionId)
	w.loginIpMap.Delete(ip)
}
func (w *WebSocketServer) SendMessage(sessionId websocket.SessionID, msg any) error {
	return w.wsSrv.SendMessage(sessionId, msg)
}
func (w *WebSocketServer) handleConnect(sessionId websocket.SessionID, register bool) {
	if register {
		global.GVA_LOG.Info("instance connected: " + string(sessionId))
	} else {
		w.RemoveClientInfo(sessionId)
		global.GVA_LOG.Info("instance disconnected: " + string(sessionId))
		return //在连接断开时也会调用这个函数，为了防止死锁，进行return
	}
}

// NewWebsocketServer create a websocket server.
func (w *WebSocketServer) NewWebsocketServer() error {
	w.wsSrv = websocket.NewServer(
		websocket.WithAddress(":7082"),
		websocket.WithPath("/"),
		websocket.WithConnectHandle(w.handleConnect),
		websocket.WithCodec("json"),
		websocket.WithMsgType(websocket.MsgTypeText),
	)

	global.GVA_MSG_TRANSFER = *gmap.NewSyncMap[websocket.SessionID, gmap.SyncMap[int64, chan any]]()

	w.messageHandlers = make(websocket.MessageHandlerMap)
	w.loginIpMap = *gmap.NewSyncMap[string, websocket.SessionID]()
	w.loginSessionIDMap = *gmap.NewSyncMap[websocket.SessionID, string]()

	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_LOGIN, w.HandleLoginMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_KEEP, w.HandleKeepMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Monitor, w.HandleMonitorMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Alarm, w.HandleAlarmMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Download_RET, w.HandleDownloadRetMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Download_Finish, w.HandleDownloadFinishMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Install_RET, w.HandleInstallRetMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_Install_Finish, w.HandleInstallFinishMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_GET_INFO_RET, w.HandleGetInfoRetMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_ADD_TASK_RET, w.HandleAddTaskRetMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_LIST_TASK_RET, w.HandleListTaskRetMsg)
	websocket.RegisterServerMessageHandler(w.wsSrv, protocol.CMD_MANAGE_TASK_RET, w.HandleManageTaskRetMsg)
	err := w.wsSrv.Start()
	if err != nil {
		global.GVA_LOG.Error("websocket server start failed", zap.Error(err))
		return err
	}
	return nil
}
