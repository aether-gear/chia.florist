package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	productRepo "service-core/internal/modules/product/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type AddItemUsecase struct {
	cartRepo      repository.CartRepository
	inventoryRepo inventoryRepo.InventoryRepository
	productRepo   productRepo.ProductRepository
	transactor    transaction.Transactor
	executor      transaction.Executor
}

func NewAddItemUsecase(
	cartRepo repository.CartRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productRepo productRepo.ProductRepository,
	transactor transaction.Transactor,
	executor transaction.Executor,
) *AddItemUsecase {
	return &AddItemUsecase{
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		transactor:    transactor,
		executor:      executor,
	}
}

type AddItemInput struct {
	UserID, ProductID, ShopID uuid.UUID
	Quantity                  int
}

func (u *AddItemUsecase) Execute(
	ctx context.Context,
	input AddItemInput,
) error {
	if input.ShopID == uuid.Nil {
		return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	if input.Quantity <= 0 {
		return apperrors.NewInvalidInput(domain.ErrInvalidQuantity.Error())
	}

	inventory, err := u.inventoryRepo.GetByProductIDAndShopID(
		ctx, u.executor, input.ProductID, input.ShopID,
	)
	if err != nil {
		return fmt.Errorf("failed to load inventory by product and shop: %w", err)
	}
	if inventory == nil {
		return apperrors.NewNotFound(domain.ErrProductNotFound.Error())
	}

	product, err := u.productRepo.GetByID(ctx, u.executor, input.ProductID)
	if err != nil {
		return fmt.Errorf("failed to load product with inventory: %w", err)
	}
	if product == nil {
		return apperrors.NewNotFound(domain.ErrProductNotFound.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(ctx, u.executor, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		cart, err = u.cartRepo.NewCart(ctx, u.executor, input.UserID)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if cart.HasProductInAnotherShop(input.ProductID, input.ShopID) {
		return apperrors.
			NewConflict(domain.ErrProductAlreadyAssignedToShop.Error())
	}

	targetQuantity := input.Quantity
	if existingItem := cart.FindItem(input.ProductID, input.ShopID); existingItem != nil && existingItem.DeletedAt == nil {
		targetQuantity += existingItem.Quantity
	}
	if targetQuantity > inventory.Available() {
		return apperrors.NewConflict(domain.ErrInsufficientStock.Error())
	}

	if err := cart.AddItem(input.ProductID, input.ShopID, input.Quantity); err != nil {
		return apperrors.NewInvalidInput(err.Error())
	}

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
				return fmt.Errorf("failed to add item: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
