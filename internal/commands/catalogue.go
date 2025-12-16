package commands

import (
	"fmt"
	"maps"
	"slices"

	"github.com/andreashoj/order-system/internal/domain"
)

func DisplayCatalogue(catalogue map[int]domain.Product) {
	sortedProducts := slices.Sorted(maps.Keys(catalogue))
	for _, key := range sortedProducts {
		item := catalogue[key]
		fmt.Printf("> %v: %s, %v$\n", key, item.Name, item.Price)
	}
}

func GetProductSelection(catalogue map[int]domain.Product) *domain.Product {
	var selection int
	var product domain.Product
	for {
		scanner.Scan()
		_, err := fmt.Sscanf(scanner.Text(), "%d", &selection)
		if err == nil {
			if _, exists := catalogue[selection]; exists {
				product = catalogue[selection]
				break
			}
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	return &product
}

func GetProductQuantity() int {
	var quantity int
	fmt.Println("> Nice choice! How many do you want?")
	for {
		scanner.Scan()
		_, err := fmt.Sscanf(scanner.Text(), "%d", &quantity)
		if err == nil {
			break
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	fmt.Printf("> Nice, you added: %v, to the cart\n", quantity)
	return quantity
}

func PromptCheckout() bool {
	var wantsToCheckout bool
	fmt.Println("> Do you wish to checkout (Y)es? or continue shopping (N)o?")
	for {
		scanner.Scan()
		input := scanner.Text()
		if input == "Y" || input == "y" {
			return true
		} else if input == "N" || input == "n" {
			return false
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	return wantsToCheckout
}
