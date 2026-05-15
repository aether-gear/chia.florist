package repository

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountProps struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
}

type CreateSessionProps struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Agent     string
	IPAddress string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
