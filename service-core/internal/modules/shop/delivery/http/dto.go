package http

import (
	"time"

	"github.com/google/uuid"
)

type CreateShopRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    string  `json:"is_active"`
}

type GetShopResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
