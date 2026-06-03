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
		method domain.PaymentMethod,
	) error

	FindByName(
		ctx context.Context,
		name string,
	) (*domain.PaymentMethod, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.PaymentMethod, error)

	ListAll(
		ctx context.Context,
	) ([]domain.PaymentMethod, error)
}

type PaymentAccountRepository interface {
	Save(
		ctx context.Context,
		paymentAccount domain.PaymentAccount,
	) error

	GetByID(
		ctx context.Context,
		paymentID uuid.UUID,
	) (*domain.PaymentAccount, error)

	AcquireLeastLoaded(
		ctx context.Context,
		exec transaction.Executor,
		methodID uuid.UUID,
	) (*domain.PaymentAccount, error)

	IncrementLoad(
		ctx context.Context,
		accountID uuid.UUID,
	) error
	DecrementLoad(
		ctx context.Context,
		accountID uuid.UUID,
	) error

	ListByMethodID(
		ctx context.Context,
		methodID uuid.UUID,
	) ([]domain.PaymentAccount, error)
	ListAll(
		ctx context.Context,
	) ([]domain.PaymentAccount, error)
}
