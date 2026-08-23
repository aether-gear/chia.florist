package http

import (
	"net/http"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/analytics/usecase"

	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	getOrderMetrics     *usecase.GetOrderMetricsUsecase
	getPaymentMetrics   *usecase.GetPaymentMetricsUsecase
	getShipmentMetrics  *usecase.GetShipmentMetricsUsecase
	getInventoryMetrics *usecase.GetInventoryMetricsUsecase
	getProductMetrics   *usecase.GetProductMetricsUsecase
	getDemandForecast   *usecase.GetDemandForecastUsecase
	getStockoutRisks    *usecase.GetStockoutRisksUsecase
}

func NewAnalyticsHandler(
	getOrderMetrics *usecase.GetOrderMetricsUsecase,
	getPaymentMetrics *usecase.GetPaymentMetricsUsecase,
	getShipmentMetrics *usecase.GetShipmentMetricsUsecase,
	getInventoryMetrics *usecase.GetInventoryMetricsUsecase,
	getProductMetrics *usecase.GetProductMetricsUsecase,
	getDemandForecast *usecase.GetDemandForecastUsecase,
	getStockoutRisks *usecase.GetStockoutRisksUsecase,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		getOrderMetrics:     getOrderMetrics,
		getPaymentMetrics:   getPaymentMetrics,
		getShipmentMetrics:  getShipmentMetrics,
		getInventoryMetrics: getInventoryMetrics,
		getProductMetrics:   getProductMetrics,
		getDemandForecast:   getDemandForecast,
		getStockoutRisks:    getStockoutRisks,
	}
}

