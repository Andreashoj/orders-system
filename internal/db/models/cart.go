package models

import "time"

type Cart struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	Quantity    int       `json:"quantity,omitempty"`
	LastUpdated time.Time `json:"last_updated"`
}
