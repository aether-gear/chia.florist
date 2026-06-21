package http

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
