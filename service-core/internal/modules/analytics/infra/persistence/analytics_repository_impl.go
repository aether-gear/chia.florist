package persistence

import (
	"context"
	"fmt"

	"service-core/internal/modules/analytics/domain"
	"service-core/internal/modules/analytics/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/jackc/pgx/v5"
)

type analyticsRepositoryImpl struct{}

func NewAnalyticsRepositoryImpl() repository.AnalyticsRepository {
	return &analyticsRepositoryImpl{}
}

func (r *analyticsRepositoryImpl) GetOrderSummary(
	ctx context.Context,
	exec transaction.Executor,
	params repository.OrderMetricsParams,
) (*domain.OrderSummary, error) {
	query := `
		SELECT
			COALESCE(COUNT(*), 0)
				AS total_orders,
			COALESCE(SUM(total), 0)
				AS total_gmv,
			COALESCE(SUM(CASE WHEN status != 'cancelled' THEN total ELSE 0 END), 0)
				AS total_revenue,
			COALESCE(SUM(shipping_fee), 0)
				AS total_shipping_fee,
			COALESCE(AVG(CASE WHEN status != 'cancelled' THEN total END), 0.0)
				AS aov,
			COALESCE(COUNT(*) FILTER (WHERE status = 'pending'), 0)
				AS pending_count,
			COALESCE(COUNT(*) FILTER (WHERE status = 'confirmed'), 0)
				AS confirmed_count,
			COALESCE(COUNT(*) FILTER (WHERE status = 'processing'), 0)
				AS processing_count,
			COALESCE(COUNT(*) FILTER (WHERE status = 'shipped'), 0)
				AS shipped_count,
			COALESCE(COUNT(*) FILTER (WHERE status = 'delivered'), 0)
				AS delivered_count,
			COALESCE(COUNT(*) FILTER (WHERE status = 'cancelled'), 0)
				AS cancelled_count
		FROM
			orders
		WHERE
			created_at BETWEEN $1 AND $2
	`

	args := []any{params.DateRange.From, params.DateRange.To}
	if params.ShopID != nil {
		query += ` AND id IN (SELECT DISTINCT order_id FROM order_items WHERE shop_id = $3)`
		args = append(args, *params.ShopID)
	}

	var summary domain.OrderSummary
	err := exec.QueryRow(ctx, query, args...).Scan(
		&summary.TotalOrders,
		&summary.TotalGMV,
		&summary.TotalRevenue,
		&summary.TotalShippingFee,
		&summary.AOV,
		&summary.PendingCount,
		&summary.ConfirmedCount,
		&summary.ProcessingCount,
		&summary.ShippedCount,
		&summary.DeliveredCount,
		&summary.CancelledCount,
	)
	if err != nil {
		return nil, fmt.Errorf("query order summary failed: %w", err)
	}

	if summary.TotalOrders > 0 {
		summary.CancellationRate = float64(summary.CancelledCount) / float64(summary.TotalOrders)
	}

	return &summary, nil
}

