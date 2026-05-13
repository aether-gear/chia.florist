package usecase

import (
	"fmt"

	"service-core/internal/infra/storage"
	inventoryD "service-core/internal/modules/inventory/domain"
	inventoryR "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	productRepo    repository.ProductRepository
	inventoryRepo  inventoryR.InventoryRepository
	productImgRepo repository.ProductImageRepository
	fileStore      storage.Provider
}

func NewGetProductUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productImgRepo repository.ProductImageRepository,
	fileStore storage.Provider,
) *GetProductUsecase {
	return &GetProductUsecase{
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		productImgRepo: productImgRepo,
		fileStore:      fileStore,
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
	ShopInventories []inventoryD.Inventory
	Images          []ImageProductDetail
}

func (u *GetProductUsecase) Execute(id uuid.UUID) (*ProductDetailResponse, error) {
	product, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load products with inventory: %w", err)
	}
	if product == nil {
		return nil, nil
	}

	inventories, err := u.inventoryRepo.ListByProductID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory by product: %w", err)
	}

	images, err := u.productImgRepo.FindByProductID(id)
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
