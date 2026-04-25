package usecase

import (
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentMethod struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewListPaymentMethod(
	pMR repository.PaymentMethodRepository,
) *ListPaymentMethod {
	return &ListPaymentMethod{
		paymentMethodRepo: pMR,
	}
}

func (u *ListPaymentMethod) ListAll() ([]domain.PaymentMethod, error) {
	methods, err := u.paymentMethodRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load payment methods: %w", err)
	}

	return methods, nil
}
