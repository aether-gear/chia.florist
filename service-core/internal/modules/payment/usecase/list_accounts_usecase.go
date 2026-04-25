package usecase

import (
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentAccount struct {
	paymentAccRepo repository.PaymentAccountRepository
}

func NewListPaymentAccount(
	pAR repository.PaymentAccountRepository,
) *ListPaymentAccount {
	return &ListPaymentAccount{
		paymentAccRepo: pAR,
	}
}

func (u *ListPaymentAccount) ListAll() ([]domain.PaymentAccount, error) {
	accounts, err := u.paymentAccRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load payment accounts: %w", err)
	}

	return accounts, nil
}
