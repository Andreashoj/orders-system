package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"

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
	productRepo := repos.NewProductRepo(DB)

	// Declare services
	registrationService := services.NewRegistrationService(userRepo)
	catalogueService := services.NewCatalogueService(productRepo)

	_, err = handleIntroduction(registrationService)
	if err != nil {
		fmt.Printf("Introduction failed: %s", err)
		return
	}

	for {
		cmd := commands.GetInput()

		switch cmd {
		case commands.Catalogue:
			err = handleCatalogue(catalogueService)
			if err != nil {
				fmt.Printf("Something went wrong while showing the catalogue: %s", err)
				break
			}
		case commands.Cart:
		case commands.Checkout:

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

func handleCatalogue(catalogueService *services.CatalogueService) error {
	products, err := catalogueService.GetAllProducts()
	if err != nil {
		return fmt.Errorf("beep boop, something went wrong - is that a you or me problem.. ?: %s", err)
	}

	// Create catalogue mapping, will be used as the input value in when showing the catalogue
	catalogue := make(map[int]domain.Product, len(products))
	for i, product := range products {
		catalogue[i+1] = product
	}

	// Print catalogue
	sortedProducts := slices.Sorted(maps.Keys(catalogue))
	for _, key := range sortedProducts {
		item := catalogue[key]
		fmt.Printf("> %v: %s, %v$\n", key, item.Name, item.Price)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var selection int

	for {
		scanner.Scan()
		_, err = fmt.Sscanf(scanner.Text(), "%d", &selection)
		if err == nil {
			if _, exists := catalogue[selection]; exists {
				break
			}
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	var quantity int
	fmt.Println("> Nice choice! How many do you want?")
	for {
		scanner.Scan()
		_, err = fmt.Sscanf(scanner.Text(), "%d", &quantity)
		if err == nil {
			break
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	fmt.Printf("> Nice, you added: %v, to the cart\n", quantity)

	// Store selected into cart

	// User can check out
	// Should see their can and get a proceed confirmation selection
	// This should trigger events for
	// - Create transaction
	// - Update inventory
	// - Create shipping - create some sort of tracking ID, that holds procees of shipment status

	return nil
}
