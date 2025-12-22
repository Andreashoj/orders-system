package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type MockOrderCreator struct {
	CreateOrderFunc func(userID string) (*domain.Order, error)
}

func (m *MockOrderCreator) CreateOrder(userID string) (*domain.Order, error) {
	return m.CreateOrderFunc(userID)
}

type MockEventPublisher struct{}

func (m *MockEventPublisher) PubOrder(orderID string) error {
	return nil
}

type mockAcknowledger struct{}

func (m *mockAcknowledger) Ack(tag uint64, multiple bool) error {
	return nil
}

func (m *mockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

func (m *mockAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func setupTestUser() *domain.User {
	return &domain.User{
		ID:        uuid.NewString(),
		Name:      "Tester",
		Balance:   1000,
		CreatedAt: time.Now(),
	}
}

func setupTestOrderCreator(orderID string) *MockOrderCreator {
	return &MockOrderCreator{
		CreateOrderFunc: func(userID string) (*domain.Order, error) {
			return &domain.Order{ID: orderID}, nil
		},
	}
}

func setupTestChannels(txChan, shipChan, invChan chan amqp091.Delivery) pubsub.ReplyChannels {
	return pubsub.ReplyChannels{
		TransactionReply: txChan,
		InventoryReply:   invChan,
		ShippingReply:    shipChan,
	}
}

func TestHandleCheckout_OrderSuccess(t *testing.T) {
	mockEventPublisher := &MockEventPublisher{}
	user := setupTestUser()
	orderID := uuid.NewString()
	mockShopping := setupTestOrderCreator(orderID)

	msg, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: orderID})
	reply := amqp091.Delivery{
		ContentType:  "application/json",
		Body:         msg,
		Acknowledger: &mockAcknowledger{},
	}

	txChan := make(chan amqp091.Delivery, 10)
	txChan <- reply
	replyChans := setupTestChannels(txChan, make(chan amqp091.Delivery, 10), make(chan amqp091.Delivery, 10))

	err := HandleCheckout(&replyChans, mockEventPublisher, mockShopping, user)
	if err != nil {
		t.Errorf("Failed running tests with error: %s", err)
	}
}

func TestHandleCheckout_OrderShipmentFailure(t *testing.T) {
	mockEventPublisher := &MockEventPublisher{}
	user := setupTestUser()
	orderID := uuid.NewString()
	mockShopping := setupTestOrderCreator(orderID)

	msg, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: orderID})
	reply := amqp091.Delivery{
		ContentType:  "application/json",
		Body:         msg,
		Acknowledger: &mockAcknowledger{},
	}

	txChan := make(chan amqp091.Delivery, 10)
	txChan <- reply
	replyChans := setupTestChannels(txChan, make(chan amqp091.Delivery, 10), make(chan amqp091.Delivery, 10))

	err := HandleCheckout(&replyChans, mockEventPublisher, mockShopping, user)
	if err != nil {
		t.Errorf("Failed running tests with error: %s", err)
	}
}
