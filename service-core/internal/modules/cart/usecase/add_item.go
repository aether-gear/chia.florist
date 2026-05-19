package usecase

import (
	"fmt"

	appErr "service-core/internal/common/errors"
	cartD "service-core/internal/modules/cart/domain"
	cartR "service-core/internal/modules/cart/repository"
	inventoryR "service-core/internal/modules/inventory/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type AddItemUsecase struct {
	cartRepo      cartR.CartRepository
	inventoryRepo inventoryR.InventoryRepository
	productRepo   productR.ProductRepository
}

func NewAddItemUsecase(
	cartRepo cartR.CartRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productRepo productR.ProductRepository,
) *AddItemUsecase {
	return &AddItemUsecase{
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

type AddItemInput struct {
	UserID, ProductID, ShopID uuid.UUID
	Quantity                  int
}

func (u *AddItemUsecase) Execute(input AddItemInput) error {
	if input.ShopID == uuid.Nil {
		return appErr.NewInvalidInput(cartD.ErrInvalidShopID.Error())
	}

	if input.Quantity <= 0 {
		return appErr.NewInvalidInput(cartD.ErrInvalidQuantity.Error())
	}

	inventory, err := u.inventoryRepo.GetByProductIDAndShopID(input.ProductID, input.ShopID)
	if err != nil {
		return fmt.Errorf("failed to load inventory by product and shop: %w", err)
	}
	if inventory == nil {
		return appErr.NewNotFound(cartD.ErrProductNotFound.Error())
	}

	product, err := u.productRepo.GetByID(input.ProductID)
	if err != nil {
		return fmt.Errorf("failed to load product with inventory: %w", err)
	}
	if product == nil {
		return appErr.NewNotFound(cartD.ErrProductNotFound.Error())
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

	if cart.HasProductInAnotherShop(input.ProductID, input.ShopID) {
		return appErr.NewConflict(cartD.ErrProductAlreadyAssignedToShop.Error())
	}

	targetQuantity := input.Quantity
	if existingItem := cart.FindItem(input.ProductID, input.ShopID); existingItem != nil && existingItem.DeletedAt == nil {
		targetQuantity += existingItem.Quantity
	}
	if targetQuantity > inventory.Available() {
		return appErr.NewConflict(cartD.ErrInsufficientStock.Error())
	}

	if err := cart.AddItem(input.ProductID, input.ShopID, input.Quantity); err != nil {
		return appErr.NewInvalidInput(err.Error())
	}

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}
