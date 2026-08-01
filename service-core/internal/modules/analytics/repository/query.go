package repository

import (
	"time"

	"github.com/google/uuid"
)

type DateRangeParams struct {
	From time.Time
	To   time.Time
}

type Granularity string

const (
	GranularityDaily   Granularity = "daily"
	GranularityWeekly  Granularity = "weekly"
	GranularityMonthly Granularity = "monthly"
)

type OrderMetricsParams struct {
	DateRange   DateRangeParams
	Granularity Granularity
	ShopID      *uuid.UUID
	TopN        int
}

type PaymentMetricsParams struct {
	DateRange DateRangeParams
}

type ShipmentMetricsParams struct {
	DateRange DateRangeParams
	TopN      int
}

type InventoryMetricsParams struct {
	ShopID            *uuid.UUID
	LowStockThreshold int
}

type ProductMetricsParams struct {
	DateRange DateRangeParams
	TopN      int
}
