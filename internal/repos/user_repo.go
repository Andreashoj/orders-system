package repos

import (
	"database/sql"
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
)

type UserRepo interface {
	Get(userID string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
}

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(DB *sql.DB) *userRepo {
	return &userRepo{
		DB: DB,
	}
}

func (u *userRepo) Get(userID string) (*domain.User, error) {
	var user domain.User
	err := u.DB.QueryRow(`SELECT id, name, balance, created_at FROM users WHERE id = $1`, userID).Scan(&user.ID, &user.Name, &user.Balance, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %s", err)
	}

	return &user, nil
}

func (u *userRepo) Create(user *domain.User) error {
	_, err := u.DB.Exec(`INSERT INTO users (id, name, balance, created_at) VALUES ($1, $2, $3, $4)`, user.ID, user.Name, user.Balance, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed creating user: %s", err)
	}

	return nil
}

func (u *userRepo) Update(user *domain.User) error {
	_, err := u.DB.Exec(`UPDATE users SET name = $1, balance = $2 WHERE id = $3`, user.Name, user.Balance, user.ID)
	if err != nil {
		return fmt.Errorf("failed updating user: %s", err)
	}

	return nil
}
