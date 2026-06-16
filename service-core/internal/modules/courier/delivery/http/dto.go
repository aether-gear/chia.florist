package http

type courierSelectionRequest struct {
	Code      string `json:"code"`
	IsEnabled bool   `json:"is_enabled"`
}

type configureCourierShopRequest struct {
	Couriers []courierSelectionRequest `json:"couriers"`
}
