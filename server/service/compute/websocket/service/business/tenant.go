package business

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	systemModel "ygpt/server/model/system"
	renterModel "ygpt/server/model/yunguan/renter"
	"ygpt/server/service/compute/websocket/protocol"
	"ygpt/server/utils"
)

type TenantService struct {
}

func (t *TenantService) GiveTenantInfo(message *protocol.TenantInfoRequest) error {
	tenantInfoMsg := protocol.TenantInfoResponse{}
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

func (t *TenantService) GiveControlTenantResult(message *protocol.TenantAccountControlRequest) error {
	controlRetMsg := protocol.TenantAccountControlResponse{}
	controlRetMsg.ID = message.ID
	controlRetMsg.Command = protocol.TenantControlRetCmd
	var err error
	//TODO: 更新数据库
	if err != nil {
		controlRetMsg.RetCode = -400
		controlRetMsg.RetMsg = "更新失败"
	} else {
		controlRetMsg.RetCode = 0
		controlRetMsg.RetMsg = "更新成功"
	}
	return global.GVA_WS.SendingWithRetry(controlRetMsg)
}

func (t *TenantService) GiveTenantToken(message *protocol.TokenReq) error {
	tokenMsg := protocol.TokenRes{
		BaseRetMsg: protocol.BaseRetMsg{
			BaseMsg: protocol.BaseMsg{
				Command: protocol.TokenRetCmd,
				ID:      message.ID,
			},
		},
		OrderComputingID: message.OrderComputingID,
	}
	var user systemModel.SysUser
	var renter renterModel.Renter
	if err := global.GVA_DB.Where("phone = ? and enable = ?", message.Username, 1).First(&user).Error; err != nil {
		tokenMsg.RetCode = -400
		tokenMsg.RetMsg = "用户不存在或该用户账号冻结"
		return global.GVA_WS.SendingWithRetry(tokenMsg)
	} else {
		if err := global.GVA_DB.Where("username = ?", message.Username).First(&renter).Error; err != nil {
			tokenMsg.RetCode = -400
			tokenMsg.RetMsg = "未找到对应的租户信息，请联系管理员"
		} else {
			token, claims, err := utils.LoginToken(&user)
			if err != nil {
				global.GVA_LOG.Error("获取token失败!", zap.Error(err))
				return err
			}
			tokenMsg.RetCode = protocol.SuccessCodeCmd
			tokenMsg.RetMsg = protocol.SuccessMsgCmd
			tokenMsg.Payload.Token = token
			tokenMsg.Payload.ExpiresAt = claims.RegisteredClaims.ExpiresAt.Unix() * 1000
		}
		return global.GVA_WS.WsCli.SendMessage(tokenMsg)
	}
}

func (t *TenantService) RequestTenantInfoList(reqInfo protocol.ReqInfo) (*protocol.GetTenantFinanceDetailResponse, error) {
	tenantInfoListMsg := protocol.GetTenantFinanceDetailRequest{}
	tenantInfoListMsg.ID = global.GVA_WS.Counter.GetID()
	tenantInfoListMsg.Command = protocol.TenantFinanceCmd
	tenantInfoListMsg.Payload = reqInfo

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(tenantInfoListMsg.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(tenantInfoListMsg.ID))
	}()
	if err := global.GVA_WS.SendingWithRetry(tenantInfoListMsg); err != nil {
		return nil, err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.GetTenantFinanceDetailResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("TenantFinance timeout")
		return nil, fmt.Errorf("TenantFinance timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return nil, err
		}
		result, ok = r.(*protocol.GetTenantFinanceDetailResponse)
		if !ok {
			global.GVA_LOG.Error("error type of TenantFinance receive from task channel")
			return nil, fmt.Errorf("error type of TenantFinance receive from task channel")
		}
		if result.RetCode != 0 {
			return nil, fmt.Errorf("TenantFinance error：%v %v", result.RetCode, result.RetMsg)
		}
		return result, nil
	}
}

func (t *TenantService) HandleTenantInfoListResponse(message *protocol.GetTenantFinanceDetailResponse) {

}
