package http

import "github.com/google/uuid"

type addItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartResponse struct {
	CartID uuid.UUID
	Items  []CartItemView
	Total  int64
}

type CartItemView struct {
	ProductID uuid.UUID
	Name      string
	Price     int64
	Subtotal  int64
	Quantity  int
}
