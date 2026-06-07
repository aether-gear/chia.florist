package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"
)

type ListPaymentAccountUsecase struct {
	paymentAccRepo repository.PaymentAccountRepository
	executor       transaction.Executor
}

func NewListPaymentAccountUsecase(
	paymentAccRepo repository.PaymentAccountRepository,
	executor transaction.Executor,
) *ListPaymentAccountUsecase {
	return &ListPaymentAccountUsecase{
		paymentAccRepo: paymentAccRepo,
		executor:       executor,
	}
}

func (u *ListPaymentAccountUsecase) ListAll(ctx context.Context) ([]domain.PaymentAccount, error) {
	accounts, err := u.paymentAccRepo.ListAll(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment accounts: %w", err)
	}

	return accounts, nil
}
