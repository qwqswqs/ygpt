package websocket

import (
	"fmt"
	"ygpt/server/global"
	"ygpt/server/service/compute/websocket/protocol"
	"ygpt/server/service/compute/websocket/service"
)

// ______________________________________________以下是对服务器的请求的处理________________________________________________________
func HandleMonitorRequest(message *protocol.Monitor) error {
	//收到消息，从发送消息队列中删除消息ID,确认收到了应答
	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandlePackageAllocateRequest(message *protocol.PackageAllocateRequest) error {
	fmt.Println("收到消息:%+v", message)

	service.SuanLiServiceGroupApp.ResourceService.AllocateResource(message)
	return nil
}

func HandleReleaseRequest(message *protocol.ReleaseRequest) error {
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.ResourceService.ReleaseResource(message)
	return nil
}
func HandleYgModifyRequest(message *protocol.YgModifyRequest) error {
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.ModifyService.YgModify(message)
	return nil
}

func HandleTenantInfoQueryRequest(message *protocol.TenantInfoRequest) error {
	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleTenantControlRequest(message *protocol.TenantAccountControlRequest) error {
	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleTicketRequest(message *protocol.SoftDownloadRequest) error {
	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandleProductStatusChangeRequest(message *protocol.ProductStatusChangeRequest) error {
	service.SuanLiServiceGroupApp.SpecUploadService.HandleProductStatusChangeRequest(message)
	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandleTokenRequest(message *protocol.TokenReq) error {
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.GiveTenantToken(message)
	return nil
}

func HandleOrderTicketRequest(message *protocol.OrderTicketRequest) error {
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.ResourceService.HandleOrderTicketRequest(message)
	return nil
}

func HandleCopyAiProductToCloudResource(message *protocol.CopyAiProductToCloudResource) error {
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.ElementService.HandleCopyAiProductToCloudResource(message)
	return nil
}

// ___________________________________________以下是对服务器返回的应答的处理________________________________________________

func HandleLoginResponse(message *protocol.LoginResponse) error {
	fmt.Println("收到消息:%+v", message)
	global.GVA_WS.HandleAck(message.ID)
	// 异地登录了
	if message.RetCode == protocol.LoginBeOccupiedCode {
		return fmt.Errorf(message.RetMsg)
	}
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	service.SuanLiServiceGroupApp.LoginService.HandleLoginResult(*message)
	return nil
}

func HandleKeepResponse(message *protocol.KeepResponse) error {
	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleLogoutResponse(message *protocol.LogoutResponse) error {
	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandlePushResponse(message *protocol.PushResponse) error {
	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	service.SuanLiServiceGroupApp.HandlePushResult(message)
	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandleOrderPushResponse(message *protocol.OrderPushResponse) error {
	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

// 供需，元素相关
func HandleSupplyInfoReleaseResponse(message *protocol.SupplyInfoReleaseResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	service.SuanLiServiceGroupApp.HandleSupplyInfoReleaseResult(message)
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleSupplyInfoManageResponse(message *protocol.SupplyInfoManageResponse) error {

	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandleDemandInfoReleaseResponse(message *protocol.DemandInfoReleaseResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	//service.SuanLiServiceGroupApp.HandleDemandInfoReleaseResult(message)
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleDemandInfoManageResponse(message *protocol.DemandInfoManageResponse) error {
	global.GVA_WS.HandleAck(message.ID)

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleElementResponse(message *protocol.ElementReleaseResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	service.SuanLiServiceGroupApp.ElementService.HandleElementReleaseResult(message)
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message
	//WriteMsgToChan(message.ID, message.RetCode, message, message.RetMsg)

	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleElementManageResponse(message *protocol.ElementManageResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleSoftwareVersionResponse(message *protocol.SoftVersionResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	service.SuanLiServiceGroupApp.SoftwareService.HandleQuerySoftwareVersionResult(message)
	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleSoftwareDownloadResponse(message *protocol.SoftDownloadResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return nil
}
func HandleNodeUploadResponse(message *protocol.UploadNodeStaticsResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return nil
}

// 获取元素类型
func HandleElementTypeResponse(message *protocol.DataSetModelRes) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return nil
}

// 获取地域信息
func HandleRegionResponse(message *protocol.RegionRes) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return nil
}

func HandleQuoteResponse(message *protocol.QuoteResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	fmt.Println("收到消息:%+v", message)
	return service.SuanLiServiceGroupApp.TicketService.HandleQuoteResponse(message)
}

func HandleTenantFinanceResponse(message *protocol.GetTenantFinanceDetailResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message
	fmt.Println("收到消息:%+v", message)
	service.SuanLiServiceGroupApp.TenantService.HandleTenantInfoListResponse(message)
	return nil
}

func HandleProductStatusChangeResponse(message *protocol.ProductStatusChangeResponse) error {
	global.GVA_WS.HandleAck(message.ID)
	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(message.ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}

	if message.RetCode < 0 {
		msgChan <- fmt.Errorf(message.RetMsg)
		return nil
	}
	msgChan <- message
	fmt.Println("收到消息:%+v", message)
	return nil
}
func WriteMsgToChan(ID, RetCode int, message interface{}, RetMsg string) error {

	msgChan, ok := global.GVA_SLMSG_TRANSFER.Load(int64(ID))
	if !ok {
		// 可能是超时了
		global.GVA_LOG.Error("not found msg transfer channel")
		return nil
	}
	if RetCode != 0 {
		msgChan <- fmt.Errorf(RetMsg)
		return nil
	}
	msgChan <- message
	return nil
}
