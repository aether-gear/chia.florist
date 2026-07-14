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

type shipmentEventRepositoryImpl struct{}

func NewShipmentEventRepositoryImpl() repository.ShipmentEventRepository {
	return &shipmentEventRepositoryImpl{}
}

const shipmentEventSelectCols = `
	id,
	shipment_id,
	status,
	description,
	location,
	timestamp
`

func (r *shipmentEventRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.ShipmentEvent, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM shipment_events
		WHERE id = $1
	`, shipmentEventSelectCols)

	e, err := r.scanShipmentEvent(exec.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shipment event by id failed: %w", err)
	}
	return e, nil
}

func (r *shipmentEventRepositoryImpl) ListByShipmentID(
	ctx context.Context,
	exec transaction.Executor,
	shipmentID uuid.UUID,
) ([]domain.ShipmentEvent, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM shipment_events
		WHERE shipment_id = $1
		ORDER BY timestamp ASC
	`, shipmentEventSelectCols)

	rows, err := exec.Query(ctx, query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("query shipment events by shipment id failed: %w", err)
	}
	defer rows.Close()

	var events []domain.ShipmentEvent
	for rows.Next() {
		e, err := r.scanShipmentEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return events, nil
}

func (r *shipmentEventRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	event domain.ShipmentEvent,
) error {
	query := `
		INSERT INTO shipment_events (
			id,
			shipment_id,
			status,
			description,
			location,
			timestamp
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := exec.Exec(ctx, query,
		event.ID,
		event.ShipmentID,
		event.Status,
		event.Description,
		event.Location,
		event.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("create shipment event failed: %w", err)
	}

	return nil
}

func (r *shipmentEventRepositoryImpl) scanShipmentEvent(row pgx.Row) (*domain.ShipmentEvent, error) {
	var e domain.ShipmentEvent
	err := row.Scan(
		&e.ID,
		&e.ShipmentID,
		&e.Status,
		&e.Description,
		&e.Location,
		&e.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("scan shipment event failed: %w", err)
	}

	return &e, nil
}
