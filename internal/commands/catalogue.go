package commands

import (
	"bufio"
	"fmt"
	"maps"
	"os"
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

func GetProductSelection(catalogue map[int]domain.Product) int {
	scanner := bufio.NewScanner(os.Stdin)
	var selection int
	for {
		scanner.Scan()
		_, err := fmt.Sscanf(scanner.Text(), "%d", &selection)
		if err == nil {
			if _, exists := catalogue[selection]; exists {
				break
			}
		}

		fmt.Println("> Sorry that's not a valid selection, try again!")
	}

	return selection
}

func GetProductQuantity() int {
	scanner := bufio.NewScanner(os.Stdin)
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
