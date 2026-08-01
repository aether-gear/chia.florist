package domain

type ShipmentSummary struct {
	Total             int
	Delivered         int
	Failed            int
	Returned          int
	Cancelled         int
	DeliveryRate      float64
	AvgFulfillmentSec float64
}

type CourierBreakdown struct {
	Courier      string
	Service      string
	Count        int
	DeliveryRate float64
	AvgCost      int64
}
