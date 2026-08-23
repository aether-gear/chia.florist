package usecase

import (
	"context"
	"errors"
	"testing"

	intelligencelayer "service-core/internal/infra/intelligence_layer"
	shipping "service-core/internal/infra/shipping"
)

type mockShippingProvider struct {
	rates []shipping.RateOption
	err   error
}

func (m *mockShippingProvider) CalculateRates(ctx context.Context, input shipping.CalculateRatesInput) ([]shipping.RateOption, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rates, nil
}

type mockAIProvider struct {
	slaResp *intelligencelayer.CourierSLAResponse
	err     error
}

func (m *mockAIProvider) HealthCheck(ctx context.Context) (*intelligencelayer.HealthData, error) {
	return nil, nil
}

func (m *mockAIProvider) PredictDemand(ctx context.Context, req intelligencelayer.DemandForecastRequest) (*intelligencelayer.DemandForecastResponse, error) {
	return nil, nil
}

func (m *mockAIProvider) PredictStockoutRisk(ctx context.Context, req intelligencelayer.StockoutRiskRequest) (*intelligencelayer.StockoutRiskResponse, error) {
	return nil, nil
}

func (m *mockAIProvider) PredictCourierSLA(ctx context.Context, req intelligencelayer.CourierSLARequest) (*intelligencelayer.CourierSLAResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.slaResp, nil
}

func (m *mockAIProvider) DetectAnomaly(ctx context.Context, req intelligencelayer.AnomalyCheckRequest) (*intelligencelayer.AnomalyCheckResponse, error) {
	return nil, nil
}

func TestEstimateShippingCost_WithAIPrediction(t *testing.T) {
	shipProv := &mockShippingProvider{
		rates: []shipping.RateOption{
			{
				Name:    "JNE",
				Code:    "jne",
				Service: "REG",
				Cost:    22000,
				Etd:     "2-3 hari",
			},
		},
	}

	aiProv := &mockAIProvider{
		slaResp: &intelligencelayer.CourierSLAResponse{
			CourierCode:            "jne",
			EstimatedDurationHours: 26.4,
			SLAConfidenceScore:     0.95,
			DeliveryStatus:         "ON_TRACK",
		},
	}

	u := NewEstimateShippingOptionsUsecase(shipProv, &mockExecutor{}, aiProv)

	res, err := u.Execute(context.Background(), EstimateShippingOptionsInput{
		Origin:      1,
		Destination: 2,
		Weight:      1000,
		Couriers:    []string{"jne"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 rate option, got %d", len(res))
	}

	if res[0].EstimatedDurationHours == nil || *res[0].EstimatedDurationHours != 26.4 {
		t.Errorf("expected duration 26.4, got %v", res[0].EstimatedDurationHours)
	}

	if res[0].DeliveryStatus == nil || *res[0].DeliveryStatus != "ON_TRACK" {
		t.Errorf("expected delivery status ON_TRACK, got %v", res[0].DeliveryStatus)
	}
}

func TestEstimateShippingCost_AIFailure_FailsOpen(t *testing.T) {
	shipProv := &mockShippingProvider{
		rates: []shipping.RateOption{
			{
				Name:    "JNE",
				Code:    "jne",
				Service: "REG",
				Cost:    22000,
				Etd:     "2-3 hari",
			},
		},
	}

	aiProv := &mockAIProvider{
		err: errors.New("connection refused to intelligence-layer"),
	}

	u := NewEstimateShippingOptionsUsecase(shipProv, &mockExecutor{}, aiProv)

	res, err := u.Execute(context.Background(), EstimateShippingOptionsInput{
		Origin:      1,
		Destination: 2,
		Weight:      1000,
		Couriers:    []string{"jne"},
	})

	if err != nil {
		t.Fatalf("expected success with fail-open, got error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 rate option, got %d", len(res))
	}

	if res[0].EstimatedDurationHours != nil {
		t.Errorf("expected nil duration when AI fails open, got %v", res[0].EstimatedDurationHours)
	}
}
