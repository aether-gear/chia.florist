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

type FindProductsUsecase struct {
	productRepo    repository.ProductRepository
	inventoryRepo  inventoryR.InventoryRepository
	productImgRepo repository.ProductImageRepository
	fileStore      storage.Provider
}

func NewFindProductsUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryR.InventoryRepository,
	productImgRepo repository.ProductImageRepository,
	fileStore storage.Provider,
) *FindProductsUsecase {
	return &FindProductsUsecase{
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		productImgRepo: productImgRepo,
		fileStore:      fileStore,
	}
}

type ProductCatalogResponse struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryD.Inventory
	Images          struct {
		Thumbnail string
	}
}

type FindProductsInput struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

func (u *FindProductsUsecase) Execute(
	input FindProductsInput,
) (
	[]ProductCatalogResponse,
	int,
	error,
) {
	products, total, err := u.productRepo.FindProducts(repository.FindProductParams{
		Page:  input.Page,
		Limit: input.Limit,
		ID:    input.ID,
		Name:  input.Name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load products: %w", err)
	}
	if len(products) == 0 {
		return []ProductCatalogResponse{}, total, nil
	}

	productIDs := make([]uuid.UUID, 0, len(products))
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	inventoryMap, err := u.inventoryRepo.ListByProductIDs(productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load inventory for products: %w", err)
	}

	imagesMap, err := u.productImgRepo.ListByProductIDs(productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load images for products: %w", err)
	}

	results := make([]ProductCatalogResponse, 0, len(products))
	for _, product := range products {
		inventories := inventoryMap[product.ID]
		images := imagesMap[product.ID]

		result := ProductCatalogResponse{
			Product:         product,
			ShopInventories: inventories,
		}

		if len(images) > 0 {
			key := images[0].Variants[domain.ResolutionThumbnail].Key
			result.Images.Thumbnail = u.fileStore.PublicURL(key, "public-assets")
		}

		for _, inventory := range inventories {
			result.Inventory.TotalStock += inventory.TotalStock
			result.Inventory.ReservedStock += inventory.ReservedStock
		}

		results = append(results, result)
	}

	return results, total, nil
}
