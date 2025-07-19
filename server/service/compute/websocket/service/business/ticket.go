package business

import (
	"fmt"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/compute/websocket/protocol"
)

type TicketService struct{}

func (d *TicketService) HandleTicket(message *protocol.TicketRequest, result string) error {
	//写入数据库
	ticketRetMsg := protocol.TicketRet{}
	ticketRetMsg.ID = message.ID
	ticketRetMsg.Command = protocol.TicketRetCmd
	ticketRetMsg.RetCode = 0
	ticketRetMsg.RetMsg = "success"
	ticketRetMsg.Payload.TicketsID = message.Payload.TicketsID
	ticketRetMsg.Payload.Result = result
	return global.GVA_WS.SendingWithRetry(ticketRetMsg)
}

func (d *TicketService) QuoteRequest(TicketID int, contract system.Contract, invoice system.Invoice, quotation float64) error {
	var quoteMsg protocol.QuoteRequest
	quoteMsg.ID = global.GVA_WS.Counter.GetID()
	quoteMsg.Command = protocol.QuoteCmd
	quoteMsg.Payload.YgTicketID = TicketID
	var ticket system.SystemTicket
	if err := global.GVA_DB.Where("id = ?", TicketID).First(&ticket).Error; err != nil {
		return err
	}
	quoteMsg.Payload.OrderID = ticket.OrderId
	quoteMsg.Payload.QuotePrice = quotation
	quoteMsg.Payload.Contract = protocol.Contract(contract)
	quoteMsg.Payload.Invoice = protocol.Invoice(invoice)

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(quoteMsg.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(quoteMsg.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(quoteMsg); err != nil {
		return err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.QuoteResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("QuoteResponse timeout")
		return fmt.Errorf("QuoteResponse timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.QuoteResponse)
		if !ok {
			global.GVA_LOG.Error("error type of QuoteResponse receive from task channel")
			return fmt.Errorf("error type of QuoteResponse receive from task channel")
		}
		if result.RetCode != 0 {
			return fmt.Errorf("QuoteResponse error：%v %v", result.RetCode, result.RetMsg)
		}
		return nil
	}
}

func (d *TicketService) HandleQuoteResponse(message *protocol.QuoteResponse) error {
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
	return nil
}
