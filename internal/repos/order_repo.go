package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

type OrderRepo interface {
	Create(order *domain.Order) error
	Get(userID string) (*domain.Order, error) // change param type if userID doesn't make sense..
}

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(DB *sql.DB) OrderRepo {
	return &orderRepo{
		DB: DB,
	}
}

func (o orderRepo) Create(order *domain.Order) error {
	tx, err := o.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed commiting order db transaction: %s", err)
	}

	_, err = tx.Exec(`INSERT INTO orders (id, user_id, complete, completed_at, created_at) VALUES ($1, $2, $3, $4, $5)`, order.ID, order.UserID, order.Complete, order.CompletedAt, order.CreatedAt)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed creating order: %s", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`INSERT INTO order_products (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed inserting product: %s into table, got error: %s", item.ProductID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed commiting order db transaction: %s", err)
	}

	return nil
}

func (o orderRepo) Get(userID string) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}
