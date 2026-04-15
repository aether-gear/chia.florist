package repository

import (
	"service-core/internal/features/cart/domain"

	"github.com/google/uuid"
)

type CartWithItems struct {
	*domain.Cart
}

type CartRepository interface {
	GetWithItemsByUserID(userID uuid.UUID) (*domain.Cart, error)
	NewCart(userID uuid.UUID) (*domain.Cart, error)
	Save(cart *domain.Cart) error
}
