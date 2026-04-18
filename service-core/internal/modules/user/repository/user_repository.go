package repository

import (
	"service-core/internal/modules/user/domain"
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

type UserRepository interface {
	FindUsers(params FindUserParams) ([]domain.User, int, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	GetUserWithAccount(id uuid.UUID) (*UserWithAccount, error)
}
