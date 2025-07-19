package business

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/service/compute/websocket/protocol"
	json "yunion.io/x/jsonutils"
)

type SupplyDemandService struct{}

func (s *SupplyDemandService) SyncSupplyInfo() error {
	var productInfo product.ProductComputing
	var productIDs []int
	if err := global.GVA_DB.Model(&productInfo).
		Where("is_synced = ?", 1).
		Pluck("id", &productIDs).Error; err != nil {
		return err
	} else {
		for _, productID := range productIDs {
			err = s.SupplyInfoRelease(productID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SupplyInfoRelease 发布供给信息
func (s *SupplyDemandService) SupplyInfoRelease(supplyID int) error {

	supplyInfoMsg := protocol.SupplyInfoReleaseRequest{}
	supplyInfoMsg.ID = global.GVA_WS.Counter.GetID()
	supplyInfoMsg.Command = protocol.SupplyInfoReleaseCmd
	supplyInfoMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	//查数据库，添加元素Json细节
	var supplyDemandInfo product.ProductSupply
	if err := global.GVA_DB.Where("id = ?", supplyID).First(&supplyDemandInfo).Error; err != nil {
		return err
	}
	jsonStr := json.Marshal(supplyDemandInfo)
	supplyInfoMsg.ProductDetailsJson = jsonStr.String()

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(supplyInfoMsg.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(supplyInfoMsg.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(supplyInfoMsg); err != nil {
		return err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.SupplyInfoReleaseResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("SupplyInfoReleaseResponse timeout")
		return fmt.Errorf("SupplyInfoReleaseResponse timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.SupplyInfoReleaseResponse)
		if !ok {
			global.GVA_LOG.Error("error type of SupplyInfoReleaseResponse receive from task channel")
			return fmt.Errorf("error type of SupplyInfoReleaseResponse receive from task channel")
		}
		if result.RetCode != 0 {
			return fmt.Errorf("SupplyInfoReleaseResponse error：%v %v", result.RetCode, result.RetMsg)
		}
		return nil
	}
}

func (s *SupplyDemandService) HandleSupplyInfoReleaseResult(message *protocol.SupplyInfoReleaseResponse) error {
	var productInfo product.ProductSupply
	if err := global.GVA_DB.Where("id = ?", message.SupplyInfoID).First(&productInfo).Error; err != nil {
		return err
	}
	if message.RetCode != 0 {
		if message.RetMsg == "产品已存在" {
			productInfo.Status = 2
			global.GVA_DB.Where("id = ?", message.SupplyInfoID).Save(&productInfo)
			return nil
		}
		productInfo.Status = -2
		global.GVA_DB.Where("id = ?", message.SupplyInfoID).Save(&productInfo)
		err := fmt.Errorf("element release ret id error：%v %v", message.ID, message.RetCode)
		return err
	} else {
		productInfo.Status = 2
		if err := global.GVA_DB.Where("id = ?", message.SupplyInfoID).Save(&productInfo).Error; err != nil {
			global.GVA_LOG.Error("error update productInfo status", zap.Error(err))
		}
	}
	return nil
}

// ManageSupplyInfo 管理供给信息
func (s *SupplyDemandService) ManageSupplyInfo(method int8, supplyID ...int) error {
	if method == 0 {
		return fmt.Errorf("报文参数不全")
	}
	manageMsg := protocol.SupplyInfoManageRequest{}
	manageMsg.ID = global.GVA_WS.Counter.GetID()
	manageMsg.Command = protocol.ManageSupplyInfoCmd
	manageMsg.Payload.NodeCode = global.GVA_CONFIG.System.NodeCode
	//1.查询，2.上架，3.下架
	manageMsg.Payload.Method = method
	//要上架、下架的产品ID
	if method == 2 || method == 3 {
		manageMsg.Payload.SupplyInfoID = supplyID[0]
	}
	return global.GVA_WS.SendingWithRetry(manageMsg)
}
func (s *SupplyDemandService) HandleManageSupplyInfoResult(message *protocol.SupplyInfoManageResponse) error {
	// 成功与否更新元素产品数据库
	var productInfo product.ProductSupply
	supplyID := message.SupplyInfoID
	global.GVA_DB.Where("id = ?", supplyID).First(&productInfo)
	if message.RetCode != 0 {
	} else {
		if message.Method == 1 {
		}
		if message.RetCode == 2 {
			productInfo.Status = 1
			global.GVA_DB.Save(&productInfo)
		}
		if message.RetCode == 2 {
			productInfo.Status = 0
			global.GVA_DB.Save(&productInfo)
		}
	}
	return nil
}
