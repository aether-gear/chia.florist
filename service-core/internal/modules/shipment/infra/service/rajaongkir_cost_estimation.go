package service

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

	"service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"
)

type rajaOngkirCostEstimation struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewRajaOngkirCostEstimation(
	apiKey string,
	baseURL string,
	client *http.Client,
) repository.ShippingCostProvider {
	return &rajaOngkirCostEstimation{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
}

type rajaOngkirRequest struct {
	Origin      int    `json:"origin"`
	Destination int    `json:"destination"`
	Weight      int    `json:"weight"`
	Courier     string `json:"courier"`
	Price       string `json:"price"`
}

type rajaOngkirResponse struct {
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

func (s *rajaOngkirCostEstimation) CalculateCost(
	ctx context.Context,
	exec transaction.Executor,
	input repository.CalculateCostInput,
) ([]repository.CostOption, error) {
	body, err := s.doRequest(ctx, http.MethodPost, input)
	if err != nil {
		return nil, err
	}

	var resp rajaOngkirResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode provider response: %w", err)
	}

	if resp.Meta.Code != 200 {
		return nil, errors.New(resp.Meta.Message)
	}

	var costOptions []repository.CostOption
	for _, cO := range resp.Data {
		costOption := repository.CostOption{
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

func (s *rajaOngkirCostEstimation) doRequest(
	ctx context.Context,
	method string,
	input repository.CalculateCostInput,
) ([]byte, error) {
	endpoint := s.baseURL

	form := url.Values{}
	form.Set("origin", strconv.Itoa(input.OriginID))
	form.Set("destination", strconv.Itoa(input.DestinationID))
	form.Set("weight", strconv.Itoa(input.Weight))
	form.Set("courier", strings.Join(input.Couriers, ":"))
	if input.PriceFilter != nil {
		form.Set("price", *input.PriceFilter)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("http error: %s", err)
	}

	req.Header.Set("Key", s.apiKey)
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
