package usecase

import (
	"context"
	"fmt"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	intelligencelayer "service-core/internal/infra/intelligence_layer"
	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetDemandForecastInput struct {
	ProductID uuid.UUID
	ShopID    *uuid.UUID
}

type GetDemandForecastUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
	aiProv        intelligencelayer.Provider
}

func NewGetDemandForecastUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
	aiProv intelligencelayer.Provider,
) *GetDemandForecastUsecase {
	return &GetDemandForecastUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
		aiProv:        aiProv,
	}
}

func (u *GetDemandForecastUsecase) Execute(
	ctx context.Context,
	input GetDemandForecastInput,
) (*domain.ProductDemandForecast, error) {
	if input.ProductID == uuid.Nil {
		return nil, apperrors.NewInvalidInput("product_id is required")
	}

	lagFeat, err := u.analyticsRepo.GetProductLagFeatures(ctx, u.exec, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product lag features: %w", err)
	}
	if lagFeat == nil {
		return nil, apperrors.NewNotFound("product not found")
	}

	now := appclock.Now()
	pythonDayOfWeek := (int(now.Weekday()) + 6) % 7
	isWeekend := 0.0
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		isWeekend = 1.0
	}

	var predictedUnits float64
	confidenceTier := "baseline"

	if u.aiProv != nil {
		req := intelligencelayer.DemandForecastRequest{
			ProductID:      lagFeat.ProductID.String(),
			GrossMarginPct: lagFeat.GrossMarginPct,
			ViewCount:      lagFeat.ViewCount,
			DayOfWeek:      pythonDayOfWeek,
			Month:          int(now.Month()),
			IsWeekend:      isWeekend,
			UnitsSoldLag1:  lagFeat.UnitsSoldLag1,
			UnitsSoldLag7:  lagFeat.UnitsSoldLag7,
			UnitsSoldLag14: lagFeat.UnitsSoldLag14,
			UnitsSoldLag30: lagFeat.UnitsSoldLag30,
			Rolling7dMean:  lagFeat.Rolling7dMean,
			Rolling7dStd:   lagFeat.Rolling7dStd,
			Rolling30dMean: lagFeat.Rolling30dMean,
			Rolling30dStd:  lagFeat.Rolling30dStd,
		}

		resp, err := u.aiProv.PredictDemand(ctx, req)
		if err == nil && resp != nil {
			predictedUnits = resp.PredictedUnitsSold7d
			confidenceTier = resp.ConfidenceTier
		} else {
			// Fallback to moving average if AI service unavailable
			predictedUnits = lagFeat.Rolling7dMean * 7.0
		}
	} else {
		predictedUnits = lagFeat.Rolling7dMean * 7.0
	}

	return &domain.ProductDemandForecast{
		ProductID:            lagFeat.ProductID,
		ProductName:          lagFeat.ProductName,
		ShopID:               input.ShopID,
		PredictedUnitsSold7d: predictedUnits,
		ConfidenceTier:       confidenceTier,
		HistoricalVelocity7d: int(lagFeat.Rolling7dMean * 7.0),
		CurrentStock:         lagFeat.CurrentStock,
		ForecastGeneratedAt:  now,
	}, nil
}
