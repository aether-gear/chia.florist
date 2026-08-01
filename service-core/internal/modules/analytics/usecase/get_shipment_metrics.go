package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetShipmentMetricsInput struct {
	From time.Time
	To   time.Time
	TopN int
}

type GetShipmentMetricsResult struct {
	Summary  *domain.ShipmentSummary
	Couriers []domain.CourierBreakdown
}

type GetShipmentMetricsUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
}

func NewGetShipmentMetricsUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
) *GetShipmentMetricsUsecase {
	return &GetShipmentMetricsUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
	}
}

func (u *GetShipmentMetricsUsecase) Execute(
	ctx context.Context,
	input GetShipmentMetricsInput,
) (*GetShipmentMetricsResult, error) {
	params := repository.ShipmentMetricsParams{
		DateRange: repository.DateRangeParams{
			From: input.From,
			To:   input.To,
		},
		TopN: input.TopN,
	}

	summary, err := u.analyticsRepo.GetShipmentSummary(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shipment summary: %w", err)
	}

	couriers, err := u.analyticsRepo.GetCourierBreakdown(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load courier breakdowns: %w", err)
	}

	result := GetShipmentMetricsResult{
		Summary:  summary,
		Couriers: couriers,
	}

	return &result, nil
}
