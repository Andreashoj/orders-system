package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
	"github.com/andreashoj/order-system/internal/services"
)

//type CheckoutHandler struct {
//	EventHandler    *pubsub.EventHandler
//	shoppingService *services.ShoppingService
//}

// TODO:
// Move publish into a service that can be mocked, so we avoid having a direct dependency to rabbitmq
// Mock shoppingService with an interface for CreateOrder
// Mock responses to the reply channel
// Assert results

func HandleCheckout(eventHandler pubsub.EventHandler, shoppingService *services.ShoppingService, user *domain.User) error {
	order, err := shoppingService.CreateOrder(user.ID)
	if err != nil {
		return fmt.Errorf("failed creating order: %s", err)
	}

	// Queue up events for transaction, shipment and inventory
	err = pubsub.NewPublish(eventHandler.GetClient(), pubsub.ExchangeOrderDirect, pubsub.TransactionKey, pubsub.PubTransaction{OrderId: order.ID})
	if err != nil {
		return fmt.Errorf("failed publishing: %s", err)
	}

	var responses pubsub.OrderProcessRequirements
	replyCounter := 0
	for {
		select {
		case tx := <-eventHandler.GetReplyChannels().TransactionReply:
			var transactionReply pubsub.TransactionReplyMessage
			err = json.Unmarshal(tx.Body, &transactionReply)
			if err != nil {
				return fmt.Errorf("failed decoding transaction payload: %s", err)
			}

			if transactionReply.CorrelationID == order.ID {
				responses.TransactionComplete = transactionReply.Success
				replyCounter++
				tx.Ack(false)
			}

			// Requeue message
			tx.Nack(false, true)
		case tx := <-eventHandler.GetReplyChannels().ShippingReply:
			fmt.Println("got message!", tx)
		case tx := <-eventHandler.GetReplyChannels().InventoryReply:
			fmt.Println("got message!", tx)
		}

		if replyCounter == 1 {
			break
		}
	}

	return nil
}
