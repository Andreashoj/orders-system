package pubsub

import (
	"context"
	"encoding/json"
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

func SetupExchange(client *amqp091.Connection) error {
	_, err := NewExchange(client, ExchangeOrderDirect, amqp091.ExchangeDirect)
	if err != nil {
		return fmt.Errorf("failed setting up exchange: %s", err)
	}

	return nil
}

func NewExchange(client *amqp091.Connection, exchange, kind string) (*amqp091.Channel, error) {
	ch, _ := client.Channel()
	err := ch.ExchangeDeclare(exchange, kind, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating exchange: %s", err)
	}

	return ch, nil
}

func NewQueue(
	client *amqp091.Connection,
	queueName, routingKey, exchange string) (<-chan amqp091.Delivery, *amqp091.Channel, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening channel: %s", err)
	}

	// Create queue
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed creating queue: %s", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed binding queue to exchange: %s", err)
	}

	// Consume queue
	transactionCh, _ := ch.Consume(queueName, "", false, false, false, false, nil)

	return transactionCh, ch, nil
}

func QueueHandler[T any](
	client *amqp091.Connection,
	ch <-chan amqp091.Delivery,
	handler func(client *amqp091.Connection, payload T) bool) {
	go func() {
		for tr := range ch {

			var payload T
			fmt.Println(payload)
			err := json.Unmarshal(tr.Body, &payload)
			if err != nil {
				fmt.Printf("failed decoding payload in transaction goroutine: %s", err)
				tr.Nack(false, true)
				return
			}

			ok := handler(client, payload)
			if !ok {
				fmt.Printf("something went wrong while handling the transaction! ARH!")
				tr.Nack(false, true)
				return
			}

			tr.Ack(false)
		}
	}()
}

func NewPublish(client *amqp091.Connection, exchange, routingKey string, data any) error {
	ch, err := client.Channel()
	if err != nil {
		return fmt.Errorf("failed starting channel in publish: %s", err)
	}
	defer ch.Close()

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed encoding pub msg: %s", err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})

	if err != nil {
		return fmt.Errorf("failed creating publish: %s", err)
	}

	return nil
}
