package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type AddCustomItemInput struct {
	CustomerID     uuid.UUID
	ShopID         uuid.UUID
	Quantity       int
	ProductName    string
	PhysicalSizeID string
	CustomDesign   json.RawMessage
}

type AddCustomItemUsecase struct {
	executor   transaction.Executor
	transactor transaction.Transactor
	cartRepo   repository.CartRepository
}

func NewAddCustomItemUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	cartRepo repository.CartRepository,
) *AddCustomItemUsecase {
	return &AddCustomItemUsecase{
		executor:   executor,
		transactor: transactor,
		cartRepo:   cartRepo,
	}
}

func (u *AddCustomItemUsecase) Execute(ctx context.Context, input AddCustomItemInput) error {
	if input.ShopID == uuid.Nil {
		return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}
	if input.Quantity <= 0 {
		return apperrors.NewInvalidInput(domain.ErrInvalidQuantity.Error())
	}
	if len(input.CustomDesign) == 0 {
		return apperrors.NewInvalidInput("custom_design is required for custom items")
	}
	if input.PhysicalSizeID == "" {
		return apperrors.NewInvalidInput("physical_size_id is required for custom items")
	}

	cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor,
		input.CustomerID,
	)
	if err != nil {
		return fmt.Errorf("failed to load cart: %w", err)
	}
	if cart == nil {
		cart, err = u.cartRepo.NewCart(ctx, u.executor,
			input.CustomerID,
		)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if err := cart.AddCustomItem(
		input.ShopID,
		input.Quantity,
		input.CustomDesign,
	); err != nil {
		return apperrors.NewInvalidInput(err.Error())
	}

	if err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
			return fmt.Errorf("failed to save cart with custom item: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
