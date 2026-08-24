package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductDemandForecast struct {
	ProductID            uuid.UUID
	ProductName          string
	ShopID               *uuid.UUID
	PredictedUnitsSold7d float64
	ConfidenceTier       string
	HistoricalVelocity7d int
	CurrentStock         int
	ForecastGeneratedAt  time.Time
}

type ProductLagFeatures struct {
	ProductID       uuid.UUID
	ProductName     string
	GrossMarginPct  float64
	ViewCount       int
	UnitsSoldLag1   float64
	UnitsSoldLag7   float64
	UnitsSoldLag14  float64
	UnitsSoldLag30  float64
	Rolling7dMean   float64
	Rolling7dStd    float64
	Rolling30dMean  float64
	Rolling30dStd   float64
	CurrentStock    int
}
