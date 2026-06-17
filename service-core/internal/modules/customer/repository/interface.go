package repository

import (
	"context"

	"service-core/internal/modules/customer/domain"
	transaction "service-core/internal/shared/transaction"
)

type CustomerRepository interface {
	FindCustomers(
		ctx context.Context,
		exec transaction.Executor,
		params FindCustomerParams,
	) ([]domain.Customer, int, error)
}
