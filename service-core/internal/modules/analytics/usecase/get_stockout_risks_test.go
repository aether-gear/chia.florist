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

type mockStockoutRiskRepo struct {
	mockAnalyticsRepo
	burnRates []domain.InventoryBurnRateData
	err       error
}

func (m *mockStockoutRiskRepo) GetInventoryBurnRates(_ context.Context, _ transaction.Executor, _ *uuid.UUID) ([]domain.InventoryBurnRateData, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.burnRates, nil
}

type mockAIStockoutProvider struct {
	stockoutResp *intelligencelayer.StockoutRiskResponse
	err          error
}

func (m *mockAIStockoutProvider) HealthCheck(ctx context.Context) (*intelligencelayer.HealthData, error) {
	return nil, nil
}
func (m *mockAIStockoutProvider) PredictDemand(ctx context.Context, req intelligencelayer.DemandForecastRequest) (*intelligencelayer.DemandForecastResponse, error) {
	return nil, nil
}
func (m *mockAIStockoutProvider) PredictStockoutRisk(ctx context.Context, req intelligencelayer.StockoutRiskRequest) (*intelligencelayer.StockoutRiskResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stockoutResp, nil
}
func (m *mockAIStockoutProvider) PredictCourierSLA(ctx context.Context, req intelligencelayer.CourierSLARequest) (*intelligencelayer.CourierSLAResponse, error) {
	return nil, nil
}
func (m *mockAIStockoutProvider) DetectAnomaly(ctx context.Context, req intelligencelayer.AnomalyCheckRequest) (*intelligencelayer.AnomalyCheckResponse, error) {
	return nil, nil
}

func TestGetStockoutRisks_Success(t *testing.T) {
	prodID := uuid.New()
	shopID := uuid.New()

	repo := &mockStockoutRiskRepo{
		burnRates: []domain.InventoryBurnRateData{
			{
				ProductID:            prodID,
				ProductName:          "Lily White",
				ShopID:               shopID,
				ShopName:             "Main Branch",
				Stock:                5,
				ReservedStock:        2,
				StockBurnRate7d:      3.5,
				SupplierLeadTimeDays: 7.0,
			},
		},
	}

	aiProv := &mockAIStockoutProvider{
		stockoutResp: &intelligencelayer.StockoutRiskResponse{
			ProductID:           prodID.String(),
			StockoutProbability: 0.9421,
			WillStockout:        true,
			RiskLevel:           "CRITICAL",
		},
	}

	uc := NewGetStockoutRisksUsecase(&mockExecutor{}, repo, aiProv)
	res, err := uc.Execute(context.Background(), GetStockoutRisksInput{ShopID: &shopID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 stockout risk item, got %d", len(res))
	}

	if res[0].RiskLevel != "CRITICAL" {
		t.Errorf("expected risk level CRITICAL, got %s", res[0].RiskLevel)
	}
	if !res[0].WillStockout {
		t.Errorf("expected will_stockout=true")
	}
}

func TestGetStockoutRisks_AIFailure_HeuristicFallback(t *testing.T) {
	prodID := uuid.New()
	shopID := uuid.New()

	repo := &mockStockoutRiskRepo{
		burnRates: []domain.InventoryBurnRateData{
			{
				ProductID:            prodID,
				ProductName:          "Lily White",
				ShopID:               shopID,
				ShopName:             "Main Branch",
				Stock:                1,
				ReservedStock:        0,
				StockBurnRate7d:      2.0,
				SupplierLeadTimeDays: 7.0,
			},
		},
	}

	aiProv := &mockAIStockoutProvider{
		err: errors.New("ai server offline"),
	}

	uc := NewGetStockoutRisksUsecase(&mockExecutor{}, repo, aiProv)
	res, err := uc.Execute(context.Background(), GetStockoutRisksInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 stockout risk item, got %d", len(res))
	}

	if res[0].RiskLevel != "CRITICAL" {
		t.Errorf("expected heuristic CRITICAL for 1 stock remaining, got %s", res[0].RiskLevel)
	}
}
