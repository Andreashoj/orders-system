package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type CatalogueService struct {
	repo repos.ProductRepo
}

func NewCatalogueService(productRepo repos.ProductRepo) *CatalogueService {
	return &CatalogueService{
		repo: productRepo,
	}
}

func (c *CatalogueService) GetAllProducts() ([]domain.Product, error) {
	products, err := c.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed getting products: %s", err)
	}

	return products, nil
}
