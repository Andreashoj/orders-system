package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type ShoppingService interface {
	GetAllProducts() ([]domain.Product, error)
	AddToCart(userID, productID string, quantity int) (*domain.Cart, error)
	GetCart(userID string) (*domain.Cart, error)
	CreateOrder(userID string) (*domain.Order, error)
	GetOrder(orderID string) (*domain.Order, error)
	ChargeUser(userID string, amount int) bool
}

type shoppingService struct {
	productRepo repos.ProductRepo
	cartRepo    repos.CartRepo
	orderRepo   repos.OrderRepo
	userRepo    repos.UserRepo
}

func NewShoppingService(productRepo repos.ProductRepo, cartRepo repos.CartRepo, orderRepo repos.OrderRepo, userRepo repos.UserRepo) ShoppingService {
	return &shoppingService{
		productRepo: productRepo,
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
		userRepo:    userRepo,
	}
}

func (c *shoppingService) GetAllProducts() ([]domain.Product, error) {
	products, err := c.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed getting products: %s", err)
	}

	return products, nil
}

func (c *shoppingService) AddToCart(userID, productID string, quantity int) (*domain.Cart, error) {
	cart, err := c.cartRepo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting the cart: %s", err)
	}

	product, err := c.productRepo.Get(productID)
	if err != nil {
		return nil, fmt.Errorf("failed getting the added product: %s", err)
	}

	cart.Add(product, quantity)
	err = c.cartRepo.Update(cart)
	if err != nil {
		return nil, fmt.Errorf("failed updating the cart: %s", err)
	}

	return cart, nil
}

func (c *shoppingService) GetCart(userID string) (*domain.Cart, error) {
	cart, err := c.cartRepo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting cart from repo: %s", err)
	}

	return cart, nil
}

func (c *shoppingService) CreateOrder(userID string) (*domain.Order, error) {
	order := domain.NewOrder(userID)
	cart, err := c.cartRepo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting cart: %s", err)
	}
	order.AddCart(cart)
	err = c.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed creating order: %s", err)
	}

	return order, nil
}

func (c *shoppingService) GetOrder(orderID string) (*domain.Order, error) {
	order, err := c.orderRepo.Get(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed getting order: %s", err)
	}

	return order, nil
}

func (c *shoppingService) ChargeUser(userID string, amount int) bool {
	user, err := c.userRepo.Get(userID)
	if err != nil {
		fmt.Printf("failed retrieving user: %s", err)
		return false
	}

	if user.Balance >= amount {
		user.Balance -= amount
		err = c.userRepo.Update(user)

		if err != nil {
			fmt.Errorf("failed updating user")
			return false
		}

		return true
	}

	return false
}
