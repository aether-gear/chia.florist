package usecase

import (
	"context"
	"fmt"

	"service-core/internal/infra/storage"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	productRepo    repository.ProductRepository
	inventoryRepo  inventoryRepo.InventoryRepository
	productImgRepo repository.ProductImageRepository
	fileStore      storage.Provider
	executor       transaction.Executor
}

func NewGetProductUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productImgRepo repository.ProductImageRepository,
	fileStore storage.Provider,
	executor transaction.Executor,
) *GetProductUsecase {
	return &GetProductUsecase{
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		productImgRepo: productImgRepo,
		fileStore:      fileStore,
		executor:       executor,
	}
}

type ImageProductDetail struct {
	Thumbnail string
	Detail    string
	Preview   string
}

type ProductDetailResponse struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryDomain.Inventory
	Images          []ImageProductDetail
}

func (u *GetProductUsecase) Execute(
	ctx context.Context,
	id uuid.UUID,
) (*ProductDetailResponse, error) {
	product, err := u.productRepo.GetByID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load products with inventory: %w", err)
	}
	if product == nil {
		return nil, nil
	}

	inventories, err := u.inventoryRepo.ListByProductID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory by product: %w", err)
	}

	images, err := u.productImgRepo.ListByProductID(ctx, u.executor, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load images by product: %w", err)
	}

	result := ProductDetailResponse{
		Product:         *product,
		ShopInventories: inventories,
	}

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

	for _, inventory := range inventories {
		result.Inventory.TotalStock += inventory.TotalStock
		result.Inventory.ReservedStock += inventory.ReservedStock
	}

	return &result, nil
}
