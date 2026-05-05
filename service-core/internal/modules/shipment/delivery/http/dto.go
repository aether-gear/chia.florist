package http

type EstimateShippingCostRequest struct {
	Origin      int      `json:"origin"`
	Destination int      `json:"destination"`
	Weight      int      `json:"weight"`
	Couriers    []string `json:"couriers"`
	PriceFilter *string  `json:"price_filter"`
}

type EstimateShippingCostResponse struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int64  `json:"cost"`
	Etd         string `json:"etd"`
}
