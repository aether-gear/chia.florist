package domain

type ProductVariant struct {
	VariantID string `json:"variant_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
}
