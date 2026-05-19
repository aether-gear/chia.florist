package usecase

import (
	"fmt"

	appErr "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	cartR "service-core/internal/modules/cart/repository"
	inventoryR "service-core/internal/modules/inventory/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type UpdateItemUsecase struct {
	cartRepo      cartR.CartRepository
	inventoryRepo inventoryR.InventoryRepository
	productRepo   productR.ProductRepository
}

func NewUpdateItemUsecase(
	cartRepo cartR.CartRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productRepo productR.ProductRepository,
) *UpdateItemUsecase {
	return &UpdateItemUsecase{
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

type UpdateItemInput struct {
	UserID, ProductID, ShopID uuid.UUID
	Quantity                  int
}

func (u *UpdateItemUsecase) Execute(input UpdateItemInput) error {
	if input.ShopID == uuid.Nil {
		return appErr.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	if input.Quantity <= 0 {
		return appErr.NewInvalidInput(domain.ErrInvalidQuantity.Error())
	}

	inventory, err := u.inventoryRepo.GetByProductIDAndShopID(input.ProductID, input.ShopID)
	if err != nil {
		return fmt.Errorf("failed to load inventory by product and shop: %w", err)
	}
	if inventory == nil {
		return appErr.NewNotFound(domain.ErrProductNotFound.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}

	if cart == nil {
		cart, err = u.cartRepo.NewCart(input.UserID)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if !cart.HasItem(input.ProductID, input.ShopID) {
		return appErr.NewNotFound(domain.ErrCartItemNotFound.Error())
	}

	product, err := u.productRepo.GetByID(input.ProductID)
	if err != nil {
		return fmt.Errorf("failed to retrieve product: %w", err)
	}
	if product == nil {
		return appErr.NewNotFound(domain.ErrProductNotFound.Error())
	}

	if input.Quantity > inventory.Available() {
		return appErr.NewConflict(domain.ErrInsufficientStock.Error())
	}

	if err := cart.SetItem(input.ProductID, input.ShopID, input.Quantity); err != nil {
		return appErr.NewInvalidInput(err.Error())
	}

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	return nil
}
