package repository

import (
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
