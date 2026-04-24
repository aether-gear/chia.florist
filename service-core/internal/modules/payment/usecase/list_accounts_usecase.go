package usecase

import (
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentAccount struct {
	paymentAccRepo repository.PaymentAccountRepository
}

func NewListPaymentAccount(pAR repository.PaymentAccountRepository) *ListPaymentAccount {
	return &ListPaymentAccount{
		paymentAccRepo: pAR,
	}
}

func (u *ListPaymentAccount) ListAll() ([]domain.PaymentAccount, error) {
	return u.paymentAccRepo.ListAll()
}
