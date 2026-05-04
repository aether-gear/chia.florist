package usecase

import (
	"errors"
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	inventoryD "service-core/internal/modules/inventory/domain"
	inventoryR "service-core/internal/modules/inventory/repository"
	productD "service-core/internal/modules/product/domain"
	productR "service-core/internal/modules/product/repository"
	shopR "service-core/internal/modules/shop/repository"

	"github.com/google/uuid"
)

type CreateInventoryUsecase struct {
	inventoryRepo inventoryR.InventoryRepository
	productRepo   productR.ProductRepository
	shopRepo      shopR.ShopRepository
}

func NewCreateInventoryUsecase(
	inventoryRepo inventoryR.InventoryRepository,
	productRepo productR.ProductRepository,
	shopRepo shopR.ShopRepository,
) *CreateInventoryUsecase {
	return &CreateInventoryUsecase{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		shopRepo:      shopRepo,
	}
}

type CreateInventoryInput struct {
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Stock     int
}

func (u *CreateInventoryUsecase) Execute(input CreateInventoryInput) error {
	product, err := u.productRepo.GetByID(input.ProductID)
	if err != nil {
		return fmt.Errorf("failed to load product: %w", err)
	}
	if product == nil {
		return appErr.NewNotFound(productD.ErrProductNotFound.Error())
	}

	shop, err := u.shopRepo.GetByID(input.ShopID)
	if err != nil {
		return fmt.Errorf("failed to load shop: %w", err)
	}
	if shop == nil {
		return appErr.NewNotFound("shop not found")
	}

	existing, err := u.inventoryRepo.GetByProductAndShop(input.ProductID, input.ShopID)
	if err != nil {
		return fmt.Errorf("failed to load inventory: %w", err)
	}
	if existing != nil {
		return appErr.NewConflict("inventory already exists for product and shop")
	}

	inventory := &inventoryD.Inventory{
		ID:        uuid.New(),
		ProductID: input.ProductID,
		ShopID:    input.ShopID,
		Stock:     input.Stock,
		Reserved:  0,
		CreatedAt: time.Now(),
	}
	if err := inventory.Validate(); err != nil {
		if errors.Is(err, inventoryD.ErrInvalidStock) || errors.Is(err, inventoryD.ErrInvalidReserved) {
			return appErr.NewInvalidInput(err.Error())
		}

		return err
	}

	if err := u.inventoryRepo.Create(inventory); err != nil {
		return fmt.Errorf("failed to save inventory: %w", err)
	}

	return nil
}
