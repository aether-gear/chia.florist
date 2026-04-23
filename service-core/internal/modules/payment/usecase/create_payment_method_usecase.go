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
	paymentMethodRepo repository.PaymentMethodRepository,
) *CreatePaymentMethod {
	return &CreatePaymentMethod{
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (u *CreatePaymentMethod) Execute(input CreatePaymentMethodInput) error {
	validTypes := map[string]bool{
		string(domain.TypeBankTransfer): true,
		string(domain.TypeEWallet):      true,
		string(domain.TypeQRCode):       true,
	}
	if !validTypes[input.Type] {
		return errors.New("invalid payment type")
	}

	paymentMethod := domain.PaymentMethod{
		ID:          uuid.New(),
		Name:        input.Name,
		Type:        input.Type,
		IsActive:    input.IsActive,
		Description: input.Description,
	}

	err := u.paymentMethodRepo.Save(paymentMethod)
	if err != nil {
		return err
	}

	return nil
}

type CreatePaymentMethodInput struct {
	Name        string
	Type        string
	IsActive    bool
	Description string
}
