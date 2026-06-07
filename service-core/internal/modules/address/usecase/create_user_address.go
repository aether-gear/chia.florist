package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateAddressUsecase struct {
	userAddressRepo repository.UserAddressRepository
	executor        transaction.Executor
}

func NewCreateAddressUsecase(
	userAddressRepo repository.UserAddressRepository,
	executor transaction.Executor,

) *CreateAddressUsecase {
	return &CreateAddressUsecase{
		userAddressRepo: userAddressRepo,
	}
}

type CreateAddressInput struct {
	UserID       uuid.UUID
	ReceiverName string
	Phone        *string
	IsDefault    *bool
	ProvinceID   string
	CityID       string
	DistrictID   string
	VillageID    string
	FullAddress  string
	PostalCode   string
}

func (u *CreateAddressUsecase) Execute(
	ctx context.Context,
	input CreateAddressInput,
) error {
	var isDefault bool
	if *input.IsDefault {
		isDefault = *input.IsDefault
	} else {
		isDefault = false
	}

	address := domain.Address{
		ID:           uuid.New(),
		UserID:       input.UserID,
		ReceiverName: input.ReceiverName,
		Phone:        input.Phone,
		IsDefault:    isDefault,
		Detail: domain.AddressDetail{
			ProvinceID:  input.ProvinceID,
			CityID:      input.CityID,
			DistrictID:  input.DistrictID,
			VillageID:   input.VillageID,
			FullAddress: input.FullAddress,
			PostalCode:  input.PostalCode,
		},
		CreatedAt: time.Now(),
	}

	err := u.userAddressRepo.Create(ctx, u.executor, address)
	if err != nil {
		return fmt.Errorf("failed to save address: %w", err)
	}

	return nil
}
