package http

import "time"

type estimateShippingOptionsRequest struct {
	Origin      int      `json:"origin"`
	Destination int      `json:"destination"`
	Weight      int      `json:"weight"`
	Couriers    []string `json:"couriers"`
	PriceFilter *string  `json:"price_filter"`
}

type estimateShippingOptionsResponse struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int64  `json:"cost"`
	Etd         string `json:"etd"`
}

type updateShipmentStatusRequest struct {
	Status      string  `json:"status"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type updateShipmentRequest struct {
	TrackingNumber *string `json:"tracking_number"`
	Courier        *string `json:"courier"`
	Service        *string `json:"service"`
}

type shipmentResponse struct {
	ID                string     `json:"id"`
	OrderID           string     `json:"order_id"`
	Status            string     `json:"status"`
	FulfillmentMethod string     `json:"fulfillment_method"`
	TrackingNumber    *string    `json:"tracking_number,omitempty"`
	Courier           string     `json:"courier"`
	Service           string     `json:"service"`
	Cost              int64      `json:"cost"`
	Weight            int        `json:"weight"`
	CreatedAt         time.Time  `json:"created_at"`
}
