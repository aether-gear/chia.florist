package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type CreatePaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewCreatePaymentMethodUsecase(
	paymentMethodRepo repository.PaymentMethodRepository,
) *CreatePaymentMethodUsecase {
	return &CreatePaymentMethodUsecase{
		paymentMethodRepo: paymentMethodRepo,
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

func (u *CreatePaymentMethodUsecase) Execute(
	ctx context.Context,
	input CreatePaymentMethodInput,
) error {
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
		return apperrors.NewInvalidInput(err.Error())
	}

	err := u.paymentMethodRepo.Save(ctx, paymentMethod)
	if err != nil {
		return fmt.Errorf("failed to save payment method: %w", err)
	}

	return nil
}
