package job

import (
	"context"
	"time"

	applogger "service-core/internal/common/logger"
	analyticsUsecase "service-core/internal/modules/analytics/usecase"
)

type StockoutRiskScanJob struct {
	stockoutRisksUC *analyticsUsecase.GetStockoutRisksUsecase
	logger          applogger.Logger
	interval        time.Duration
}

func NewStockoutRiskScanJob(
	stockoutRisksUC *analyticsUsecase.GetStockoutRisksUsecase,
	interval time.Duration,
	logger applogger.Logger,
) *StockoutRiskScanJob {
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	return &StockoutRiskScanJob{
		stockoutRisksUC: stockoutRisksUC,
		logger:          logger,
		interval:        interval,
	}
}

func (j *StockoutRiskScanJob) Start(ctx context.Context) {
	if j.logger != nil {
		j.logger.Info(ctx, "stockout risk scan job: started",
			applogger.Field{Key: "interval", Value: j.interval.String()},
		)
	}

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if j.logger != nil {
				j.logger.Info(ctx, "stockout risk scan job: stopped")
			}
			return
		case <-ticker.C:
			j.Scan(ctx)
		}
	}
}

func (j *StockoutRiskScanJob) Scan(ctx context.Context) {
	if j.logger != nil {
		j.logger.Debug(ctx, "stockout risk scan job: tick — evaluating inventory stockout probabilities")
	}

	risks, err := j.stockoutRisksUC.Execute(ctx, analyticsUsecase.GetStockoutRisksInput{})
	if err != nil {
		if j.logger != nil {
			j.logger.Error(ctx, "stockout risk scan job: failed to scan inventory",
				applogger.Field{Key: "error", Value: err.Error()},
			)
		}
		return
	}

	criticalCount := 0
	for _, item := range risks {
		if item.RiskLevel == "CRITICAL" {
			criticalCount++
			if j.logger != nil {
				j.logger.Warn(ctx, "critical stockout risk detected",
					applogger.Field{Key: "product_id", Value: item.ProductID.String()},
					applogger.Field{Key: "product_name", Value: item.ProductName},
					applogger.Field{Key: "shop_id", Value: item.ShopID.String()},
					applogger.Field{Key: "stockout_probability", Value: item.StockoutProbability},
					applogger.Field{Key: "days_until_stockout", Value: item.EstimatedDaysToStockout},
				)
			}
		}
	}

	if j.logger != nil && (len(risks) > 0 || criticalCount > 0) {
		j.logger.Info(ctx, "stockout risk scan job: completed",
			applogger.Field{Key: "total_items_scanned", Value: len(risks)},
			applogger.Field{Key: "critical_risks_count", Value: criticalCount},
		)
	}
}
