package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type SavePaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
	executor          transaction.Executor
}

func NewSavePaymentMethodUsecase(
	paymentMethodRepo repository.PaymentMethodRepository,
	executor transaction.Executor,
) *SavePaymentMethodUsecase {
	return &SavePaymentMethodUsecase{
		paymentMethodRepo: paymentMethodRepo,
		executor:          executor,
	}
}

type CreatePaymentMethodInput struct {
	ID            *uuid.UUID
	Name          string
	Code          string
	Provider      string
	Type          string
	IsActive      bool
	Description   string
	FeeType       string
	FeeFixed      int64
	FeePercentage float64
}

func (u *SavePaymentMethodUsecase) Execute(
	ctx context.Context,
	input CreatePaymentMethodInput,
) error {
	var id uuid.UUID
	if input.ID != nil && *input.ID != uuid.Nil {
		id = *input.ID
	} else {
		id = uuid.New()
	}

	paymentMethod := domain.PaymentMethod{
		ID:            id,
		Name:          input.Name,
		Code:          domain.PaymentMethodCode(input.Code),
		Provider:      input.Provider,
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

	err := u.paymentMethodRepo.
		Save(ctx, u.executor, paymentMethod)
	if err != nil {
		return fmt.Errorf("failed to save payment method: %w", err)
	}

	return nil
}
