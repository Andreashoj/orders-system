package repos

import (
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

type ProductRepo interface {
	GetAll() ([]domain.Product, error)
}

type productRepo struct {
	DB *sql.DB
}

func NewProductRepo(DB *sql.DB) ProductRepo {
	return productRepo{
		DB: DB,
	}
}

func (p productRepo) GetAll() ([]domain.Product, error) {
	rows, err := p.DB.Query(`SELECT id, name, price FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving products: %s", err)
	}
	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed scanning product: %s", err)
		}

		products = append(products, product)
	}

	return products, nil
}
