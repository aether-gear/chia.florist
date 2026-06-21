package repository

import (
	"context"

	transaction "service-core/internal/shared/transaction"
)

type PricingService interface {
	Calculate(
		ctx context.Context,
		exec transaction.Executor,
		input PricingInput,
	) (*PricingResult, error)
}
