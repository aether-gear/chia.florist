package intelligencelayer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	applogger "service-core/internal/common/logger"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     applogger.Logger
}

func NewClient(baseURL string, timeout time.Duration, logger applogger.Logger) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if timeout <= 0 {
		timeout = 500 * time.Millisecond
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (c *Client) HealthCheck(ctx context.Context) (*HealthData, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create health request: %w", err)
	}

	var resp APIResponse[HealthData]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("intelligence-layer health check returned failure: %v", resp.Error)
	}

	return resp.Data, nil
}

func (c *Client) PredictDemand(ctx context.Context, payload DemandForecastRequest) (*DemandForecastResponse, error) {
	url := fmt.Sprintf("%s/predict/demand", c.baseURL)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal demand payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create demand request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp APIResponse[DemandForecastResponse]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("intelligence-layer demand prediction failed: %v", resp.Error)
	}

	return resp.Data, nil
}

func (c *Client) PredictStockoutRisk(ctx context.Context, payload StockoutRiskRequest) (*StockoutRiskResponse, error) {
	url := fmt.Sprintf("%s/predict/stockout-risk", c.baseURL)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stockout payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create stockout request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp APIResponse[StockoutRiskResponse]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("intelligence-layer stockout prediction failed: %v", resp.Error)
	}

	return resp.Data, nil
}

func (c *Client) PredictCourierSLA(ctx context.Context, payload CourierSLARequest) (*CourierSLAResponse, error) {
	url := fmt.Sprintf("%s/predict/courier-sla", c.baseURL)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal courier SLA payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create courier SLA request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp APIResponse[CourierSLAResponse]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("intelligence-layer courier SLA prediction failed: %v", resp.Error)
	}

	return resp.Data, nil
}

func (c *Client) DetectAnomaly(ctx context.Context, payload AnomalyCheckRequest) (*AnomalyCheckResponse, error) {
	url := fmt.Sprintf("%s/predict/anomaly", c.baseURL)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal anomaly payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create anomaly check request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var resp APIResponse[AnomalyCheckResponse]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("intelligence-layer anomaly detection failed: %v", resp.Error)
	}

	return resp.Data, nil
}

func (c *Client) doRequest(req *http.Request, target any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("intelligence-layer HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("intelligence-layer returned HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("failed to decode intelligence-layer JSON response: %w (raw: %s)", err, string(bodyBytes))
	}

	return nil
}
