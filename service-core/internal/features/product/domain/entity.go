package domain

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       *int
	Category    string
	Stock       int

	Variants []ProductVariant
	Images   []string

	Rating       *float64
	ReviewsCount *int

	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
	DeletedAt  *time.Time
}

type ProductVariant struct {
	VariantID string `json:"variant_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
}
