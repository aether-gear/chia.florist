package http

import (
	"time"

	"github.com/google/uuid"
)

type orderSummaryResponse struct {
	TotalOrders      int     `json:"total_orders"`
	TotalGMV         int64   `json:"total_gmv"`
	TotalRevenue     int64   `json:"total_revenue"`
	TotalShippingFee int64   `json:"total_shipping_fee"`
	AOV              float64 `json:"aov"`
	CancellationRate float64 `json:"cancellation_rate"`
	PendingCount     int     `json:"pending_count"`
	ConfirmedCount   int     `json:"confirmed_count"`
	ProcessingCount  int     `json:"processing_count"`
	ShippedCount     int     `json:"shipped_count"`
	DeliveredCount   int     `json:"delivered_count"`
	CancelledCount   int     `json:"cancelled_count"`
}

type orderTimeSeriesResponse struct {
	Date       time.Time `json:"date"`
	OrderCount int       `json:"order_count"`
	GMV        int64     `json:"gmv"`
	AOV        float64   `json:"aov"`
}

type topProductResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Revenue     int64     `json:"revenue"`
}

type topShopResponse struct {
	ShopID   uuid.UUID `json:"shop_id"`
	ShopName string    `json:"shop_name"`
	Revenue  int64     `json:"revenue"`
	Orders   int       `json:"orders"`
}

type getOrderMetricsResponse struct {
	Summary     orderSummaryResponse      `json:"summary"`
	TimeSeries  []orderTimeSeriesResponse `json:"time_series"`
	TopProducts []topProductResponse      `json:"top_products"`
	TopShops    []topShopResponse         `json:"top_shops"`
}

type paymentSummaryResponse struct {
	TotalPaid          int64   `json:"total_paid"`
	TotalPending       int64   `json:"total_pending"`
	TotalExpired       int64   `json:"total_expired"`
	TotalRefunded      int64   `json:"total_refunded"`
	PaymentSuccessRate float64 `json:"payment_success_rate"`
	AvgTimeToPay       float64 `json:"avg_time_to_pay"`
}

type paymentMethodBreakdownResponse struct {
	MethodID    uuid.UUID `json:"method_id"`
	MethodName  string    `json:"method_name"`
	MethodType  string    `json:"method_type"`
	Count       int       `json:"count"`
	Amount      int64     `json:"amount"`
	SuccessRate float64   `json:"success_rate"`
}

type getPaymentMetricsResponse struct {
	Summary   paymentSummaryResponse           `json:"summary"`
	Breakdown []paymentMethodBreakdownResponse `json:"breakdown"`
}

type shipmentSummaryResponse struct {
	Total             int     `json:"total"`
	Delivered         int     `json:"delivered"`
	Failed            int     `json:"failed"`
	Returned          int     `json:"returned"`
	Cancelled         int     `json:"cancelled"`
	DeliveryRate      float64 `json:"delivery_rate"`
	AvgFulfillmentSec float64 `json:"avg_fulfillment_sec"`
}

type courierBreakdownResponse struct {
	Courier      string  `json:"courier"`
	Service      string  `json:"service"`
	Count        int     `json:"count"`
	DeliveryRate float64 `json:"delivery_rate"`
	AvgCost      int64   `json:"avg_cost"`
}

type getShipmentMetricsResponse struct {
	Summary  shipmentSummaryResponse    `json:"summary"`
	Couriers []courierBreakdownResponse `json:"couriers"`
}

type inventorySummaryResponse struct {
	TotalProducts  int `json:"total_products"`
	TotalStock     int `json:"total_stock"`
	TotalReserved  int `json:"total_reserved"`
	TotalAvailable int `json:"total_available"`
	StockoutCount  int `json:"stockout_count"`
	LowStockCount  int `json:"low_stock_count"`
}

type topProductStatResponse struct {
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	Revenue          int64     `json:"revenue"`
	UnitsSold        int       `json:"units_sold"`
	ConversionRate   float64   `json:"conversion_rate"`
	ReturnRate       *float64  `json:"return_rate,omitempty"`
	GrossMarginPct   *float64  `json:"gross_margin_pct,omitempty"`
	SalesVelocity7d  int       `json:"sales_velocity_7d"`
	SalesVelocity30d int       `json:"sales_velocity_30d"`
}

type productMetricsSummaryResponse struct {
	TopByRevenue    []topProductStatResponse `json:"top_by_revenue"`
	TopByVolume     []topProductStatResponse `json:"top_by_volume"`
	AvgConversion   float64                  `json:"avg_conversion"`
	AvgReturnRate   float64                  `json:"avg_return_rate"`
	InvoiceVoidRate float64                  `json:"invoice_void_rate"`
}

type demandForecastResponse struct {
	ProductID            uuid.UUID  `json:"product_id"`
	ProductName          string     `json:"product_name"`
	ShopID               *uuid.UUID `json:"shop_id,omitempty"`
	PredictedUnitsSold7d float64    `json:"predicted_units_sold_7d"`
	ConfidenceTier       string     `json:"confidence_tier"`
	HistoricalVelocity7d int        `json:"historical_velocity_7d"`
	CurrentStock         int        `json:"current_stock"`
	ForecastGeneratedAt  time.Time  `json:"forecast_generated_at"`
}

type stockoutRiskItemResponse struct {
	ProductID               uuid.UUID `json:"product_id"`
	ProductName             string    `json:"product_name"`
	ShopID                  uuid.UUID `json:"shop_id"`
	ShopName                string    `json:"shop_name"`
	Stock                   int       `json:"stock"`
	ReservedStock           int       `json:"reserved_stock"`
	AvailableStock          int       `json:"available_stock"`
	StockBurnRate7d         float64   `json:"stock_burn_rate_7d"`
	SupplierLeadTimeDays    float64   `json:"supplier_lead_time_days"`
	EstimatedDaysToStockout float64   `json:"estimated_days_to_stockout"`
	ReorderUrgencyRatio     float64   `json:"reorder_urgency_ratio"`
	StockoutProbability     float64   `json:"stockout_probability"`
	WillStockout            bool      `json:"will_stockout"`
	RiskLevel               string    `json:"risk_level"`
	EvaluatedAt             time.Time `json:"evaluated_at"`
}
