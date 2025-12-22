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

type mockAcknowledger struct {
	AckCalled  bool
	NackCalled bool
}

func (m *mockAcknowledger) Ack(tag uint64, multiple bool) error {
	m.AckCalled = true
	return nil
}

func (m *mockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	m.NackCalled = true
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

	msg, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: orderID, Success: true})
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

func TestHandleCheckout_OrderTransactionCorrectAcksCalled(t *testing.T) {
	mockEventPublisher := &MockEventPublisher{}
	user := setupTestUser()
	orderID := uuid.NewString()
	mockShopping := setupTestOrderCreator(orderID)

	msgWrongCorrelationID, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: "wrong-id"})
	replyWrongMockAcknowledger := &mockAcknowledger{}
	replyWrong := amqp091.Delivery{
		ContentType:  "application/json",
		Body:         msgWrongCorrelationID,
		Acknowledger: replyWrongMockAcknowledger,
	}

	msgCorrectCorrelationID, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: orderID})
	replyCorrectMockAcknowledger := &mockAcknowledger{}
	replyCorrect := amqp091.Delivery{
		ContentType:  "application/json",
		Body:         msgCorrectCorrelationID,
		Acknowledger: replyCorrectMockAcknowledger,
	}

	txChan := make(chan amqp091.Delivery, 10)
	txChan <- replyWrong
	txChan <- replyCorrect
	replyChans := setupTestChannels(txChan, make(chan amqp091.Delivery, 10), make(chan amqp091.Delivery, 10))

	HandleCheckout(&replyChans, mockEventPublisher, mockShopping, user)

	if !replyWrongMockAcknowledger.NackCalled {
		t.Errorf("Unrelated order was not requeued after being handled")
	}

	if !replyCorrectMockAcknowledger.AckCalled {
		t.Errorf("Did not acknowledge the correct order")
	}
}

func TestHandleCheckout_OrderTransactionFailed(t *testing.T) {
	mockEventPublisher := &MockEventPublisher{}
	user := setupTestUser()
	orderID := uuid.NewString()
	mockShopping := setupTestOrderCreator(orderID)

	msg, _ := json.Marshal(pubsub.TransactionReplyMessage{CorrelationID: orderID, Success: false})
	reply := amqp091.Delivery{
		ContentType:  "application/json",
		Body:         msg,
		Acknowledger: &mockAcknowledger{},
	}

	txChan := make(chan amqp091.Delivery, 10)
	txChan <- reply
	replyChans := setupTestChannels(txChan, make(chan amqp091.Delivery, 10), make(chan amqp091.Delivery, 10))

	err := HandleCheckout(&replyChans, mockEventPublisher, mockShopping, user)
	if err == nil {
		t.Errorf("Failed running tests with error: %s", err)
	}
}
