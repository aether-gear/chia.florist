package repository

import (
	"time"

	query "service-core/internal/shared/query"

	"github.com/google/uuid"
)

type FindUserParams struct {
	ID       *uuid.UUID
	Name     *string
	Username *string
	Email    *string

	query.Pagination
	query.Sorts
}

type CreateUserProps struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Phone     *string
	AvatarURL *string
	CreatedAt time.Time
}

type SaveProfileProps struct {
	UserID    uuid.UUID
	Name      *string
	Phone     *string
	AvatarURL *string
	UpdatedAt time.Time
}
