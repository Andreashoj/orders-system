package main

import (
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

	if err = pubsub.NewQueue(rclient, pubsub.QueueTransaction, pubsub.TransactionKey, pubsub.ExchangeOrderDirect, handleTransaction); err != nil {
		fmt.Printf("failed creating transaction queue: %s", err)
		return
	}

	if err = pubsub.NewQueue(rclient, pubsub.QueueShipping, pubsub.ShippingKey, pubsub.ExchangeOrderDirect, handleShipping); err != nil {
		fmt.Printf("failed creating shipping queue: %s", err)
		return
	}

	if err = pubsub.NewQueue(rclient, pubsub.QueueInventory, pubsub.InventoryKey, pubsub.ExchangeOrderDirect, handleInventory); err != nil {
		fmt.Printf("failed creating inventory queue: %s", err)
		return
	}

	// Repos
	userRepo := repos.NewUserRepo(DB)
	cartRepo := repos.NewCartRepo(DB)
	productRepo := repos.NewProductRepo(DB)

	// Services
	shoppingService := services.NewShoppingService(productRepo, cartRepo)
	registrationService := services.NewRegistrationService(userRepo, cartRepo)

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
				err = handleCheckout(rclient, shoppingService, user)
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
			err = handleCheckout(rclient, shoppingService, user)
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

func handleCheckout(rclient *amqp091.Connection, shoppingService *services.ShoppingService, user *domain.User) error {
	// check cart for items
	_, err := shoppingService.GetCart(user.ID)
	if err != nil {
		return fmt.Errorf("failed retrieving cart: %s", err)
	}

	// >>>>>> create order <<<<<<
	err = pubsub.NewPublish(rclient, pubsub.ExchangeOrderDirect, pubsub.TransactionKey, map[string]string{"hello": "there"})
	if err != nil {
		return fmt.Errorf("failed publishing: %s", err)
	}

	// Create 3 subs [x]
	// Create 3 pubs [x]
	// start transaction
	// start shipping
	// start inventory

	// pub all 3 events
	// sub all 3 events and await for response from replyQ
	// If any of the events fail, create event handlers to rollback

	return nil
}

func handleTransaction(payload map[string]string) bool {
	fmt.Println("handling transaction!")
	return true
}

func handleShipping(payload string) bool {
	fmt.Println("handling shipping!")
	return false
}

func handleInventory(payload string) bool {
	fmt.Println("handling inventory!")
	return false
}
