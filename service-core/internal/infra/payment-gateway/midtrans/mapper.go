package midtrans

import (
	"fmt"
	"strings"
)

// MapPaymentType converts an internal payment method name
// into the Midtrans Core API payment_type string required
// by ChargeRequest.PaymentType.
//
// Returns an error for method names that are not supported
// by this provider, so the caller can reject the request early
// rather than sending an invalid charge to the gateway.
func MapPaymentType(methodName string) (string, error) {
	methodName = strings.ToLower(methodName)

	switch methodName {
	case "qris":
		return "qris", nil
	case "gopay":
		return "gopay", nil
	case "shopeepay":
		return "shopeepay", nil
	case "mandiri":
		return "echannel", nil
	default:
		return "", fmt.Errorf("midtrans: unsupported payment method %q", methodName)
	}
}
