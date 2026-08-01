package domain

type InventorySummary struct {
	TotalProducts  int
	TotalStock     int
	TotalReserved  int
	TotalAvailable int
	StockoutCount  int
	LowStockCount  int
}
