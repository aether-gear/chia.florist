package usecase

import (
	"context"
	"fmt"

	"service-core/internal/infra/storage"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetCartUsecase struct {
	cartRepo       repository.CartRepository
	inventoryRepo  inventoryRepo.InventoryRepository
	productRepo    productRepo.ProductRepository
	productImgRepo productRepo.ProductImageRepository
	fileStore      storage.Provider
}

func NewGetCartUsecase(
	cartRepo repository.CartRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productRepo productRepo.ProductRepository,
	productImgRepo productRepo.ProductImageRepository,
	fileStore storage.Provider,
) *GetCartUsecase {
	return &GetCartUsecase{
		cartRepo:       cartRepo,
		inventoryRepo:  inventoryRepo,
		productRepo:    productRepo,
		productImgRepo: productImgRepo,
		fileStore:      fileStore,
	}
}

type ProductCartResponse struct {
	Product   productDomain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryDomain.Inventory
	Images          struct {
		Thumbnail string
	}
}

type GetCartResult struct {
	Cart     *domain.Cart
	Products map[uuid.UUID]ProductCartResponse
}

func (u *GetCartUsecase) Execute(
	ctx context.Context,
	userID uuid.UUID,
) (*GetCartResult, error) {
	cart, err := u.cartRepo.GetWithItemsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cart: %w", err)
	}

	if cart == nil {
		cart, err = u.cartRepo.NewCart(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if len(cart.Items) == 0 {
		return &GetCartResult{
			Cart:     cart,
			Products: map[uuid.UUID]ProductCartResponse{},
		}, nil
	}

	productIDs := make([]uuid.UUID, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := u.productRepo.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load cart with products: %w", err)
	}

	inventoryMap, err := u.inventoryRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for cart products: %w", err)
	}

	imagesMap, err := u.productImgRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load images for products: %w", err)
	}

	productMap := make(map[uuid.UUID]ProductCartResponse)
	for _, p := range products {
		inventories := inventoryMap[p.ID]
		images := imagesMap[p.ID]

		result := ProductCartResponse{
			Product:         p,
			ShopInventories: inventories,
		}

		if len(images) > 0 {
			key := images[0].Variants[productDomain.ResolutionThumbnail].Key
			result.Images.Thumbnail = u.fileStore.
				PublicURL(key, "public-assets")
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
