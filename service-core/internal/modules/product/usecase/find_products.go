package usecase

import (
	"context"
	"fmt"
	"strings"

	storage "service-core/internal/infra/storage"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type FindProductsUsecase struct {
	productRepo    repository.ProductRepository
	inventoryRepo  inventoryRepo.InventoryRepository
	productImgRepo repository.ProductImageRepository
	fileStore      storage.Provider
	executor       transaction.Executor
}

func NewFindProductsUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryRepo.InventoryRepository,
	productImgRepo repository.ProductImageRepository,
	fileStore storage.Provider,
	executor transaction.Executor,
) *FindProductsUsecase {
	return &FindProductsUsecase{
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		productImgRepo: productImgRepo,
		fileStore:      fileStore,
		executor:       executor,
	}
}

type ProductCatalogResponse struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryDomain.Inventory
	Images          struct {
		Thumbnail string
	}
}

type FindProductsInput struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
	Sort  string
}

func (u *FindProductsUsecase) Execute(
	ctx context.Context,
	input FindProductsInput,
) (
	[]ProductCatalogResponse,
	int,
	error,
) {
	var productSortKeys = map[string]query.SortKey{
		"latest":   repository.ProductSortLatest,
		"date":     repository.ProductSortLatest,
		"name":     repository.ProductSortName,
		"price":    repository.ProductSortPrice,
		"weight":   repository.ProductSortWeight,
		"status":   repository.ProductSortStatus,
		"modified": repository.ProductSortModified,
		"archived": repository.ProductSortArchived,
	}

	var sorts query.Sorts
	if input.Sort != "" {
		parts := strings.SplitSeq(input.Sort, ",")
		for part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			subparts := strings.Split(part, ":")
			key := strings.TrimSpace(subparts[0])

			var dir query.SortDirection = query.SortDesc
			if len(subparts) > 1 {
				d := strings.ToLower(strings.TrimSpace(subparts[1]))
				if d == "asc" {
					dir = query.SortAsc
				}
			}

			sortKey, exists := productSortKeys[key]
			if exists {
				sorts = append(sorts, query.Sort{
					By:        sortKey,
					Direction: dir,
				})
			}
		}
	}

	if len(sorts) == 0 {
		sorts = query.Sorts{
			{
				By:        repository.ProductSortLatest,
				Direction: query.SortDesc,
			},
		}
	}

	params := repository.FindProductParams{
		ID:   input.ID,
		Name: input.Name,
		Pagination: query.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
		},
		Sorts: sorts,
	}

	products, total, err := u.productRepo.
		FindProducts(ctx, u.executor, params)
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

	inventoryMap, err := u.inventoryRepo.
		ListByProductIDs(ctx, u.executor, productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load inventory for products: %w", err)
	}

	imagesMap, err := u.productImgRepo.
		ListByProductIDs(ctx, u.executor, productIDs)
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
