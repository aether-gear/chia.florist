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

func (pM *PaymentMethod) Validate() error {
	if pM.Name == "" {
		return ErrInvalidName
	}

	if !pM.Type.IsValid() {
		return ErrInvalidType
	}

	if !pM.FeeType.IsValid() {
		return ErrInvalidFeeType
	}

	switch pM.FeeType {
	case FeeTypeFlat:
		if pM.FeeFixed < 0 {
			return ErrInvalidFeeFixed
		}

	case FeeTypePercentage:
		if pM.FeePercentage < 0 || pM.FeePercentage > 1 {
			return ErrInvalidFeePercentage
		}

	case FeeTypeMixed:
		if pM.FeeFixed < 0 {
			return ErrInvalidFeeFixed
		}
		if pM.FeePercentage < 0 || pM.FeePercentage > 1 {
			return ErrInvalidFeePercentage
		}
	}

	return nil
}

func (t PaymentMethodType) IsValid() bool {
	switch t {
	case TypeBankTransfer, TypeEWallet, TypeQRCode:
		return true
	default:
		return false
	}
}

func (t PaymentFeeType) IsValid() bool {
	switch t {
	case FeeTypeFlat, FeeTypePercentage, FeeTypeMixed:
		return true
	default:
		return false
	}
}
