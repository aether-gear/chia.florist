package usecase

import (
	"context"
	"errors"
	"testing"

	intelligencelayer "service-core/internal/infra/intelligence_layer"
	"service-core/internal/modules/analytics/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockDemandForecastRepo struct {
	mockAnalyticsRepo
	feat    *domain.ProductLagFeatures
	featErr error
}

func (m *mockDemandForecastRepo) GetProductLagFeatures(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*domain.ProductLagFeatures, error) {
	if m.featErr != nil {
		return nil, m.featErr
	}
	return m.feat, nil
}

type mockAIDemandProvider struct {
	demandResp *intelligencelayer.DemandForecastResponse
	err        error
}

func (m *mockAIDemandProvider) HealthCheck(ctx context.Context) (*intelligencelayer.HealthData, error) {
	return nil, nil
}
func (m *mockAIDemandProvider) PredictDemand(ctx context.Context, req intelligencelayer.DemandForecastRequest) (*intelligencelayer.DemandForecastResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.demandResp, nil
}
func (m *mockAIDemandProvider) PredictStockoutRisk(ctx context.Context, req intelligencelayer.StockoutRiskRequest) (*intelligencelayer.StockoutRiskResponse, error) {
	return nil, nil
}
func (m *mockAIDemandProvider) PredictCourierSLA(ctx context.Context, req intelligencelayer.CourierSLARequest) (*intelligencelayer.CourierSLAResponse, error) {
	return nil, nil
}
func (m *mockAIDemandProvider) DetectAnomaly(ctx context.Context, req intelligencelayer.AnomalyCheckRequest) (*intelligencelayer.AnomalyCheckResponse, error) {
	return nil, nil
}

func TestGetDemandForecast_Success(t *testing.T) {
	prodID := uuid.New()
	repo := &mockDemandForecastRepo{
		feat: &domain.ProductLagFeatures{
			ProductID:      prodID,
			ProductName:    "Rose Red",
			GrossMarginPct: 0.45,
			ViewCount:      850,
			UnitsSoldLag1:  12.0,
			UnitsSoldLag7:  15.0,
			Rolling7dMean:  13.5,
			CurrentStock:   100,
		},
	}

	aiProv := &mockAIDemandProvider{
		demandResp: &intelligencelayer.DemandForecastResponse{
			ProductID:            prodID.String(),
			PredictedUnitsSold7d: 98.5,
			ConfidenceTier:       "high",
		},
	}

	uc := NewGetDemandForecastUsecase(&mockExecutor{}, repo, aiProv)
	res, err := uc.Execute(context.Background(), GetDemandForecastInput{ProductID: prodID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ProductID != prodID {
		t.Errorf("expected product ID %v, got %v", prodID, res.ProductID)
	}
	if res.PredictedUnitsSold7d != 98.5 {
		t.Errorf("expected predicted units 98.5, got %f", res.PredictedUnitsSold7d)
	}
	if res.ConfidenceTier != "high" {
		t.Errorf("expected confidence high, got %s", res.ConfidenceTier)
	}
}

func TestGetDemandForecast_AIFailure_FallbackMovingAverage(t *testing.T) {
	prodID := uuid.New()
	repo := &mockDemandForecastRepo{
		feat: &domain.ProductLagFeatures{
			ProductID:      prodID,
			ProductName:    "Rose Red",
			Rolling7dMean:  10.0,
			CurrentStock:   50,
		},
	}

	aiProv := &mockAIDemandProvider{
		err: errors.New("ai service down"),
	}

	uc := NewGetDemandForecastUsecase(&mockExecutor{}, repo, aiProv)
	res, err := uc.Execute(context.Background(), GetDemandForecastInput{ProductID: prodID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFallback := 70.0 // 10.0 * 7
	if res.PredictedUnitsSold7d != expectedFallback {
		t.Errorf("expected fallback units %f, got %f", expectedFallback, res.PredictedUnitsSold7d)
	}
}
