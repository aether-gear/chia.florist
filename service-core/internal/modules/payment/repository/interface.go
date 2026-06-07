package repository

import (
	"context"

	"service-core/internal/modules/payment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type PaymentMethodRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		method domain.PaymentMethod,
	) error

	FindByName(
		ctx context.Context,
		exec transaction.Executor,
		name string,
	) (*domain.PaymentMethod, error)

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.PaymentMethod, error)

	ListAll(
		ctx context.Context,
		exec transaction.Executor,
	) ([]domain.PaymentMethod, error)
}

type PaymentAccountRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		paymentAccount domain.PaymentAccount,
	) error

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		paymentID uuid.UUID,
	) (*domain.PaymentAccount, error)

	AcquireLeastLoaded(
		ctx context.Context,
		exec transaction.Executor,
		methodID uuid.UUID,
	) (*domain.PaymentAccount, error)

	IncrementLoad(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) error
	DecrementLoad(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) error

	ListByMethodID(
		ctx context.Context,
		exec transaction.Executor,
		methodID uuid.UUID,
	) ([]domain.PaymentAccount, error)
	ListAll(
		ctx context.Context,
		exec transaction.Executor,
	) ([]domain.PaymentAccount, error)
}
