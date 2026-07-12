package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/inventory/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DeleteInventoryUsecase struct {
	inventoryRepo repository.InventoryRepository
	executor      transaction.Executor
}

func NewDeleteInventoryUsecase(
	inventoryRepo repository.InventoryRepository,
	executor transaction.Executor,
) *DeleteInventoryUsecase {
	return &DeleteInventoryUsecase{
		inventoryRepo: inventoryRepo,
		executor:      executor,
	}
}

type DeleteInventoryInput struct {
	ProductID uuid.UUID
	ShopID    uuid.UUID
}

func (u *DeleteInventoryUsecase) Execute(
	ctx context.Context,
	input DeleteInventoryInput,
) error {
	existing, err := u.inventoryRepo.
		GetByProductIDAndShopID(ctx, u.executor,
			input.ProductID,
			input.ShopID,
		)
	if err != nil {
		return fmt.Errorf("failed to load inventory: %w", err)
	}

	if existing == nil {
		return apperrors.NewNotFound("inventory not found")
	}

	if existing.ReservedStock > 0 {
		return apperrors.NewConflict("cannot delete inventory with active reservations")
	}

	if err := u.inventoryRepo.
		Delete(ctx, u.executor,
			input.ProductID,
			input.ShopID,
		); err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	return nil
}
