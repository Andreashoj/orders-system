package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type ShoppingService struct {
	productRepo repos.ProductRepo
	cartRepo    repos.CartRepo
}

func NewShoppingService(productRepo repos.ProductRepo, cartRepo repos.CartRepo) *ShoppingService {
	return &ShoppingService{
		productRepo: productRepo,
		cartRepo:    cartRepo,
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
