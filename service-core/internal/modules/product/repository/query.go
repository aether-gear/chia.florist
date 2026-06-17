package repository

import (
	"io"

	"service-core/internal/modules/product/domain"
	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

var (
	ProductSortLatest   query.SortKey = "latest"
	ProductSortName     query.SortKey = "name"
	ProductSortPrice    query.SortKey = "price"
	ProductSortWeight   query.SortKey = "weight"
	ProductSortStatus   query.SortKey = "status"
	ProductSortModified query.SortKey = "modified"
	ProductSortArchived query.SortKey = "archived"
)

type FindProductParams struct {
	ID   *string
	Name *string

	Pagination query.Pagination
	Sorts      query.Sorts
}

type UploadProductImagesParams struct {
	ProductID uuid.UUID
	Files     []ImageFile
}

type ImageFile struct {
	File io.Reader
	domain.ProductImageMetadata
}

type UploadedProductImage struct {
	Sequence int

	CatalogURL string
	CartURL    string
	DetailURL  string

	IsPrimary bool
}
