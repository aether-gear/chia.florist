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

type UpdateShopAddressUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
	executor        transaction.Executor
	transactor      transaction.Transactor
}

func NewUpdateShopAddressUsecase(
	shopAddressRepo repository.ShopAddressRepository,
	executor transaction.Executor,
	transactor transaction.Transactor,
) *UpdateShopAddressUsecase {
	return &UpdateShopAddressUsecase{
		shopAddressRepo: shopAddressRepo,
		executor:        executor,
		transactor:      transactor,
	}
}

type UpdateShopAddressInput struct {
	ID          uuid.UUID
	ShopID      uuid.UUID
	Label       string
	Phone       *string
	IsActive    *bool
	ProvinceID  string
	CityID      string
	DistrictID  string
	VillageID   string
	FullAddress string
	PostalCode  string
}

func (u *UpdateShopAddressUsecase) Execute(
	ctx context.Context,
	input UpdateShopAddressInput,
) error {
	existing, err := u.shopAddressRepo.GetByID(ctx, u.executor, input.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve address: %w", err)
	}

	if existing == nil ||
		existing.ShopID != input.ShopID ||
		existing.DeletedAt != nil {

		return apperrors.NewNotFound(domain.ErrAddressNotFound.Error())
	}

	isDefault := existing.IsActive
	if input.IsActive != nil {
		isDefault = *input.IsActive
	}

	now := appclock.Now()
	address := domain.ShopAddress{
		ID:       input.ID,
		ShopID:   input.ShopID,
		Label:    input.Label,
		Phone:    input.Phone,
		IsActive: isDefault,
		Detail: domain.AddressDetail{
			ProvinceID:  input.ProvinceID,
			CityID:      input.CityID,
			DistrictID:  input.DistrictID,
			VillageID:   input.VillageID,
			FullAddress: input.FullAddress,
			PostalCode:  input.PostalCode,
		},
		CreatedAt: existing.CreatedAt,
		UpdatedAt: &now,
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if isDefault {
			if err := u.shopAddressRepo.UnsetActiveByShopID(ctx, exec,
				input.ShopID,
			); err != nil {
				return fmt.Errorf("failed to unset active address: %w", err)
			}
		}

		if err := u.shopAddressRepo.Update(ctx, exec,
			address,
		); err != nil {
			return fmt.Errorf("failed to update address: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
