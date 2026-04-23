package usecase

import (
	"errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	"time"

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
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, errors.New("no active payment accounts")
	}

	selected := accounts[0]
	for _, acc := range accounts {
		if acc.CurrentLoad < selected.CurrentLoad {
			selected = acc
		}
	}

	selected.CurrentLoad += 1
	now := time.Now()
	selected.LastUsedAt = &now

	err = u.paymentAccountRepo.Save(selected)
	if err != nil {
		return nil, err
	}

	return &selected, nil
}
