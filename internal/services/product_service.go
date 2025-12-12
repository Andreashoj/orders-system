package services

import (
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/db/models"
)

type ProductService interface {
	Get() ([]models.Product, error)
}

type productService struct {
	DB *sql.DB
}

func StartNewProductService(DB *sql.DB) ProductService {
	return &productService{
		DB: DB,
	}
}

func (p *productService) Get() ([]models.Product, error) {
	rows, err := p.DB.Query(`SELECT id, name, price FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving products: %s", err)
	}
	var products []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed scanning product: %s", err)
		}

		products = append(products, product)
	}

	return products, nil
}
