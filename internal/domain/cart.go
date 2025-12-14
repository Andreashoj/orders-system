package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID          string     `json:"id,omitempty"`
	UserID      string     `json:"user_id,omitempty"`
	Items       []CartItem `json:"quantity,omitempty"` // Is most likely not needed
	LastUpdated time.Time  `json:"last_updated"`
}

type CartItem struct {
	ProductID string
	Quantity  int
}

func NewCart(user *User) *Cart {
	return &Cart{
		ID:          uuid.New().String(),
		UserID:      user.ID,
		Items:       []CartItem{},
		LastUpdated: time.Now(),
	}
}

func (c *Cart) UpdateCart() {

}
