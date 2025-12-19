package main

import (
	"encoding/json"
	"fmt"

	"github.com/andreashoj/order-system/internal/commands"
	"github.com/andreashoj/order-system/internal/db"
	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
	"github.com/andreashoj/order-system/internal/repos"
	"github.com/andreashoj/order-system/internal/services"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	DB, err := db.NewDB()
	if err != nil {
		fmt.Printf("Failed starting the DB: %s", err)
		return
	}

	// Repos
	userRepo := repos.NewUserRepo(DB)
	cartRepo := repos.NewCartRepo(DB)
	productRepo := repos.NewProductRepo(DB)
	orderRepo := repos.NewOrderRepo(DB)

	// Services
	shoppingService := services.NewShoppingService(productRepo, cartRepo, orderRepo, userRepo)
	registrationService := services.NewRegistrationService(userRepo, cartRepo)

	// Message Broker
	rclient, err := pubsub.NewRabbitMqClient()
	if err != nil {
		fmt.Printf("Failed starting rabbit client: %s", err)
		return
	}

	err = pubsub.SetupExchange(rclient)
	if err != nil {
		fmt.Printf("Failed setting up exchange: %s", err)
		return
	}

	// Queues
	txQueue, _, err := pubsub.NewQueue(rclient, pubsub.QueueTransaction, pubsub.TransactionKey, pubsub.ExchangeOrderDirect)
	if err != nil {
		fmt.Printf("failed creating transaction queue: %s", err)
		return
	}

	txCh, txConn, err := pubsub.NewQueue(rclient, pubsub.QueueTransactionReply, pubsub.TransactionKeyReply, pubsub.ExchangeOrderDirect)
	if err != nil {
		fmt.Printf("failed creating inventory queue: %s", err)
		return
	}
	fmt.Println(txConn)

	// Event/Queue Handlers
	eventHandler := pubsub.NewEventHandler(rclient, shoppingService)
	eventHandler.ReplyChannels.TransactionReply = txCh

	pubsub.QueueHandler(txQueue, func(payload pubsub.PubTransaction) bool {
		fmt.Println("payload", payload)
		return eventHandler.HandleTransaction(payload)
	})

	user, err := handleIntroduction(registrationService)
	if err != nil {
		fmt.Printf("Introduction failed: %s", err)
		return
	}

	for {
		cmd := commands.GetMenu()

		switch cmd {
		case commands.Catalogue:
			err = handleCatalogue(shoppingService, user)
			if err != nil {
				fmt.Printf("Something went wrong while showing the catalogue: %s", err)
				break
			}

			if userWantsCheckout := commands.PromptCheckout(); userWantsCheckout {
				err = handleCheckout(rclient, eventHandler, shoppingService, user)
				if err != nil {
					fmt.Printf("Failed checkout: %s", err)
					return
				}
			}
		case commands.Cart:
			err = handleCart(shoppingService, user)
			if err != nil {
				fmt.Printf("Failed showing cart: %s", err)
				return
			}
		case commands.Checkout:
			err = handleCheckout(rclient, eventHandler, shoppingService, user)
			if err != nil {
				fmt.Printf("Failed creating order: %s", err)
				return
			}
		case commands.Exit:
			fmt.Println("Thanks for stopping by, cya next time!")
			return
		}
	}
}

func handleIntroduction(registrationService *services.RegistrationService) (*domain.User, error) {
	username, err := commands.WelcomeMessage()
	if err != nil {
		return nil, fmt.Errorf("failed getting username: %s", err)
	}

	user, err := registrationService.CreateUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %s", err)
	}

	return user, nil
}

func handleCatalogue(shoppingService *services.ShoppingService, user *domain.User) error {
	products, err := shoppingService.GetAllProducts()
	if err != nil {
		return fmt.Errorf("beep boop, something went wrong - is that a you or me problem.. ?: %s", err)
	}

	// Create catalogue mapping, will be used as the display number the user can select when browsing the catalogue
	catalogue := make(map[int]domain.Product, len(products))
	for i, product := range products {
		catalogue[i+1] = product
	}

	commands.DisplayCatalogue(catalogue)

	// Awaits user inputs
	selection := commands.GetProductSelection(catalogue)
	quantity := commands.GetProductQuantity()

	// Add the selected product + quantity to cart
	_, err = shoppingService.AddToCart(user.ID, selection.ID, quantity)
	if err != nil {
		return fmt.Errorf("failed adding product to cart: %s", err)
	}

	return nil
}

func handleCart(shoppingService *services.ShoppingService, user *domain.User) error {
	cart, err := shoppingService.GetCart(user.ID)
	if err != nil {
		return fmt.Errorf("failed retrieving cart: %s", err)
	}

	commands.DisplayCart(cart)
	return nil
}

func handleCheckout(rclient *amqp091.Connection, eventHandler *pubsub.EventHandler, shoppingService *services.ShoppingService, user *domain.User) error {
	order, err := shoppingService.CreateOrder(user.ID)
	if err != nil {
		return fmt.Errorf("failed creating order: %s", err)
	}

	// Queue up events for transaction, shipment and inventory
	err = pubsub.NewPublish(rclient, pubsub.ExchangeOrderDirect, pubsub.TransactionKey, pubsub.PubTransaction{OrderId: order.ID})
	if err != nil {
		return fmt.Errorf("failed publishing: %s", err)
	}

	var responses pubsub.OrderProcessRequirements
	replyCounter := 0
	for {
		select {
		// GOT ABSOLUTELY REKT BY OLD QUEUE MESSAGES WITHOUT PAYLOAD. - YUP HAPPENED AGAIN
		// NEXT STEP IS ACTUALLY DOING TRANSACTION FUNCTIONALITY IN THE HANDLER.
		// SO WHAT SHOULD HAPPEN IN THE TRANSACTION HANDLER?
		// User should PAY - need to create a balance ? no pay
		// AFTER THAT IS DONE, SAME FOR SHIPPING AND INVENTORY
		// ALSO, FIGURE OUT HOW TO DI IN THE HANDLERS, STRUCT FOR THE EVENT HANDLERS.
		case tx := <-eventHandler.ReplyChannels.TransactionReply:
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
		case tx := <-eventHandler.ReplyChannels.ShippingReply:
			fmt.Println("got message!", tx)
		case tx := <-eventHandler.ReplyChannels.InventoryReply:
			fmt.Println("got message!", tx)
		}

		if replyCounter == 1 {
			break
		}
	}

	//for response := range responses {
	//	// Rollbacks:il
	//	// Transaction => status refund
	//	// Inventory => check order and re-add products to product_inventory table
	//	// Shipping => delete entry
	//	// These should also be created as queues and published to
	//	if !responses.TransactionComplete {
	//		// Rollback strategy
	//		return fmt.Errorf("order failed, transaction was not completed: %s", err)
	//	}
	//}

	// Mark order as complete here ()

	// Create 3 subs [x]
	// Create 3 pubs [x]
	// start transaction
	// start shipping
	// start inventory

	// Handlers should communicate to reply queues when they are done
	// Compare reply queues responses here (order_id), if all 3 are done then we can mark order as complete

	// pub all 3 events
	// sub all 3 events and await for response from replyQ
	// If any of the events fail, create event handlers to rollback

	return nil
}
