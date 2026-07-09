package repository

import (
	"context"

	"service-core/internal/modules/customer/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(
		ctx context.Context,
		exec transaction.Executor,
		customer domain.Customer,
	) error

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Customer, error)

	GetByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*domain.Customer, error)

	GetProfileByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*domain.CustomerProfile, error)

	FindCustomers(
		ctx context.Context,
		exec transaction.Executor,
		params FindCustomerParams,
	) ([]domain.CustomerProfile, int, error)

	Delete(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) error
}
