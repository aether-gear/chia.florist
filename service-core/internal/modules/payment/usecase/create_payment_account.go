package usecase

import (
	"errors"
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
)

type CreatePaymentAccountUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
	paymentAccRepo    repository.PaymentAccountRepository
}

func NewCreatePaymentAccountUsecase(
	paymentAccRepo repository.PaymentAccountRepository,
	paymentMethodRepo repository.PaymentMethodRepository,
) *CreatePaymentAccountUsecase {
	return &CreatePaymentAccountUsecase{
		paymentAccRepo:    paymentAccRepo,
		paymentMethodRepo: paymentMethodRepo,
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

func (u *CreatePaymentAccountUsecase) Execute(input CreatePaymentAccountInput) error {
	method, err := u.paymentMethodRepo.GetByID(input.MethodID)
	if err != nil {
		return fmt.Errorf("failed to retrieve payment account: %w", err)
	}
	if method == nil {
		return appErr.NewNotFound(domain.ErrPaymentMethodNotFound.Error())
	}

	paymentAccount := domain.PaymentAccount{
		ID:            uuid.New(),
		MethodID:      method.ID,
		AccountName:   input.AccountName,
		AccountNumber: input.AccountNumber,
		PhoneNumber:   input.PhoneNumber,
		QRString:      input.QRString,
		IsActive:      input.IsActive,
		CurrentLoad:   0,
		CreatedAt:     time.Now(),
	}

	if err := paymentAccount.ValidateForMethod(method.Type); err != nil {
		if errors.Is(err, domain.ErrUnsupportedPaymentMethod) {
			return appErr.NewBadRequest(err.Error())
		}

		return appErr.NewInvalidInput(err.Error())
	}

	err = u.paymentAccRepo.Save(paymentAccount)
	if err != nil {
		return fmt.Errorf("failed to save payment account: %w", err)
	}

	return nil
}
