package server

//var (
//	UdpServer *UDPServer
//)
//
//type Any interface{}
//type Message Any
//type MessageHandler func(Message, *net.UDPAddr) error
//
//type HandlerData struct {
//	Handler MessageHandler
//	Binder  Binder
//}
//type MessageHandlerMap map[websocket.MessageCmd]*HandlerData
//type Binder func() Message
//
//// UDPServer represents a UDP server instance.
//type UDPServer struct {
//	conn            *net.UDPConn
//	addr            *net.UDPAddr
//	messageHandlers MessageHandlerMap
//	loginMap        gmap.SyncMap[string, string]
//}
//
//func InitUDPServer() *UDPServer {
//	global.GVA_MSG_TRANSFER = *gmap.NewSyncMap[string, gmap.SyncMap[int64, chan any]]()
//	u := NewUDPServer()
//	go u.Run()
//	return u
//}
//
//// NewUDPServer initializes and starts a new UDP server.
//func NewUDPServer() *UDPServer {
//	laddr, err := net.ResolveUDPAddr("udp", ":7082")
//	if err != nil {
//		global.GVA_LOG.Error("Failed to resolve address: " + err.Error())
//		panic(err)
//	}
//
//	conn, err := net.ListenUDP("udp", laddr)
//	if err != nil {
//		global.GVA_LOG.Debug("Failed to listen on UDP: " + err.Error())
//		panic(err)
//	}
//
//	global.GVA_LOG.Debug("UDP server started on: " + conn.LocalAddr().String())
//
//	server := &UDPServer{
//		conn:            conn,
//		addr:            laddr,
//		messageHandlers: make(MessageHandlerMap),
//		loginMap:        *gmap.NewSyncMap[string, string](),
//	}
//
//	// Register message handlers
//	RegisterServerMessageHandler(server, protocol.CMD_LOGIN, server.HandleLoginMsg)
//	//RegisterServerMessageHandler(server, "Keep", server.HandleKeepMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Monitor_RET, server.HandleMonitorRetMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Alarm, server.HandleAlarmMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Download_RET, server.HandleDownloadRetMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Download_Finish, server.HandleDownloadFinishMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Install_RET, server.HandleInstallRetMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_Install_Finish, server.HandleInstallFinishMsg)
//	RegisterServerMessageHandler(server, protocol.CMD_GET_INFO_RET, server.HandleGetInfoRetMsg)
//
//	return server
//}
//
//func (u *UDPServer) RegisterMessageHandler(cmd protocol.MessageCmd, handler MessageHandler, binder Binder) {
//	if _, ok := u.messageHandlers[cmd]; ok {
//		return
//	}
//
//	u.messageHandlers[cmd] = &HandlerData{handler, binder}
//}
//
//func RegisterServerMessageHandler[T any](srv *UDPServer, cmd protocol.MessageCmd, handler func(*T, *net.UDPAddr) error) {
//	srv.RegisterMessageHandler(cmd,
//		func(message Message, remoteAddr *net.UDPAddr) error {
//			switch t := message.(type) {
//			case *T:
//				return handler(t, remoteAddr)
//			default:
//				return errors.New("invalid payload struct type")
//			}
//		},
//		func() Message {
//			var t T
//			return &t
//		},
//	)
//}
//
//func (u *UDPServer) DeregisterMessageHandler(cmd protocol.MessageCmd) {
//	delete(u.messageHandlers, cmd)
//}
//
//func (u *UDPServer) Run() {
//	for {
//		data := make([]byte, 10240)
//		n, remoteAddr, err := u.conn.ReadFromUDP(data)
//		if err != nil {
//			global.GVA_LOG.Error("Error receiving UDP message: " + err.Error())
//			continue
//		}
//
//		fmt.Println("Received UDP message from " + remoteAddr.String() + ": " + string(data[:n]))
//		go func() {
//			err := u.messageHandler(data[:n], remoteAddr)
//			if err != nil {
//				global.GVA_LOG.Error("Error processing UDP message: " + err.Error())
//			}
//		}()
//	}
//}
//
//func (u *UDPServer) messageHandler(buf []byte, remoteAddr *net.UDPAddr) error {
//	var err error
//	var handler *HandlerData
//	var message Message
//
//	if handler, message, err = u.unmarshalMessage(buf); err != nil {
//		return err
//	}
//
//	if err = handler.Handler(message, remoteAddr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (u *UDPServer) marshalMessage(message Message) ([]byte, error) {
//	var (
//		codecJsonMsg []byte
//		err          error
//	)
//
//	codecJsonMsg, err = json.Marshal(message)
//	if err != nil {
//		return nil, err
//	}
//
//	//msgWithLength, err = LengthMarshal(codecJsonMsg, c.msgType)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//global.GVA_LOG.DebugInfo("msgWithLength:", string(msgWithLength))
//
//	return codecJsonMsg, nil
//}
//
//func (u *UDPServer) unmarshalMessage(msg []byte) (*HandlerData, Message, error) {
//	var (
//		handler *HandlerData
//		message Message
//		baseMsg protocol.BaseMsg
//		ok      bool
//		err     error
//	)
//
//	//msgWithoutLength, length, err = LengthUnmarshal(msgWithLength, c.msgType)
//	//if err != nil {
//	//	return nil, nil, fmt.Errorf("lengthUnmarshal message exception:%v", err)
//	//}
//	//
//	//if int(length) != len(msgWithoutLength) {
//	//	return nil, nil, fmt.Errorf("incomplete message")
//	//}
//
//	err = json.Unmarshal(msg, &baseMsg)
//	if err != nil {
//		return nil, nil, fmt.Errorf("parse the Json command failed:%v", err)
//	}
//
//	handler, ok = u.messageHandlers[baseMsg.Command]
//	if !ok {
//		return nil, nil, errors.New("message handler not found")
//	}
//
//	if handler.Binder != nil {
//		message = handler.Binder()
//		err = json.Unmarshal(msg, &message)
//		if err != nil {
//			return nil, nil, fmt.Errorf("parse the Json failed:%v", err)
//		}
//	} else {
//		return nil, nil, errors.New("message Binder not found")
//	}
//
//	return handler, message, nil
//}
//
//func (u *UDPServer) SendMessage(message interface{}, remoteAddr *net.UDPAddr) error {
//	msgByte, _ := json.Marshal(message)
//	_, err := u.conn.WriteTo(msgByte, remoteAddr)
//	return err
//}
