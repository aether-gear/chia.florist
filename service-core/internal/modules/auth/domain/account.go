package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Email    string
	Password string

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
