package repository

import (
	"context"

	"service-core/internal/modules/cart/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CartRepository interface {
	GetWithItemsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*domain.Cart, error)

	NewCart(
		ctx context.Context,
		userID uuid.UUID,
	) (*domain.Cart, error)

	Save(
		ctx context.Context,
		exec transaction.Executor,
		cart *domain.Cart,
	) error
}
