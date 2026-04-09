package domain

type Product struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Category  string  `json:"category"`
	Thumbnail string  `json:"thumbnail"`
	Rating    float64 `json:"rating"`
}
