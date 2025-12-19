package pubsub

import "github.com/rabbitmq/amqp091-go"

type ReplyChannels struct {
	TransactionReply <-chan amqp091.Delivery
	InventoryReply   <-chan amqp091.Delivery
	ShippingReply    <-chan amqp091.Delivery
}

type TransactionReplyMessage struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	Success       bool   `json:"success,omitempty"`
}

type OrderProcessRequirements struct {
	TransactionComplete bool
	ShipmentComplete    bool
	InventoryComplete   bool
}

type PubTransaction struct {
	OrderId string `json:"order_id"`
}
