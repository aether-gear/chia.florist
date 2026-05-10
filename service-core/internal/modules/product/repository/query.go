package repository

import (
	"io"
	inventoryD "service-core/internal/modules/inventory/domain"
	"service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type FindProductParams struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

type UploadProductImageParams struct {
	ProductID uuid.UUID
	Metadata  domain.ProductImageMetadata
	File      io.Reader
}

type ProductWithInventory struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryD.Inventory
}
