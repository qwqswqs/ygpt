package maintenance

import (
	"github.com/pkg/errors"
	"ygpt/server/global"
	"ygpt/server/service/compute/websocket/protocol"
)

type MonitorService struct{}

func (m *MonitorService) GiveTenantInfo(message *protocol.Monitor) error {
	tenantInfoMsg := protocol.MonitorResponse{}
	tenantInfoMsg.ID = message.ID
	tenantInfoMsg.Command = protocol.TenantInfoQueryRetCmd
	var err error
	//TODO: 查询数据库
	if err != nil {
		tenantInfoMsg.RetCode = -400
		tenantInfoMsg.RetMsg = "查询失败"
	} else {
		tenantInfoMsg.RetCode = 0
		tenantInfoMsg.RetMsg = "查询成功"
		//TODO:粘贴数据
	}
	return global.GVA_WS.SendingWithRetry(tenantInfoMsg)
}

func (m *MonitorService) NodeStaticsUpload(staticsJson string) error {
	if global.GVA_WS == nil {
		return errors.New("websocket is nil")
	}
	msg := protocol.UploadNodeStaticsRequest{}
	msg.ID = global.GVA_WS.Counter.GetID()
	msg.Command = protocol.NodeStaticsCmd
	msg.NodeCode = global.GVA_CONFIG.System.NodeCode
	msg.StaticsJson = staticsJson
	return global.GVA_WS.SendingWithRetry(msg)
}

func (m *MonitorService) HandleNodeStaticsUploadResponse(message *protocol.UploadNodeStaticsResponse) error {
	if message.RetCode == 0 {
		return nil
	} else {
		return errors.New(message.RetMsg)
	}
}
