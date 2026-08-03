package paymentgateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChargeItem struct {
	ID       string
	Name     string
	Quantity int
	Price    int64
}

type ChargeRequest struct {
	// PaymentID is the internal Payment UUID
	// — used as the idempotency key external
	// order-reference sent to the gateway
	PaymentID uuid.UUID

	// OrderID is the staff's order reference
	// exposed to the gateway
	OrderID uuid.UUID

	// Amount is the gross amount
	// in the smallest currency unit
	// (e.g. IDR cents)
	Amount int64

	// PaymentType is the gateway-specific
	// payment channel identifier
	// (e.g. "bank_transfer", "gopay", "qris").
	PaymentType string

	// BankCode is only required when
	// PaymentType is "bank_transfer"
	BankCode string

	// CustomerName is optional;
	// included in payment instructions
	// when available.
	CustomerName string

	// optional
	CustomerEmail string
	CustomerPhone string

	// ExpiresAt is the desired
	// expiry time for the payment
	//
	// When zero the gateway will apply
	// its own default
	ExpiresAt time.Time

	// Items lists the individual products, shipping fees, or
	// adjustments that comprise the total transaction amount.
	Items []ChargeItem
}

// PaymentInstruction is a single actionable
// piece of payment guidance returned by the gateway
// (e.g. a virtual-account number or a QR string)
type PaymentInstruction struct {
	Type  string // "bank_transfer" | "qris" | "ewallet" | …
	Label string // human-readable label, e.g. "BCA Virtual Account"
	Value string // the account number, QR string, or deep-link URL
}

type ChargeResponse struct {
	// GatewayTransactionID is the unique
	// transaction identifier assigned by the
	// payment gateway
	// (e.g. Midtrans transaction_id)
	GatewayTransactionID string

	// GatewayOrderID is the order identifier
	// echoed back by the gateway
	GatewayOrderID string

	// PaymentType is the resolved payment-channel type
	// (as returned by the gateway,
	// may differ from the requested type for aliases)
	PaymentType string

	// GrossAmount is the total amount
	// as confirmed by the gateway
	GrossAmount int64

	// Status is the initial gateway
	// status of the transaction
	Status string

	// Instructions carries human-readable
	// payment instructions
	// (VA numbers, QR strings, deep-link URLs, etc.)
	Instructions []PaymentInstruction

	// ExpiresAt is the transaction expiry time
	// reported by the gateway
	//
	// May be zero when the gateway does not
	// return one
	ExpiresAt time.Time
}

type NotificationStatus string

const (
	NotificationStatusPending    NotificationStatus = "pending"
	NotificationStatusSettlement NotificationStatus = "settlement"
	NotificationStatusExpire     NotificationStatus = "expire"
	NotificationStatusCancel     NotificationStatus = "cancel"
	NotificationStatusDeny       NotificationStatus = "deny"
	NotificationStatusRefund     NotificationStatus = "refund"
	NotificationStatusChallenge  NotificationStatus = "challenge"
)

// NotificationResult is the normalised result
// of parsing a gateway webhook
type NotificationResult struct {
	// GatewayTransactionID is the gateway-side
	// transaction identifier
	GatewayTransactionID string

	// GatewayOrderID maps back to the internal PaymentID
	// sent as the order reference at charge time
	GatewayOrderID string

	// Status is the normalised transaction status
	Status NotificationStatus

	// GrossAmount is the amount confirmed
	// by the gateway
	GrossAmount int64

	// FraudStatus is the anti-fraud result
	// Empty when not applicable
	// e.g. ("accept", "challenge", "deny")
	FraudStatus string

	// RawStatus is the unmodified status string
	// from the gateway
	RawStatus string
}

// NotificationPayload is the raw webhook
// body forwarded from the gateway
//
// It is intentionally generic
// so each provider can deserialise it itself
type NotificationPayload map[string]any

type AllowedPaymentMethod struct {
	Code          string
	Name          string
	Type          string // "bank_transfer" | "ewallet" | "qr_code"
	FeeType       string // "flat" | "percentage" | "mixed"
	FeeFixed      int64
	FeePercentage float64
	Description   string
}

type Provider interface {
	// Name returns the unique provider identifier.
	Name() string

	// AllowedPaymentMethods returns the list of payment methods supported by this provider.
	AllowedPaymentMethods() []AllowedPaymentMethod

	// Supports returns true if the gateway provider is configured
	// and capable of handling the given payment channel code.
	Supports(code string) bool

	// Charge creates a new payment transaction
	// with the gateway and returns the instructions
	// the customer needs to complete the payment
	Charge(
		ctx context.Context,
		req ChargeRequest,
	) (*ChargeResponse, error)

	// ParseNotification deserialises and validates an
	// inbound webhook payload from the gateway
	// and returns a normalised NotificationResult.
	ParseNotification(
		ctx context.Context,
		payload NotificationPayload,
	) (*NotificationResult, error)

	// GetTransactionStatus fetches the current status
	// of a transaction directly from the gateway using
	// its gateway-side order ID.
	//
	// Used by the reconciliation job and the customer-triggered
	// sync when no webhook was received (e.g. service was unreachable
	// during Midtrans's delivery attempts).
	GetTransactionStatus(
		ctx context.Context,
		gatewayOrderID string,
	) (*NotificationResult, error)

	// CancelTransaction requests the gateway to
	// cancel / void a pending transaction
	// identified by its gateway-side order ID.
	CancelTransaction(
		ctx context.Context,
		gatewayOrderID string,
	) error
}
