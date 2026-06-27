package http

import "time"

type createOrderItemRequest struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"name"`
	Quantity    int    `json:"quantity"`
}

type createOrderPaymentRequest struct {
	ID       string `json:"id"`
	IsManual bool   `json:"is_manual"`
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

type createOrderPaymentAccountResponse struct {
	AccountName   string  `json:"account_name"`
	AccountNumber *string `json:"account_number,omitempty"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	QRString      *string `json:"qr_string,omitempty"`
}

type createOrderResponse struct {
	OrderID        string                             `json:"order_id"`
	Instruction    *string                            `json:"instruction"`
	PaymentAccount *createOrderPaymentAccountResponse `json:"payment_account,omitempty"`
}

type orderItemResponse struct {
	ID               string  `json:"id"`
	ProductID        string  `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         int     `json:"quantity"`
	UnitPrice        int64   `json:"unit_price"`
	Subtotal         int64   `json:"subtotal"`
	ShopID           string  `json:"shop_id"`
	ShopName         string  `json:"shop_name"`
	CourierCode      *string `json:"courier_code,omitempty"`
	CourierService   *string `json:"courier_service,omitempty"`
	ShippingFeeTotal int64   `json:"shipping_fee"`
}

type orderResponse struct {
	ID          string              `json:"id"`
	Number      string              `json:"number"`
	CustomerID  string              `json:"customer_id"`
	AddressID   string              `json:"address_id"`
	Status      string              `json:"status"`
	Subtotal    int64               `json:"subtotal"`
	ShippingFee int64               `json:"shipping_fee"`
	Total       int64               `json:"total"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty"`
	Items       []orderItemResponse `json:"items"`
}
