package usecase

import (
	"context"
	"errors"
	"testing"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockInventoryRepo struct {
	mockAnalyticsRepo
	summary    *domain.InventorySummary
	summaryErr error
	captured   repository.InventoryMetricsParams
}

func (m *mockInventoryRepo) GetInventorySummary(_ context.Context, _ transaction.Executor, params repository.InventoryMetricsParams) (*domain.InventorySummary, error) {
	m.captured = params
	if m.summaryErr != nil {
		return nil, m.summaryErr
	}
	return m.summary, nil
}

func TestGetInventoryMetricsUsecase_Execute_Success(t *testing.T) {
	exec := &mockExecutor{}
	shopID := uuid.New()

	expectedSummary := &domain.InventorySummary{
		TotalProducts:  50,
		TotalStock:     1000,
		TotalReserved:  100,
		TotalAvailable: 900,
		StockoutCount:  2,
		LowStockCount:  5,
	}

	repo := &mockInventoryRepo{
		summary: expectedSummary,
	}

	uc := NewGetInventoryMetricsUsecase(exec, repo)

	input := GetInventoryMetricsInput{
		ShopID:            &shopID,
		LowStockThreshold: 10,
	}

	res, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.TotalProducts != 50 {
		t.Errorf("expected TotalProducts=50, got %d", res.TotalProducts)
	}
	if repo.captured.LowStockThreshold != 10 {
		t.Errorf("expected LowStockThreshold=10, got %d", repo.captured.LowStockThreshold)
	}
	if repo.captured.ShopID == nil || *repo.captured.ShopID != shopID {
		t.Errorf("expected ShopID=%v, got %v", shopID, repo.captured.ShopID)
	}
}

func TestGetInventoryMetricsUsecase_Execute_Error(t *testing.T) {
	exec := &mockExecutor{}
	dbErr := errors.New("inventory db error")

	repo := &mockInventoryRepo{summaryErr: dbErr}
	uc := NewGetInventoryMetricsUsecase(exec, repo)
	_, err := uc.Execute(context.Background(), GetInventoryMetricsInput{})
	if !errors.Is(err, dbErr) {
		t.Errorf("expected dbErr, got %v", err)
	}
}
