package repository

import (
	"service-core/internal/modules/user/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindUsers(params FindUserParams) ([]domain.User, int, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	GetUserWithAccount(id uuid.UUID) (*UserWithAccount, error)
}
