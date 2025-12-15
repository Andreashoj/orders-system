package main

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/commands"
	"github.com/andreashoj/order-system/internal/db"
	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
	"github.com/andreashoj/order-system/internal/repos"
	"github.com/andreashoj/order-system/internal/services"
)

func main() {
	DB, err := db.NewDB()
	if err != nil {
		fmt.Printf("Failed starting the DB: %s", err)
		return
	}

	rclient, err := pubsub.NewRabbitMqClient()
	if err != nil {
		fmt.Printf("Failed starting rabbit client: %s", err)
		return
	}

	err = pubsub.SetupSubs(rclient)
	if err != nil {
		fmt.Printf("Failed setting up subscriptions: %s", err)
		return
	}

	// Repos
	userRepo := repos.NewUserRepo(DB)
	cartRepo := repos.NewCartRepo(DB)
	productRepo := repos.NewProductRepo(DB)

	// Declare services
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
		case commands.Cart:
		case commands.Checkout:
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
	_, err := shoppingService.AddToCart(user.ID, selection.ID, quantity)
	if err != nil {
		return fmt.Errorf("failed adding product to cart: %s", err)
	}

	// Prompt user to continue shopping or to check out
	// Continue ? Print Display Catalogue
	// Check out ? Trigger check out flow (somehow)

	// User can check out
	// Should see their can and get a proceed confirmation selection
	// This should trigger events for
	// - Create transaction
	// - Update inventory
	// - Create shipping - create some sort of tracking ID, that holds procees of shipment status

	return nil
}
