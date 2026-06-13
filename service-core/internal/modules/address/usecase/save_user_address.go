package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type SaveUserAddressUsecase struct {
	executor        transaction.Executor
	transactor      transaction.Transactor
	userAddressRepo repository.UserAddressRepository
}

func NewSaveUserAddressUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	userAddressRepo repository.UserAddressRepository,

) *SaveUserAddressUsecase {
	return &SaveUserAddressUsecase{
		executor:        executor,
		transactor:      transactor,
		userAddressRepo: userAddressRepo,
	}
}

type SaveUserAddressInput struct {
	ID           *uuid.UUID
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

func (u *SaveUserAddressUsecase) Execute(
	ctx context.Context,
	input SaveUserAddressInput,
) error {
	var addressID uuid.UUID

	isCreate := input.ID == nil
	isDefault := input.IsDefault != nil && *input.IsDefault

	if !isCreate {
		addressID = *input.ID
	}
	if isCreate {
		addressID = uuid.New()

		count, err := u.userAddressRepo.
			CountByUserID(ctx, u.executor, input.UserID)
		if err != nil {
			return fmt.Errorf("failed to count addresses: %w", err)
		}

		if count != nil {
			if *count >= 10 {
				return apperrors.NewConflict(
					domain.ErrAddressLimitReached.Error(),
				)
			}

			if *count == 0 {
				isDefault = true
			}
		}
	}

	address := domain.Address{
		ID:           addressID,
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

	err := u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if isDefault {
				if err := u.userAddressRepo.
					UnsetDefaultByUserID(
						ctx,
						exec,
						input.UserID,
					); err != nil {
					return fmt.Errorf("failed to unset default address: %w", err)
				}
			}

			if err := u.userAddressRepo.
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
