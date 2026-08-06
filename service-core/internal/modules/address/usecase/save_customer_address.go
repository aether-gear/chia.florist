package usecase

import (
	"context"
	"fmt"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type SaveCustomerAddressUsecase struct {
	executor            transaction.Executor
	transactor          transaction.Transactor
	customerAddressRepo repository.CustomerAddressRepository
}

func NewSaveCustomerAddressUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	customerAddressRepo repository.CustomerAddressRepository,
) *SaveCustomerAddressUsecase {
	return &SaveCustomerAddressUsecase{
		executor:            executor,
		transactor:          transactor,
		customerAddressRepo: customerAddressRepo,
	}
}

type SaveCustomerAddressInput struct {
	ID           *uuid.UUID
	CustomerID   uuid.UUID
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

func (u *SaveCustomerAddressUsecase) Execute(
	ctx context.Context,
	input SaveCustomerAddressInput,
) error {
	var addressID uuid.UUID

	isCreate := input.ID == nil
	isDefault := input.IsDefault != nil && *input.IsDefault

	if isCreate {
		addressID = uuid.New()

		count, err := u.customerAddressRepo.
			CountByCustomerID(ctx, u.executor, input.CustomerID)
		if err != nil {
			return fmt.Errorf("failed to count addresses: %w", err)
		}

		if count != nil {
			if *count >= 10 {
				return apperrors.NewConflict(domain.ErrAddressLimitReached.Error())
			}

			if *count == 0 {
				isDefault = true
			}
		}
	} else {
		addressID = *input.ID
	}

	address := domain.CustomerAddress{
		ID:           addressID,
		CustomerID:   input.CustomerID,
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
		CreatedAt: appclock.Now(),
	}

	err := u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if isDefault {
				if err := u.customerAddressRepo.
					UnsetDefaultByCustomerID(
						ctx,
						exec,
						input.CustomerID,
					); err != nil {
					return fmt.Errorf("failed to unset default address: %w", err)
				}
			}

			if err := u.customerAddressRepo.
				Save(
					ctx,
					exec,
					address,
				); err != nil {
				return fmt.Errorf("failed to save address: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
