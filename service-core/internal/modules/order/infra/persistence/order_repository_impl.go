package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/repository"
	query "service-core/internal/shared/query"
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
			id,
			number,
			user_id,
			address_id,
			status,
			subtotal,
			shipping_fee,
			total,
			created_at,
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
			id,
			number,
			user_id,
			address_id,
			status,
			subtotal,
			shipping_fee,
			total,
			created_at,
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

func (r *orderRepositoryImpl) FindOrders(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindOrderParams,
) ([]domain.Order, int, error) {
	baseQuery := `
		FROM orders o
	`

	selectQuery := `
		SELECT
			o.id,
			o.number,
			o.user_id,
			o.address_id,
			o.status,
			o.subtotal,
			o.shipping_fee,
			o.total,
			o.created_at,
			o.updated_at
	`

	whereClause := ""
	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("o.id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Number != nil {
		conditions = append(conditions, fmt.Sprintf("o.number ILIKE $%d", argPos))
		args = append(args, "%"+*params.Number+"%")
		argPos++
	}

	if params.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("o.user_id = $%d", argPos))
		args = append(args, *params.UserID)
		argPos++
	}

	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("o.status = $%d", argPos))
		args = append(args, *params.Status)
		argPos++
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count matching orders
	countQuery := `
		SELECT COUNT(DISTINCT o.id)
	` + baseQuery + whereClause

	countArgs := append([]any{}, args...)
	var total int
	err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count orders failed: %w", err)
	}

	// Build sorting expressions
	var orderSortKeys = map[query.SortKey]string{
		repository.OrderSortLatest: "o.created_at",
		repository.OrderSortNumber: "o.number",
		repository.OrderSortTotal:  "o.total",
		repository.OrderSortStatus: "o.status",
		repository.OrderSortModify: "o.updated_at",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := orderSortKeys[sort.By]
		if !exists {
			continue
		}

		direction := "DESC"
		if sort.Direction == query.SortAsc {
			direction = "ASC"
		}

		sortClauses = append(
			sortClauses,
			fmt.Sprintf("%s %s", colName, direction),
		)
	}

	orderBy := "ORDER BY o.created_at DESC"
	if len(sortClauses) > 0 {
		orderBy = "ORDER BY " + strings.Join(sortClauses, ", ")
	}

	// Apply pagination
	limitPos := argPos
	offsetPos := argPos + 1

	limit := params.Pagination.Limit
	if limit <= 0 {
		limit = 10
	}

	page := params.Pagination.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	// The execution
	queryStr := selectQuery +
		baseQuery +
		whereClause +
		" " +
		orderBy +
		fmt.Sprintf(
			" LIMIT $%d OFFSET $%d",
			limitPos,
			offsetPos,
		)

	rows, err := exec.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query orders failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Order
	for rows.Next() {
		var item domain.Order
		err := rows.Scan(
			&item.ID,
			&item.Number,
			&item.UserID,
			&item.AddressID,
			&item.Status,
			&item.Subtotal,
			&item.ShippingFee,
			&item.Total,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("mapping order model to domain failed: %w", err)
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate orders failed: %w", err)
	}

	return results, total, nil
}

