package usecase

import (
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
)

type ListPaymentAccountUsecase struct {
	paymentAccRepo repository.PaymentAccountRepository
}

func NewListPaymentAccountUsecase(
	paymentAccRepo repository.PaymentAccountRepository,
) *ListPaymentAccountUsecase {
	return &ListPaymentAccountUsecase{
		paymentAccRepo: paymentAccRepo,
	}
}

func (u *ListPaymentAccountUsecase) ListAll() ([]domain.PaymentAccount, error) {
	accounts, err := u.paymentAccRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load payment accounts: %w", err)
	}

	return accounts, nil
}
