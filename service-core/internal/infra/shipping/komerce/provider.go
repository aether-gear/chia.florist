package komerce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	shipping "service-core/internal/infra/shipping"
	config "service-core/internal/shared/config"
)

type komerceProvider struct {
	client       *http.Client
	orderBaseURL string
	trackBaseURL string
	apiKey       string
	shippingKey  string
}

func NewKomerceProvider(
	cfg config.KomerceConfig,
) (shipping.LogisticsProvider, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("komerce: api key is required")
	}
	if strings.TrimSpace(cfg.OrderBaseURL) == "" {
		return nil, fmt.Errorf("komerce: order base URL is required")
	}
	if strings.TrimSpace(cfg.TrackBaseURL) == "" {
		return nil, fmt.Errorf("komerce: track base URL is required")
	}

	return &komerceProvider{
		client:       &http.Client{Timeout: 30 * time.Second},
		orderBaseURL: strings.TrimRight(cfg.OrderBaseURL, "/"),
		trackBaseURL: strings.TrimRight(cfg.TrackBaseURL, "/"),
		apiKey:       cfg.APIKey,
		shippingKey:  cfg.ShippingKey,
	}, nil
}

type createOrderRequestBody struct {
	OriginAreaID      int    `json:"origin_area_id"`
	DestinationAreaID int    `json:"destination_area_id"`
	Courier           string `json:"courier"`
	CourierService    string `json:"courier_service_code"`
	Weight            int    `json:"weight"`
	UniqueOrderID     string `json:"unique_order_id"`
	ItemName          string `json:"item_name"`
	ItemPrice         int64  `json:"item_price"`
	Quantity          int    `json:"quantity"`
	IsCOD             int    `json:"is_cod"`
	Note              string `json:"note"`
	ShipperName       string `json:"shipper_name"`
	ShipperPhone      string `json:"shipper_phone"`
	ShipperAddress    string `json:"shipper_address"`
	ReceiverName      string `json:"receiver_name"`
	ReceiverPhone     string `json:"receiver_phone"`
	ReceiverAddress   string `json:"receiver_address"`
}

type createOrderResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data struct {
		OrderNo    string `json:"order_no"`
		AirwayBill string `json:"airway_bill"`
	} `json:"data"`
}

func (p *komerceProvider) CreateOrder(
	ctx context.Context,
	input shipping.CreateOrderInput,
) (*shipping.CreateOrderResult, error) {
	endpoint := "/order/api/v1/orders/store"
	reqBody := createOrderRequestBody{
		OriginAreaID:      input.OriginAreaID,
		DestinationAreaID: input.DestinationAreaID,
		Courier:           input.CourierCode,
		CourierService:    input.CourierService,
		Weight:            input.Weight,
		UniqueOrderID:     input.UniqueOrderID,
		ItemName:          input.ItemName,
		ItemPrice:         input.ItemPrice,
		Quantity:          input.ItemQty,
		IsCOD:             0,
		Note:              "",
		ShipperName:       input.ShipperName,
		ShipperPhone:      input.ShipperPhone,
		ShipperAddress:    input.ShipperAddress,
		ReceiverName:      input.ReceiverName,
		ReceiverPhone:     input.ReceiverPhone,
		ReceiverAddress:   input.ReceiverAddress,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := p.orderBaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, fmt.Errorf("one or more request parameters are invalid")

		case http.StatusUnauthorized:
			return nil, fmt.Errorf("invalid or missing API key")

		case http.StatusUnprocessableEntity:
			return nil, fmt.Errorf("missing required parameter")

		case http.StatusInternalServerError:
			return nil, fmt.Errorf("provider service unavailable")

		default:
			return nil, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result createOrderResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if result.Meta.Code != 200 {
		return nil, fmt.Errorf("rejected (%d): %s", result.Meta.Code, result.Meta.Message)
	}

	return &shipping.CreateOrderResult{
		KomerceOrderNo: result.Data.OrderNo,
		TrackingNumber: result.Data.AirwayBill,
	}, nil
}

type trackWaybillRequestBody struct {
	AWB       string  `json:"awb"`
	Courier   string  `json:"courier"`
	LastPhone *string `json:"last_phone_number,omitempty"`
}

type trackWaybillResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		Summary struct {
			Status string `json:"status"`
		} `json:"summary"`
		Manifest []struct {
			ManifestDate        string `json:"manifest_date"`
			ManifestDescription string `json:"manifest_description"`
			CityName            string `json:"city_name"`
			ManifestCode        string `json:"manifest_code"`
		} `json:"manifest"`
	} `json:"data"`
}

func (p *komerceProvider) TrackShipment(
	ctx context.Context,
	input shipping.TrackShipmentInput,
) ([]shipping.TrackingEvent, error) {
	endpoint := "/api/v1/track/waybill"
	reqBody := trackWaybillRequestBody{
		AWB:       input.TrackingNumber,
		Courier:   input.Courier,
		LastPhone: input.LastPhone,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := p.trackBaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Key", p.shippingKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result trackWaybillResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if result.Meta.Code != 200 {
		return nil, fmt.Errorf("rejected (%d): %s", result.Meta.Code, result.Meta.Message)
	}

	var events []shipping.TrackingEvent
	for _, m := range result.Data.Manifest {
		t, _ := time.Parse("2006-01-02 15:04:05", m.ManifestDate)
		events = append(events, shipping.TrackingEvent{
			Status:      m.ManifestCode,
			Description: m.ManifestDescription,
			Location:    m.CityName,
			Timestamp:   t,
		})
	}

	return events, nil
}
