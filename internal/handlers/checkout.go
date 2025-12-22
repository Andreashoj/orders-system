package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
)

type OrderCreator interface {
	CreateOrder(userID string) (*domain.Order, error)
}

var (
	CheckoutTransactionError = errors.New("failed transaction")
	CheckoutInventoryError   = errors.New("failed inventory")
	CheckoutShipmentError    = errors.New("failed shipment")
)

func HandleCheckout(replyChannels *pubsub.ReplyChannels, eventPublisher pubsub.EventPublisher, shoppingService OrderCreator, user *domain.User) error {
	order, err := shoppingService.CreateOrder(user.ID)
	if err != nil {
		return fmt.Errorf("failed creating order: %s", err)
	}

	err = eventPublisher.PubOrder(order.ID)
	if err != nil {
		return fmt.Errorf("failed publishing event: %s", err)
	}

	// TODO: Remove after all cases have been written
	responses := pubsub.OrderProcessRequirements{
		ShipmentComplete:  true,
		InventoryComplete: true,
	}
	replyCounter := 0
	for {
		select {
		case tx := <-replyChannels.TransactionReply:
			var transactionReply pubsub.TransactionReplyMessage
			err = json.Unmarshal(tx.Body, &transactionReply)
			if err != nil {
				return fmt.Errorf("failed decoding transaction payload: %s", err)
			}

			if transactionReply.CorrelationID == order.ID {
				responses.TransactionComplete = transactionReply.Success
				replyCounter++
				if err = tx.Ack(false); err != nil {
					return fmt.Errorf("failed acknowledgement: %s", err)
				}
			}

			// Requeue message
			tx.Nack(false, true)
		case tx := <-replyChannels.ShippingReply:
			fmt.Println("got message!", tx)
		case tx := <-replyChannels.InventoryReply:
			fmt.Println("got message!", tx)
		}

		if replyCounter == 1 {
			break
		}
	}

	if !responses.TransactionComplete {
		return CheckoutTransactionError
	} else if !responses.ShipmentComplete {
		return CheckoutShipmentError
	} else if !responses.InventoryComplete {
		return CheckoutInventoryError
	}

	return nil
}
