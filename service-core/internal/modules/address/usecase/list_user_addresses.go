package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type ListUserAddressUsecase struct {
	userAddressRepo repository.UserAddressRepository
}

func NewListUserAddressUsecase(
	userAddressRepo repository.UserAddressRepository,
) *ListUserAddressUsecase {
	return &ListUserAddressUsecase{
		userAddressRepo: userAddressRepo,
	}
}

func (u *ListUserAddressUsecase) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Address, error) {
	res, err := u.userAddressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
