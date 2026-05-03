package repository

import (
	"time"

	"github.com/google/uuid"
)

type FindUserParams struct {
	Page     int
	Limit    int
	ID       *uuid.UUID
	Name     *string
	Username *string
	Email    *string
}

type UserWithAccount struct {
	ID       uuid.UUID
	Name     string
	Username string
	Email    string
	Phone    *string

	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	LastLoginAt *time.Time
}
