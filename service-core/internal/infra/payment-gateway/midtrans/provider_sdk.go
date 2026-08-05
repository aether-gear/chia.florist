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

type midtransSDKProvider struct {
	cfg    config.MidTransConfig
	client coreapi.Client
}

func NewMidtransSDKProvider(
	cfg config.MidTransConfig,
) (*midtransSDKProvider, error) {
	if strings.TrimSpace(cfg.ServerKey) == "" {
		return nil, fmt.Errorf("midtrans: server key is required")
	}

	env := midtranssdk.Sandbox
	if cfg.IsProduction {
		env = midtranssdk.Production
	}

	c := coreapi.Client{}
	c.New(cfg.ServerKey, env)

	return &midtransSDKProvider{
		cfg:    cfg,
		client: c,
	}, nil
}

func (p *midtransSDKProvider) Charge(
	ctx context.Context,
	req paymentgateway.ChargeRequest,
) (*paymentgateway.ChargeResponse, error) {
	chargeReq, err := p.buildChargeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans: build charge request: %w", err)
	}

	resp, midErr := p.client.ChargeTransaction(chargeReq)
	if midErr != nil {
		return nil, fmt.Errorf("midtrans: charge transaction: %s (status %d)",
			midErr.Message, midErr.StatusCode)
	}

	result, err := p.mapChargeResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *midtransSDKProvider) ParseNotification(
	_ context.Context,
	payload paymentgateway.NotificationPayload,
) (*paymentgateway.NotificationResult, error) {
	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		return nil, fmt.Errorf("midtrans: parse notification: missing order_id in payload")
	}

	result, midErr :=
		p.client.CheckTransaction(orderID)
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
func (p *midtransSDKProvider) CancelTransaction(
	_ context.Context,
	gatewayOrderID string,
) error {
	_, midErr :=
		p.client.CancelTransaction(gatewayOrderID)
	if midErr != nil {
		return fmt.Errorf("midtrans: cancel transaction %q: %s (status %d)",
			gatewayOrderID, midErr.Message, midErr.StatusCode)
	}
	return nil
}

// RefundTransaction requests Midtrans
// to refund a paid transaction identified by its gateway-side order ID
func (p *midtransSDKProvider) RefundTransaction(
	_ context.Context,
	req paymentgateway.RefundRequest,
) (*paymentgateway.RefundResponse, error) {
	refundReq := &coreapi.RefundReq{
		RefundKey: fmt.Sprintf("refund-%s-%d", req.GatewayOrderID, time.Now().Unix()),
		Amount:    req.RefundAmount,
		Reason:    req.Reason,
	}
	resp, midErr := p.client.RefundTransaction(req.GatewayOrderID, refundReq)
	if midErr != nil {
		return nil, fmt.Errorf("midtrans: refund transaction %q: %s (status %d)",
			req.GatewayOrderID, midErr.Message, midErr.StatusCode)
	}

	grossAmount, _ := parseAmount(resp.GrossAmount)
	return &paymentgateway.RefundResponse{
		GatewayTransactionID: resp.TransactionID,
		GatewayOrderID:       resp.OrderID,
		RefundAmount:         grossAmount,
		Status:               resp.TransactionStatus,
	}, nil
}

// buildChargeRequest converts the app
// eneric ChargeRequest into a Midtrans
// coreapi.ChargeReq.
func (p *midtransSDKProvider) buildChargeRequest(
	req paymentgateway.ChargeRequest,
) (*coreapi.ChargeReq, error) {
	txDetails := midtranssdk.TransactionDetails{
		OrderID:  req.PaymentID.String(),
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
	case "mandiri":
		bankCode := "Mandiri"
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
func (p *midtransSDKProvider) mapChargeResponse(
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
// or deep-link actions from the raw Midtrans SDK response
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
