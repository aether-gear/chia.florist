package domain

import (
	"time"

	"github.com/google/uuid"
)

// PaymentChannelData holds the gateway-returned payment channel details
// (QR string, VA number, deep-link URL) persisted after a charge so they
// can be retrieved on any subsequent request — e.g. after a page refresh.
type PaymentChannelData struct {
	ID uuid.UUID

	PaymentID uuid.UUID

	// ChannelType is the payment method type, using the domain constants:
	// TypeBankTransfer | TypeEWallet | TypeQRCode
	ChannelType PaymentMethodType

	// DisplayName is the human-readable label returned by the gateway,
	// e.g. "QRIS", "GoPay", "BCA Virtual Account".
	DisplayName string

	// ActionURL is the actionable value: QR string, deep-link URL, or VA number.
	// Nil when the gateway returns no instruction value for this channel.
	ActionURL *string

	ExpiresAt *time.Time

	CreatedAt time.Time
}
