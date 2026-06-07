package repository

import (
	"context"

	"service-core/internal/modules/merchant/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type MerchantRepository interface {
	GetByAccountID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) (*domain.Merchant, error)
}
