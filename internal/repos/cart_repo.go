package repos

import "github.com/andreashoj/order-system/internal/domain"

type CartRepo interface {
	Create() error
	Update(cart *domain.Cart) error
}
