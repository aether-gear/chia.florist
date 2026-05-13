package usecase

import (
	"fmt"

	"service-core/internal/infra/storage"
	cartD "service-core/internal/modules/cart/domain"
	cartR "service-core/internal/modules/cart/repository"
	inventoryD "service-core/internal/modules/inventory/domain"
	inventoryR "service-core/internal/modules/inventory/repository"
	productD "service-core/internal/modules/product/domain"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetCartUsecase struct {
	cartRepo       cartR.CartRepository
	inventoryRepo  inventoryR.InventoryRepository
	productRepo    productR.ProductRepository
	productImgRepo productR.ProductImageRepository
	fileStore      storage.Provider
}

func NewGetCartUsecase(
	cartRepo cartR.CartRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productRepo productR.ProductRepository,
	productImgRepo productR.ProductImageRepository,
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
	Product   productD.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryD.Inventory
	Images          struct {
		Thumbnail string
	}
}

type GetCartResult struct {
	Cart     *cartD.Cart
	Products map[uuid.UUID]ProductCartResponse
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
			Products: map[uuid.UUID]ProductCartResponse{},
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

	imagesMap, err := u.productImgRepo.ListByProductIDs(productIDs)
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
			key := images[0].Variants[productD.ResolutionThumbnail].Key
			result.Images.Thumbnail = u.fileStore.PublicURL(key, "public-assets")
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
