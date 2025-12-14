package domain

import "time"

type CartProduct struct {
	CardID    int       `json:"card_id,omitempty"`
	ProductID int       `json:"product_id,omitempty"`
	Quantity  int       `json:"quantity,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
