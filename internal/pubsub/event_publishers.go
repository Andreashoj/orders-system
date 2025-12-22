package pubsub

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	PubOrder(orderID string) error
}

type eventPublisher struct {
	client *amqp091.Connection
}

func NewEventPublisher(client *amqp091.Connection) EventPublisher {
	return &eventPublisher{
		client: client,
	}
}

func (e *eventPublisher) PubOrder(orderID string) error {
	err := NewPublish(e.client, ExchangeOrderDirect, TransactionKey, PubTransaction{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("failed publishing event: %s", err)
	}

	return nil
}
