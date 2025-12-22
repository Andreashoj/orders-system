package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
)

type OrderCreator interface {
	CreateOrder(userID string) (*domain.Order, error)
}

func HandleCheckout(replyChannels *pubsub.ReplyChannels, eventPublisher pubsub.EventPublisher, shoppingService OrderCreator, user *domain.User) error {
	order, err := shoppingService.CreateOrder(user.ID)
	if err != nil {
		return fmt.Errorf("failed creating order: %s", err)
	}

	err = eventPublisher.PubOrder(order.ID)
	if err != nil {
		return fmt.Errorf("failed publishing event: %s", err)
	}

	var responses pubsub.OrderProcessRequirements
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

	return nil
}
