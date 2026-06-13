package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ListUserAddressUsecase struct {
	userAddressRepo repository.UserAddressRepository
	executor        transaction.Executor
}

func NewListUserAddressUsecase(
	userAddressRepo repository.UserAddressRepository,
	executor transaction.Executor,
) *ListUserAddressUsecase {
	return &ListUserAddressUsecase{
		userAddressRepo: userAddressRepo,
		executor:        executor,
	}
}

func (u *ListUserAddressUsecase) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Address, error) {
	res, err := u.userAddressRepo.ListByUserID(ctx, u.executor, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
