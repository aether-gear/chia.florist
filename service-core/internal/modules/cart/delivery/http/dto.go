package http

import "github.com/google/uuid"

type addItemRequest struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemRequest struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Quantity  int    `json:"quantity"`
}

type cartResponse struct {
	CartID uuid.UUID      `json:"cart_id"`
	Items  []cartItemView `json:"items"`
	Total  int64          `json:"total"`
}

type productImageResponse struct {
	Thumbnail *string `json:"thumbnail,omitempty"`
	Preview   *string `json:"preview,omitempty"`
	Detail    *string `json:"detail,omitempty"`
}

type cartItemView struct {
	ProductID uuid.UUID            `json:"product_id"`
	ShopID    uuid.UUID            `json:"shop_id"`
	Name      string               `json:"name"`
	Price     int64                `json:"price"`
	Subtotal  int64                `json:"subtotal"`
	Quantity  int                  `json:"quantity"`
	Image     productImageResponse `json:"images"`
}
