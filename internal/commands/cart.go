package commands

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

func DisplayCart(cart *domain.Cart) {
	var totalItems int
	var totalPrice int

	fmt.Println("------------------------------------")
	for i, item := range cart.Items {
		totalPrice += item.Quantity * item.Price
		totalItems += item.Quantity
		fmt.Printf("> %v: %s $%v - x%v\n", i+1, item.Name, item.Price, item.Quantity)
	}
	fmt.Println("------------------------------------")

	// Display total price and item amount
	fmt.Printf("> You have: %v items\n", totalItems)
	fmt.Printf("> The total cost of the items in your cart: %v$\n", totalPrice)
}
