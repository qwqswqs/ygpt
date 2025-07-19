package websocket

import (
	"github.com/CDUwenbojin/websocket"
	gmap "github.com/doraemonkeys/sync-gmap"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/service/compute/websocket/protocol"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/service/compute/websocket/ws"
	"ygpt/server/task"
)

func NewWebsocketClient(endpoint string) *ws.WsClient {
	WsCli := websocket.NewClient(
		websocket.WithEndpoint(endpoint),
		websocket.WithClientCodec("json"),
		websocket.WithClientPayloadType(websocket.MsgTypeBinary),
	)

	c := &ws.WsClient{
		WsCli:   WsCli,
		Counter: ws.Counter{},
	}
	//创建一个资源信息发送通道
	global.GVA_SLMSG_TRANSFER = *gmap.NewSyncMap[int64, chan any]()
	//注册请求处理函数
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.MonitorCmd, HandleMonitorRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.AllocateCmd, HandlePackageAllocateRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ReleaseCmd, HandleReleaseRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.YgModifyCmd, HandleYgModifyRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.TicketCmd, HandleTicketRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.TenantInfoQueryCmd, HandleTenantInfoQueryRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.TenantControlCmd, HandleTenantControlRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ProductStatusChangeCmd, HandleProductStatusChangeRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.TokenReqCmd, HandleTokenRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.OrderTicketCmd, HandleOrderTicketRequest)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.CopyAiCmd, HandleCopyAiProductToCloudResource)
	//注册响应处理函数
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.LoginRetCmd, HandleLoginResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.LogoutRetCmd, HandleLogoutResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.KeepRetCmd, HandleKeepResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.PushRetCmd, HandlePushResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.OrderPushRetCmd, HandleOrderPushResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.SupplyInfoReleaseRetCmd, HandleSupplyInfoReleaseResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ManageSupplyInfoRetCmd, HandleSupplyInfoManageResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.DemandInfoReleaseRetCmd, HandleDemandInfoReleaseResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ManageDemandInfoRetCmd, HandleDemandInfoManageResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ElementUploadRetCmd, HandleElementResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ElementManageRetCmd, HandleElementManageResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.SoftVersionRetCmd, HandleSoftwareVersionResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.SoftDownloadRetCmd, HandleSoftwareDownloadResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.NodeStaticsRetCmd, HandleNodeUploadResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.QuoteRetCmd, HandleQuoteResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.TenantFinanceRetCmd, HandleTenantFinanceResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ProductStatusChangeRetCmd, HandleProductStatusChangeResponse)

	websocket.RegisterClientMessageHandler(c.WsCli, protocol.RegionRetCmd, HandleRegionResponse)
	websocket.RegisterClientMessageHandler(c.WsCli, protocol.ThirdTagRetCmd, HandleElementTypeResponse)

	//websocket.RegisterClientMessageHandler(c.WsCli, protocol.PackageUploadRetCmd, HandlePackageUploadResponse)
	//websocket.RegisterClientMessageHandler(c.WsCli, protocol.AlertRetCmd, HandleAlertResponse)
	return c
}

func InitWebsocketClient(endpoint string) *ws.WsClient {
	webs := NewWebsocketClient(endpoint)
	err := webs.Connect()
	if err != nil {
		global.GVA_LOG.Error("连接平台失败", zap.Error(err))
	}
	return webs
}
func KeepLogin(wc *ws.WsClient) {
	//maxRetries := 5        // 最大重试次数
	initialRetryDelay := 2 // 初始重试延迟（秒）
	backoffFactor := 2     // 指数退避因子
	maxRetryDelay := 60    // 最大重试延迟（秒）

	retryDelay := time.Duration(initialRetryDelay) * time.Second
	for {
		if !wc.WsCli.GetIsConnected() {
			err := wc.Connect() // 尝试连接平台,如果无法连接，无响应的情况本身会阻塞若干秒（默认超时时间）
			if err != nil {
				time.Sleep(retryDelay)
				retryDelay *= time.Duration(backoffFactor)
				if retryDelay > time.Duration(maxRetryDelay)*time.Second {
					retryDelay = time.Duration(maxRetryDelay) * time.Second
				}
				continue
			}

			// 成功连接后，尝试登录平台,并且同步一次资源信息
			err = service.SuanLiServiceGroupApp.Login()
			if err != nil {
				global.GVA_LOG.Error("发送节点登录请求失败", zap.Error(err))
				// 不进行重试，直接返回
				return
			} else {
				task.PushResInfo()
			}
			// 成功连接并登录后，重置重试计数和延迟
			retryDelay = time.Duration(initialRetryDelay) * time.Second
		}
		time.Sleep(time.Second * 5)
	}
}
