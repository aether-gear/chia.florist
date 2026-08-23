package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockoutRiskItem struct {
	ProductID               uuid.UUID
	ProductName             string
	ShopID                  uuid.UUID
	ShopName                string
	Stock                   int
	ReservedStock           int
	AvailableStock          int
	StockBurnRate7d         float64
	SupplierLeadTimeDays    float64
	EstimatedDaysToStockout float64
	ReorderUrgencyRatio     float64
	StockoutProbability     float64
	WillStockout            bool
	RiskLevel               string
	EvaluatedAt             time.Time
}

type InventoryBurnRateData struct {
	ProductID            uuid.UUID
	ProductName          string
	ShopID               uuid.UUID
	ShopName             string
	Stock                int
	ReservedStock        int
	UnitsSold7d          int
	StockBurnRate7d      float64
	SupplierLeadTimeDays float64
}
