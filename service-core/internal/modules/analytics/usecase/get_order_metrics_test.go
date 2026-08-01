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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockExecutor struct{}

func (m *mockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *mockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *mockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type mockAnalyticsRepo struct {
	summary       *domain.OrderSummary
	summaryErr    error
	timeSeries    []domain.OrderTimeSeries
	timeSeriesErr error
	topProducts   []domain.TopProduct
	topProductsErr error
	topShops      []domain.TopShop
	topShopsErr   error

	capturedOrderParams repository.OrderMetricsParams
}

func (m *mockAnalyticsRepo) GetOrderSummary(_ context.Context, _ transaction.Executor, params repository.OrderMetricsParams) (*domain.OrderSummary, error) {
	m.capturedOrderParams = params
	if m.summaryErr != nil {
		return nil, m.summaryErr
	}
	return m.summary, nil
}

func (m *mockAnalyticsRepo) GetOrderTimeSeries(_ context.Context, _ transaction.Executor, _ repository.OrderMetricsParams) ([]domain.OrderTimeSeries, error) {
	if m.timeSeriesErr != nil {
		return nil, m.timeSeriesErr
	}
	return m.timeSeries, nil
}

func (m *mockAnalyticsRepo) GetTopProducts(_ context.Context, _ transaction.Executor, _ repository.OrderMetricsParams) ([]domain.TopProduct, error) {
	if m.topProductsErr != nil {
		return nil, m.topProductsErr
	}
	return m.topProducts, nil
}

func (m *mockAnalyticsRepo) GetTopShops(_ context.Context, _ transaction.Executor, _ repository.OrderMetricsParams) ([]domain.TopShop, error) {
	if m.topShopsErr != nil {
		return nil, m.topShopsErr
	}
	return m.topShops, nil
}

func (m *mockAnalyticsRepo) GetPaymentSummary(_ context.Context, _ transaction.Executor, _ repository.PaymentMetricsParams) (*domain.PaymentSummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetPaymentMethodBreakdown(_ context.Context, _ transaction.Executor, _ repository.PaymentMetricsParams) ([]domain.PaymentMethodBreakdown, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetShipmentSummary(_ context.Context, _ transaction.Executor, _ repository.ShipmentMetricsParams) (*domain.ShipmentSummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetCourierBreakdown(_ context.Context, _ transaction.Executor, _ repository.ShipmentMetricsParams) ([]domain.CourierBreakdown, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetInventorySummary(_ context.Context, _ transaction.Executor, _ repository.InventoryMetricsParams) (*domain.InventorySummary, error) {
	return nil, nil
}

func (m *mockAnalyticsRepo) GetProductMetricsSummary(_ context.Context, _ transaction.Executor, _ repository.ProductMetricsParams) (*domain.ProductMetricsSummary, error) {
	return nil, nil
}

func TestGetOrderMetricsUsecase_Execute_Success(t *testing.T) {
	exec := &mockExecutor{}
	shopID := uuid.New()
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	expectedSummary := &domain.OrderSummary{
		TotalOrders: 10,
		TotalGMV:    500000,
	}
	expectedTS := []domain.OrderTimeSeries{
		{Date: from, OrderCount: 5, GMV: 250000},
	}
	expectedProducts := []domain.TopProduct{
		{ProductID: uuid.New(), ProductName: "Rose", Revenue: 200000},
	}
	expectedShops := []domain.TopShop{
		{ShopID: shopID, ShopName: "Florist A", Revenue: 500000},
	}

	repo := &mockAnalyticsRepo{
		summary:     expectedSummary,
		timeSeries:  expectedTS,
		topProducts: expectedProducts,
		topShops:    expectedShops,
	}

	uc := NewGetOrderMetricsUsecase(exec, repo)

	input := GetOrderMetricsInput{
		From:        from,
		To:          to,
		Granularity: "daily",
		ShopID:      &shopID,
		TopN:        5,
	}

	res, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Summary.TotalOrders != 10 {
		t.Errorf("expected TotalOrders=10, got %d", res.Summary.TotalOrders)
	}
	if len(res.TimeSeries) != 1 {
		t.Errorf("expected 1 time series item, got %d", len(res.TimeSeries))
	}
	if len(res.TopProducts) != 1 {
		t.Errorf("expected 1 top product, got %d", len(res.TopProducts))
	}
	if len(res.TopShops) != 1 {
		t.Errorf("expected 1 top shop, got %d", len(res.TopShops))
	}

	if repo.capturedOrderParams.Granularity != repository.GranularityDaily {
		t.Errorf("expected granularity daily, got %s", repo.capturedOrderParams.Granularity)
	}
}

func TestGetOrderMetricsUsecase_Execute_Errors(t *testing.T) {
	exec := &mockExecutor{}
	dbErr := errors.New("db error")

	t.Run("Summary error", func(t *testing.T) {
		repo := &mockAnalyticsRepo{summaryErr: dbErr}
		uc := NewGetOrderMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetOrderMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})

	t.Run("TimeSeries error", func(t *testing.T) {
		repo := &mockAnalyticsRepo{
			summary:       &domain.OrderSummary{},
			timeSeriesErr: dbErr,
		}
		uc := NewGetOrderMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetOrderMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})

	t.Run("TopProducts error", func(t *testing.T) {
		repo := &mockAnalyticsRepo{
			summary:        &domain.OrderSummary{},
			timeSeries:     []domain.OrderTimeSeries{},
			topProductsErr: dbErr,
		}
		uc := NewGetOrderMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetOrderMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})

	t.Run("TopShops error", func(t *testing.T) {
		repo := &mockAnalyticsRepo{
			summary:     &domain.OrderSummary{},
			timeSeries:  []domain.OrderTimeSeries{},
			topProducts: []domain.TopProduct{},
			topShopsErr: dbErr,
		}
		uc := NewGetOrderMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetOrderMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})
}
