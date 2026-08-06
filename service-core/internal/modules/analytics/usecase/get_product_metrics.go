package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetProductMetricsInput struct {
	From time.Time
	To   time.Time
	TopN int
}

type GetProductMetricsUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
}

func NewGetProductMetricsUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
) *GetProductMetricsUsecase {
	return &GetProductMetricsUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
	}
}

func (u *GetProductMetricsUsecase) Execute(
	ctx context.Context,
	input GetProductMetricsInput,
) (*domain.ProductMetricsSummary, error) {
	params := repository.ProductMetricsParams{
		DateRange: repository.DateRangeParams{
			From: input.From,
			To:   input.To,
		},
		TopN: input.TopN,
	}

	productMetrics, err := u.analyticsRepo.GetProductMetricsSummary(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product metrics: %w", err)
	}

	return productMetrics, nil
}
