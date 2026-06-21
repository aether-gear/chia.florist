package midtrans

import "fmt"

// MapPaymentType converts an internal payment method name
// into the Midtrans Core API payment_type string required
// by ChargeRequest.PaymentType.
//
// Returns an error for method names that are not supported
// by this provider, so the caller can reject the request early
// rather than sending an invalid charge to the gateway.
func MapPaymentType(methodName string) (string, error) {
	switch methodName {
	case "QRIS":
		return "qris", nil
	case "GoPay":
		return "gopay", nil
	case "DANA":
		return "dana", nil
	case "ShopeePay":
		return "shopeepay", nil
	case "Mandiri":
		return "echannel", nil
	default:
		return "", fmt.Errorf("midtrans: unsupported payment method %q", methodName)
	}
}
