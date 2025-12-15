package services

import (
	"fmt"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/repos"
)

type RegistrationService struct {
	userRepo repos.UserRepo
	cartRepo repos.CartRepo
}

func NewRegistrationService(userRepo repos.UserRepo, cartRepo repos.CartRepo) *RegistrationService {
	return &RegistrationService{
		userRepo: userRepo,
		cartRepo: cartRepo,
	}
}

func (r *RegistrationService) CreateUser(username string) (*domain.User, error) {
	user, err := domain.NewUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed creating domain user: %s", err)
	}

	if err = r.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed creating the user in the database: %s", err)
	}

	cart := domain.NewCart(user)
	if err = r.cartRepo.Create(cart); err != nil {
		return nil, fmt.Errorf("failed creating cart: %s", err)
	}

	return user, nil
}
