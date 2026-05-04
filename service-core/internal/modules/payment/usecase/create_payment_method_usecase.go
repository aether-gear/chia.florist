package usecase

import (
	"fmt"

	appErr "service-core/internal/common/errors"
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

type CreatePaymentMethodInput struct {
	Name          string
	Type          string
	IsActive      bool
	Description   string
	FeeType       string
	FeeFixed      int64
	FeePercentage float64
}

func (u *CreatePaymentMethod) Execute(input CreatePaymentMethodInput) error {
	paymentMethod := domain.PaymentMethod{
		ID:            uuid.New(),
		Name:          input.Name,
		Type:          domain.PaymentMethodType(input.Type),
		IsActive:      input.IsActive,
		FeeType:       domain.PaymentFeeType(input.FeeType),
		FeeFixed:      input.FeeFixed,
		FeePercentage: input.FeePercentage,
		Description:   input.Description,
	}

	if err := paymentMethod.Validate(); err != nil {
		return appErr.NewInvalidInput(err.Error())
	}

	err := u.paymentMethodRepo.Save(paymentMethod)
	if err != nil {
		return fmt.Errorf("failed to save payment method: %w", err)
	}

	return nil
}
