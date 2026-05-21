package http

type courierSelectionRequest struct {
	Code      string `json:"code" validate:"required"`
	IsEnabled bool   `json:"is_enabled"`
}

type configureCourierShopRequest struct {
	ShopID   string                    `json:"shop_id"`
	Couriers []courierSelectionRequest `json:"couriers" validate:"required,min=1,dive"`
}
