package midtrans

import (
	"context"
	"fmt"
	"strings"
	"time"

	paymentgateway "service-core/internal/infra/payment-gateway"
	config "service-core/internal/shared/config"

	midtranssdk "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// MidtransProvider implements paymentgateway.Provider
// using the Midtrans Core API
//
// It supports bank-transfer (VA), QRIS and e-wallet charges
type MidtransProvider struct {
	cfg    config.MidTransConfig
	client coreapi.Client
}

// NewMidtransProvider constructs a ready-to-use MidtransProvider
//
// It validates that a ServerKey is present before returning
func NewMidtransProvider(
	cfg config.MidTransConfig,
) (*MidtransProvider, error) {
	if strings.TrimSpace(cfg.ServerKey) == "" {
		return nil, fmt.Errorf("midtrans: server key is required")
	}

	env := midtranssdk.Sandbox
	if cfg.IsProduction {
		env = midtranssdk.Production
	}

	c := coreapi.Client{}
	c.New(cfg.ServerKey, env)

	return &MidtransProvider{
		cfg:    cfg,
		client: c,
	}, nil
}

// Charge creates a new Midtrans Core API transaction
// and returns the payment instructions the customer
// must follow to complete the payment
func (p *MidtransProvider) Charge(
	ctx context.Context,
	req paymentgateway.ChargeRequest,
) (*paymentgateway.ChargeResponse, error) {
	chargeReq, err := p.buildChargeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans: build charge request: %w", err)
	}

	resp, midErr := p.client.ChargeTransaction(chargeReq)
	if midErr != nil {
		return nil, fmt.Errorf("midtrans: charge transaction: %s (status %d)", midErr.Message, midErr.StatusCode)
	}

	return p.mapChargeResponse(resp)
}

// ParseNotification validates and normalises an inbound
// Midtrans notification webhook payload into a
// provider-agnostic NotificationResult.
//
// The payload is the raw JSON body decoded into a map
// (order_id, transaction_id, transaction_status,
// gross_amount, fraud_status, …)
//
// Re-query the gateway to verify authenticity
// before trusting the values (prevents spoofing)
func (p *MidtransProvider) ParseNotification(
	_ context.Context,
	payload paymentgateway.NotificationPayload,
) (*paymentgateway.NotificationResult, error) {
	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		return nil, fmt.Errorf("midtrans: parse notification: missing order_id in payload")
	}

	result, midErr := p.client.CheckTransaction(orderID)
	if midErr != nil {
		return nil, fmt.Errorf("midtrans: check transaction %q: %s (status %d)",
			orderID, midErr.Message, midErr.StatusCode)
	}

	grossAmount, err := parseAmount(result.GrossAmount)
	if err != nil {
		return nil, fmt.Errorf("midtrans: parse notification gross_amount: %w", err)
	}

	return &paymentgateway.NotificationResult{
		GatewayTransactionID: result.TransactionID,
		GatewayOrderID:       result.OrderID,
		Status:               mapNotificationStatus(result.TransactionStatus, result.FraudStatus),
		GrossAmount:          grossAmount,
		FraudStatus:          result.FraudStatus,
		RawStatus:            result.TransactionStatus,
	}, nil
}

// CancelTransaction requests Midtrans
// to cancel a pending / authorised transaction
// identified by its gateway-side order ID
func (p *MidtransProvider) CancelTransaction(
	_ context.Context,
	gatewayOrderID string,
) error {
	_, midErr := p.client.CancelTransaction(gatewayOrderID)
	if midErr != nil {
		return fmt.Errorf("midtrans: cancel transaction %q: %s (status %d)",
			gatewayOrderID, midErr.Message, midErr.StatusCode)
	}
	return nil
}

// buildChargeRequest converts our generic ChargeRequest into a Midtrans
// coreapi.ChargeReq. Add more payment-type branches here as needed
func (p *MidtransProvider) buildChargeRequest(
	req paymentgateway.ChargeRequest,
) (*coreapi.ChargeReq, error) {
	txDetails := midtranssdk.TransactionDetails{
		OrderID:  req.PaymentID.String(), // payment UUID as the idempotency key
		GrossAmt: req.Amount,
	}

	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: txDetails,
		CustomerDetails: &midtranssdk.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
	}

	// Attach expiry when explicitly set.
	if !req.ExpiresAt.IsZero() {
		dur := time.Until(req.ExpiresAt)
		if dur > 0 {
			minutes := int(dur.Minutes())
			if minutes < 1 {
				minutes = 1
			}
			chargeReq.CustomExpiry = &coreapi.CustomExpiry{
				ExpiryDuration: minutes,
				Unit:           "minute",
			}
		}
	}

	paymentType := strings.ToLower(req.PaymentType)

	switch paymentType {
	case "bank_transfer":
		bankCode := strings.ToUpper(req.BankCode)
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtranssdk.Bank(bankCode),
		}

	case "qris", "qr_code":
		chargeReq.PaymentType = coreapi.PaymentTypeQris

	case "gopay":
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
		chargeReq.Gopay = &coreapi.GopayDetails{
			EnableCallback: true,
		}

	case "shopeepay":
		chargeReq.PaymentType = coreapi.PaymentTypeShopeepay
		chargeReq.ShopeePay = &coreapi.ShopeePayDetails{
			CallbackUrl: "",
		}

	default:
		return nil, fmt.Errorf("unsupported payment type %q", req.PaymentType)
	}

	return chargeReq, nil
}

