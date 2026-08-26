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
	executor   transaction.Executor
	transactor transaction.Transactor
	cartRepo   repository.CartRepository
}

func NewRemoveItemUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	cartRepo repository.CartRepository,
) *RemoveItemUsecase {
	return &RemoveItemUsecase{
		executor:   executor,
		transactor: transactor,
		cartRepo:   cartRepo,
	}
}

type RemoveItemInput struct {
	CustomerID, ProductID, ShopID uuid.UUID
	ItemOptions                   *domain.ItemOptions
}

func (u *RemoveItemUsecase) Execute(
	ctx context.Context,
	input RemoveItemInput,
) error {
	if input.ShopID == uuid.Nil {
		return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor,
		input.CustomerID,
	)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return apperrors.NewNotFound(domain.ErrCartNotFound.Error())
	}

	var opts []domain.ItemOptions
	if input.ItemOptions != nil {
		opts = append(opts, *input.ItemOptions)
	}

	if cart.FindItem(input.ProductID, input.ShopID, opts...) == nil {
		return apperrors.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

	cart.RemoveItem(input.ProductID, input.ShopID, opts...)

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

type RemoveItemByIDInput struct {
	CustomerID uuid.UUID
	CartItemID uuid.UUID
}

func (u *RemoveItemUsecase) ExecuteByID(
	ctx context.Context,
	input RemoveItemByIDInput,
) error {
	if input.CartItemID == uuid.Nil {
		return apperrors.NewInvalidInput("invalid cart item id")
	}

	cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor,
		input.CustomerID,
	)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return apperrors.NewNotFound(domain.ErrCartNotFound.Error())
	}

	if !cart.RemoveItemByID(input.CartItemID) {
		return apperrors.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

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
