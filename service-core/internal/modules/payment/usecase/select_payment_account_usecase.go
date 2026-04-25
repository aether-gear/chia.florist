package usecase

import (
	"fmt"
	"time"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type SelectPayment struct {
	paymentAccountRepo repository.PaymentAccountRepository
}

func NewSelectPayment(
	paymentAccountRepo repository.PaymentAccountRepository,
) *SelectPayment {
	return &SelectPayment{
		paymentAccountRepo: paymentAccountRepo,
	}
}

func (u *SelectPayment) Execute(methodID uuid.UUID) (*domain.PaymentAccount, error) {
	accounts, err := u.paymentAccountRepo.ListByMethodID(methodID)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment methods: %w", err)
	}
	if len(accounts) == 0 {
		return nil, fmt.Errorf("failed to load payment methods: no payment accounts active")
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

	if err := u.paymentAccountRepo.Save(selected); err != nil {
		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	return &selected, nil
}
