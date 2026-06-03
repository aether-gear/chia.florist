package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type SelectPayment struct {
	paymentAccRepo repository.PaymentAccountRepository
}

func NewSelectPayment(
	paymentAccRepo repository.PaymentAccountRepository,
) *SelectPayment {
	return &SelectPayment{
		paymentAccRepo: paymentAccRepo,
	}
}

func (u *SelectPayment) Execute(ctx context.Context, methodID uuid.UUID) (*domain.PaymentAccount, error) {
	accounts, err := u.paymentAccRepo.ListByMethodID(ctx, methodID)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment methods: %w", err)
	}
	if len(accounts) == 0 {
		return nil, apperrors.NewNotFound(domain.ErrNoActivePaymentAccount.Error())
	}

	selected := accounts[0]
	for _, acc := range accounts {
		if acc.CurrentLoad < selected.CurrentLoad {
			selected = acc
		}
	}

	now := time.Now()
	selected.LastUsedAt = &now
	selected.CurrentLoad += 1

	if err := u.paymentAccRepo.Save(ctx, selected); err != nil {
		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	return &selected, nil
}
