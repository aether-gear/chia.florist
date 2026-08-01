package repository

import (
	"context"

	"service-core/internal/modules/analytics/domain"
	transaction "service-core/internal/shared/transaction"
)

type AnalyticsRepository interface {
	GetOrderSummary(
		ctx context.Context,
		exec transaction.Executor,
		params OrderMetricsParams,
	) (*domain.OrderSummary, error)
	GetOrderTimeSeries(
		ctx context.Context,
		exec transaction.Executor,
		params OrderMetricsParams,
	) ([]domain.OrderTimeSeries, error)
	GetTopProducts(
		ctx context.Context,
		exec transaction.Executor,
		params OrderMetricsParams,
	) ([]domain.TopProduct, error)
	GetTopShops(
		ctx context.Context,
		exec transaction.Executor,
		params OrderMetricsParams,
	) ([]domain.TopShop, error)

	GetPaymentSummary(
		ctx context.Context,
		exec transaction.Executor,
		params PaymentMetricsParams,
	) (*domain.PaymentSummary, error)
	GetPaymentMethodBreakdown(
		ctx context.Context,
		exec transaction.Executor,
		params PaymentMetricsParams,
	) ([]domain.PaymentMethodBreakdown, error)

	GetShipmentSummary(
		ctx context.Context,
		exec transaction.Executor,
		params ShipmentMetricsParams,
	) (*domain.ShipmentSummary, error)
	GetCourierBreakdown(
		ctx context.Context,
		exec transaction.Executor,
		params ShipmentMetricsParams,
	) ([]domain.CourierBreakdown, error)

	GetInventorySummary(
		ctx context.Context,
		exec transaction.Executor,
		params InventoryMetricsParams,
	) (*domain.InventorySummary, error)

	GetProductMetricsSummary(
		ctx context.Context,
		exec transaction.Executor,
		params ProductMetricsParams,
	) (*domain.ProductMetricsSummary, error)
}
