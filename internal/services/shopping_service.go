package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type ShoppingService struct {
	productRepo repos.ProductRepo
	cartRepo    repos.CartRepo
	orderRepo   repos.OrderRepo
}

func NewShoppingService(productRepo repos.ProductRepo, cartRepo repos.CartRepo, orderRepo repos.OrderRepo) *ShoppingService {
	return &ShoppingService{
		productRepo: productRepo,
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
	}
}

func (c *ShoppingService) GetAllProducts() ([]domain.Product, error) {
	products, err := c.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed getting products: %s", err)
	}

	return products, nil
}

func (c *ShoppingService) AddToCart(userID, productID string, quantity int) (*domain.Cart, error) {
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

func (c *ShoppingService) GetCart(userID string) (*domain.Cart, error) {
	cart, err := c.cartRepo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting cart from repo: %s", err)
	}

	return cart, nil
}

func (c *ShoppingService) CreateOrder(userID string) (*domain.Order, error) {
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
