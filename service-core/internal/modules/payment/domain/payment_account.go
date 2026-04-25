package domain

import (
	"fmt"
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
			return fmt.Errorf("invalid payment account: account number is required")
		}
		if pA.AccountName == "" {
			return fmt.Errorf("invalid payment account: account name is required")
		}

	case TypeEWallet:
		if pA.PhoneNumber == "" {
			return fmt.Errorf("invalid payment account: phone number is required")
		}

	case TypeQRCode:
		if pA.QRString == nil || *pA.QRString == "" {
			return fmt.Errorf("invalid payment account: qr string is required")
		}

	default:
		return fmt.Errorf("invalid payment account: unsupported payment method type")
	}

	return nil
}
