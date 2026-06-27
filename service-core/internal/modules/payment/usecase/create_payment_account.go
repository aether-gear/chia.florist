package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreatePaymentAccountUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
	paymentAccRepo    repository.PaymentAccountRepository
	executor          transaction.Executor
}

func NewCreatePaymentAccountUsecase(
	paymentAccRepo repository.PaymentAccountRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
	executor transaction.Executor,
) *CreatePaymentAccountUsecase {
	return &CreatePaymentAccountUsecase{
		paymentAccRepo:    paymentAccRepo,
		paymentMethodRepo: paymentMethodRepo,
		executor:          executor,
	}
}

type CreatePaymentAccountInput struct {
	MethodID      uuid.UUID
	AccountName   string
	AccountNumber *string
	PhoneNumber   string
	QRString      *string
	IsActive      bool
}

func (u *CreatePaymentAccountUsecase) Execute(
	ctx context.Context,
	input CreatePaymentAccountInput,
) error {
	method, err := u.paymentMethodRepo.GetByID(ctx, u.executor, input.MethodID)
	if err != nil {
		return fmt.Errorf("failed to retrieve payment account: %w", err)
	}
	if method == nil {
		return apperrors.NewNotFound(domain.ErrPaymentMethodNotFound.Error())
	}

	var phonePtr *string
	if input.PhoneNumber != "" {
		phonePtr = &input.PhoneNumber
	}

	paymentAccount := domain.PaymentAccount{
		ID:            uuid.New(),
		MethodID:      method.ID,
		AccountName:   input.AccountName,
		AccountNumber: input.AccountNumber,
		PhoneNumber:   phonePtr,
		QRString:      input.QRString,
		IsActive:      input.IsActive,
		CurrentLoad:   0,
		CreatedAt:     time.Now(),
	}

	if err := paymentAccount.ValidateForMethod(method.Type); err != nil {
		if errors.Is(err, domain.ErrUnsupportedPaymentMethod) {
			return apperrors.NewBadRequest(err.Error())
		}

		return apperrors.NewInvalidInput(err.Error())
	}

	err = u.paymentAccRepo.Save(ctx, u.executor, paymentAccount)
	if err != nil {
		return fmt.Errorf("failed to save payment account: %w", err)
	}

	return nil
}
