package komerce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	apperrors "service-core/internal/common/errors"
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
	if !isValidKomerceCourier(input.CourierCode) {
		return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")
	}

	endpoint := "/order/api/v1/orders/store"
	reqBody := createOrderRequestBody{
		OriginAreaID:      input.OriginAreaID,
		DestinationAreaID: input.DestinationAreaID,
		Courier:           normalizeCourierCode(input.CourierCode),
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
		case http.StatusBadRequest, http.StatusUnprocessableEntity:
			return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")

		case http.StatusUnauthorized:
			return nil, apperrors.NewUnauthorized("invalid or missing API key")

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
		return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")
	}

	if result.Meta.Code != 200 {
		return nil, mapKomerceError(result.Meta.Code, result.Meta.Message)
	}

	return &shipping.CreateOrderResult{
		KomerceOrderNo: result.Data.OrderNo,
		TrackingNumber: result.Data.AirwayBill,
	}, nil
}

func (p *komerceProvider) CancelOrder(
	ctx context.Context,
	komerceOrderNo string,
) error {
	if strings.TrimSpace(komerceOrderNo) == "" {
		return nil
	}

	endpoint := "/order/api/v1/orders/cancel"
	reqBody := map[string]string{
		"order_no": komerceOrderNo,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal cancel request: %w", err)
	}

	url := p.orderBaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build cancel request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("cancel request failed: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

type trackWaybillRequestBody struct {
	AWB             string  `json:"awb"`
	Courier         string  `json:"courier"`
	LastPhone       *string `json:"last_phone,omitempty"`
	LastPhoneNumber *string `json:"last_phone_number,omitempty"`
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
	if !isValidKomerceCourier(input.Courier) {
		return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")
	}

	endpoint := "/api/v1/track/waybill"
	reqURL, err := url.Parse(p.trackBaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := reqURL.Query()
	q.Set("courier", normalizeCourierCode(input.Courier))
	q.Set("awb", input.TrackingNumber)
	if input.LastPhone != nil && *input.LastPhone != "" {
		q.Set("last_phone_number", *input.LastPhone)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("key", p.shippingKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		switch resp.StatusCode {
		case http.StatusBadRequest, http.StatusUnprocessableEntity:
			return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")

		case http.StatusUnauthorized:
			return nil, apperrors.NewUnauthorized("invalid or missing API key")

		case http.StatusNotFound:
			return nil, apperrors.NewNotFound("tracking number not found or not yet scanned by courier")

		case http.StatusTooManyRequests:
			return nil, apperrors.NewTooManyRequests("Komerce API rate limit exceeded")

		default:
			return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result trackWaybillResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, apperrors.NewBadRequest("external service is having a problem with processing current request")
	}

	if result.Meta.Code != 200 {
		appErr := mapKomerceError(result.Meta.Code, result.Meta.Message)
		return nil, appErr
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

func isValidKomerceCourier(courier string) bool {
	c := normalizeCourierCode(courier)
	switch c {
	case "jne", "jnt", "ninja", "tiki", "pos", "anteraja", "sap", "lion", "wahana", "first", "ide":
		return true
	default:
		return false
	}
}

func normalizeCourierCode(courier string) string {
	c := strings.ToLower(strings.TrimSpace(courier))
	switch {
	case strings.Contains(c, "jne"):
		return "jne"
	case strings.Contains(c, "j&t") || strings.Contains(c, "jnt"):
		return "jnt"
	case strings.Contains(c, "ninja"):
		return "ninja"
	case strings.Contains(c, "tiki"):
		return "tiki"
	case strings.Contains(c, "pos"):
		return "pos"
	case strings.Contains(c, "anteraja"):
		return "anteraja"
	case strings.Contains(c, "sap"):
		return "sap"
	case strings.Contains(c, "lion"):
		return "lion"
	case strings.Contains(c, "wahana"):
		return "wahana"
	case strings.Contains(c, "first"):
		return "first"
	case strings.Contains(c, "idexpress") || strings.Contains(c, "ide"):
		return "ide"
	default:
		fields := strings.Fields(c)
		if len(fields) > 0 {
			return fields[0]
		}
		return c
	}
}
