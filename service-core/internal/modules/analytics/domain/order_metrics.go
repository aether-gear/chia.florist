package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderSummary struct {
	TotalOrders      int
	TotalGMV         int64
	TotalRevenue     int64
	TotalShippingFee int64
	AOV              float64
	CancellationRate float64
	PendingCount     int
	ConfirmedCount   int
	ProcessingCount  int
	ShippedCount     int
	DeliveredCount   int
	CancelledCount   int
}

type OrderTimeSeries struct {
	Date       time.Time
	OrderCount int
	GMV        int64
	AOV        float64
}

type TopProduct struct {
	ProductID   uuid.UUID
	ProductName string
	Quantity    int
	Revenue     int64
}

type TopShop struct {
	ShopID   uuid.UUID
	ShopName string
	Revenue  int64
	Orders   int
}
