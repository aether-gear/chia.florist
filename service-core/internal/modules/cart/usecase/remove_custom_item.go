package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type RemoveCustomItemInput struct {
	CustomerID uuid.UUID
	CartItemID uuid.UUID // the cart_items.id of the custom item to remove
}

type RemoveCustomItemUsecase struct {
	executor   transaction.Executor
	transactor transaction.Transactor
	cartRepo   repository.CartRepository
}

func NewRemoveCustomItemUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	cartRepo repository.CartRepository,
) *RemoveCustomItemUsecase {
	return &RemoveCustomItemUsecase{
		executor:   executor,
		transactor: transactor,
		cartRepo:   cartRepo,
	}
}

func (u *RemoveCustomItemUsecase) Execute(ctx context.Context, input RemoveCustomItemInput) error {
	cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor, input.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to load cart: %w", err)
	}
	if cart == nil {
		return apperrors.NewNotFound(domain.ErrCartNotFound.Error())
	}

	if !cart.RemoveCustomItem(input.CartItemID) {
		return apperrors.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
			return fmt.Errorf("failed to update cart after removing custom item: %w", err)
		}
		return nil
	})
}
