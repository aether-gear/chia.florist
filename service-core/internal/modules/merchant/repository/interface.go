package repository

import (
	"context"

	"service-core/internal/modules/merchant/domain"

	"github.com/google/uuid"
)

type MerchantRepository interface {
	GetByAccountID(
		ctx context.Context,
		accountID uuid.UUID,
	) (*domain.Merchant, error)
}
