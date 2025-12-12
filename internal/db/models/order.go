package models

import "time"

type Order struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	Complete    bool      `json:"complete,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}
