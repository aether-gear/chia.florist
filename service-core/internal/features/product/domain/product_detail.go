package domain

type ProductDetail struct {
	ProductID    string           `json:"product_id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Price        int              `json:"price"`
	Category     string           `json:"category"`
	Stock        int              `json:"stock"`
	Variants     []ProductVariant `json:"variants"`
	Images       []string         `json:"images"`
	Rating       float64          `json:"rating"`
	ReviewsCount int              `json:"reviews_count"`
	CreatedAt    string           `json:"created_at"`
}
