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

type UploadProductImagesParams struct {
	ProductID uuid.UUID
	Files     []ImageFile
}

type ImageFile struct {
	File io.Reader
	domain.ProductImageMetadata
}

type ProductWithInventory struct {
	Product   domain.Product
	Inventory struct {
		TotalStock    int
		ReservedStock int
	}
	ShopInventories []inventoryD.Inventory
}

type UploadedProductImage struct {
	Sequence int

	CatalogURL string
	CartURL    string
	DetailURL  string

	IsPrimary bool
}
