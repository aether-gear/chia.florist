package repository

import (
	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

var (
	CustomerSortLatest    query.SortKey = "latest"
	CustomerSortName      query.SortKey = "name"
	CustomerSortUsername  query.SortKey = "username"
	CustomerSortPhone     query.SortKey = "phone"
	CustomerSortModify    query.SortKey = "modify"
	CustomerSortLastLogin query.SortKey = "last_login"
)

type FindCustomerParams struct {
	ID       *uuid.UUID
	Name     *string
	Username *string
	Email    *string

	query.Pagination
	query.Sorts
}
