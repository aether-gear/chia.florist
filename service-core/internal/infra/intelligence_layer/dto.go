package intelligencelayer

// APIResponse represents the standard JSON response envelope from the AI Lab FastAPI server.
type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    *T   `json:"data"`
	Error   any  `json:"error"`
	Meta    any  `json:"meta"`
}

// HealthData represents the system status and loaded models.
type HealthData struct {
	Status       string            `json:"status"`
	Version      string            `json:"version"`
	LoadedModels []string          `json:"loaded_models"`
	ModelDetails map[string]string `json:"model_details"`
}

// DemandForecastRequest matches POST /api/v1/predict/demand
type DemandForecastRequest struct {
	ProductID      string  `json:"product_id"`
	GrossMarginPct float64 `json:"gross_margin_pct"`
	ViewCount      int     `json:"view_count"`
	DayOfWeek      int     `json:"day_of_week"`
	Month          int     `json:"month"`
	IsWeekend      float64 `json:"is_weekend"`
	UnitsSoldLag1  float64 `json:"units_sold_lag_1"`
	UnitsSoldLag7  float64 `json:"units_sold_lag_7"`
	UnitsSoldLag14 float64 `json:"units_sold_lag_14"`
	UnitsSoldLag30 float64 `json:"units_sold_lag_30"`
	Rolling7dMean  float64 `json:"rolling_7d_mean"`
	Rolling7dStd   float64 `json:"rolling_7d_std"`
	Rolling30dMean float64 `json:"rolling_30d_mean"`
	Rolling30dStd  float64 `json:"rolling_30d_std"`
}

// DemandForecastResponse matches response from POST /api/v1/predict/demand
type DemandForecastResponse struct {
	ProductID            string  `json:"product_id"`
	PredictedUnitsSold7d float64 `json:"predicted_units_sold_7d"`
	ConfidenceTier       string  `json:"confidence_tier"`
}

// StockoutRiskRequest matches POST /api/v1/predict/stockout-risk
type StockoutRiskRequest struct {
	ProductID               string   `json:"product_id"`
	Stock                   float64  `json:"stock"`
	ReservedStock           float64  `json:"reserved_stock"`
	StockBurnRate7d         float64  `json:"stock_burn_rate_7d"`
	SupplierLeadTimeDays    float64  `json:"supplier_lead_time_days"`
	EstimatedDaysToStockout *float64 `json:"estimated_days_to_stockout,omitempty"`
	ReorderUrgencyRatio     *float64 `json:"reorder_urgency_ratio,omitempty"`
}

// StockoutRiskResponse matches response from POST /api/v1/predict/stockout-risk
type StockoutRiskResponse struct {
	ProductID           string  `json:"product_id"`
	StockoutProbability float64 `json:"stockout_probability"`
	WillStockout        bool    `json:"will_stockout"`
	RiskLevel           string  `json:"risk_level"`
}

// CourierSLARequest matches POST /api/v1/predict/courier-sla
type CourierSLARequest struct {
	CourierCode       string  `json:"courier_code"`
	ShippingCost      float64 `json:"shipping_cost"`
	DispatchDayOfWeek int     `json:"dispatch_day_of_week"`
	DispatchHour      int     `json:"dispatch_hour"`
	DispatchIsWeekend float64 `json:"dispatch_is_weekend"`
}

// CourierSLAResponse matches response from POST /api/v1/predict/courier-sla
type CourierSLAResponse struct {
	CourierCode            string  `json:"courier_code"`
	EstimatedDurationHours float64 `json:"estimated_duration_hours"`
	SLAConfidenceScore     float64 `json:"sla_confidence_score"`
	DeliveryStatus         string  `json:"delivery_status"`
}

// AnomalyCheckRequest matches POST /api/v1/predict/anomaly
type AnomalyCheckRequest struct {
	Amount           float64 `json:"amount"`
	TimeToPaySec     float64 `json:"time_to_pay_sec"`
	IsFailedStatus   float64 `json:"is_failed_status"`
	IsManualTransfer float64 `json:"is_manual_transfer"`
}

// AnomalyCheckResponse matches response from POST /api/v1/predict/anomaly
type AnomalyCheckResponse struct {
	IsAnomaly    bool     `json:"is_anomaly"`
	AnomalyScore float64  `json:"anomaly_score"`
	Severity     string   `json:"severity"`
	Reasons      []string `json:"reasons"`
}
