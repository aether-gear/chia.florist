package http

import (
	"time"

	"github.com/google/uuid"
)

type createPaymentAccountRequest struct {
	MethodID      string  `json:"method_id"`
	AccountName   string  `json:"account_name"`
	AccountNumber *string `json:"account_number"`
	PhoneNumber   string  `json:"phone_number"`
	QRString      *string `json:"qr_string"`
	IsActive      string  `json:"is_active"`
}

type savePaymentMethodRequest struct {
	ID            *string `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Provider      string  `json:"provider"`
	Type          string  `json:"type"`
	IsActive      string  `json:"is_active"`
	Description   string  `json:"description"`
	FeeType       string  `json:"fee_type"`
	FeeFixed      *string `json:"fee_amount"`
	FeePercentage *string `json:"fee_percentage"`
}

type paymentInstructionResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type paymentMethodResponse struct {
	ID            uuid.UUID                   `json:"id"`
	Name          string                      `json:"name"`
	Code          string                      `json:"code"`
	Provider      string                      `json:"provider"`
	Type          string                      `json:"type"`
	IsActive      bool                        `json:"is_active"`
	Description   string                      `json:"description"`
	FeeType       string                      `json:"fee_type"`
	FeeFixed      int64                       `json:"fee_fixed"`
	FeePercentage float64                     `json:"fee_percentage"`
	Instruction   *paymentInstructionResponse `json:"instruction,omitempty"`
}

type paymentAccountResponse struct {
	ID            uuid.UUID `json:"id"`
	MethodID      uuid.UUID `json:"method_id"`
	AccountName   string    `json:"account_name"`
	AccountNumber *string   `json:"account_number"`
	PhoneNumber   *string   `json:"phone_number"`
	QRString      *string   `json:"qr_string"`
}

type manualPaymentActionRequest struct {
	Action string `json:"action"`
}

type savePaymentInstructionRequest struct {
	Content string `json:"content"`
}

type getPaymentDetailResponse struct {
	PaymentID   string     `json:"payment_id"`
	Status      string     `json:"status"`
	Amount      int64      `json:"amount"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ChannelType *string    `json:"channel_type,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
	ActionURL   *string    `json:"action_url,omitempty"`
	// For manual payments:
	AccountName   *string `json:"account_name,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	QRString      *string `json:"qr_string,omitempty"`
	// Rendered instruction markdown:
	Instruction *string `json:"instruction,omitempty"`
}

type checkPaymentStatusResponse struct {
	Status string `json:"status"`
	Synced bool   `json:"synced"`
}
