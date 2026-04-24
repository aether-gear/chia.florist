package usecase

import (
	"errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type CreatePaymentMethod struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewCreatePaymentMethod(
	pMR repository.PaymentMethodRepository,
) *CreatePaymentMethod {
	return &CreatePaymentMethod{
		paymentMethodRepo: pMR,
	}
}

func (u *CreatePaymentMethod) Execute(input CreatePaymentMethodInput) error {
	validMethodTypes := map[string]domain.PaymentMethodType{
		string(domain.TypeBankTransfer): domain.TypeBankTransfer,
		string(domain.TypeEWallet):      domain.TypeEWallet,
		string(domain.TypeQRCode):       domain.TypeQRCode,
	}
	methodType, ok := validMethodTypes[input.Type]
	if !ok {
		return errors.New("invalid payment type")
	}

	validFeeTypes := map[string]domain.PaymentFeeType{
		string(domain.FeeTypeFlat):       domain.FeeTypeFlat,
		string(domain.FeeTypePercentage): domain.FeeTypePercentage,
		string(domain.FeeTypeMixed):      domain.FeeTypeMixed,
	}
	feeType, ok := validFeeTypes[input.FeeType]
	if !ok {
		return errors.New("invalid payment type")
	}

	if input.FeeFixed < 0 || input.FeePercentage > 1 || input.FeePercentage < 0 {
		return errors.New("Invalid fee amount")
	}

	paymentMethod := domain.PaymentMethod{
		ID:            uuid.New(),
		Name:          input.Name,
		Type:          methodType,
		IsActive:      input.IsActive,
		FeeType:       feeType,
		FeeFixed:      input.FeeFixed,
		FeePercentage: input.FeePercentage,
		Description:   input.Description,
	}

	err := u.paymentMethodRepo.Save(paymentMethod)
	if err != nil {
		return err
	}

	return nil
}

type CreatePaymentMethodInput struct {
	Name          string
	Type          string
	IsActive      bool
	Description   string
	FeeType       string
	FeeFixed      int64
	FeePercentage float64
}
