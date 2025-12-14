package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type RegistrationService struct {
	repo repos.UserRepo
}

func NewRegistrationService(userRepo repos.UserRepo) *RegistrationService {
	return &RegistrationService{
		repo: userRepo,
	}
}

func (r *RegistrationService) CreateUser(username string) (*domain.User, error) {
	user, err := domain.NewUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed creating domain user: %s", err)
	}

	if err = r.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed creating the user in the database: %s", err)
	}

	return user, nil
}
