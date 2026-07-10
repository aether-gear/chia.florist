package midtrans

import (
	"encoding/json"
	"strings"

	paymentgateway "service-core/internal/infra/payment-gateway"
)

type channelStrategy interface {
	BuildPayload(req paymentgateway.ChargeRequest, body map[string]any)
	ParseInstructions(respBody []byte) ([]paymentgateway.PaymentInstruction, error)
}

// gopayStrategy handles GoPay payment method.
type gopayStrategy struct{}

func (s *gopayStrategy) BuildPayload(req paymentgateway.ChargeRequest, body map[string]any) {
	body["payment_type"] = "gopay"
	body["gopay"] = map[string]any{
		"enable_callback": true,
	}
}

func (s *gopayStrategy) ParseInstructions(respBody []byte) ([]paymentgateway.PaymentInstruction, error) {
	var resp gopayChargeAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	var instructions []paymentgateway.PaymentInstruction
	for _, act := range resp.Actions {
		if strings.EqualFold(act.Name, "deeplink-redirect") ||
			strings.EqualFold(act.Name, "generate-qr-code") {
			instructions = append(instructions, paymentgateway.PaymentInstruction{
				Type:  "ewallet",
				Label: act.Name,
				Value: act.URL,
			})
		}
	}
	return instructions, nil
}

// shopeepayStrategy handles ShopeePay payment method.
type shopeepayStrategy struct{}

func (s *shopeepayStrategy) BuildPayload(req paymentgateway.ChargeRequest, body map[string]any) {
	body["payment_type"] = "shopeepay"
}

func (s *shopeepayStrategy) ParseInstructions(respBody []byte) ([]paymentgateway.PaymentInstruction, error) {
	var resp shopeepayChargeAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	var instructions []paymentgateway.PaymentInstruction
	for _, act := range resp.Actions {
		if strings.EqualFold(act.Name, "deeplink-redirect") ||
			strings.EqualFold(act.Name, "generate-qr-code") {
			instructions = append(instructions, paymentgateway.PaymentInstruction{
				Type:  "ewallet",
				Label: act.Name,
				Value: act.URL,
			})
		}
	}
	return instructions, nil
}

// qrisStrategy handles QRIS and QR Code payment methods.
type qrisStrategy struct{}

func (s *qrisStrategy) BuildPayload(req paymentgateway.ChargeRequest, body map[string]any) {
	body["payment_type"] = "qris"
}

func (s *qrisStrategy) ParseInstructions(respBody []byte) ([]paymentgateway.PaymentInstruction, error) {
	var resp qrisChargeAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	var instructions []paymentgateway.PaymentInstruction
	if resp.QRString != "" {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "qris",
			Label: "QRIS",
			Value: resp.QRString,
		})
	}
	for _, act := range resp.Actions {
		if strings.EqualFold(act.Name, "deeplink-redirect") ||
			strings.EqualFold(act.Name, "generate-qr-code") {
			instructions = append(instructions, paymentgateway.PaymentInstruction{
				Type:  "ewallet",
				Label: act.Name,
				Value: act.URL,
			})
		}
	}
	return instructions, nil
}

// mandiriStrategy handles Mandiri virtual account/bill payment.
type mandiriStrategy struct{}

func (s *mandiriStrategy) BuildPayload(req paymentgateway.ChargeRequest, body map[string]any) {
	body["payment_type"] = "bank_transfer"
	body["bank_transfer"] = map[string]any{
		"bank": "mandiri",
	}
}

func (s *mandiriStrategy) ParseInstructions(respBody []byte) ([]paymentgateway.PaymentInstruction, error) {
	var resp bankTransferChargeAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	var instructions []paymentgateway.PaymentInstruction
	if resp.PermataVANumber != "" {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "bank_transfer",
			Label: "PERMATA Virtual Account",
			Value: resp.PermataVANumber,
		})
	}
	for _, va := range resp.VaNumbers {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "bank_transfer",
			Label: strings.ToUpper(va.Bank) + " Virtual Account",
			Value: va.VANumber,
		})
	}
	if resp.BillKey != "" {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "bank_transfer",
			Label: "MANDIRI Bill Key",
			Value: resp.BillKey,
		})
	}
	if resp.BillerCode != "" {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "bank_transfer",
			Label: "MANDIRI Biller Code",
			Value: resp.BillerCode,
		})
	}
	return instructions, nil
}
