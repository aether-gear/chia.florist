package repository

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetByEmail(email string) (*domain.Account, error)
	GetByID(id uuid.UUID) (*domain.Account, error)

	Create(account domain.Account) error
	// UpdateLastLogin(id uuid.UUID) error
}

type SessionRepository interface {
	Create(session domain.Session) error
}
