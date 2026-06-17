package repository

import (
	"time"

	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

var (
	UserSortLatest    query.SortKey = "latest"
	UserSortName      query.SortKey = "name"
	UserSortUsername  query.SortKey = "username"
	UserSortPhone     query.SortKey = "phone"
	UserSortModify    query.SortKey = "modify"
	UserSortLastLogin query.SortKey = "last_login"
)

type FindUserParams struct {
	ID       *uuid.UUID
	Name     *string
	Username *string
	Email    *string

	query.Pagination
	query.Sorts
}

type UserWithAccount struct {
	ID          uuid.UUID
	Name        string
	Username    string
	Email       string
	Phone       *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	LastLoginAt *time.Time
}

type CreateUserProps struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Phone     *string
	CreatedAt time.Time
}
