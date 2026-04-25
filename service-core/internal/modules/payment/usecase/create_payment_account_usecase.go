package usecase

import (
	"fmt"
	"time"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type CreatePaymentAccount struct {
	paymentMethodRepo  repository.PaymentMethodRepository
	paymentAccountRepo repository.PaymentAccountRepository
}

func NewCreatePaymentAccount(
	pAR repository.PaymentAccountRepository,
	pMR repository.PaymentMethodRepository,
) *CreatePaymentAccount {
	return &CreatePaymentAccount{
		paymentMethodRepo:  pMR,
		paymentAccountRepo: pAR,
	}
}

func (u *CreatePaymentAccount) Execute(input CreatePaymentAccountInput) error {
	method, err := u.paymentMethodRepo.GetByID(input.MethodID)
	if err != nil {
		return fmt.Errorf("failed to retrieve payment account: %w", err)
	}
	if method == nil {
		return fmt.Errorf("failed to create payment account: payment method not found")
	}

	paymentAccount := domain.PaymentAccount{
		ID:            uuid.New(),
		MethodID:      method.ID,
		AccountName:   input.AccountName,
		AccountNumber: input.AccountNumber,
		PhoneNumber:   input.PhoneNumber,
		QRString:      input.QRString,
		IsActive:      input.IsActive,
		CurrentLoad:   0,
		CreatedAt:     time.Now(),
	}

	if err := paymentAccount.ValidateForMethod(method.Type); err != nil {
		return fmt.Errorf("failed to create payment account: %w", err)
	}

	err = u.paymentAccountRepo.Save(paymentAccount)
	if err != nil {
		return fmt.Errorf("failed to save payment account: %w", err)
	}

	return nil
}

type CreatePaymentAccountInput struct {
	MethodID      uuid.UUID
	AccountName   string
	AccountNumber *string
	PhoneNumber   string
	QRString      *string
	IsActive      bool
}
