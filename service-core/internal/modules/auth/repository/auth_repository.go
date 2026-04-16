package repository

import (
	"service-core/internal/modules/auth/domain"
	"time"

	"github.com/google/uuid"
)

type CreateAccountProps struct {
	ID           uuid.UUID
	Name         string
	Username     string
	Email        string
	PasswordHash string
	Phone        *string
	CreatedAt    time.Time
}

type AuthRepository interface {
	GetByEmail(email string) (*domain.Account, error)
	GetByID(id uuid.UUID) (*domain.Account, error)

	Create(account CreateAccountProps) error
	// UpdateLastLogin(id uuid.UUID) error
}