func (r *analyticsRepositoryImpl) GetOrderTimeSeries(
	ctx context.Context,
	exec transaction.Executor,
	params repository.OrderMetricsParams,
) ([]domain.OrderTimeSeries, error) {
	truncUnit := "day"
	switch params.Granularity {
	case repository.GranularityWeekly:
		truncUnit = "week"
	case repository.GranularityMonthly:
		truncUnit = "month"
	}

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at)
				AS date_bucket,
			COALESCE(COUNT(*), 0)
				AS order_count,
			COALESCE(SUM(total), 0)
				AS gmv,
			COALESCE(AVG(total), 0.0)
				AS aov
		FROM
			orders
		WHERE
			created_at BETWEEN $1 AND $2
	`, truncUnit)

	args := []any{params.DateRange.From, params.DateRange.To}
	if params.ShopID != nil {
		query += ` AND id IN (SELECT DISTINCT order_id FROM order_items WHERE shop_id = $3)`
		args = append(args, *params.ShopID)
	}

	query += fmt.Sprintf(` GROUP BY DATE_TRUNC('%s', created_at) ORDER BY date_bucket ASC`, truncUnit)

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query order time series failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.OrderTimeSeries, error) {
		var item domain.OrderTimeSeries
		err := row.Scan(
			&item.Date,
			&item.OrderCount,
			&item.GMV,
			&item.AOV,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan order time series failed: %w", err)
	}

	return items, nil
}

func (r *analyticsRepositoryImpl) GetTopProducts(
	ctx context.Context,
	exec transaction.Executor,
	params repository.OrderMetricsParams,
) ([]domain.TopProduct, error) {
	limit := params.TopN
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT
			oi.product_id,
			oi.product_name,
			COALESCE(SUM(oi.quantity), 0)
				AS quantity,
			COALESCE(SUM(oi.subtotal), 0)
				AS revenue
		FROM
			order_items oi
		JOIN
			orders o ON o.id = oi.order_id
		WHERE
			o.created_at BETWEEN $1 AND $2
			AND o.status != 'cancelled'
	`

	args := []any{params.DateRange.From, params.DateRange.To}
	if params.ShopID != nil {
		query += ` AND oi.shop_id = $3`
		args = append(args, *params.ShopID)
	}

	query += fmt.Sprintf(` GROUP BY oi.product_id, oi.product_name ORDER BY revenue DESC LIMIT %d`, limit)

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query top products failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.TopProduct, error) {
		var item domain.TopProduct
		err := row.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Revenue,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan top products failed: %w", err)
	}

	return items, nil
}

func (r *analyticsRepositoryImpl) GetTopShops(
	ctx context.Context,
	exec transaction.Executor,
	params repository.OrderMetricsParams,
) ([]domain.TopShop, error) {
	limit := params.TopN
	if limit <= 0 {
		limit = 10
	}

	query := fmt.Sprintf(`
		SELECT
			oi.shop_id,
			oi.shop_name,
			COALESCE(SUM(oi.subtotal), 0)
				AS revenue,
			COALESCE(COUNT(DISTINCT oi.order_id), 0)
				AS orders
		FROM
			order_items oi
		JOIN
			orders o ON o.id = oi.order_id
		WHERE
			o.created_at BETWEEN $1 AND $2
			AND o.status != 'cancelled'
		GROUP BY
			oi.shop_id,
			oi.shop_name
		ORDER BY revenue DESC
		LIMIT %d
	`, limit)

	args := []any{params.DateRange.From, params.DateRange.To}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query top shops failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.TopShop, error) {
		var item domain.TopShop
		err := row.Scan(
			&item.ShopID,
			&item.ShopName,
			&item.Revenue,
			&item.Orders,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan top shops failed: %w", err)
	}

	return items, nil
}

