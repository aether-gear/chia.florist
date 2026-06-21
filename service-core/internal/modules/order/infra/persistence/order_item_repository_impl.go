package persistence

import (
	"context"
	"fmt"

	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type orderItemRepositoryImpl struct{}

func NewOrderItemRepositoryImpl() repository.OrderItemRepository {
	return &orderItemRepositoryImpl{}
}

func (r *orderItemRepositoryImpl) ListByOrderID(
	ctx context.Context,
	exec transaction.Executor,
	orderID uuid.UUID,
) ([]domain.OrderItem, error) {
	query := `
		SELECT
			id,
			order_id,
			shop_id,
			shop_name,
			product_id,
			product_name,
			quantity,
			unit_price,
			subtotal,
			courier_code,
			courier_service,
			shipping_fee_total
		FROM
			order_items
		WHERE
			order_id = $1
	`

	rows, err := exec.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order items by order id failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.OrderItem, error) {
		var item domain.OrderItem
		err := row.Scan(
			&item.ID,
			&item.OrderID,
			&item.ShopID,
			&item.ShopName,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.CourierCode,
			&item.CourierService,
			&item.ShippingFee,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan order items failed: %w", err)
	}

	return items, nil
}

func (r *orderItemRepositoryImpl) SaveBulk(
	ctx context.Context,
	exec transaction.Executor,
	items []domain.OrderItem,
) error {
	query := `
		INSERT INTO order_items (
			id,
			order_id,
			shop_id,
			shop_name,
			product_id,
			product_name,
			quantity,
			unit_price,
			subtotal,
			courier_code,
			courier_service,
			shipping_fee_total
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id)
		DO UPDATE SET
			order_id = EXCLUDED.order_id,
			shop_id = EXCLUDED.shop_id,
			shop_name = EXCLUDED.shop_name,
			product_id = EXCLUDED.product_id,
			product_name = EXCLUDED.product_name,
			quantity = EXCLUDED.quantity,
			unit_price = EXCLUDED.unit_price,
			subtotal = EXCLUDED.subtotal,
			courier_code = EXCLUDED.courier_code,
			courier_service = EXCLUDED.courier_service,
			shipping_fee_total = EXCLUDED.shipping_fee_total
	`

	for _, item := range items {
		_, err := exec.Exec(ctx, query,
			item.ID,
			item.OrderID,
			item.ShopID,
			item.ShopName,
			item.ProductID,
			item.ProductName,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
			item.CourierCode,
			item.CourierService,
			item.ShippingFee,
		)
		if err != nil {
			return fmt.Errorf("query to save order item: %w", err)
		}
	}

	return nil
}

func (r *orderItemRepositoryImpl) ListByOrderIDs(
	ctx context.Context,
	exec transaction.Executor,
	orderIDs []uuid.UUID,
) ([]domain.OrderItem, error) {
	if len(orderIDs) == 0 {
		return []domain.OrderItem{}, nil
	}

	query := `
		SELECT
			id,
			order_id,
			shop_id,
			shop_name,
			product_id,
			product_name,
			quantity,
			unit_price,
			subtotal,
			courier_code,
			courier_service,
			shipping_fee_total
		FROM
			order_items
		WHERE
			order_id = ANY($1::uuid[])
	`

	orderIDStrings := make([]string, len(orderIDs))
	for i, id := range orderIDs {
		orderIDStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, orderIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query order items by order ids failed: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.OrderItem, error) {
		var item domain.OrderItem
		err := row.Scan(
			&item.ID,
			&item.OrderID,
			&item.ShopID,
			&item.ShopName,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.CourierCode,
			&item.CourierService,
			&item.ShippingFee,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan order items failed: %w", err)
	}

	return items, nil
}

