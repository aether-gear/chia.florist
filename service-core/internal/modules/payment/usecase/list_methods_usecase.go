package usecase

import (
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentMethod struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewListPaymentMethod(pMR repository.PaymentMethodRepository) *ListPaymentMethod {
	return &ListPaymentMethod{
		paymentMethodRepo: pMR,
	}
}

func (u *ListPaymentMethod) ListAll() ([]domain.PaymentMethod, error) {
	return u.paymentMethodRepo.ListAll()
}
