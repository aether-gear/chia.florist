package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/shipment/domain"
	"service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type shipmentRepositoryImpl struct{}

func NewShipmentRepositoryImpl() repository.ShipmentRepository {
	return &shipmentRepositoryImpl{}
}

const shipmentSelectCols = `
	id,
	order_id,
	status,
	tracking_number,
	courier_name,
	service,
	shipping_cost,
	weight,
	origin_id,
	destination_id,
	created_at
`

func (r *shipmentRepositoryImpl) scanShipment(row pgx.Row) (*domain.Shipment, error) {
	var s domain.Shipment
	err := row.Scan(
		&s.ID,
		&s.OrderID,
		&s.Status,
		&s.TrackingNumber,
		&s.Courier,
		&s.Service,
		&s.Cost,
		&s.Weight,
		&s.OriginID,
		&s.DestinationID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan shipment failed: %w", err)
	}
	return &s, nil
}

func (r *shipmentRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.Shipment, error) {
	query := `
		SELECT id,
			order_id,
			status,
			tracking_number,
			courier_name,
			service,
			shipping_cost,
			weight,
			origin_id,
			destination_id,
			created_at
		FROM shipments
		WHERE id = $1
	`

	s, err := r.scanShipment(exec.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shipment by id failed: %w", err)
	}
	return s, nil
}

func (r *shipmentRepositoryImpl) GetByOrderID(
	ctx context.Context,
	exec transaction.Executor,
	orderID uuid.UUID,
) (*domain.Shipment, error) {
	query := `
		SELECT
			id,
			order_id,
			status,
			tracking_number,
			courier_name,
			service,
			shipping_cost,
			weight,
			origin_id,
			destination_id,
			created_at
		FROM shipments
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	s, err := r.scanShipment(exec.QueryRow(ctx, query, orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shipment by order id failed: %w", err)
	}
	return s, nil
}

func (r *shipmentRepositoryImpl) ListByOrderIDs(
	ctx context.Context,
	exec transaction.Executor,
	orderIDs []uuid.UUID,
) ([]domain.Shipment, error) {
	if len(orderIDs) == 0 {
		return []domain.Shipment{}, nil
	}

	query := `
		SELECT DISTINCT ON (order_id)
			id,
			order_id,
			status,
			tracking_number,
			courier_name,
			service,
			shipping_cost,
			weight,
			origin_id,
			destination_id,
			created_at
		FROM shipments
		WHERE order_id = ANY($1::uuid[])
		ORDER BY order_id, created_at DESC
	`

	orderIDStrings := make([]string, len(orderIDs))
	for i, id := range orderIDs {
		orderIDStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, orderIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query shipments by order ids failed: %w", err)
	}
	defer rows.Close()

	shipments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Shipment, error) {
		var s domain.Shipment
		err := row.Scan(
			&s.ID,
			&s.OrderID,
			&s.Status,
			&s.TrackingNumber,
			&s.Courier,
			&s.Service,
			&s.Cost,
			&s.Weight,
			&s.OriginID,
			&s.DestinationID,
			&s.CreatedAt,
		)
		return s, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan shipments failed: %w", err)
	}

	return shipments, nil
}

func (r *shipmentRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	shipment domain.Shipment,
) error {
	query := `
		INSERT INTO shipments (
			id,
			order_id,
			status,
			tracking_number,
			courier_name,
			service,
			shipping_cost,
			weight,
			origin_id,
			destination_id,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := exec.Exec(ctx, query,
		shipment.ID,
		shipment.OrderID,
		shipment.Status,
		shipment.TrackingNumber,
		shipment.Courier,
		shipment.Service,
		shipment.Cost,
		shipment.Weight,
		shipment.OriginID,
		shipment.DestinationID,
		shipment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create shipment failed: %w", err)
	}
	return nil
}

func (r *shipmentRepositoryImpl) Update(
	ctx context.Context,
	exec transaction.Executor,
	shipment domain.Shipment,
) error {
	query := `
		UPDATE shipments
		SET
			status          = $2,
			tracking_number = $3,
			courier_name    = $4,
			service         = $5,
			shipping_cost   = $6,
			weight          = $7
		WHERE id = $1
	`

	_, err := exec.Exec(ctx, query,
		shipment.ID,
		shipment.Status,
		shipment.TrackingNumber,
		shipment.Courier,
		shipment.Service,
		shipment.Cost,
		shipment.Weight,
	)
	if err != nil {
		return fmt.Errorf("update shipment failed: %w", err)
	}
	return nil
}
