package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"
)

type ListPaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
	executor          transaction.Executor
}

func NewListPaymentMethodUsecase(
	paymentMethodRepo repository.PaymentMethodRepository,
	executor transaction.Executor,
) *ListPaymentMethodUsecase {
	return &ListPaymentMethodUsecase{
		paymentMethodRepo: paymentMethodRepo,
		executor:          executor,
	}
}

func (u *ListPaymentMethodUsecase) ListAll(ctx context.Context) ([]domain.PaymentMethod, error) {
	methods, err := u.paymentMethodRepo.ListAll(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment methods: %w", err)
	}

	return methods, nil
}
