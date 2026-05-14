package repository

import (
	"service-core/internal/modules/user/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindUsers(params FindUserParams) ([]domain.User, int, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	CreateUser(props CreateUserProps) error
}
