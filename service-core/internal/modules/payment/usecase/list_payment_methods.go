package usecase

import (
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewListPaymentMethodUsecase(
	paymentMethodRepo repository.PaymentMethodRepository,
) *ListPaymentMethodUsecase {
	return &ListPaymentMethodUsecase{
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (u *ListPaymentMethodUsecase) ListAll() ([]domain.PaymentMethod, error) {
	methods, err := u.paymentMethodRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load payment methods: %w", err)
	}

	return methods, nil
}
