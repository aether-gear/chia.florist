package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethodType string
type PaymentFeeType string

const (
	TypeBankTransfer PaymentMethodType = "bank_transfer"
	TypeEWallet      PaymentMethodType = "ewallet"
	TypeQRCode       PaymentMethodType = "qr_code"

	FeeTypeFlat       PaymentFeeType = "flat"
	FeeTypePercentage PaymentFeeType = "percentage"
	FeeTypeMixed      PaymentFeeType = "mixed"
)

type PaymentMethod struct {
	ID uuid.UUID

	Name        string
	Type        PaymentMethodType
	IsActive    bool
	Description string

	FeeType       PaymentFeeType
	FeeFixed      int64
	FeePercentage float64

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
