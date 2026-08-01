package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockProductRepo struct {
	mockAnalyticsRepo
	summary    *domain.ProductMetricsSummary
	summaryErr error
	captured   repository.ProductMetricsParams
}

func (m *mockProductRepo) GetProductMetricsSummary(_ context.Context, _ transaction.Executor, params repository.ProductMetricsParams) (*domain.ProductMetricsSummary, error) {
	m.captured = params
	if m.summaryErr != nil {
		return nil, m.summaryErr
	}
	return m.summary, nil
}

func TestGetProductMetricsUsecase_Execute_Success(t *testing.T) {
	exec := &mockExecutor{}
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	margin := 35.5
	expectedSummary := &domain.ProductMetricsSummary{
		TopByRevenue: []domain.TopProductStat{
			{
				ProductID:        uuid.New(),
				ProductName:      "Sunflower Bouquet",
				Revenue:          1500000,
				UnitsSold:        30,
				GrossMarginPct:   &margin,
				SalesVelocity7d:  10,
				SalesVelocity30d: 30,
			},
		},
		TopByVolume: []domain.TopProductStat{
			{
				ProductID:        uuid.New(),
				ProductName:      "Single Rose",
				Revenue:          500000,
				UnitsSold:        100,
				SalesVelocity7d:  25,
				SalesVelocity30d: 100,
			},
		},
		InvoiceVoidRate: 0.02,
	}

	repo := &mockProductRepo{
		summary: expectedSummary,
	}

	uc := NewGetProductMetricsUsecase(exec, repo)

	input := GetProductMetricsInput{
		From: from,
		To:   to,
		TopN: 10,
	}

	res, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.TopByRevenue) != 1 {
		t.Errorf("expected 1 product in TopByRevenue, got %d", len(res.TopByRevenue))
	}
	if len(res.TopByVolume) != 1 {
		t.Errorf("expected 1 product in TopByVolume, got %d", len(res.TopByVolume))
	}
	if res.InvoiceVoidRate != 0.02 {
		t.Errorf("expected InvoiceVoidRate=0.02, got %f", res.InvoiceVoidRate)
	}
	if repo.captured.TopN != 10 {
		t.Errorf("expected TopN=10, got %d", repo.captured.TopN)
	}
}

func TestGetProductMetricsUsecase_Execute_Error(t *testing.T) {
	exec := &mockExecutor{}
	dbErr := errors.New("product metrics db error")

	repo := &mockProductRepo{summaryErr: dbErr}
	uc := NewGetProductMetricsUsecase(exec, repo)
	_, err := uc.Execute(context.Background(), GetProductMetricsInput{})
	if !errors.Is(err, dbErr) {
		t.Errorf("expected dbErr, got %v", err)
	}
}
