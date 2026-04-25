package http

import (
	"github.com/google/uuid"
)

type CreatePaymentAccountRequest struct {
	MethodID      string  `json:"method_id"`
	AccountName   string  `json:"account_name"`
	AccountNumber *string `json:"account_number"`
	PhoneNumber   string  `json:"phone_number"`
	QRString      *string `json:"qr_string"`
	IsActive      string  `json:"is_active"`
}

type CreatePaymentMethodRequest struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	IsActive      string  `json:"is_active"`
	Description   string  `json:"description"`
	FeeType       string  `json:"fee_type"`
	FeeFixed      *string `json:"fee_amount"`
	FeePercentage *string `json:"fee_percentage"`
}

type PaymentMethodResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	IsActive      bool      `json:"is_active"`
	Description   string    `json:"description"`
	FeeType       string    `json:"fee_type"`
	FeeFixed      int64     `json:"fee_fixed"`
	FeePercentage float64   `json:"fee_percentage"`
}

type PaymentAccountResponse struct {
	ID            uuid.UUID `json:"id"`
	MethodID      uuid.UUID `json:"method_id"`
	AccountName   string    `json:"account_name"`
	AccountNumber *string   `json:"account_number"`
	PhoneNumber   string    `json:"phone_number"`
	QRString      *string   `json:"qr_string"`
}
