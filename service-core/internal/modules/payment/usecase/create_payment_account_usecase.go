package usecase

import (
	"errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	"time"

	"github.com/google/uuid"
)

type CreatePaymentAccount struct {
	paymentMethodRepo  repository.PaymentMethodRepository
	paymentAccountRepo repository.PaymentAccountRepository
}

func NewCreatePaymentAccount(
	paymentMethodRepo repository.PaymentMethodRepository,
	paymentAccountRepo repository.PaymentAccountRepository,
) *CreatePaymentAccount {
	return &CreatePaymentAccount{
		paymentMethodRepo:  paymentMethodRepo,
		paymentAccountRepo: paymentAccountRepo,
	}
}

func (u *CreatePaymentAccount) Execute(input CreatePaymentAccountInput) error {
	method, err := u.paymentMethodRepo.GetByID(input.MethodID)
	if err != nil {
		return err
	}
	if method == nil {
		return errors.New("payment method not found")
	}

	switch method.Type {

	case string(domain.TypeBankTransfer):
		if input.AccountNumber == nil || *input.AccountNumber == "" {
			return errors.New("account number required")
		}
		if input.AccountName == "" {
			return errors.New("account name required")
		}

	case string(domain.TypeEWallet):
		if input.PhoneNumber == "" {
			return errors.New("phone number required")
		}

	case string(domain.TypeQRCode):
		if input.QRString == nil || *input.QRString == "" {
			return errors.New("qr string required")
		}

	default:
		return errors.New("unsupported payment type")
	}

	paymentAccount := domain.PaymentAccount{
		ID:            uuid.New(),
		MethodID:      method.ID,
		AccountName:   input.AccountName,
		AccountNumber: input.AccountNumber,
		PhoneNumber:   input.PhoneNumber,
		QRString:      input.QRString,
		IsActive:      input.IsActive,
		CreatedAt:     time.Now(),
	}

	err = u.paymentAccountRepo.Save(paymentAccount)
	if err != nil {
		return err
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
