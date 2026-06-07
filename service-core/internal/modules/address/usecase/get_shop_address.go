package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetShopAddressUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
	executor        transaction.Executor
}

func NewGetShopAddressUsecase(
	shopAddressRepo repository.ShopAddressRepository,
	executor transaction.Executor,
) *GetShopAddressUsecase {
	return &GetShopAddressUsecase{
		shopAddressRepo: shopAddressRepo,
		executor:        executor,
	}
}

func (u *GetShopAddressUsecase) GetByID(
	ctx context.Context,
	addressID uuid.UUID,
) (*domain.ShopAddress, error) {
	res, err := u.shopAddressRepo.GetByID(ctx, u.executor, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
