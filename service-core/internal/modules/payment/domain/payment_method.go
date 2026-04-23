package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethodType string

const (
	TypeBankTransfer PaymentMethodType = "BANK_TRANSFER"
	TypeEWallet      PaymentMethodType = "EWALLET"
	TypeQRCode       PaymentMethodType = "QR_CODE"
)

type PaymentMethod struct {
	ID          uuid.UUID
	Name        string
	Type        string
	IsActive    bool
	Description string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
