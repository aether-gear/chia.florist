package repository

import (
	"context"

	"service-core/internal/modules/merchant/domain"
	transaction "service-core/internal/shared/transaction"
)

type MerchantRepository interface {
	Create(
		ctx context.Context,
		exec transaction.Executor,
		merchant domain.Merchant,
	) error
}
