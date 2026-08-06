package http

import (
	"encoding/json"
	"time"
)

// ---- Request DTOs ----

type createOrderItemRequest struct {
	ProductID          *string         `json:"product_id,omitempty"`
	CartItemID         *string         `json:"cart_item_id,omitempty"`
	ProductVariantType string          `json:"product_variant_type,omitempty"`
	ItemType           string          `json:"item_type,omitempty"`
	IsCustom           *bool           `json:"is_custom,omitempty"`
	ProductName        string          `json:"name"`
	Quantity           int             `json:"quantity"`
	CustomDesign       json.RawMessage `json:"custom_design,omitempty"`
}

type createOrderPaymentRequest struct {
	ID string `json:"id"`
}

type createOrderCourierRequest struct {
	Code    string `json:"code"`
	Service string `json:"service"`
}

type createOrderShopRequest struct {
	ShopID   string                     `json:"shop_id"`
	ShopName string                     `json:"name"`
	Courier  *createOrderCourierRequest `json:"selected_courier"`
	Items    []createOrderItemRequest   `json:"items"`
}

type createOrderRequest struct {
	AddressID       string                    `json:"address_id"`
	SelectedPayment createOrderPaymentRequest `json:"selected_payment"`
	Shops           []createOrderShopRequest  `json:"shops"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status"`

	// TrackingNumber is optional. It is used in manual logistics mode to
	// pre-set the shipment tracking number. Ignored when the server is
	// configured with an automated provider (e.g. Komerce).
	TrackingNumber *string `json:"tracking_number"`

	// FulfillmentMethod is optional. Chooses who delivers the order: "courier"
	// or "self_delivery". If not specified, it defaults to "courier".
	FulfillmentMethod *string `json:"fulfillment_method"`
}

// ---- Response DTOs ----

type createOrderPaymentAccountResponse struct {
	AccountName   string  `json:"account_name"`
	AccountNumber *string `json:"account_number,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	QRString      *string `json:"qr_string,omitempty"`
}

type createOrderResponse struct {
	OrderID        string                             `json:"order_id"`
	Instruction    *string                            `json:"instruction"`
	PaymentAccount *createOrderPaymentAccountResponse `json:"payment_account,omitempty"`
	ChannelData    *paymentChannelDataResponse        `json:"channel_data,omitempty"`
}

type orderItemResponse struct {
	ID                 string  `json:"id"`
	ProductID          *string `json:"product_id,omitempty"`
	ProductVariantType string  `json:"product_variant_type"`
	IsCustom           bool    `json:"is_custom"`
	ProductName        string  `json:"product_name"`
	Quantity           int     `json:"quantity"`
	UnitPrice          int64   `json:"unit_price"`
	Subtotal           int64   `json:"subtotal"`
	ShopID             string  `json:"shop_id"`
	ShopName           string  `json:"shop_name"`
	CourierCode        *string `json:"courier_code,omitempty"`
	CourierService     *string `json:"courier_service,omitempty"`
	ShippingFeeTotal   int64   `json:"shipping_fee"`
}

type paymentChannelDataResponse struct {
	ChannelType string     `json:"channel_type"`
	DisplayName string     `json:"display_name"`
	ActionURL   *string    `json:"action_url,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type paymentDetailResponse struct {
	ID          string                      `json:"id"`
	Status      string                      `json:"status"`
	Provider    string                      `json:"provider"`
	Amount      int64                       `json:"amount"`
	ExpiresAt   *time.Time                  `json:"expires_at,omitempty"`
	ChannelData *paymentChannelDataResponse `json:"channel_data,omitempty"`
	CreatedAt   time.Time                   `json:"created_at"`
}

type shipmentEventResponse struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
}

type shipmentDetailResponse struct {
	ID                string                  `json:"id"`
	Status            string                  `json:"status"`
	FulfillmentMethod string                  `json:"fulfillment_method"`
	Courier           string                  `json:"courier"`
	Service           string                  `json:"service"`
	TrackingNumber    *string                 `json:"tracking_number,omitempty"`
	Cost              int64                   `json:"cost"`
	CreatedAt         time.Time               `json:"created_at"`
	Events            []shipmentEventResponse `json:"events,omitempty"`
}

type orderAddressResponse struct {
	ID           string  `json:"id"`
	CustomerID   string  `json:"customer_id"`
	ReceiverName string  `json:"receiver_name"`
	Phone        *string `json:"phone,omitempty"`
	IsDefault    bool    `json:"is_default"`
	ProvinceID   string  `json:"province_id"`
	CityID       string  `json:"city_id"`
	DistrictID   string  `json:"district_id"`
	VillageID    string  `json:"village_id"`
	FullAddress  string  `json:"full_address"`
	PostalCode   string  `json:"postal_code"`
}

type orderResponse struct {
	ID          string                  `json:"id"`
	Number      string                  `json:"number"`
	CustomerID  string                  `json:"customer_id"`
	AddressID   string                  `json:"address_id"`
	Status      string                  `json:"status"`
	Subtotal    int64                   `json:"subtotal"`
	ShippingFee int64                   `json:"shipping_fee"`
	Total       int64                   `json:"total"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   *time.Time              `json:"updated_at,omitempty"`
	Items       []orderItemResponse     `json:"items"`
	Payment     *paymentDetailResponse  `json:"payment,omitempty"`
	Shipment    *shipmentDetailResponse `json:"shipment,omitempty"`
	Address     *orderAddressResponse   `json:"address,omitempty"`
}

type trackingTimelineEventResponse struct {
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
}

type orderTrackingResponse struct {
	OrderID        string                          `json:"order_id"`
	ShipmentID     string                          `json:"shipment_id"`
	Courier        string                          `json:"courier"`
	TrackingNumber *string                         `json:"tracking_number,omitempty"`
	Timeline       []trackingTimelineEventResponse `json:"timeline"`
}
