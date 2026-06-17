package persistence

import (
	"context"
	"fmt"
	"strings"

	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

type merchantRepositoryImpl struct{}

func NewMerchantRepositoryImpl() repository.MerchantRepository {
	return &merchantRepositoryImpl{}
}

func (r *merchantRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	merchant domain.Merchant,
) error {
	query := `
		INSERT INTO merchants (
			id,
			name,
			description,
			logo_url,
			banner_url,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := exec.Exec(ctx, query,
		merchant.ID,
		merchant.Name,
		merchant.Description,
		merchant.LogoUrl,
		merchant.BannerUrl,
		merchant.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert merchant failed: %w", err)
	}
	return nil
}

func (r *merchantRepositoryImpl) FindMerchants(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindMerchantParams,
) ([]domain.Merchant, int, error) {
	baseQuery := `
		FROM merchants m
	`

	selectQuery := `
		SELECT
			m.id,
			m.name,
			m.description,
			m.logo_url,
			m.banner_url,
			m.created_at,
			m.updated_at,
			m.deleted_at
	`

	// Build filters
	// Apply search criteria and soft-delete constraints
	whereClause := ""
	notDeletedCondition := "m.deleted_at IS NULL"

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	conditions = append(conditions, notDeletedCondition)

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("m.id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("m.name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Name+"%")
		argPos++
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count matching products
	// Used for pagination metadata
	countQuery := `
		SELECT COUNT(*)
	` + baseQuery + whereClause

	countArgs := append([]any{}, args...)

	var total int
	err := exec.
		QueryRow(ctx, countQuery, countArgs...).
		Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count merchants failed: %w", err)
	}

	// Build sorting expressions
	// Convert requested sort keys into SQL ORDER BY clauses
	var merchantSortKeys = map[query.SortKey]string{
		repository.MerchantSortLatest: "m.created_at",
		repository.MerchantSortName:   "m.name",
		repository.MerchantSortModify: "m.updated_at",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := merchantSortKeys[sort.By]
		if !exists {
			continue
		}

		dir := "DESC"
		if sort.Direction == query.SortAsc {
			dir = "ASC"
		}

		sortClauses = append(
			sortClauses,
			fmt.Sprintf("%s %s", colName, dir),
		)
	}

	orderBy := "ORDER BY m.created_at DESC"
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
	queryStr := selectQuery + baseQuery + whereClause + " " + orderBy +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	rows, err := exec.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query list merchants failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Merchant
	for rows.Next() {
		var m domain.Merchant
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.LogoUrl,
			&m.BannerUrl,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("mapping merchant model to domain failed: %w", err)
		}
		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate merchants failed: %w", err)
	}

	return results, total, nil
}
