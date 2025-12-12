package pubsub

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMqClient() (*amqp091.Connection, error) {
	conn := "amqp://guest:guest@localhost:5672/"
	client, err := amqp091.Dial(conn)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbitmq: %s", err)
	}

	return client, nil
}
