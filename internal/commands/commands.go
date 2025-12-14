package commands

import "fmt"

type Command string

const (
	Catalogue Command = "1"
	Cart      Command = "2"
	Checkout  Command = "3"
)

func WelcomeMessage() (string, error) {
	fmt.Println("Welcome to the online CLI shopper, only chads allowed here. \nWhat do you prefer to be addressed as?")
	var username string
	_, err := fmt.Scan(&username)
	if err != nil {
		return "", fmt.Errorf("failed getting username: %s", err)
	}

	return username, nil
}

func GetInput() Command {
	fmt.Println(">1: To see catalogue")
	fmt.Println(">2: See cart")
	fmt.Println(">3: Checkout")
	fmt.Println(">Q: Exit")

	var input Command
	fmt.Scan(&input)

	return input
}

func GetCatalogue() {
	// Get DB products here
	fmt.Println(">YOOOO")
	fmt.Println(">YOOOO")
	fmt.Println(">YOOOO")
	fmt.Println(">YOOOO")
}
