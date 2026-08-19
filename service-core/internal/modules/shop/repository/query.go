package repository

import (
	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

var (
	ShopSortLatest query.SortKey = "latest"
	ShopSortName   query.SortKey = "name"
	ShopSortActive query.SortKey = "active"
	ShopSortModify query.SortKey = "modify"
)

type FindShopsParams struct {
	ID             *string
	ShopIDs        []uuid.UUID
	Name           *string
	IsActive       *bool
	ApprovalStatus *string

	Pagination query.Pagination
	Sorts      query.Sorts
}


