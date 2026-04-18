package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Username string
	Phone    *string

	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	LastLoginAt *time.Time
}
