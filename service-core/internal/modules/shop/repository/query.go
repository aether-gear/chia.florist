package repository

import query "service-core/internal/shared/query"

var (
	ShopSortLatest query.SortKey = "latest"
	ShopSortName   query.SortKey = "name"
	ShopSortActive query.SortKey = "active"
	ShopSortModify query.SortKey = "modify"
)

type FindShopsParams struct {
	ID       *string
	Name     *string
	IsActive *bool

	Pagination query.Pagination
	Sorts      query.Sorts
}
