package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

type OrderRepo interface {
	Create(order *domain.Order) error
	Get(orderID string) (*domain.Order, error) // change param type if userID doesn't make sense..
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

func (o orderRepo) Get(orderID string) (*domain.Order, error) {
	var order domain.Order
	err := o.DB.QueryRow(`SELECT id, user_id, complete, completed_at, created_at FROM orders`).Scan(&order.ID, &order.UserID, &order.Complete, &order.CompletedAt, &order.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving order: %s", err)
	}

	itemsQuery := `
		SELECT products.name, products.id, order_products.price, order_products.quantity FROM order_products
			LEFT JOIN products
				ON order_products.product_id = products.id
					WHERE order_products.order_id = $1
	`
	rows, err := o.DB.Query(itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed querying order items: %s", err)
	}

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		err = rows.Scan(&item.Name, &item.ProductID, &item.Price, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed scanning order item")
		}

		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}
