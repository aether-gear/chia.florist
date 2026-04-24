package usecase

import (
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type AllocatePaymentByMethod struct {
	paymentMethodRepo  repository.PaymentMethodRepository
	paymentAccountRepo repository.PaymentAccountRepository
}

func NewAllocatePaymentByMethod(
	pMR repository.PaymentMethodRepository,
	pAR repository.PaymentAccountRepository,
) *AllocatePaymentByMethod {
	return &AllocatePaymentByMethod{
		paymentMethodRepo:  pMR,
		paymentAccountRepo: pAR,
	}
}

func (u *AllocatePaymentByMethod) Execute(methodID uuid.UUID) error {
	return nil
}

type AllocatePaymentByMethodResponse struct {
	AccountID     uuid.UUID
	MethodID      uuid.UUID
	AccountName   string
	AccountNumber *string
	PhoneNumber   *string
	QRString      *string
}