func (r *analyticsRepositoryImpl) GetPaymentSummary(
	ctx context.Context,
	exec transaction.Executor,
	params repository.PaymentMetricsParams,
) (*domain.PaymentSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE status = 'paid'), 0)
				AS total_paid,
			COALESCE(SUM(amount) FILTER (WHERE status = 'pending'), 0)
				AS total_pending,
			COALESCE(SUM(amount) FILTER (WHERE status = 'expired'), 0)
				AS total_expired,
			COALESCE(SUM(amount) FILTER (WHERE status = 'refunded'), 0)
				AS total_refunded,
			COALESCE(COUNT(*) FILTER (WHERE status = 'paid'), 0)
				AS paid_count,
			COALESCE(COUNT(*), 0)
				AS total_count,
			COALESCE(
			    AVG(
			        EXTRACT(EPOCH FROM (paid_at - created_at))
			    ) FILTER (WHERE
			        status = 'paid'
			        AND paid_at IS NOT NULL
			    ),
			    0.0
			) AS avg_time_to_pay
		FROM payments
		WHERE created_at BETWEEN $1 AND $2
	`

	var summary domain.PaymentSummary
	var paidCount, totalCount int

	err := exec.QueryRow(ctx, query, params.DateRange.From, params.DateRange.To).Scan(
		&summary.TotalPaid,
		&summary.TotalPending,
		&summary.TotalExpired,
		&summary.TotalRefunded,
		&paidCount,
		&totalCount,
		&summary.AvgTimeToPay,
	)
	if err != nil {
		return nil, fmt.Errorf("query payment summary failed: %w", err)
	}

	if totalCount > 0 {
		summary.PaymentSuccessRate = float64(paidCount) / float64(totalCount)
	}

	return &summary, nil
}

func (r *analyticsRepositoryImpl) GetPaymentMethodBreakdown(
	ctx context.Context,
	exec transaction.Executor,
	params repository.PaymentMetricsParams,
) ([]domain.PaymentMethodBreakdown, error) {
	query := `
		SELECT
			pm.id
				AS method_id,
			pm.name
				AS method_name,
			pm.type
				AS method_type,
			COALESCE(COUNT(p.id), 0)
				AS count,
			COALESCE(SUM(p.amount) FILTER (WHERE p.status = 'paid'), 0)
				AS amount,
			COALESCE(COUNT(p.id) FILTER (WHERE p.status = 'paid'), 0)
				AS paid_count
		FROM
			payment_methods pm
		LEFT JOIN
			payments p ON p.method_id = pm.id
			AND p.created_at BETWEEN $1 AND $2
		GROUP BY pm.id, pm.name, pm.type
		ORDER BY amount DESC
	`

	rows, err := exec.Query(ctx, query, params.DateRange.From, params.DateRange.To)
	if err != nil {
		return nil, fmt.Errorf("query payment method breakdown failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PaymentMethodBreakdown, error) {
		var item domain.PaymentMethodBreakdown
		var paidCount int
		err := row.Scan(
			&item.MethodID,
			&item.MethodName,
			&item.MethodType,
			&item.Count,
			&item.Amount,
			&paidCount,
		)
		if err == nil && item.Count > 0 {
			item.SuccessRate = float64(paidCount) / float64(item.Count)
		}
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan payment method breakdown failed: %w", err)
	}

	return items, nil
}

func (r *analyticsRepositoryImpl) GetShipmentSummary(
	ctx context.Context,
	exec transaction.Executor,
	params repository.ShipmentMetricsParams,
) (*domain.ShipmentSummary, error) {
	query := `
		SELECT
			COALESCE(COUNT(*), 0)
				AS total,
			COALESCE(COUNT(*) FILTER (WHERE s.status = 'delivered'), 0)
				AS delivered,
			COALESCE(COUNT(*) FILTER (WHERE s.status = 'failed'), 0)
				AS failed,
			COALESCE(COUNT(*) FILTER (WHERE s.status = 'returned'), 0)
				AS returned,
			COALESCE(COUNT(*) FILTER (WHERE s.status = 'cancelled'), 0)
				AS cancelled,
			COALESCE(
				AVG(
					EXTRACT(EPOCH FROM (s.created_at - o.created_at))
				) FILTER (WHERE
					s.status = 'delivered'
				), 0.0
			) AS avg_fulfillment_sec
		FROM
			shipments s
		JOIN
			orders o ON o.id = s.order_id
		WHERE
			o.created_at BETWEEN $1 AND $2
	`

	args := []any{params.DateRange.From, params.DateRange.To}

	var summary domain.ShipmentSummary
	err := exec.QueryRow(ctx, query, args...).Scan(
		&summary.Total,
		&summary.Delivered,
		&summary.Failed,
		&summary.Returned,
		&summary.Cancelled,
		&summary.AvgFulfillmentSec,
	)
	if err != nil {
		return nil, fmt.Errorf("query shipment summary failed: %w", err)
	}

	if summary.Total > 0 {
		summary.DeliveryRate = float64(summary.Delivered) / float64(summary.Total)
	}

	return &summary, nil
}

func (r *analyticsRepositoryImpl) GetCourierBreakdown(
	ctx context.Context,
	exec transaction.Executor,
	params repository.ShipmentMetricsParams,
) ([]domain.CourierBreakdown, error) {
	limit := params.TopN
	if limit <= 0 {
		limit = 10
	}

	query := fmt.Sprintf(`
		SELECT
			COALESCE(s.courier_name, '')
				AS courier,
			COALESCE(s.service, '')
				AS service,
			COALESCE(COUNT(*), 0)
				AS count,
			COALESCE(COUNT(*) FILTER (WHERE s.status = 'delivered'), 0)
				AS delivered_count,
			COALESCE(AVG(s.shipping_cost), 0)
				AS avg_cost
		FROM
			shipments s
		JOIN
			orders o ON o.id = s.order_id
		WHERE
			o.created_at BETWEEN $1 AND $2
		GROUP BY
			s.courier_name,
			s.service
		ORDER BY count DESC
		LIMIT %d
	`, limit)

	args := []any{params.DateRange.From, params.DateRange.To}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query courier breakdown failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.CourierBreakdown, error) {
		var item domain.CourierBreakdown
		var deliveredCount int
		var avgCost float64
		err := row.Scan(
			&item.Courier,
			&item.Service,
			&item.Count,
			&deliveredCount,
			&avgCost,
		)
		if err == nil {
			item.AvgCost = int64(avgCost)
			if item.Count > 0 {
				item.DeliveryRate = float64(deliveredCount) / float64(item.Count)
			}
		}
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan courier breakdown failed: %w", err)
	}

	return items, nil
}

func (r *analyticsRepositoryImpl) GetInventorySummary(
	ctx context.Context,
	exec transaction.Executor,
	params repository.InventoryMetricsParams,
) (*domain.InventorySummary, error) {
	threshold := params.LowStockThreshold
	if threshold <= 0 {
		threshold = 5
	}

	query := `
		SELECT
			COALESCE(COUNT(DISTINCT product_id), 0)
				AS total_products,
			COALESCE(SUM(stock), 0)
				AS total_stock,
			COALESCE(SUM(reserved_stock), 0)
				AS total_reserved,
			COALESCE(SUM(stock - reserved_stock), 0)
				AS total_available,
			COALESCE(COUNT(*) FILTER (WHERE stock - reserved_stock <= 0), 0)
				AS stockout_count,
			COALESCE(COUNT(*) FILTER (WHERE
		        stock - reserved_stock <= $1
		        AND stock - reserved_stock > 0
			)) AS low_stock_count
		FROM inventory
	`

	args := []any{threshold}
	if params.ShopID != nil {
		query += ` WHERE shop_id = $2`
		args = append(args, *params.ShopID)
	}

	var summary domain.InventorySummary
	err := exec.QueryRow(ctx, query, args...).Scan(
		&summary.TotalProducts,
		&summary.TotalStock,
		&summary.TotalReserved,
		&summary.TotalAvailable,
		&summary.StockoutCount,
		&summary.LowStockCount,
	)
	if err != nil {
		return nil, fmt.Errorf("query inventory summary failed: %w", err)
	}

	return &summary, nil
}

func (r *analyticsRepositoryImpl) GetProductMetricsSummary(
	ctx context.Context,
	exec transaction.Executor,
	params repository.ProductMetricsParams,
) (*domain.ProductMetricsSummary, error) {
	limit := params.TopN
	if limit <= 0 {
		limit = 10
	}

	// 1. Calculate invoice void rate
	invoiceQuery := `
		SELECT
			COALESCE(COUNT(*) FILTER (WHERE status = 'void'), 0)
				AS void_count,
			COALESCE(COUNT(*), 0)
				AS total_count
		FROM invoices
		WHERE created_at BETWEEN $1 AND $2
	`
	var voidCount, totalInvoices int
	err := exec.QueryRow(ctx, invoiceQuery, params.DateRange.From, params.DateRange.To).Scan(&voidCount, &totalInvoices)
	if err != nil {
		return nil, fmt.Errorf("query invoice void rate failed: %w", err)
	}

	var invoiceVoidRate float64
	if totalInvoices > 0 {
		invoiceVoidRate = float64(voidCount) / float64(totalInvoices)
	}

	// 2. Fetch top products stat
	productQuery := fmt.Sprintf(`
		SELECT
			p.id
				AS product_id,
			p.name
				AS product_name,
			COALESCE(SUM(oi.subtotal), 0)
				AS revenue,
			COALESCE(SUM(oi.quantity), 0)
				AS units_sold,
			pp.gross_margin_pct
				AS gross_margin_pct,
			COALESCE(SUM(oi.quantity) FILTER (WHERE o.created_at >= NOW() - INTERVAL '7 days'), 0)
				AS velocity_7d,
			COALESCE(SUM(oi.quantity) FILTER (WHERE o.created_at >= NOW() - INTERVAL '30 days'), 0)
				AS velocity_30d
		FROM
			products p
		LEFT JOIN
			order_items oi ON oi.product_id = p.id
		LEFT JOIN
			orders o ON o.id = oi.order_id
			AND o.created_at BETWEEN $1 AND $2
			AND o.status != 'cancelled'
		LEFT JOIN
			product_performance pp ON pp.product_id = p.id
		WHERE p.deleted_at IS NULL
		GROUP BY
			p.id,
			p.name,
			pp.gross_margin_pct
		ORDER BY revenue DESC
		LIMIT %d
	`, limit)

	rows, err := exec.Query(ctx, productQuery, params.DateRange.From, params.DateRange.To)
	if err != nil {
		return nil, fmt.Errorf("query product metrics failed: %w", err)
	}
	defer rows.Close()

	topByRevenue, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.TopProductStat, error) {
		var item domain.TopProductStat
		var marginPct *float64
		err := row.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Revenue,
			&item.UnitsSold,
			&marginPct,
			&item.SalesVelocity7d,
			&item.SalesVelocity30d,
		)
		if err == nil {
			item.GrossMarginPct = marginPct
		}
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan product metrics failed: %w", err)
	}

	// Fetch top by volume
	productVolumeQuery := fmt.Sprintf(`
		SELECT
			p.id
				AS product_id,
			p.name
				AS product_name,
			COALESCE(SUM(oi.subtotal), 0)
				AS revenue,
			COALESCE(SUM(oi.quantity), 0)
				AS units_sold,
			pp.gross_margin_pct
				AS gross_margin_pct,
			COALESCE(SUM(oi.quantity) FILTER (WHERE o.created_at >= NOW() - INTERVAL '7 days'), 0)
				AS velocity_7d,
			COALESCE(SUM(oi.quantity) FILTER (WHERE o.created_at >= NOW() - INTERVAL '30 days'), 0)
				AS velocity_30d
		FROM
			products p
		LEFT JOIN
			order_items oi ON oi.product_id = p.id
		LEFT JOIN
			orders o ON o.id = oi.order_id
			AND o.created_at BETWEEN $1 AND $2
			AND o.status != 'cancelled'
		LEFT JOIN
			product_performance pp ON pp.product_id = p.id
		WHERE p.deleted_at IS NULL
		GROUP BY
			p.id,
			p.name,
			pp.gross_margin_pct
		ORDER BY units_sold DESC
		LIMIT %d
	`, limit)

	rowsVol, err := exec.Query(ctx, productVolumeQuery, params.DateRange.From, params.DateRange.To)
	if err != nil {
		return nil, fmt.Errorf("query product volume metrics failed: %w", err)
	}
	defer rowsVol.Close()

	topByVolume, err := pgx.CollectRows(rowsVol, func(row pgx.CollectableRow) (domain.TopProductStat, error) {
		var item domain.TopProductStat
		var marginPct *float64
		err := row.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Revenue,
			&item.UnitsSold,
			&marginPct,
			&item.SalesVelocity7d,
			&item.SalesVelocity30d,
		)
		if err == nil {
			item.GrossMarginPct = marginPct
		}
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan product volume metrics failed: %w", err)
	}

	pMS := domain.ProductMetricsSummary{
		TopByRevenue:    topByRevenue,
		TopByVolume:     topByVolume,
		InvoiceVoidRate: invoiceVoidRate,
	}

	return &pMS, nil
}
