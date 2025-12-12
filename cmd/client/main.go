package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/andreashoj/order-system/internal/commands"
	"github.com/andreashoj/order-system/internal/db"
	"github.com/andreashoj/order-system/internal/db/models"
	"github.com/andreashoj/order-system/internal/pubsub"
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

	// Declare services
	productService := services.StartNewProductService(DB)

	// Create products
	// Create cart
	// Allow adding product(s) to cart

	commands.WelcomeMessage()

	for {
		cmd := commands.GetInput()

		switch cmd {
		case commands.Catalogue:
			products, err := productService.Get()
			if err != nil {
				fmt.Println("Beep boop, something went wrong - is that a you or me problem.. ?")
				break
			}

			// Create catalogue mapping
			catalogue := make(map[int]models.Product, len(products))
			for i, product := range products {
				catalogue[i+1] = product
			}

			// Print catalogue
			sortedProducts := slices.Sorted(maps.Keys(catalogue))
			for _, key := range sortedProducts {
				item := catalogue[key]
				fmt.Printf("> %v: %s, %v$\n", key, item.Name, item.Price)
			}

			// Get user input
			// Prompt quantity ?
			// Store selected into cart

			// Remember to prompt user for username and create user

			// Show catalogue
			// How do we do that?
			// Service => repo => products
			// Pass products to catalogue
			// User can add product to inventory upon selection
			// Nothing happens at this point
			// User can checkout
			//
		case commands.Cart:
		case commands.Checkout:

		}
	}
}
