package repository

import (
	"service-core/internal/modules/auth/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetByEmail(email string) (*domain.Account, error)
	GetByID(id uuid.UUID) (*domain.Account, error)

	Create(account CreateAccountProps) error
	// UpdateLastLogin(id uuid.UUID) error
}
