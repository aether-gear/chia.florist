package repository

import (
	"context"

	"service-core/internal/modules/payment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Payment, error)

	GetByOrderID(
		ctx context.Context,
		exec transaction.Executor,
		orderID uuid.UUID,
	) (*domain.Payment, error)

	ListByOrderIDs(
		ctx context.Context,
		exec transaction.Executor,
		orderIDs []uuid.UUID,
	) ([]domain.Payment, error)

	UpdateStatus(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
		status domain.PaymentStatus,
	) error

	Save(
		ctx context.Context,
		exec transaction.Executor,
		payment domain.Payment,
	) error
}

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

	// RetrieveLeastLoaded returns the active account with the
	// lowest current load for the specified payment method
	RetrieveLeastLoaded(
		ctx context.Context,
		exec transaction.Executor,
		methodID uuid.UUID,
	) (*domain.PaymentAccount, error)

	// IncrementLoad increases the account's current load after
	// it has been assigned to a payment
	IncrementLoad(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) error

	// DecrementLoad decreases the account's current load when
	// the payment assignment is released or no longer active
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

type PaymentEventRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.PaymentEvent, error)

	ListByPaymentID(
		ctx context.Context,
		exec transaction.Executor,
		paymentID uuid.UUID,
	) ([]domain.PaymentEvent, error)

	Create(
		ctx context.Context,
		exec transaction.Executor,
		event domain.PaymentEvent,
	) error
}

type PaymentInstructionRepository interface {
	GetByPaymentMethodID(
		ctx context.Context,
		exec transaction.Executor,
		methodID uuid.UUID,
	) (*domain.PaymentInstruction, error)

	Save(
		ctx context.Context,
		exec transaction.Executor,
		instruction domain.PaymentInstruction,
	) error
}
