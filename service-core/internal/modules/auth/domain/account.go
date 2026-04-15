package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Email    string
	Password string

	LastLoginAt *time.Time
}
