package usecase

import (
	"context"
	"fmt"
	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DeleteUserAddressUsecase struct {
	executor        transaction.Executor
	userAddressRepo repository.UserAddressRepository
}

func NewDeleteUserAddressUsecase(
	executor transaction.Executor,
	userAddressRepo repository.UserAddressRepository,
) *DeleteUserAddressUsecase {
	return &DeleteUserAddressUsecase{
		executor:        executor,
		userAddressRepo: userAddressRepo,
	}
}

func (u *DeleteUserAddressUsecase) Execute(
	ctx context.Context,
	addressID uuid.UUID,
) error {
	address, err := u.userAddressRepo.
		GetByID(ctx, u.executor, addressID)
	if err != nil {
		return fmt.Errorf("failed to retrieve address: %w", err)
	}
	if address == nil {
		return apperrors.NewNotFound(domain.ErrAddressNotFound.Error())
	}

	if address.IsDefault {
		return apperrors.NewConflict(domain.ErrCannotDeleteDefaultAddress.Error())
	}

	if err := u.userAddressRepo.
		Delete(ctx, u.executor, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	return nil
}
