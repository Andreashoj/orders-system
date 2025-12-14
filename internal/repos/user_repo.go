package repos

import (
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

type UserRepo interface {
	Create(user *domain.User) error
}

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(DB *sql.DB) *userRepo {
	return &userRepo{
		DB: DB,
	}
}

func (u *userRepo) Create(user *domain.User) error {
	_, err := u.DB.Exec(`INSERT INTO users (id, name, created_at) VALUES ($1, $2, $3)`, user.ID, user.Name, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed creating user: %s", err)
	}

	return nil
}
