package domain

import (
	"fmt"
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
		return fmt.Errorf("failed to validate payment method: name is required")
	}

	if !pM.Type.IsValid() {
		return fmt.Errorf("failed to validate payment method: invalid type")
	}

	if !pM.FeeType.IsValid() {
		return fmt.Errorf("failed to validate payment method: invalid fee type")
	}

	switch pM.FeeType {
	case FeeTypeFlat:
		if pM.FeeFixed < 0 {
			return fmt.Errorf("failed to validate payment method: fixed fee must be >= 0")
		}

	case FeeTypePercentage:
		if pM.FeePercentage < 0 || pM.FeePercentage > 1 {
			return fmt.Errorf("failed to validate payment method: percentage must be between 0 and 1")
		}

	case FeeTypeMixed:
		if pM.FeeFixed < 0 {
			return fmt.Errorf("failed to validate payment method: fixed fee must be >= 0")
		}
		if pM.FeePercentage < 0 || pM.FeePercentage > 1 {
			return fmt.Errorf("failed to validate payment method: percentage must be between 0 and 1")
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
