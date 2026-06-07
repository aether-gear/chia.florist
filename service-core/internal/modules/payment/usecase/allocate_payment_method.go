package usecase

import (
	"context"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type AllocatePaymentByMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
	paymentAccRepo    repository.PaymentAccountRepository
}

func NewAllocatePaymentByMethodUsecase(
	paymentMethodRepo repository.PaymentMethodRepository,
	paymentAccRepo repository.PaymentAccountRepository,
) *AllocatePaymentByMethodUsecase {
	return &AllocatePaymentByMethodUsecase{
		paymentMethodRepo: paymentMethodRepo,
		paymentAccRepo:    paymentAccRepo,
	}
}

type AllocatePaymentByMethodResponse struct {
	AccountID     uuid.UUID
	MethodID      uuid.UUID
	AccountName   string
	AccountNumber *string
	PhoneNumber   *string
	QRString      *string
}

func (u *AllocatePaymentByMethodUsecase) Execute(ctx context.Context, methodID uuid.UUID) error {
	return nil
}
