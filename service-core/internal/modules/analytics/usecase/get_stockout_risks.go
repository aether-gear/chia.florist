package usecase

import (
	"context"
	"fmt"
	"sort"

	appclock "service-core/internal/common/clock"
	intelligencelayer "service-core/internal/infra/intelligence_layer"
	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetStockoutRisksInput struct {
	ShopID *uuid.UUID
}

type GetStockoutRisksUsecase struct {
	exec          transaction.Executor
	analyticsRepo repository.AnalyticsRepository
	aiProv        intelligencelayer.Provider
}

func NewGetStockoutRisksUsecase(
	exec transaction.Executor,
	analyticsRepo repository.AnalyticsRepository,
	aiProv intelligencelayer.Provider,
) *GetStockoutRisksUsecase {
	return &GetStockoutRisksUsecase{
		exec:          exec,
		analyticsRepo: analyticsRepo,
		aiProv:        aiProv,
	}
}

func (u *GetStockoutRisksUsecase) Execute(
	ctx context.Context,
	input GetStockoutRisksInput,
) ([]domain.StockoutRiskItem, error) {
	burnRates, err := u.analyticsRepo.GetInventoryBurnRates(ctx, u.exec, input.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inventory burn rates: %w", err)
	}

	now := appclock.Now()
	results := make([]domain.StockoutRiskItem, 0, len(burnRates))

	for _, item := range burnRates {
		burnRate := item.StockBurnRate7d
		if burnRate < 0.01 {
			burnRate = 0.01
		}

		availableStock := item.Stock - item.ReservedStock
		if availableStock < 0 {
			availableStock = 0
		}

		estDays := float64(availableStock) / burnRate
		urgencyRatio := (burnRate * item.SupplierLeadTimeDays) / float64(max(1, availableStock))

		riskLevel := "NORMAL"
		willStockout := false
		stockoutProb := 0.05

		if u.aiProv != nil {
			req := intelligencelayer.StockoutRiskRequest{
				ProductID:               item.ProductID.String(),
				Stock:                   float64(availableStock),
				ReservedStock:           float64(item.ReservedStock),
				StockBurnRate7d:         item.StockBurnRate7d,
				SupplierLeadTimeDays:    item.SupplierLeadTimeDays,
				EstimatedDaysToStockout: &estDays,
				ReorderUrgencyRatio:     &urgencyRatio,
			}

			resp, err := u.aiProv.PredictStockoutRisk(ctx, req)
			if err == nil && resp != nil {
				riskLevel = resp.RiskLevel
				willStockout = resp.WillStockout
				stockoutProb = resp.StockoutProbability
			} else {
				// Heuristic fallback
				if estDays <= 2.0 || availableStock <= 2 {
					riskLevel = "CRITICAL"
					willStockout = true
					stockoutProb = 0.90
				} else if estDays <= item.SupplierLeadTimeDays {
					riskLevel = "WARNING"
					willStockout = true
					stockoutProb = 0.60
				}
			}
		} else {
			if estDays <= 2.0 || availableStock <= 2 {
				riskLevel = "CRITICAL"
				willStockout = true
				stockoutProb = 0.90
			} else if estDays <= item.SupplierLeadTimeDays {
				riskLevel = "WARNING"
				willStockout = true
				stockoutProb = 0.60
			}
		}

		results = append(results, domain.StockoutRiskItem{
			ProductID:               item.ProductID,
			ProductName:             item.ProductName,
			ShopID:                  item.ShopID,
			ShopName:                item.ShopName,
			Stock:                   item.Stock,
			ReservedStock:           item.ReservedStock,
			AvailableStock:          availableStock,
			StockBurnRate7d:         item.StockBurnRate7d,
			SupplierLeadTimeDays:    item.SupplierLeadTimeDays,
			EstimatedDaysToStockout: estDays,
			ReorderUrgencyRatio:     urgencyRatio,
			StockoutProbability:     stockoutProb,
			WillStockout:            willStockout,
			RiskLevel:               riskLevel,
			EvaluatedAt:             now,
		})
	}

	// Sort by urgency ratio descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].ReorderUrgencyRatio > results[j].ReorderUrgencyRatio
	})

	return results, nil
}
