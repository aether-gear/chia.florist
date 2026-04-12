package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description *string
	Status      ProductStatus

	Price  int64
	Weight *float64

	CreatedAt  time.Time
	UpdatedAt  *time.Time
	ArchivedAt *time.Time
	DeletedAt  *time.Time
}
