package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type shopRepositoryImpl struct{}

func NewShopRepositoryImpl() repository.ShopRepository {
	return &shopRepositoryImpl{}
}

func (r *shopRepositoryImpl) FindByParams(
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

	// Build filters
	// Apply search criteria and soft-delete constraints
	whereClause := ""
	notDeletedCondition := "s.deleted_at IS NULL"

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	conditions = append(conditions, notDeletedCondition)

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

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count matching products
	// Used for pagination metadata
	countQuery := `
		SELECT COUNT(DISTINCT s.id)
	` + baseQuery + whereClause

	countArgs := append([]any{}, args...)

	var total int
	err := exec.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count shops failed: %w", err)
	}

	// Build sorting expressions
	// Convert requested sort keys into SQL ORDER BY clauses
	var shopSortKeys = map[query.SortKey]string{
		repository.ShopSortName:   "name",
		repository.ShopSortActive: "is_active",
		repository.ShopSortLatest: "created_at",
		repository.ShopSortModify: "updated_at",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := shopSortKeys[sort.By]
		if !exists {
			continue
		}

		dir := "DESC"
		if sort.Direction == query.SortAsc {
			dir = "ASC"
		}

		sortClauses = append(
			sortClauses,
			fmt.Sprintf("s.%s %s", colName, dir),
		)
	}

	orderBy := "ORDER BY s.created_at DESC"
	if len(sortClauses) > 0 {
		orderBy = "ORDER BY " + strings.Join(sortClauses, ", ")
	}

	// Apply pagination
	// Calculate limit and offset values
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
	query := selectQuery + baseQuery + whereClause + " " + orderBy +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)

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

func (r *shopRepositoryImpl) FindByIDs(
	ctx context.Context,
	exec transaction.Executor,
	ids []uuid.UUID,
) ([]domain.Shop, error) {
	if len(ids) == 0 {
		return []domain.Shop{}, nil
	}

	query := `
		SELECT
			s.id,
			s.name,
			s.slug,
			s.description,
			s.is_active,
			s.created_at,
			s.updated_at
		FROM shops s
		WHERE s.id = ANY($1::uuid[])
	`

	shopIDStrings := make([]string, len(ids))
	for i, id := range ids {
		shopIDStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, shopIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query shops by many ids failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Shop
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
			return nil, fmt.Errorf("mapping shop model to domain failed: %w", err)
		}

		results = append(results, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shops failed: %w", err)
	}

	return results, nil
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
