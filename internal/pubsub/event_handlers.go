package pubsub

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/services"
	"github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	client          *amqp091.Connection
	shoppingService *services.ShoppingService
	ReplyChannels   *ReplyChannels
}

type ReplyChannels struct {
	TransactionReply <-chan amqp091.Delivery
	InventoryReply   <-chan amqp091.Delivery
	ShippingReply    <-chan amqp091.Delivery
}

func NewEventHandler(client *amqp091.Connection, shoppingService *services.ShoppingService) *EventHandler {
	return &EventHandler{
		client:          client,
		shoppingService: shoppingService,
		ReplyChannels:   &ReplyChannels{},
	}
}

func (e *EventHandler) HandleTransaction(payload PubTransaction) bool {
	order, err := e.shoppingService.GetOrder(payload.OrderId)
	if err != nil {
		fmt.Printf("failed getting order: %s", err)
		return false
	}

	ok := e.shoppingService.ChargeUser(order.UserID, order.GetTotal())
	replyMsg := TransactionReplyMessage{CorrelationID: payload.OrderId, Success: ok}
	err = NewPublish(e.client, ExchangeOrderDirect, TransactionKeyReply, replyMsg)
	if err != nil {
		fmt.Printf("failed publishing to reply queue")
		return false
	}
	return true
}
