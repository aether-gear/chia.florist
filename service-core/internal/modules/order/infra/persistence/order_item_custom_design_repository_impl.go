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

type orderItemCustomDesignRepositoryImpl struct{}

func NewOrderItemCustomDesignRepositoryImpl() repository.OrderItemCustomDesignRepository {
	return &orderItemCustomDesignRepositoryImpl{}
}

func (r *orderItemCustomDesignRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	design domain.OrderItemCustomDesign,
) error {
	query := `
		INSERT INTO order_item_custom_designs (
			id,
			order_item_id,
			version,
			physical_size_id,
			preview_url,
			header_text_upper,
			body_text_upper,
			header_text_lower,
			body_text_lower,
			design_snapshot,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_item_id)
		DO UPDATE SET
			version = EXCLUDED.version,
			physical_size_id = EXCLUDED.physical_size_id,
			preview_url = EXCLUDED.preview_url,
			header_text_upper = EXCLUDED.header_text_upper,
			body_text_upper = EXCLUDED.body_text_upper,
			header_text_lower = EXCLUDED.header_text_lower,
			body_text_lower = EXCLUDED.body_text_lower,
			design_snapshot = EXCLUDED.design_snapshot
	`

	_, err := exec.Exec(
		ctx,
		query,
		design.ID,
		design.OrderItemID,
		design.Version,
		design.PhysicalSizeID,
		design.PreviewURL,
		design.HeaderTextUpper,
		design.BodyTextUpper,
		design.HeaderTextLower,
		design.BodyTextLower,
		design.DesignSnapshot,
		design.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save order item custom design: %w", err)
	}

	return nil
}

func (r *orderItemCustomDesignRepositoryImpl) SaveBulk(
	ctx context.Context,
	exec transaction.Executor,
	designs []domain.OrderItemCustomDesign,
) error {
	if len(designs) == 0 {
		return nil
	}

	for _, d := range designs {
		if err := r.Save(ctx, exec, d); err != nil {
			return err
		}
	}

	return nil
}

func (r *orderItemCustomDesignRepositoryImpl) GetByOrderItemID(
	ctx context.Context,
	exec transaction.Executor,
	orderItemID uuid.UUID,
) (*domain.OrderItemCustomDesign, error) {
	query := `
		SELECT
			id,
			order_item_id,
			version,
			physical_size_id,
			preview_url,
			header_text_upper,
			body_text_upper,
			header_text_lower,
			body_text_lower,
			design_snapshot,
			created_at
		FROM
			order_item_custom_designs
		WHERE
			order_item_id = $1
	`

	row := exec.QueryRow(ctx, query, orderItemID)

	var design domain.OrderItemCustomDesign
	err := row.Scan(
		&design.ID,
		&design.OrderItemID,
		&design.Version,
		&design.PhysicalSizeID,
		&design.PreviewURL,
		&design.HeaderTextUpper,
		&design.BodyTextUpper,
		&design.HeaderTextLower,
		&design.BodyTextLower,
		&design.DesignSnapshot,
		&design.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan order item custom design: %w", err)
	}

	return &design, nil
}

func (r *orderItemCustomDesignRepositoryImpl) ListByOrderItemIDs(
	ctx context.Context,
	exec transaction.Executor,
	orderItemIDs []uuid.UUID,
) (map[uuid.UUID]domain.OrderItemCustomDesign, error) {
	if len(orderItemIDs) == 0 {
		return map[uuid.UUID]domain.OrderItemCustomDesign{}, nil
	}

	query := `
		SELECT
			id,
			order_item_id,
			version,
			physical_size_id,
			preview_url,
			header_text_upper,
			body_text_upper,
			header_text_lower,
			body_text_lower,
			design_snapshot,
			created_at
		FROM
			order_item_custom_designs
		WHERE
			order_item_id = ANY($1::uuid[])
	`

	idStrings := make([]string, len(orderItemIDs))
	for i, id := range orderItemIDs {
		idStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, idStrings)
	if err != nil {
		return nil, fmt.Errorf("query custom designs by order item ids failed: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]domain.OrderItemCustomDesign, len(orderItemIDs))
	for rows.Next() {
		var design domain.OrderItemCustomDesign
		err := rows.Scan(
			&design.ID,
			&design.OrderItemID,
			&design.Version,
			&design.PhysicalSizeID,
			&design.PreviewURL,
			&design.HeaderTextUpper,
			&design.BodyTextUpper,
			&design.HeaderTextLower,
			&design.BodyTextLower,
			&design.DesignSnapshot,
			&design.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan custom design failed: %w", err)
		}
		result[design.OrderItemID] = design
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func (r *orderItemCustomDesignRepositoryImpl) ListByOrderID(
	ctx context.Context,
	exec transaction.Executor,
	orderID uuid.UUID,
) ([]domain.OrderItemCustomDesign, error) {
	query := `
		SELECT
			cd.id,
			cd.order_item_id,
			cd.version,
			cd.physical_size_id,
			cd.preview_url,
			cd.header_text_upper,
			cd.body_text_upper,
			cd.header_text_lower,
			cd.body_text_lower,
			cd.design_snapshot,
			cd.created_at
		FROM
			order_item_custom_designs cd
		JOIN
			order_items oi ON oi.id = cd.order_item_id
		WHERE
			oi.order_id = $1
	`

	rows, err := exec.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query custom designs by order id failed: %w", err)
	}
	defer rows.Close()

	designs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.OrderItemCustomDesign, error) {
		var design domain.OrderItemCustomDesign
		err := row.Scan(
			&design.ID,
			&design.OrderItemID,
			&design.Version,
			&design.PhysicalSizeID,
			&design.PreviewURL,
			&design.HeaderTextUpper,
			&design.BodyTextUpper,
			&design.HeaderTextLower,
			&design.BodyTextLower,
			&design.DesignSnapshot,
			&design.CreatedAt,
		)
		return design, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan custom designs failed: %w", err)
	}

	return designs, nil
}
