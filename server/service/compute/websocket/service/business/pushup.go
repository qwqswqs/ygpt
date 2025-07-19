package business

import (
	encodingJson "encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"ygpt/server/global"
	product2 "ygpt/server/model/yunguan/product"
	"ygpt/server/service/compute/websocket/protocol"
)

type SpecUploadService struct{}

func (w *SpecUploadService) SyncCloudComputeInfo() error {
	var productInfo product2.ProductComputing
	var productIDs []int
	if err := global.GVA_DB.Model(&productInfo).
		Where("is_synced = ?", 1).
		Pluck("id", &productIDs).Error; err != nil {
		return err
	} else {
		for _, productID := range productIDs {
			err = w.Push(productID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Push 上传云算力产品
func (w *SpecUploadService) Push(productID int) error {
	if productID == 0 {
		return fmt.Errorf("productID is empty")
	}
	spec := protocol.PushRequest{}
	spec.ID = global.GVA_WS.Counter.GetID()
	spec.Command = protocol.PushCmd
	spec.NodeCode = global.GVA_CONFIG.System.NodeCode
	//根据ID查云算力产品表
	var product product2.ProductComputing
	if err := global.GVA_DB.Where("id = ?", productID).First(&product).Error; err != nil {
		return err
	} else {
		if err := global.GVA_DB.Where("id = ? and is_synced = ?", productID, 1).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 处理未找到记录的情况
				return errors.New("产品已同步")
			}
			// 其他错误处理
			return err
		}
	}
	jsonStr, _ := encodingJson.Marshal(product)
	spec.ProductDetailsJson = string(jsonStr)

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(spec.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(spec.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(spec); err != nil {
		return err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.PushResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("PushResponse timeout")
		return fmt.Errorf("PushResponse timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.PushResponse)
		if !ok {
			global.GVA_LOG.Error("error type of PushResponse receive from task channel")
			return fmt.Errorf("error type of PushResponse receive from task channel")
		}
		if result.RetCode != 0 {
			return fmt.Errorf("PushResponse error：%v %v", result.RetCode, result.RetMsg)
		}
		return nil
	}
}

func (w *SpecUploadService) HandlePushResult(message *protocol.PushResponse) error {
	var product product2.ProductComputing
	if message.RetCode != 0 {
		if message.RetCode == protocol.ExistCode {
			product.IsSynced = 2
			global.GVA_DB.Where("id = ?", message.ProductID).Updates(&product)
			return nil
		}
		product.IsSynced = 1
		global.GVA_DB.Where("id = ?", message.ProductID).Updates(&product)
		err := fmt.Errorf("element release ret id error：%v %v", message.ID, message.RetCode)
		//TODO 失败，返回消息到前端
		return err
	} else {
		//TODO 成功，返回消息到前端
		global.GVA_DB.Where("id = ?", message.ProductID).First(&product)
		product.IsSynced = 2
		global.GVA_DB.Save(&product)
	}
	return nil
}

// OrderPush 上传规格, 容器，裸金属，云主机
func (w *SpecUploadService) OrderPush() error {
	OrderPushMsg := protocol.OrderPushRequest{}
	OrderPushMsg.ID = global.GVA_WS.Counter.GetID()
	OrderPushMsg.Command = protocol.OrderPushCmd
	OrderPushMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	OrderPushMsg.Contact = ""
	OrderPushMsg.Information = ""
	//var resInso res.res_info
	//TODO: 应当在资源信息表查询空闲资源信息
	return global.GVA_WS.SendingWithRetry(OrderPushMsg)
}

func (w *SpecUploadService) HandleOrderPushResult(message *protocol.OrderPushResponse) error {
	if message.RetCode != 0 {
		err := fmt.Errorf("element release ret id error：%v %v", message.ID, message.RetCode)
		return err
	} else {
		//TODO: 成功，返回消息到前端
	}
	return nil
}

func (w *SpecUploadService) SyncProductStatusChangeToPlatform(productId int, productType int, status int) error {
	// 云管删除产品要同步到算力平台
	var req protocol.ProductStatusChangeRequest
	req.ID = global.GVA_WS.Counter.GetID()
	req.Command = protocol.ProductStatusChangeCmd
	req.YgProductID = productId
	//一般只同步删除产品.，状态为5
	req.Status = status
	req.ProductType = productType

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(req.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(req.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(req); err != nil {
		return err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.ProductStatusChangeResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("ProductStatusChangeRequest timeout")
		return fmt.Errorf("ProductStatusChangeRequest timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.ProductStatusChangeResponse)
		if !ok {
			global.GVA_LOG.Error("error type of ProductStatusChangeResponse receive from task channel")
			return fmt.Errorf("error type of ProductStatusChangeResponse receive from task channel")
		}
		if result.RetCode != 0 {
			return fmt.Errorf("ProductStatusChangeResponse error：%v %v", result.RetCode, result.RetMsg)
		}
		return nil
	}
}

func (w *SpecUploadService) HandleProductStatusChangeRequest(message *protocol.ProductStatusChangeRequest) error {
	var res protocol.ProductStatusChangeResponse
	res.ID = message.ID
	res.Command = protocol.ProductStatusChangeRetCmd
	res.RetCode = protocol.SuccessCodeCmd
	res.RetMsg = protocol.SuccessMsgCmd
	res.YgProductID = message.YgProductID
	switch message.ProductType {
	case 1:
		var product product2.ProductComputing
		if message.Status == 5 {
			//删除云算力产品
			if err := global.GVA_DB.Where("id = ?", message.YgProductID).Delete(&product).Error; err != nil {
				return err
			}
		} else {
			err := global.GVA_DB.Where("id = ?", message.YgProductID).First(&product).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 处理未找到记录的情况
					res.RetCode = protocol.WrongProductIDErrorCode
					res.RetMsg = protocol.WrongProductIDErrorMsg
				} else {
					return err
				}
			} else {
				product.Status = message.Status
			}
			if err := global.GVA_DB.Save(&product).Error; err != nil {
				return err
			}
		}
	//元素
	case 2:
		var element product2.ProductElem
		if message.Status == 5 {
			//删除云算力产品
			if err := global.GVA_DB.Where("id = ?", message.YgProductID).Delete(&element).Error; err != nil {
				return err
			}
		} else {
			err := global.GVA_DB.Where("id = ?", message.YgProductID).First(&element).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 处理未找到记录的情况
					res.RetCode = protocol.WrongProductIDErrorCode
					res.RetMsg = protocol.WrongProductIDErrorMsg
				} else {
					return err
				}
			} else {
				//TODO: 待完善
				//if message.Status != 3 {
				//	element.AuditStatus = 1
				//} else {
				//	element.AuditStatus = 0
				//}
			}
			if err := global.GVA_DB.Save(&element).Error; err != nil {
				return err
			}
		}
	//供需
	case 3:
		var supplyDemand product2.ProductSupply
		if message.Status == 5 {
			//删除供需产品
			if err := global.GVA_DB.Where("id = ?", message.YgProductID).Delete(&supplyDemand).Error; err != nil {
				return err
			}
		} else {
			err := global.GVA_DB.Where("id = ?", message.YgProductID).First(&supplyDemand).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 处理未找到记录的情况
					res.RetCode = protocol.WrongProductIDErrorCode
					res.RetMsg = protocol.WrongProductIDErrorMsg
				} else {
					return err
				}
			} else {
				supplyDemand.Status = int8(message.Status)
			}
			//如果time.Time为空
			if supplyDemand.ValidTime.IsZero() {
				supplyDemand.ValidTime = time.Now()
			}
			if err := global.GVA_DB.Save(&supplyDemand).Error; err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("product type error：%v", message.ProductType)
	}
	return global.GVA_WS.WsCli.SendMessage(&res)
}
