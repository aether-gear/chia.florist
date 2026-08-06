package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetOrderMetricsInput struct {
	From        time.Time
	To          time.Time
	Granularity string
	ShopID      *uuid.UUID
	TopN        int
}

type GetOrderMetricsResult struct {
	Summary     *domain.OrderSummary
	TimeSeries  []domain.OrderTimeSeries
	TopProducts []domain.TopProduct
	TopShops    []domain.TopShop
}

type GetOrderMetricsUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
}

func NewGetOrderMetricsUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
) *GetOrderMetricsUsecase {
	return &GetOrderMetricsUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
	}
}

func (u *GetOrderMetricsUsecase) Execute(
	ctx context.Context,
	input GetOrderMetricsInput,
) (*GetOrderMetricsResult, error) {
	params := repository.OrderMetricsParams{
		DateRange: repository.DateRangeParams{
			From: input.From,
			To:   input.To,
		},
		Granularity: repository.Granularity(input.Granularity),
		ShopID:      input.ShopID,
		TopN:        input.TopN,
	}

	summary, err := u.analyticsRepo.GetOrderSummary(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order summary: %w", err)
	}

	timeSeries, err := u.analyticsRepo.GetOrderTimeSeries(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve time series: %w", err)
	}

	topProducts, err := u.analyticsRepo.GetTopProducts(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve top products: %w", err)
	}

	topShops, err := u.analyticsRepo.GetTopShops(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve top shops: %w", err)
	}

	result := GetOrderMetricsResult{
		Summary:     summary,
		TimeSeries:  timeSeries,
		TopProducts: topProducts,
		TopShops:    topShops,
	}

	return &result, nil
}
