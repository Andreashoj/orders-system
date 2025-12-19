package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          string `json:"id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Complete    bool   `json:"complete,omitempty"`
	Items       []OrderItem
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItem struct {
	ProductID string    `json:"product_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Price     int       `json:"price,omitempty"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

func NewOrder(userID string) *Order {
	return &Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Complete:  false,
		CreatedAt: time.Now(),
	}
}

func (o *Order) AddCart(cart *Cart) {
	var items []OrderItem
	for _, item := range cart.Items {
		orderItem := OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			CreatedAt: time.Now(),
		}
		items = append(items, orderItem)
	}
	o.Items = items
}

func (o *Order) GetTotal() int {
	sum := 0
	for _, item := range o.Items {
		sum += item.Price * item.Quantity
	}
	return sum
}
