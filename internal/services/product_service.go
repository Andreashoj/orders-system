package services

import "github.com/andreashoj/order-system/internal/db"

type ProductService interface {
	Get() ([]db.Product, error)
}

type productService struct {
}

func StartNewProductService() ProductService {
	return &productService{}
}
