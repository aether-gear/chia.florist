package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentAccount struct {
	ID       uuid.UUID
	MethodID uuid.UUID

	AccountName   string
	AccountNumber *string
	PhoneNumber   string
	QRString      *string

	IsActive bool

	CurrentLoad int
	LastUsedAt  *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (pA *PaymentAccount) ValidateForMethod(methodType PaymentMethodType) error {
	switch methodType {

	case TypeBankTransfer:
		if pA.AccountNumber == nil || *pA.AccountNumber == "" {
			return ErrInvalidAccountName
		}
		if pA.AccountName == "" {
			return ErrInvalidAccountNumber
		}

	case TypeEWallet:
		if pA.PhoneNumber == "" {
			return ErrInvalidPhoneNumber
		}

	case TypeQRCode:
		if pA.QRString == nil || *pA.QRString == "" {
			return ErrInvalidQRString
		}

	default:
		return ErrUnsupportedPaymentMethod
	}

	return nil
}
