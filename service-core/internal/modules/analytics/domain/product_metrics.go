package domain

import (
	"github.com/google/uuid"
)

type ProductMetricsSummary struct {
	TopByRevenue    []TopProductStat
	TopByVolume     []TopProductStat
	AvgConversion   float64
	AvgReturnRate   float64
	InvoiceVoidRate float64
}

type TopProductStat struct {
	ProductID        uuid.UUID
	ProductName      string
	Revenue          int64
	UnitsSold        int
	ConversionRate   float64
	ReturnRate       *float64
	GrossMarginPct   *float64
	SalesVelocity7d  int
	SalesVelocity30d int
}
