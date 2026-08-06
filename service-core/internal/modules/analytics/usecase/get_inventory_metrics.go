package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetInventoryMetricsInput struct {
	ShopID            *uuid.UUID
	LowStockThreshold int
}

type GetInventoryMetricsUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
}

func NewGetInventoryMetricsUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
) *GetInventoryMetricsUsecase {
	return &GetInventoryMetricsUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
	}
}

func (u *GetInventoryMetricsUsecase) Execute(
	ctx context.Context,
	input GetInventoryMetricsInput,
) (*domain.InventorySummary, error) {
	params := repository.InventoryMetricsParams{
		ShopID:            input.ShopID,
		LowStockThreshold: input.LowStockThreshold,
	}

	result, err := u.analyticsRepo.GetInventorySummary(ctx, u.exec,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory summary: %w", err)
	}

	return result, nil
}
