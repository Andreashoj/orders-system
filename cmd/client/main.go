package main

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/commands"
	"github.com/andreashoj/order-system/internal/db"
)

func main() {
	_, err := db.StartNewDB()
	if err != nil {
		fmt.Printf("Failed starting the DB: %s", err)
		return
	}

	commands.WelcomeMessage()

	for {
		cmd := commands.GetInput()

		switch cmd {
		case commands.Catalogue:
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
