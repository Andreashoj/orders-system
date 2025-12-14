package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username string) (*User, error) {
	if len(username) <= 3 {
		return nil, fmt.Errorf("Username must be longer than 3 characters")
	}

	return &User{
		ID:        uuid.New().String(),
		Name:      username,
		CreatedAt: time.Now(),
	}, nil
}
