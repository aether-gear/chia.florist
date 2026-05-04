package http

import "github.com/google/uuid"

type addItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Quantity  int    `json:"quantity"`
}

type CartResponse struct {
	CartID uuid.UUID      `json:"cart_id"`
	Items  []CartItemView `json:"items"`
	Total  int64          `json:"total"`
}

type CartItemView struct {
	ProductID uuid.UUID `json:"product_id"`
	ShopID    uuid.UUID `json:"shop_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Subtotal  int64     `json:"subtotal"`
	Quantity  int       `json:"quantity"`
}
