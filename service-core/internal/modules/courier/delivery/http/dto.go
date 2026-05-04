package http

type CourierSelectionRequest struct {
	Code      string `json:"code" validate:"required"`
	IsEnabled bool   `json:"is_enabled"`
}

type ConfigureCourierShopRequest struct {
	ShopID   string                    `json:"shop_id"`
	Couriers []CourierSelectionRequest `json:"couriers" validate:"required,min=1,dive"`
}
