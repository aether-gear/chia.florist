package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ChangeItemShopInput struct {
	CustomerID uuid.UUID
	CartItemID uuid.UUID
	NewShopID  uuid.UUID
}

type ChangeItemShopUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	cartRepo      repository.CartRepository
	shopRepo      shopRepo.ShopRepository
	inventoryRepo inventoryRepo.InventoryRepository
}

func NewChangeItemShopUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	cartRepo repository.CartRepository,
	shopRepo shopRepo.ShopRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
) *ChangeItemShopUsecase {
	return &ChangeItemShopUsecase{
		executor:      executor,
		transactor:    transactor,
		cartRepo:      cartRepo,
		shopRepo:      shopRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *ChangeItemShopUsecase) Execute(
	ctx context.Context,
	input ChangeItemShopInput,
) error {
	if input.NewShopID == uuid.Nil {
		return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	shop, err := u.shopRepo.GetByID(ctx, u.executor, input.NewShopID)
	if err != nil {
		return fmt.Errorf("failed to retrieve shop: %w", err)
	}
	if shop == nil || !shop.IsActive {
		return apperrors.NewNotFound("target shop not found or inactive")
	}

	cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor, input.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to load cart: %w", err)
	}
	if cart == nil {
		return apperrors.NewNotFound(domain.ErrCartNotFound.Error())
	}

	var targetItem *domain.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == input.CartItemID && cart.Items[i].DeletedAt == nil {
			targetItem = &cart.Items[i]
			break
		}
	}

	if targetItem == nil {
		return apperrors.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

	if targetItem.ShopID == input.NewShopID {
		return nil
	}

	// Standard catalog product validation: check inventory stock at target shop
	if targetItem.ProductVariantType != domain.ProductVariantTypeCustom && targetItem.ProductID != nil {
		inventory, err := u.inventoryRepo.GetByProductIDAndShopID(ctx, u.executor, *targetItem.ProductID, input.NewShopID)
		if err != nil {
			return fmt.Errorf("failed to load inventory for target shop: %w", err)
		}
		if inventory == nil || inventory.Available() < targetItem.Quantity {
			return apperrors.NewConflict(domain.ErrInsufficientStock.Error())
		}
	}

	if !cart.ChangeItemShop(input.CartItemID, input.NewShopID) {
		return apperrors.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

	if err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
			return fmt.Errorf("failed to save cart item shop change: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
