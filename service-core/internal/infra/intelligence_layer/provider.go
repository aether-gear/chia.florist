package intelligencelayer

import (
	"context"
)

// Provider defines the contract for interacting with the AI Lab intelligence-layer microservice.
type Provider interface {
	// HealthCheck returns the health status and loaded ML models.
	HealthCheck(ctx context.Context) (*HealthData, error)

	// PredictDemand forecasts SKU unit sales for the next 7 days.
	PredictDemand(ctx context.Context, req DemandForecastRequest) (*DemandForecastResponse, error)

	// PredictStockoutRisk evaluates the probability and risk level of an inventory stockout.
	PredictStockoutRisk(ctx context.Context, req StockoutRiskRequest) (*StockoutRiskResponse, error)

	// PredictCourierSLA estimates delivery duration in hours and SLA confidence.
	PredictCourierSLA(ctx context.Context, req CourierSLARequest) (*CourierSLAResponse, error)

	// DetectAnomaly detects operational anomalies and payment delay patterns.
	DetectAnomaly(ctx context.Context, req AnomalyCheckRequest) (*AnomalyCheckResponse, error)
}
