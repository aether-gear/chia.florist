package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"
)

type mockShipmentRepo struct {
	mockAnalyticsRepo
	shipmentSummary *domain.ShipmentSummary
	summaryErr      error
	couriers        []domain.CourierBreakdown
	couriersErr     error
}

func (m *mockShipmentRepo) GetShipmentSummary(_ context.Context, _ transaction.Executor, _ repository.ShipmentMetricsParams) (*domain.ShipmentSummary, error) {
	if m.summaryErr != nil {
		return nil, m.summaryErr
	}
	return m.shipmentSummary, nil
}

func (m *mockShipmentRepo) GetCourierBreakdown(_ context.Context, _ transaction.Executor, _ repository.ShipmentMetricsParams) ([]domain.CourierBreakdown, error) {
	if m.couriersErr != nil {
		return nil, m.couriersErr
	}
	return m.couriers, nil
}

func TestGetShipmentMetricsUsecase_Execute_Success(t *testing.T) {
	exec := &mockExecutor{}
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	expectedSummary := &domain.ShipmentSummary{
		Total:        20,
		Delivered:    18,
		DeliveryRate: 0.9,
	}
	expectedCouriers := []domain.CourierBreakdown{
		{
			Courier:      "JNE",
			Service:      "REG",
			Count:        15,
			DeliveryRate: 0.93,
			AvgCost:      15000,
		},
	}

	repo := &mockShipmentRepo{
		shipmentSummary: expectedSummary,
		couriers:        expectedCouriers,
	}

	uc := NewGetShipmentMetricsUsecase(exec, repo)

	input := GetShipmentMetricsInput{
		From: from,
		To:   to,
		TopN: 5,
	}

	res, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Summary.Total != 20 {
		t.Errorf("expected Total=20, got %d", res.Summary.Total)
	}
	if len(res.Couriers) != 1 {
		t.Errorf("expected 1 courier breakdown item, got %d", len(res.Couriers))
	}
}

func TestGetShipmentMetricsUsecase_Execute_Errors(t *testing.T) {
	exec := &mockExecutor{}
	dbErr := errors.New("shipment db error")

	t.Run("ShipmentSummary error", func(t *testing.T) {
		repo := &mockShipmentRepo{summaryErr: dbErr}
		uc := NewGetShipmentMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetShipmentMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})

	t.Run("CourierBreakdown error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			shipmentSummary: &domain.ShipmentSummary{},
			couriersErr:     dbErr,
		}
		uc := NewGetShipmentMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetShipmentMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})
}
