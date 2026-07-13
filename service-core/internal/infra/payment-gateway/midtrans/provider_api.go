package midtrans

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	paymentgateway "service-core/internal/infra/payment-gateway"
	config "service-core/internal/shared/config"
)

type midtransAPIProvider struct {
	cfg        config.MidTransConfig
	client     *http.Client
	baseURL    string
	authToken  string
	strategies map[string]channelStrategy
}

func NewMidtransAPIProvider(
	cfg config.MidTransConfig,
) (*midtransAPIProvider, error) {
	if strings.TrimSpace(cfg.ServerKey) == "" {
		return nil, fmt.Errorf("midtrans: server key is required")
	}

	// Allow an explicit URL override
	// (useful for tests / local proxies).
	var baseURL string
	if strings.TrimSpace(cfg.URL) != "" {
		cleanedURL := strings.TrimRight(cfg.URL, "/")
		cleanedURL = strings.TrimSuffix(cleanedURL, "/v2/charge")
		cleanedURL = strings.TrimSuffix(cleanedURL, "/v2")
		baseURL = cleanedURL
	}

	authToken := "Basic " +
		base64.StdEncoding.EncodeToString(
			[]byte(cfg.ServerKey+":"),
		)

	p := &midtransAPIProvider{
		cfg:       cfg,
		client:    &http.Client{Timeout: 30 * time.Second},
		baseURL:   baseURL,
		authToken: authToken,
	}

	p.strategies = map[string]channelStrategy{
		"gopay":     &gopayStrategy{},
		"shopeepay": &shopeepayStrategy{},
		"qris":      &qrisStrategy{},
		"qr_code":   &qrisStrategy{},
		"mandiri":   &mandiriStrategy{},
	}

	return p, nil
}

func (p *midtransAPIProvider) getStrategy(code string) channelStrategy {
	code = strings.ToLower(code)
	if p.strategies == nil {
		p.strategies = map[string]channelStrategy{
			"gopay":     &gopayStrategy{},
			"shopeepay": &shopeepayStrategy{},
			"qris":      &qrisStrategy{},
			"qr_code":   &qrisStrategy{},
			"mandiri":   &mandiriStrategy{},
		}
	}
	return p.strategies[code]
}

func (p *midtransAPIProvider) Supports(code string) bool {
	return p.getStrategy(code) != nil
}

type vaNumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

type action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

type chargeAPIResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	StaffID           string `json:"staff_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	ExpiryTime        string `json:"expiry_time"`
}

type bankTransferChargeAPIResponse struct {
	chargeAPIResponse
	PermataVANumber string     `json:"permata_va_number"`
	VaNumbers       []vaNumber `json:"va_numbers"`
	BillKey         string     `json:"bill_key"`
	BillerCode      string     `json:"biller_code"`
}

type gopayChargeAPIResponse struct {
	chargeAPIResponse

	Actions []action `json:"actions"`
}

type qrisChargeAPIResponse struct {
	chargeAPIResponse

	Acquirer string   `json:"acquirer"`
	Actions  []action `json:"actions"`
	QRString string   `json:"qr_string"`
}

type shopeepayChargeAPIResponse struct {
	chargeAPIResponse

	Actions                []action `json:"actions"`
	ChannelResponseCode    string   `json:"channel_response_code"`
	ChannelResponseMessage string   `json:"channel_response_message"`
}

type statusAPIResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

type cancelAPIResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

func (p *midtransAPIProvider) Charge(
	ctx context.Context,
	req paymentgateway.ChargeRequest,
) (*paymentgateway.ChargeResponse, error) {
	strategy := p.getStrategy(req.PaymentType)
	if strategy == nil {
		return nil, fmt.Errorf("midtrans: unsupported payment method %q", req.PaymentType)
	}

	body, err := p.buildChargeBody(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans: build charge request: %w", err)
	}

	url := p.baseURL + "/v2/charge"
	respBody, err := p.doRequest(
		ctx,
		http.MethodPost,
		url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("midtrans: charge transaction: %w", err)
	}

	var baseResp chargeAPIResponse
	if err := json.Unmarshal(respBody, &baseResp); err != nil {
		return nil, fmt.Errorf("midtrans: unmarshal base charge response: %w", err)
	}

	if !strings.HasPrefix(baseResp.StatusCode, "2") {
		return nil, fmt.Errorf("midtrans: charge failed: %s (status %s)",
			baseResp.StatusMessage, baseResp.StatusCode)
	}

	instructions, err := strategy.ParseInstructions(respBody)
	if err != nil {
		return nil, fmt.Errorf("midtrans: parse instructions: %w", err)
	}

	grossAmount, err := parseAmount(baseResp.GrossAmount)
	if err != nil {
		return nil, fmt.Errorf("parse gross_amount: %w", err)
	}

	var expiresAt time.Time
	if baseResp.ExpiryTime != "" {
		wib, _ := time.LoadLocation("Asia/Jakarta")
		t, parseErr := time.ParseInLocation("2006-01-02 15:04:05", baseResp.ExpiryTime, wib)
		if parseErr == nil {
			expiresAt = t
		}
	}

	result := paymentgateway.ChargeResponse{
		GatewayTransactionID: baseResp.TransactionID,
		GatewayOrderID:       baseResp.OrderID,
		PaymentType:          baseResp.PaymentType,
		GrossAmount:          grossAmount,
		Status:               baseResp.TransactionStatus,
		ExpiresAt:            expiresAt,
		Instructions:         instructions,
	}

	return &result, nil
}

func (p *midtransAPIProvider) ParseNotification(
	ctx context.Context,
	payload paymentgateway.NotificationPayload,
) (*paymentgateway.NotificationResult, error) {
	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		return nil, fmt.Errorf("midtrans: parse notification: missing order_id in payload")
	}

	return p.fetchTransactionStatus(ctx, orderID)
}

func (p *midtransAPIProvider) GetTransactionStatus(
	ctx context.Context,
	gatewayOrderID string,
) (*paymentgateway.NotificationResult, error) {
	if gatewayOrderID == "" {
		return nil, fmt.Errorf("midtrans: GetTransactionStatus: gatewayOrderID is required")
	}
	return p.fetchTransactionStatus(ctx, gatewayOrderID)
}

// fetchTransactionStatus calls the Midtrans
// GET /v2/{orderID}/status endpoint
// and returns a normalised NotificationResult.
//
// It is the shared implementation for both ParseNotification
// and GetTransactionStatus.
func (p *midtransAPIProvider) fetchTransactionStatus(
	ctx context.Context,
	orderID string,
) (*paymentgateway.NotificationResult, error) {
	url := p.baseURL + "/v2/" + orderID + "/status"
	respBody, err := p.doRequest(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("midtrans: check transaction %q: %w",
			orderID, err)
	}

	var result statusAPIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("midtrans: unmarshal status response: %w", err)
	}

	if !strings.HasPrefix(result.StatusCode, "2") {
		return nil, fmt.Errorf("midtrans: check transaction %q: %s (status %s)",
			orderID, result.StatusMessage, result.StatusCode)
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

func (p *midtransAPIProvider) CancelTransaction(
	ctx context.Context,
	gatewayOrderID string,
) error {
	url := p.baseURL + "/v2/" + gatewayOrderID + "/cancel"
	respBody, err := p.doRequest(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("midtrans: cancel transaction %q: %w",
			gatewayOrderID, err)
	}

	var resp cancelAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("midtrans: unmarshal cancel response: %w", err)
	}

	if !strings.HasPrefix(resp.StatusCode, "2") {
		return fmt.Errorf("midtrans: cancel transaction %q: %s (status %s)",
			gatewayOrderID, resp.StatusMessage, resp.StatusCode,
		)
	}

	return nil
}

// doRequest executes an authenticated HTTP request
// against the Midtrans Core API and returns the
// raw response body.
func (p *midtransAPIProvider) doRequest(
	ctx context.Context,
	method, url string,
	body []byte,
) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx, method, url, bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", p.authToken)

	httpResp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return respBody, nil
}

// buildChargeBody converts the app-generic
// ChargeRequest into a JSON byte slice matching
// the Midtrans Core API v2 /v2/charge request body.
func (p *midtransAPIProvider) buildChargeBody(
	req paymentgateway.ChargeRequest,
) ([]byte, error) {
	body := map[string]any{
		"transaction_details": map[string]any{
			"order_id":     req.OrderID.String(),
			"gross_amount": req.Amount,
		},
		"customer_details": map[string]any{
			"first_name": req.CustomerName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
	}

	if len(req.Items) > 0 {
		itemDetails := make([]map[string]any, 0, len(req.Items))
		for _, item := range req.Items {
			itemDetails = append(itemDetails, map[string]any{
				"id":       item.ID,
				"price":    item.Price,
				"quantity": item.Quantity,
				"name":     item.Name,
			})
		}
		body["item_details"] = itemDetails
	}

	if !req.ExpiresAt.IsZero() {
		dur := time.Until(req.ExpiresAt)
		if dur > 0 {
			minutes := int(dur.Minutes())
			if minutes < 1 {
				minutes = 1
			}
			body["custom_expiry"] = map[string]any{
				"expiry_duration": minutes,
				"unit":            "minute",
			}
		}
	}

	strategy := p.getStrategy(req.PaymentType)
	if strategy == nil {
		return nil, fmt.Errorf("unsupported payment type %q", req.PaymentType)
	}

	strategy.BuildPayload(req, body)

	return json.Marshal(body)
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
