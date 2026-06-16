package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type shopRepositoryImpl struct{}

func NewShopRepositoryImpl() repository.ShopRepository {
	return &shopRepositoryImpl{}
}

func (r *shopRepositoryImpl) List(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindShopsParams,
) ([]domain.Shop, int, error) {
	baseQuery := `
		FROM shops s
	`

	selectQuery := `
		SELECT
			s.id,
			s.name,
			s.slug,
			s.description,
			s.is_active,
			s.created_at,
			s.updated_at
	`

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	conditions = append(conditions, "s.deleted_at IS NULL")

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("s.id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("s.name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Name+"%")
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countArgs := append([]any{}, args...)
	countQuery := "SELECT COUNT(DISTINCT s.id) " + baseQuery + whereClause

	var total int
	err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count shops failed: %w", err)
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	limitPos := argPos
	offsetPos := argPos + 1

	query := selectQuery + baseQuery + whereClause +
		fmt.Sprintf(" ORDER BY s.created_at DESC LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query list shops failed: %w", err)
	}
	defer rows.Close()

	var shops []domain.Shop
	for rows.Next() {
		var s domain.Shop
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Slug,
			&s.Description,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("mapping shop model to domain failed: %w", err)
		}

		shops = append(shops, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate shops failed: %w", err)
	}

	return shops, total, nil
}

func (r *shopRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	shopID uuid.UUID,
) (*domain.Shop, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			is_active,
			created_at,
			updated_at
		FROM shops
		WHERE id = $1
		LIMIT 1
	`

	var s domain.Shop
	err := exec.QueryRow(ctx, query, shopID).Scan(
		&s.ID,
		&s.Name,
		&s.Slug,
		&s.Description,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query shop by id failed: %w", err)
	}

	return &s, nil
}

func (r *shopRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	shop domain.Shop,
) error {
	query := `
		INSERT INTO shops (
			id,
			name,
			slug,
			description,
			is_active,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := exec.Exec(ctx, query,
		shop.ID,
		shop.Name,
		shop.Slug,
		shop.Description,
		shop.IsActive,
		shop.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert shop failed: %w", err)
	}
	return nil
}
