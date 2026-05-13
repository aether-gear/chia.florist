package usecase

import (
	"fmt"

	cartD "service-core/internal/modules/cart/domain"
	cartR "service-core/internal/modules/cart/repository"
	inventoryR "service-core/internal/modules/inventory/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetCartUsecase struct {
	cartRepo      cartR.CartRepository
	inventoryRepo inventoryR.InventoryRepository
	productRepo   productR.ProductRepository
}

func NewGetCartUsecase(
	cartRepo cartR.CartRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productRepo productR.ProductRepository,
) *GetCartUsecase {
	return &GetCartUsecase{
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

type GetCartResult struct {
	Cart     *cartD.Cart
	Products map[uuid.UUID]productR.ProductWithInventory
}

func (u *GetCartUsecase) Execute(userID uuid.UUID) (*GetCartResult, error) {
	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cart: %w", err)
	}

	if cart == nil {
		cart, err = u.cartRepo.NewCart(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if len(cart.Items) == 0 {
		return &GetCartResult{
			Cart:     cart,
			Products: map[uuid.UUID]productR.ProductWithInventory{},
		}, nil
	}

	productIDs := make([]uuid.UUID, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := u.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load cart with products: %w", err)
	}

	inventoryMap, err := u.inventoryRepo.ListByProductIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for cart products: %w", err)
	}

	productMap := make(map[uuid.UUID]productR.ProductWithInventory)
	for _, p := range products {
		inventories := inventoryMap[p.ID]

		result := productR.ProductWithInventory{
			Product:         p,
			ShopInventories: inventories,
		}

		for _, inventory := range inventories {
			result.Inventory.TotalStock += inventory.TotalStock
			result.Inventory.ReservedStock += inventory.ReservedStock
		}

		productMap[p.ID] = result
	}

	return &GetCartResult{
		Cart:     cart,
		Products: productMap,
	}, nil
}
