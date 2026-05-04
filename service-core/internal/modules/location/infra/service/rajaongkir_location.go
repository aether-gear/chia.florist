package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"service-core/internal/modules/location/domain"
	"service-core/internal/modules/location/repository"
)

type rajaOngkirService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewRajaOngkirService(
	apiKey string,
	baseURL string,
	client *http.Client,
) repository.LocationRepository {
	return &rajaOngkirService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
}

type rajaOngkirResponse struct {
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

func (s *rajaOngkirService) ListProvinces() ([]domain.Province, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := s.doRequest(ctx, http.MethodGet, "/province", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load provinces: %w", err)
	}

	var resp rajaOngkirResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("mapping province response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("failed to load provinces: %s", resp.Meta.Message)
	}

	var provinces []domain.Province
	for _, p := range resp.Data {
		provinces = append(provinces, domain.Province{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return provinces, nil
}

func (s *rajaOngkirService) ListCitiesByProvince(provinceID string) ([]domain.City, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := fmt.Sprintf("/city/%s", provinceID)
	body, err := s.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %w", err)
	}

	var resp rajaOngkirResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("mapping city response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("failed to load cities: %s", resp.Meta.Message)
	}

	var cities []domain.City
	for _, p := range resp.Data {
		cities = append(cities, domain.City{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return cities, nil
}

func (s *rajaOngkirService) ListDistrictsByCity(cityID string) ([]domain.District, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := fmt.Sprintf("/district/%s", cityID)
	body, err := s.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load districts: %w", err)
	}

	var resp rajaOngkirResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("mapping district response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("failed to load districts: %s", resp.Meta.Message)
	}

	var districts []domain.District
	for _, p := range resp.Data {
		districts = append(districts, domain.District{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return districts, nil
}

func (s *rajaOngkirService) ListVillagesByDistrict(districtID string) ([]domain.Village, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	path := fmt.Sprintf("/sub-district/%s", districtID)
	body, err := s.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load villages: %w", err)
	}

	var resp rajaOngkirResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("mapping village response failed: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, fmt.Errorf("failed to load villages: %s", resp.Meta.Message)
	}

	var villages []domain.Village
	for _, p := range resp.Data {
		villages = append(villages, domain.Village{
			ID:   strconv.Itoa(p.ID),
			Name: p.Name,
		})
	}

	return villages, nil
}

func (s *rajaOngkirService) doRequest(
	ctx context.Context,
	method string,
	path string,
	query map[string]string,
) ([]byte, error) {
	url := s.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", s.apiKey)
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
