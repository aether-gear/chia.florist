package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID uuid.UUID

	Name        string
	Slug        string
	Description *string

	IsActive bool

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
