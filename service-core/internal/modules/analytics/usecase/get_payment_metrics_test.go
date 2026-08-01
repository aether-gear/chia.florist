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

type mockPaymentRepo struct {
	mockAnalyticsRepo
	paymentSummary   *domain.PaymentSummary
	paymentSumErr    error
	paymentBreakdown []domain.PaymentMethodBreakdown
	breakdownErr     error
}

func (m *mockPaymentRepo) GetPaymentSummary(_ context.Context, _ transaction.Executor, _ repository.PaymentMetricsParams) (*domain.PaymentSummary, error) {
	if m.paymentSumErr != nil {
		return nil, m.paymentSumErr
	}
	return m.paymentSummary, nil
}

func (m *mockPaymentRepo) GetPaymentMethodBreakdown(_ context.Context, _ transaction.Executor, _ repository.PaymentMetricsParams) ([]domain.PaymentMethodBreakdown, error) {
	if m.breakdownErr != nil {
		return nil, m.breakdownErr
	}
	return m.paymentBreakdown, nil
}

func TestGetPaymentMetricsUsecase_Execute_Success(t *testing.T) {
	exec := &mockExecutor{}
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	expectedSummary := &domain.PaymentSummary{
		TotalPaid:          1000000,
		PaymentSuccessRate: 0.95,
	}
	expectedBreakdown := []domain.PaymentMethodBreakdown{
		{
			MethodID:    uuid.New(),
			MethodName:  "BCA VA",
			MethodType:  "bank_transfer",
			Count:       10,
			Amount:      1000000,
			SuccessRate: 1.0,
		},
	}

	repo := &mockPaymentRepo{
		paymentSummary:   expectedSummary,
		paymentBreakdown: expectedBreakdown,
	}

	uc := NewGetPaymentMetricsUsecase(exec, repo)

	input := GetPaymentMetricsInput{
		From: from,
		To:   to,
	}

	res, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Summary.TotalPaid != 1000000 {
		t.Errorf("expected TotalPaid=1000000, got %d", res.Summary.TotalPaid)
	}
	if len(res.Breakdown) != 1 {
		t.Errorf("expected 1 breakdown item, got %d", len(res.Breakdown))
	}
}

func TestGetPaymentMetricsUsecase_Execute_Errors(t *testing.T) {
	exec := &mockExecutor{}
	dbErr := errors.New("payment db error")

	t.Run("PaymentSummary error", func(t *testing.T) {
		repo := &mockPaymentRepo{paymentSumErr: dbErr}
		uc := NewGetPaymentMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetPaymentMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})

	t.Run("PaymentMethodBreakdown error", func(t *testing.T) {
		repo := &mockPaymentRepo{
			paymentSummary: &domain.PaymentSummary{},
			breakdownErr:   dbErr,
		}
		uc := NewGetPaymentMetricsUsecase(exec, repo)
		_, err := uc.Execute(context.Background(), GetPaymentMetricsInput{})
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})
}
