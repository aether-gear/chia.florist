package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type orderRepositoryImpl struct{}

func NewOrderRepositoryImpl() repository.OrderRepository {
	return &orderRepositoryImpl{}
}

func (r *orderRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.Order, error) {
	query := `
		SELECT
			id
			number
			user_id
			address_id
			status
			subtotal
			shipping_fee
			total
			created_at
			updated_at
		FROM
			orders
		WHERE
			id = $1
		LIMIT 1
	`

	var order domain.Order
	err := exec.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.Number,
		&order.UserID,
		&order.AddressID,
		&order.Status,
		&order.Subtotal,
		&order.ShippingFee,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query order by id failed: %w", err)
	}

	return &order, nil
}

func (r *orderRepositoryImpl) GetByNumber(
	ctx context.Context,
	exec transaction.Executor,
	number string,
) (*domain.Order, error) {
	query := `
		SELECT
			id
			number
			user_id
			address_id
			status
			subtotal
			shipping_fee
			total
			created_at
			updated_at
		FROM
			orders
		WHERE
			number = $1
		LIMIT 1
	`

	var order domain.Order
	err := exec.QueryRow(ctx, query, number).Scan(
		&order.ID,
		&order.Number,
		&order.UserID,
		&order.AddressID,
		&order.Status,
		&order.Subtotal,
		&order.ShippingFee,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query order by number failed: %w", err)
	}

	return &order, nil
}

func (r *orderRepositoryImpl) UpdateStatus(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
	status domain.OrderStatus,
) error {
	query := `
		UPDATE orders
		SET
			status = $2
		WHERE
			id = $1
	`

	_, err := exec.Exec(ctx, query,
		id,
		status,
	)
	if err != nil {
		return fmt.Errorf("query to update status: %w", err)
	}

	return nil
}

func (r *orderRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	order domain.Order,
) error {
	query := `
		INSERT INTO orders (
			id,
			number,
			user_id,
			address_id,
			status,
			subtotal,
			shipping_fee,
			total,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id)
		DO UPDATE SET
			number = EXCLUDED.number,
			user_id = EXCLUDED.user_id,
			address_id = EXCLUDED.address_id,
			status = EXCLUDED.status,
			subtotal = EXCLUDED.subtotal,
			shipping_fee = EXCLUDED.shipping_fee,
			total = EXCLUDED.total
	`

	_, err := exec.Exec(ctx, query,
		order.ID,
		order.Number,
		order.UserID,
		order.AddressID,
		order.Status,
		order.Subtotal,
		order.ShippingFee,
		order.Total,
		order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("query to save order: %w", err)
	}

	return nil
}
