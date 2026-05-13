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

func (u *AddItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID, shopID uuid.UUID, quantity int) error {
	if shopID == uuid.Nil {
		return appErr.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	if quantity <= 0 {
		return appErr.NewInvalidInput(domain.ErrInvalidQuantity.Error())
	}

	inventory, err := u.inventoryRepo.GetByProductIDAndShopID(productID, shopID)
	if err != nil {
		return fmt.Errorf("failed to load inventory by product and shop: %w", err)
	}
	if inventory == nil {
		return appErr.NewNotFound(domain.ErrProductNotFound.Error())
	}

	product, err := u.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to load product with inventory: %w", err)
	}
	if product == nil {
		return appErr.NewNotFound(domain.ErrProductNotFound.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		cart, err = u.cartRepo.NewCart(userID)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if cart.HasProductInAnotherShop(productID, shopID) {
		return appErr.NewConflict(domain.ErrProductAlreadyAssignedToShop.Error())
	}

	targetQuantity := quantity
	if existingItem := cart.FindItem(productID, shopID); existingItem != nil && existingItem.DeletedAt == nil {
		targetQuantity += existingItem.Quantity
	}
	if targetQuantity > inventory.Available() {
		return appErr.NewConflict(domain.ErrInsufficientStock.Error())
	}

	if err := cart.AddItem(productID, shopID, quantity); err != nil {
		return appErr.NewInvalidInput(err.Error())
	}

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}
