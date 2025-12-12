package pubsub

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func SetupSubs(client *amqp091.Connection) error {
	// Setup up subs here

	return nil
}

func declareAndBind(client *amqp091.Connection, exchange, key, queueName string) error {
	ch, _ := client.Channel()

	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed declaring queue: %s", err)
	}

	err = ch.QueueBind(queueName, exchange, key, false, nil)
	if err != nil {
		return fmt.Errorf("failed binding queue: %s", err)
	}

	return nil
}
