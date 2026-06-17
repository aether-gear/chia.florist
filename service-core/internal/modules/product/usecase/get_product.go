package usecase

import (
	"context"
	"fmt"

	"service-core/internal/infra/storage"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	executor       transaction.Executor
	fileStore      storage.Provider
	productRepo    repository.ProductRepository
	inventoryRepo  inventoryRepo.InventoryRepository
	productImgRepo repository.ProductImageRepository
	shopRepo       shopRepo.ShopRepository
}

func NewGetProductUsecase(
	executor transaction.Executor,
	fileStore storage.Provider,
	productRepo repository.ProductRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productImgRepo repository.ProductImageRepository,
	shopRepo shopRepo.ShopRepository,
) *GetProductUsecase {
	return &GetProductUsecase{
		executor:       executor,
		fileStore:      fileStore,
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		productImgRepo: productImgRepo,
		shopRepo:       shopRepo,
	}
}

type ImageProductDetail struct {
	Thumbnail string
	Detail    string
	Preview   string
}

type ProductDetailResult struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryDomain.Inventory
	Images          []ImageProductDetail
	Availability    []ShopAvailabilityResult
}

func (u *GetProductUsecase) Execute(
	ctx context.Context,
	id uuid.UUID,
) (*ProductDetailResult, error) {
	product, err := u.productRepo.GetByID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load products with inventory: %w", err)
	}
	if product == nil {
		return nil, nil
	}

	inventories, err := u.inventoryRepo.
		ListByProductID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory by product: %w", err)
	}

	var shopIDs []uuid.UUID
	for _, inv := range inventories {
		shopIDs = append(shopIDs, inv.ShopID)
	}

	shops, err := u.shopRepo.
		FindByIDs(ctx, u.executor, shopIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load shops: %w", err)
	}

	shopsMap := make(map[uuid.UUID]shopDomain.Shop)
	for _, s := range shops {
		shopsMap[s.ID] = s
	}

	images, err := u.productImgRepo.
		ListByProductID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load images by product: %w", err)
	}

	result := ProductDetailResult{
		Product:         *product,
		ShopInventories: inventories,
	}

	var (
		totalStock    = 0
		reservedStock = 0
		availability  []ShopAvailabilityResult
	)

	for _, inventory := range inventories {
		totalStock += inventory.TotalStock
		reservedStock += inventory.ReservedStock

		shop, ok := shopsMap[inventory.ShopID]
		if !ok {
			continue
		}

		availability = append(availability, ShopAvailabilityResult{
			ShopName: shop.Name,
			ShopSlug: shop.Slug,
			Stock:    inventory.TotalStock,
		})
	}

	result.Inventory.TotalStock = totalStock
	result.Inventory.ReservedStock = reservedStock
	result.Availability = availability

	if len(images) > 0 {
		result.Images = make([]ImageProductDetail, 0, len(images))

		for _, img := range images {
			imageDetail := ImageProductDetail{}

			if variant, ok := img.Variants[domain.ResolutionThumbnail]; ok {
				imageDetail.Thumbnail = u.fileStore.PublicURL(
					variant.Key,
					"public-assets",
				)
			}

			if variant, ok := img.Variants[domain.ResolutionPreview]; ok {
				imageDetail.Preview = u.fileStore.PublicURL(
					variant.Key,
					"public-assets",
				)
			}

			if variant, ok := img.Variants[domain.ResolutionDetail]; ok {
				imageDetail.Detail = u.fileStore.PublicURL(
					variant.Key,
					"public-assets",
				)
			}

			result.Images = append(result.Images, imageDetail)
		}
	}

	return &result, nil
}
