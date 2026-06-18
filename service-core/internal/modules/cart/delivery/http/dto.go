package http

import "github.com/google/uuid"

type addItemRequest struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemRequest struct {
	Quantity int `json:"quantity"`
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

type checkoutItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type selectedCourierRequest struct {
	Code    string `json:"code"`
	Service string `json:"service"`
}

type checkoutShopRequest struct {
	ShopID  string                  `json:"shop_id"`
	Items   []checkoutItemRequest   `json:"items"`
	Courier *selectedCourierRequest `json:"courier"`
}

type checkoutRequest struct {
	AddressID *string               `json:"address_id"`
	Shops     []checkoutShopRequest `json:"shops"`
}

type checkoutAddressResponse struct {
	ID            uuid.UUID `json:"id"`
	RecipientName string    `json:"recipient_name"`
	Phone         *string   `json:"phone"`
	FullAddress   string    `json:"full_address"`
}

type checkoutCouriersResponse struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Service string `json:"service"`
	ETD     string `json:"etd"`
	Fee     int64  `json:"fee"`
}

type checkoutItemResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	ShopID    uuid.UUID `json:"shop_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  int64     `json:"subtotal"`
}

type selectedCourierResponse struct {
	Code    string `json:"code"`
	Service string `json:"service"`
	Fee     int64  `json:"fee"`
}

type shopResponse struct {
	ShopID          uuid.UUID                  `json:"shop_id"`
	ShopName        string                     `json:"name"`
	ShopSlug        string                     `json:"slug"`
	Subtotal        int64                      `json:"subtotal"`
	Total           *int64                     `json:"total"`
	SelectedCourier *selectedCourierResponse   `json:"selected_courier"`
	Items           []checkoutItemResponse     `json:"items"`
	CostCouriers    []checkoutCouriersResponse `json:"cost_couriers"`
}

type checkoutResponse struct {
	Address       checkoutAddressResponse `json:"address"`
	Shops         []shopResponse          `json:"shops"`
	Subtotal      int64                   `json:"subtotal"`
	TotalShipping int64                   `json:"total_shipping"`
	TotalAll      *int64                  `json:"total"`
}
