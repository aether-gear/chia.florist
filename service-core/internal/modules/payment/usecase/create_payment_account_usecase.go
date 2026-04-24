package usecase

import (
	"errors"
	"fmt"
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
	pAR repository.PaymentAccountRepository,
	pMR repository.PaymentMethodRepository,
) *CreatePaymentAccount {
	return &CreatePaymentAccount{
		paymentMethodRepo:  pMR,
		paymentAccountRepo: pAR,
	}
}

func (u *CreatePaymentAccount) Execute(input CreatePaymentAccountInput) error {
	fmt.Println("kdlneaopdeadpaedpiay")
	method, err := u.paymentMethodRepo.GetByID(input.MethodID)
	if err != nil {
		return err
	}
	if method == nil {
		return errors.New("payment method not found")
	}
	fmt.Println("kdlneaopdeadpaedpiay")

	switch method.Type {

	case domain.TypeBankTransfer:
		if input.AccountNumber == nil || *input.AccountNumber == "" {
			return errors.New("account number required")
		}
		if input.AccountName == "" {
			return errors.New("account name required")
		}

	case domain.TypeEWallet:
		if input.PhoneNumber == "" {
			return errors.New("phone number required")
		}

	case domain.TypeQRCode:
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
		CurrentLoad:   0,
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
