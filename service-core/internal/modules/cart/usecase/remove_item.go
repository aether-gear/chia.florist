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

type RemoveItemUsecase struct {
	cartRepo   repository.CartRepository
	transactor transaction.Transactor
}

func NewRemoveItemUsecase(
	cartRepo repository.CartRepository,
	transactor transaction.Transactor,
) *RemoveItemUsecase {
	return &RemoveItemUsecase{
		cartRepo:   cartRepo,
		transactor: transactor,
	}
}

type RemoveItemInput struct {
	UserID, ProductID, ShopID uuid.UUID
}

func (u *RemoveItemUsecase) Execute(
	ctx context.Context,
	input RemoveItemInput,
) error {
	if input.ShopID == uuid.Nil {
		return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return apperrors.NewNotFound(domain.ErrCartNotFound.Error())
	}

	cart.RemoveItem(input.ProductID, input.ShopID)

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
				return fmt.Errorf("failed to update cart: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
