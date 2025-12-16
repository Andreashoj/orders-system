package repos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andreashoj/order-system/internal/domain"
)

type CartRepo interface {
	Create(cart *domain.Cart) error
	Update(cart *domain.Cart) error
	Get(userID string) (*domain.Cart, error)
}

type cartRepo struct {
	DB *sql.DB
}

func NewCartRepo(DB *sql.DB) CartRepo {
	return &cartRepo{
		DB: DB,
	}
}

func (c *cartRepo) Create(cart *domain.Cart) error {
	_, err := c.DB.Exec(`INSERT INTO carts (id, user_id, last_updated) VALUES ($1, $2, $3)`, cart.ID, cart.UserID, cart.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed creating cart for user: %s", err)
	}
	return nil
}

func (c *cartRepo) Update(cart *domain.Cart) error {
	_, err := c.DB.Exec(`DELETE FROM cart_products WHERE cart_id = $1`, cart.ID)
	if err != nil {
		return fmt.Errorf("failed deleting cart items: %s", err)
	}

	for _, item := range cart.Items {
		_, err = c.DB.Exec(`INSERT INTO cart_products (cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4)`, cart.ID, item.ProductID, item.Quantity, time.Now())
		if err != nil {
			return fmt.Errorf("failed inserting product into cart: %s", err)
		}
	}

	return nil
}

func (c *cartRepo) Get(userID string) (*domain.Cart, error) {
	cart := domain.Cart{
		UserID: userID,
	}

	// Get cart id separately to ensure it's set
	if err := c.DB.QueryRow(`SELECT id FROM carts WHERE user_id = $1`, cart.UserID).Scan(&cart.ID); err != nil {
		return nil, fmt.Errorf("failed getting cart: %s", err)
	}

	// Get cart items
	query := `SELECT cart_products.quantity, products.id, products.name, products.price FROM cart_products
					INNER JOIN products
						ON products.id = cart_products.product_id
							WHERE cart_products.cart_id = $1`

	rows, err := c.DB.Query(query, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving cart for user: %s, error: %s", userID, err)
	}

	for rows.Next() {
		var cartItem domain.CartItem
		if err = rows.Scan(&cartItem.Quantity, &cartItem.ProductID, &cartItem.Name, &cartItem.Price); err != nil {
			return nil, fmt.Errorf("failed scanning cartItem: %s", err)
		}

		cart.Items = append(cart.Items, cartItem)
	}

	return &cart, nil
}
