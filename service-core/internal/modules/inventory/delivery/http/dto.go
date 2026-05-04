package http

type CreateInventoryRequest struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Stock     int    `json:"stock"`
}
