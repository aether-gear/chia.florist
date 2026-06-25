package rajaongkir

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/infra/shipping"
	locationDomain "service-core/internal/modules/location/domain"
	config "service-core/internal/shared/config"
)

type rajaOngkirProvider struct {
	client      *http.Client
	baseURL     string
	shippingKey string
	paymentKey  string
	qrislyKey   string
}

func NewRajaOngkirProvider(
	cfg config.RajaOngkirConfig,
) (shipping.Provider, error) {
	if strings.TrimSpace(cfg.URL) == "" ||
		strings.TrimSpace(cfg.QRISLYKey) == "" ||
		strings.TrimSpace(cfg.ShippingKey) == "" ||
		strings.TrimSpace(cfg.PaymentKey) == "" {

		return nil, fmt.Errorf("midtrans: server key is required")
	}

	// Allow an explicit URL override
	// (useful for tests / local proxies).
	var baseURL string
	if strings.TrimSpace(cfg.URL) != "" {
		cleanedURL := strings.TrimRight(cfg.URL, "/")
		cleanedURL = strings.TrimSuffix(cleanedURL, "/api/v1/")
		cleanedURL = strings.TrimSuffix(cleanedURL, "/api")
		baseURL = cleanedURL
	}

	return &rajaOngkirProvider{
		client:      &http.Client{Timeout: 30 * time.Second},
		baseURL:     baseURL,
		shippingKey: cfg.ShippingKey,
		paymentKey:  cfg.PaymentKey,
		qrislyKey:   cfg.QRISLYKey,
	}, nil
}

type calculateRatesResponse struct {
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        int64  `json:"cost"`
		Etd         string `json:"etd"`
	} `json:"data"`
}

func (s *rajaOngkirProvider) CalculateRates(
	ctx context.Context,
	input shipping.CalculateRatesInput,
) ([]shipping.RateOption, error) {
	body, err := s.doCalculateRatesRequest(ctx, http.MethodPost, input)
	if err != nil {
		return nil, err
	}

	var resp calculateRatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode provider response: %w", err)
	}

	if resp.Meta.Code != 200 {
		if resp.Meta.Code == 422 {
			return nil, apperrors.NewInvalidInput("The selected courier option is invalid or no longer supported")
		}
		if resp.Meta.Code == 404 {
			return nil, apperrors.NewBadRequest("Shipping is currently unavailable for this destination using the selected courier")
		}
		if resp.Meta.Code == 400 ||
			resp.Meta.Message == "Missing Params" {
			return nil, apperrors.NewInternal(errors.New("Missing params"))
		}

		return nil, errors.New(resp.Meta.Message)
	}

	var costOptions []shipping.RateOption
	for _, cO := range resp.Data {
		costOption := shipping.RateOption{
			Name:        cO.Name,
			Code:        cO.Code,
			Service:     cO.Service,
			Description: cO.Description,
			Cost:        cO.Cost,
			Etd:         cO.Etd,
		}

		costOptions = append(costOptions, costOption)
	}

	return costOptions, nil
}

type destinationResponse struct {
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (s *rajaOngkirProvider) ListProvinces(
	ctx context.Context,
) ([]locationDomain.Province, error) {
	body, err := s.doLocationRequest(ctx, http.MethodGet, "/province", nil)
	if err != nil {
		return nil, fmt.Errorf("request for provinces: %w", err)
	}

	var resp destinationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode province response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("provinces rejected: %s", resp.Meta.Message)
	}

	var provinces []locationDomain.Province
	for _, p := range resp.Data {
		provinces = append(provinces, locationDomain.Province{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return provinces, nil
}

func (s *rajaOngkirProvider) ListCities(
	ctx context.Context,
	provinceID string,
) ([]locationDomain.City, error) {
	path := fmt.Sprintf("/city/%s", provinceID)
	body, err := s.doLocationRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("request for cities: %w", err)
	}

	var resp destinationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode city response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("cities rejected: %s", resp.Meta.Message)
	}

	var cities []locationDomain.City
	for _, p := range resp.Data {
		cities = append(cities, locationDomain.City{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return cities, nil
}

func (s *rajaOngkirProvider) ListDistricts(
	ctx context.Context,
	cityID string,
) ([]locationDomain.District, error) {
	path := fmt.Sprintf("/district/%s", cityID)
	body, err := s.doLocationRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("request for districts: %w", err)
	}

	var resp destinationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode district response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("districts rejected: %s", resp.Meta.Message)
	}

	var districts []locationDomain.District
	for _, p := range resp.Data {
		districts = append(districts, locationDomain.District{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return districts, nil
}

func (s *rajaOngkirProvider) ListVillages(
	ctx context.Context,
	districtID string,
) ([]locationDomain.Village, error) {
	path := fmt.Sprintf("/sub-district/%s", districtID)
	body, err := s.doLocationRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("request for villages: %w", err)
	}

	var resp destinationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode village response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("villages rejected: %s", resp.Meta.Message)
	}

	var villages []locationDomain.Village
	for _, p := range resp.Data {
		villages = append(villages, locationDomain.Village{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return villages, nil
}

func (s *rajaOngkirProvider) doCalculateRatesRequest(
	ctx context.Context,
	method string,
	input shipping.CalculateRatesInput,
) ([]byte, error) {
	form := url.Values{}
	form.Set("origin", strconv.Itoa(input.OriginID))
	form.Set("destination", strconv.Itoa(input.DestinationID))
	form.Set("weight", strconv.Itoa(input.Weight))
	form.Set("courier", strings.Join(input.Couriers, ":"))
	if input.PriceFilter != nil {
		form.Set("price", *input.PriceFilter)
	}

	endpoint := s.baseURL + "/api/v1/calculate/district/domestic-cost"
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("http error: %s", err)
	}

	req.Header.Set("Key", s.shippingKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

func (s *rajaOngkirProvider) doLocationRequest(
	ctx context.Context,
	method string,
	path string,
	query map[string]string,
) ([]byte, error) {
	url := s.baseURL + "/api/v1/destination" + path

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", s.shippingKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rajaongkir error: %s", string(body))
	}

	return io.ReadAll(resp.Body)
}
