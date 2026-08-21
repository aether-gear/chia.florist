package job

import (
	"context"
	"time"

	applogger "service-core/internal/common/logger"
	analyticsUsecase "service-core/internal/modules/analytics/usecase"
)

type StockoutRiskScanJob struct {
	stockoutRisksUC *analyticsUsecase.GetStockoutRisksUsecase
	auditLogger     applogger.AuditLogger
	logger          applogger.Logger
	interval        time.Duration
}

func NewStockoutRiskScanJob(
	stockoutRisksUC *analyticsUsecase.GetStockoutRisksUsecase,
	auditLogger applogger.AuditLogger,
	interval time.Duration,
	logger applogger.Logger,
) *StockoutRiskScanJob {
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	return &StockoutRiskScanJob{
		stockoutRisksUC: stockoutRisksUC,
		auditLogger:     auditLogger,
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
		j.logger.Info(ctx, "stockout risk scan job: tick — evaluating inventory stockout probabilities")
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
			if j.auditLogger != nil {
				j.auditLogger.Log(ctx, applogger.AuditEvent{
					Category:   "inventory",
					Action:     "critical_stockout_risk_alert",
					Resource:   "inventory",
					ResourceID: item.ProductID.String(),
					Outcome:    applogger.OutcomeSuccess,
					Metadata: map[string]any{
						"product_name":               item.ProductName,
						"shop_id":                    item.ShopID.String(),
						"shop_name":                  item.ShopName,
						"stock":                      item.Stock,
						"reserved_stock":             item.ReservedStock,
						"available_stock":            item.AvailableStock,
						"stock_burn_rate_7d":         item.StockBurnRate7d,
						"estimated_days_to_stockout": item.EstimatedDaysToStockout,
						"reorder_urgency_ratio":      item.ReorderUrgencyRatio,
						"stockout_probability":       item.StockoutProbability,
					},
				})
			}
		}
	}

	if j.logger != nil {
		j.logger.Info(ctx, "stockout risk scan job: completed",
			applogger.Field{Key: "total_items_scanned", Value: len(risks)},
			applogger.Field{Key: "critical_risks_count", Value: criticalCount},
		)
	}
}
