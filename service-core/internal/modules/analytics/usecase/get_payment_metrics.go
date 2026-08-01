package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetPaymentMetricsInput struct {
	From time.Time
	To   time.Time
}

type GetPaymentMetricsResult struct {
	Summary   *domain.PaymentSummary
	Breakdown []domain.PaymentMethodBreakdown
}

type GetPaymentMetricsUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
}

func NewGetPaymentMetricsUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
) *GetPaymentMetricsUsecase {
	return &GetPaymentMetricsUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
	}
}

func (u *GetPaymentMetricsUsecase) Execute(
	ctx context.Context,
	input GetPaymentMetricsInput,
) (*GetPaymentMetricsResult, error) {
	params := repository.PaymentMetricsParams{
		DateRange: repository.DateRangeParams{
			From: input.From,
			To:   input.To,
		},
	}

	summary, err := u.analyticsRepo.GetPaymentSummary(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment summary: %w", err)
	}

	breakdown, err := u.analyticsRepo.GetPaymentMethodBreakdown(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to payment method breakdown: %w", err)
	}

	result := GetPaymentMetricsResult{
		Summary:   summary,
		Breakdown: breakdown,
	}

	return &result, nil
}
