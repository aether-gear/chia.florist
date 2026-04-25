package domain

import (
	"fmt"
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

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("invalid product: name is required")
	}

	if p.Price <= 0 {
		return fmt.Errorf("invalid product: price must be greater than 0")
	}

	return nil
}