// mapChargeResponse converts a Midtrans coreapi.ChargeResponse into our
// provider-agnostic ChargeResponse
func (p *MidtransProvider) mapChargeResponse(
	resp *coreapi.ChargeResponse,
) (*paymentgateway.ChargeResponse, error) {
	grossAmount, err := parseAmount(resp.GrossAmount)
	if err != nil {
		return nil, fmt.Errorf("parse gross_amount: %w", err)
	}

	var expiresAt time.Time
	if resp.ExpiryTime != "" {
		// Midtrans returns expiry in "2006-01-02 15:04:05" format (WIB / UTC+7).
		wib, _ := time.LoadLocation("Asia/Jakarta")
		t, parseErr := time.ParseInLocation("2006-01-02 15:04:05", resp.ExpiryTime, wib)
		if parseErr == nil {
			expiresAt = t
		}
	}

	out := &paymentgateway.ChargeResponse{
		GatewayTransactionID: resp.TransactionID,
		GatewayOrderID:       resp.OrderID,
		PaymentType:          resp.PaymentType,
		GrossAmount:          grossAmount,
		Status:               resp.TransactionStatus,
		ExpiresAt:            expiresAt,
	}

	// Extract payment instructions depending on the channel.
	out.Instructions = extractInstructions(resp)

	return out, nil
}

// extractInstructions pulls VA numbers, QR strings,
// or deep-link actions from the raw Midtrans response
// and normalises them into PaymentInstruction slices
func extractInstructions(resp *coreapi.ChargeResponse) []paymentgateway.PaymentInstruction {
	var instructions []paymentgateway.PaymentInstruction

	if resp.QRString != "" {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "qris",
			Label: "QRIS",
			Value: resp.QRString,
		})
	}

	for _, va := range resp.VaNumbers {
		instructions = append(instructions, paymentgateway.PaymentInstruction{
			Type:  "bank_transfer",
			Label: strings.ToUpper(va.Bank) + " Virtual Account",
			Value: va.VANumber,
		})
	}

	for _, action := range resp.Actions {
		if strings.EqualFold(action.Name, "deeplink-redirect") ||
			strings.EqualFold(action.Name, "generate-qr-code") {
			instructions = append(instructions, paymentgateway.PaymentInstruction{
				Type:  "ewallet",
				Label: action.Name,
				Value: action.URL,
			})
		}
	}

	return instructions
}

// mapNotificationStatus converts Midtrans transaction_status + fraud_status
// into our normalised NotificationStatus.
func mapNotificationStatus(txStatus, fraudStatus string) paymentgateway.NotificationStatus {
	switch txStatus {
	case "capture":
		if fraudStatus == "challenge" {
			return paymentgateway.NotificationStatusChallenge
		}
		return paymentgateway.NotificationStatusSettlement
	case "settlement":
		return paymentgateway.NotificationStatusSettlement
	case "pending":
		return paymentgateway.NotificationStatusPending
	case "deny":
		return paymentgateway.NotificationStatusDeny
	case "expire":
		return paymentgateway.NotificationStatusExpire
	case "cancel":
		return paymentgateway.NotificationStatusCancel
	case "refund", "partial_refund":
		return paymentgateway.NotificationStatusRefund
	default:
		return paymentgateway.NotificationStatus(txStatus)
	}
}

// parseAmount converts the Midtrans gross_amount string (e.g. "150000.00")
// into an int64 representing the amount in the smallest unit.
func parseAmount(raw string) (int64, error) {
	if raw == "" {
		return 0, nil
	}
	// Strip decimal part if present (IDR has no sub-unit).
	if idx := strings.Index(raw, "."); idx != -1 {
		raw = raw[:idx]
	}
	var amount int64
	_, err := fmt.Sscanf(raw, "%d", &amount)
	if err != nil {
		return 0, fmt.Errorf("parse %q as int64: %w", raw, err)
	}
	return amount, nil
}
