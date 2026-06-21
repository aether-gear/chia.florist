package repository

import (
	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

var (
	MerchantSortLatest query.SortKey = "latest"
	MerchantSortName   query.SortKey = "name"
	MerchantSortModify query.SortKey = "modified"
)

type FindMerchantParams struct {
	ID   *uuid.UUID
	Name *string

	query.Pagination
	query.Sorts
}
