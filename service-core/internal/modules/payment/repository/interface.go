package repository

import (
	"context"

	"service-core/internal/modules/payment/domain"
	query "service-core/internal/shared/query"
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

var (
	PaymentMethodSortLatest query.SortKey = "latest"
	PaymentMethodSortName   query.SortKey = "name"
	PaymentMethodSortCode   query.SortKey = "code"
	PaymentMethodSortType   query.SortKey = "type"
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
		sorts query.Sorts,
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

// PaymentChannelDataRepository persists gateway-returned payment channel
// details (QR strings, VA numbers, deep links) so they can be retrieved
// after the initial checkout response is discarded.
type PaymentChannelDataRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		data domain.PaymentChannelData,
	) error

	GetByPaymentID(
		ctx context.Context,
		exec transaction.Executor,
		paymentID uuid.UUID,
	) (*domain.PaymentChannelData, error)

	// ListByPaymentIDs returns channel data records indexed by payment ID.
	ListByPaymentIDs(
		ctx context.Context,
		exec transaction.Executor,
		paymentIDs []uuid.UUID,
	) (map[uuid.UUID]*domain.PaymentChannelData, error)
}