func (h *AnalyticsHandler) GetOrderMetrics(w http.ResponseWriter, r *http.Request) error {
	from, to, err := parseDateRange(r)
	if err != nil {
		return err
	}

	granularity := parseGranularity(r)
	topN := apphttp.QueryIntDefault(r, "top_n", 10)
	shopID, _ := apphttp.QueryUUID(r, "shop_id")

	input := usecase.GetOrderMetricsInput{
		From:        from,
		To:          to,
		Granularity: granularity,
		ShopID:      shopID,
		TopN:        topN,
	}

	result, err := h.getOrderMetrics.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	timeSeriesResp := make([]orderTimeSeriesResponse, 0, len(result.TimeSeries))
	for _, ts := range result.TimeSeries {
		timeSeriesResp = append(timeSeriesResp, orderTimeSeriesResponse{
			Date:       ts.Date,
			OrderCount: ts.OrderCount,
			GMV:        ts.GMV,
			AOV:        ts.AOV,
		})
	}

	topProductsResp := make([]topProductResponse, 0, len(result.TopProducts))
	for _, tp := range result.TopProducts {
		topProductsResp = append(topProductsResp, topProductResponse{
			ProductID:   tp.ProductID,
			ProductName: tp.ProductName,
			Quantity:    tp.Quantity,
			Revenue:     tp.Revenue,
		})
	}

	topShopsResp := make([]topShopResponse, 0, len(result.TopShops))
	for _, ts := range result.TopShops {
		topShopsResp = append(topShopsResp, topShopResponse{
			ShopID:   ts.ShopID,
			ShopName: ts.ShopName,
			Revenue:  ts.Revenue,
			Orders:   ts.Orders,
		})
	}

	resp := getOrderMetricsResponse{
		Summary: orderSummaryResponse{
			TotalOrders:      result.Summary.TotalOrders,
			TotalGMV:         result.Summary.TotalGMV,
			TotalRevenue:     result.Summary.TotalRevenue,
			TotalShippingFee: result.Summary.TotalShippingFee,
			AOV:              result.Summary.AOV,
			CancellationRate: result.Summary.CancellationRate,
			PendingCount:     result.Summary.PendingCount,
			ConfirmedCount:   result.Summary.ConfirmedCount,
			ProcessingCount:  result.Summary.ProcessingCount,
			ShippedCount:     result.Summary.ShippedCount,
			DeliveredCount:   result.Summary.DeliveredCount,
			CancelledCount:   result.Summary.CancelledCount,
		},
		TimeSeries:  timeSeriesResp,
		TopProducts: topProductsResp,
		TopShops:    topShopsResp,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetPaymentMetrics(w http.ResponseWriter, r *http.Request) error {
	from, to, err := parseDateRange(r)
	if err != nil {
		return err
	}

	input := usecase.GetPaymentMetricsInput{
		From: from,
		To:   to,
	}

	result, err := h.getPaymentMetrics.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	breakdownResp := make([]paymentMethodBreakdownResponse, 0, len(result.Breakdown))
	for _, b := range result.Breakdown {
		breakdownResp = append(breakdownResp, paymentMethodBreakdownResponse{
			MethodID:    b.MethodID,
			MethodName:  b.MethodName,
			MethodType:  b.MethodType,
			Count:       b.Count,
			Amount:      b.Amount,
			SuccessRate: b.SuccessRate,
		})
	}

	resp := getPaymentMetricsResponse{
		Summary: paymentSummaryResponse{
			TotalPaid:          result.Summary.TotalPaid,
			TotalPending:       result.Summary.TotalPending,
			TotalExpired:       result.Summary.TotalExpired,
			TotalRefunded:      result.Summary.TotalRefunded,
			PaymentSuccessRate: result.Summary.PaymentSuccessRate,
			AvgTimeToPay:       result.Summary.AvgTimeToPay,
		},
		Breakdown: breakdownResp,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetShipmentMetrics(w http.ResponseWriter, r *http.Request) error {
	from, to, err := parseDateRange(r)
	if err != nil {
		return err
	}

	topN := apphttp.QueryIntDefault(r, "top_n", 10)

	input := usecase.GetShipmentMetricsInput{
		From: from,
		To:   to,
		TopN: topN,
	}

	result, err := h.getShipmentMetrics.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	couriersResp := make([]courierBreakdownResponse, 0, len(result.Couriers))
	for _, c := range result.Couriers {
		couriersResp = append(couriersResp, courierBreakdownResponse{
			Courier:      c.Courier,
			Service:      c.Service,
			Count:        c.Count,
			DeliveryRate: c.DeliveryRate,
			AvgCost:      c.AvgCost,
		})
	}

	resp := getShipmentMetricsResponse{
		Summary: shipmentSummaryResponse{
			Total:             result.Summary.Total,
			Delivered:         result.Summary.Delivered,
			Failed:            result.Summary.Failed,
			Returned:          result.Summary.Returned,
			Cancelled:         result.Summary.Cancelled,
			DeliveryRate:      result.Summary.DeliveryRate,
			AvgFulfillmentSec: result.Summary.AvgFulfillmentSec,
		},
		Couriers: couriersResp,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetInventoryMetrics(w http.ResponseWriter, r *http.Request) error {
	var shopID *uuid.UUID
	shopIDStr := apphttp.Query(r, "shop_id")
	if shopIDStr != "" {
		parsed, err := uuid.Parse(shopIDStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid shop_id")
		}
		shopID = &parsed
	}

	lowStockThreshold := apphttp.QueryIntDefault(r, "low_stock_threshold", 5)

	input := usecase.GetInventoryMetricsInput{
		ShopID:            shopID,
		LowStockThreshold: lowStockThreshold,
	}

	summary, err := h.getInventoryMetrics.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	resp := inventorySummaryResponse{
		TotalProducts:  summary.TotalProducts,
		TotalStock:     summary.TotalStock,
		TotalReserved:  summary.TotalReserved,
		TotalAvailable: summary.TotalAvailable,
		StockoutCount:  summary.StockoutCount,
		LowStockCount:  summary.LowStockCount,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetProductMetrics(w http.ResponseWriter, r *http.Request) error {
	from, to, err := parseDateRange(r)
	if err != nil {
		return err
	}

	topN := apphttp.QueryIntDefault(r, "top_n", 10)

	input := usecase.GetProductMetricsInput{
		From: from,
		To:   to,
		TopN: topN,
	}

	summary, err := h.getProductMetrics.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	topByRevResp := make([]topProductStatResponse, 0, len(summary.TopByRevenue))
	for _, p := range summary.TopByRevenue {
		topByRevResp = append(topByRevResp, topProductStatResponse{
			ProductID:        p.ProductID,
			ProductName:      p.ProductName,
			Revenue:          p.Revenue,
			UnitsSold:        p.UnitsSold,
			ConversionRate:   p.ConversionRate,
			ReturnRate:       p.ReturnRate,
			GrossMarginPct:   p.GrossMarginPct,
			SalesVelocity7d:  p.SalesVelocity7d,
			SalesVelocity30d: p.SalesVelocity30d,
		})
	}

	topByVolResp := make([]topProductStatResponse, 0, len(summary.TopByVolume))
	for _, p := range summary.TopByVolume {
		topByVolResp = append(topByVolResp, topProductStatResponse{
			ProductID:        p.ProductID,
			ProductName:      p.ProductName,
			Revenue:          p.Revenue,
			UnitsSold:        p.UnitsSold,
			ConversionRate:   p.ConversionRate,
			ReturnRate:       p.ReturnRate,
			GrossMarginPct:   p.GrossMarginPct,
			SalesVelocity7d:  p.SalesVelocity7d,
			SalesVelocity30d: p.SalesVelocity30d,
		})
	}

	resp := productMetricsSummaryResponse{
		TopByRevenue:    topByRevResp,
		TopByVolume:     topByVolResp,
		AvgConversion:   summary.AvgConversion,
		AvgReturnRate:   summary.AvgReturnRate,
		InvoiceVoidRate: summary.InvoiceVoidRate,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetDemandForecast(w http.ResponseWriter, r *http.Request) error {
	productIDStr := apphttp.Query(r, "product_id")
	if productIDStr == "" {
		return apperrors.NewBadRequest("product_id query param is required")
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid product_id format")
	}

	var shopID *uuid.UUID
	shopIDStr := apphttp.Query(r, "shop_id")
	if shopIDStr != "" {
		parsed, err := uuid.Parse(shopIDStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid shop_id format")
		}
		shopID = &parsed
	}

	res, err := h.getDemandForecast.Execute(r.Context(), usecase.GetDemandForecastInput{
		ProductID: productID,
		ShopID:    shopID,
	})
	if err != nil {
		return err
	}

	resp := demandForecastResponse{
		ProductID:            res.ProductID,
		ProductName:          res.ProductName,
		ShopID:               res.ShopID,
		PredictedUnitsSold7d: res.PredictedUnitsSold7d,
		ConfidenceTier:       res.ConfidenceTier,
		HistoricalVelocity7d: res.HistoricalVelocity7d,
		CurrentStock:         res.CurrentStock,
		ForecastGeneratedAt:  res.ForecastGeneratedAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *AnalyticsHandler) GetStockoutRisks(w http.ResponseWriter, r *http.Request) error {
	var shopID *uuid.UUID
	shopIDStr := apphttp.Query(r, "shop_id")
	if shopIDStr != "" {
		parsed, err := uuid.Parse(shopIDStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid shop_id format")
		}
		shopID = &parsed
	}

	items, err := h.getStockoutRisks.Execute(r.Context(), usecase.GetStockoutRisksInput{
		ShopID: shopID,
	})
	if err != nil {
		return err
	}

	respList := make([]stockoutRiskItemResponse, 0, len(items))
	for _, item := range items {
		respList = append(respList, stockoutRiskItemResponse{
			ProductID:               item.ProductID,
			ProductName:             item.ProductName,
			ShopID:                  item.ShopID,
			ShopName:                item.ShopName,
			Stock:                   item.Stock,
			ReservedStock:           item.ReservedStock,
			AvailableStock:          item.AvailableStock,
			StockBurnRate7d:         item.StockBurnRate7d,
			SupplierLeadTimeDays:    item.SupplierLeadTimeDays,
			EstimatedDaysToStockout: item.EstimatedDaysToStockout,
			ReorderUrgencyRatio:     item.ReorderUrgencyRatio,
			StockoutProbability:     item.StockoutProbability,
			WillStockout:            item.WillStockout,
			RiskLevel:               item.RiskLevel,
			EvaluatedAt:             item.EvaluatedAt,
		})
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{
		"risks": respList,
		"count": len(respList),
	})
	return nil
}

func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	now := appclock.Now()
	from := now.AddDate(0, 0, -30)
	to := now

	fromStr := apphttp.Query(r, "from")
	if fromStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			parsedDate, errDate := time.Parse("2006-01-02", fromStr)
			if errDate != nil {
				return from, to, apperrors.NewBadRequest("invalid 'from' date format, must be RFC3339 or YYYY-MM-DD")
			}
			from = parsedDate
		} else {
			from = parsed
		}
	}

	toStr := apphttp.Query(r, "to")
	if toStr != "" {
		parsed, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			parsedDate, errDate := time.Parse("2006-01-02", toStr)
			if errDate != nil {
				return from, to, apperrors.NewBadRequest("invalid 'to' date format, must be RFC3339 or YYYY-MM-DD")
			}
			to = parsedDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		} else {
			to = parsed
		}
	}

	return from, to, nil
}

func parseGranularity(r *http.Request) string {
	granularity := apphttp.Query(r, "granularity")
	switch granularity {
	case "weekly", "monthly":
		return granularity
	default:
		return "daily"
	}
}
