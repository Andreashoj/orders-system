package main

import (
	"encoding/json"
	"fmt"

	"github.com/andreashoj/order-system/internal/commands"
	"github.com/andreashoj/order-system/internal/db"
	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/handlers"
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
				err = handlers.HandleCheckout(rclient, eventHandler, shoppingService, user)
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
