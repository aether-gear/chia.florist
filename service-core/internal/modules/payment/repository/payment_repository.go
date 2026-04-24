package repository

import (
	"service-core/internal/modules/payment/domain"

	"github.com/google/uuid"
)

type PaymentMethodRepository interface {
	Save(method domain.PaymentMethod) error
	FindByName(name string) (*domain.PaymentMethod, error)
	GetByID(id uuid.UUID) (*domain.PaymentMethod, error)
	ListAll() ([]domain.PaymentMethod, error)
}

type PaymentAccountRepository interface {
	Save(paymentAccount domain.PaymentAccount) error
	GetByID(paymentID uuid.UUID) (*domain.PaymentAccount, error)

	AcquireLeastLoaded(methodID uuid.UUID) (*domain.PaymentAccount, error)

	IncrementLoad(accountID uuid.UUID) error
	DecrementLoad(accountID uuid.UUID) error

	ListByMethodID(methodID uuid.UUID) ([]domain.PaymentAccount, error)
	ListAll() ([]domain.PaymentAccount, error)
}
