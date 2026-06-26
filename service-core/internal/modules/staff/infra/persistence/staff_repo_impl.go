package persistence

import (
	"context"
	"fmt"
	"strings"

	"service-core/internal/modules/staff/domain"
	"service-core/internal/modules/staff/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

type staffRepositoryImpl struct{}

func NewStaffRepositoryImpl() repository.StaffRepository {
	return &staffRepositoryImpl{}
}

func (r *staffRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	staff domain.Staff,
) error {
	query := `
		INSERT INTO staff (
			id,
			user_id,
			created_at
		) VALUES ($1,$2,$3)
	`

	_, err := exec.Exec(ctx, query,
		staff.ID,
		staff.UserID,
		staff.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert staff failed: %w", err)
	}
	return nil
}

func (r *staffRepositoryImpl) FindStaff(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindStaffParams,
) ([]domain.Staff, int, error) {
	baseQuery := `
		FROM staff m
	`

	selectQuery := `
		SELECT
			m.id,
			m.user_id,
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

	// if params.Name != nil {
	// 	conditions = append(conditions, fmt.Sprintf("m.name ILIKE $%d", argPos))
	// 	args = append(args, "%"+*params.Name+"%")
	// 	argPos++
	// }

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
		return nil, 0, fmt.Errorf("query count staff failed: %w", err)
	}

	// Build sorting expressions
	// Convert requested sort keys into SQL ORDER BY clauses
	var staffSortKeys = map[query.SortKey]string{
		repository.StaffSortLatest: "m.created_at",
		repository.StaffSortModify: "m.updated_at",
		// repository.StaffSortName:   "m.name",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := staffSortKeys[sort.By]
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
		return nil, 0, fmt.Errorf("query list staff failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Staff
	for rows.Next() {
		var m domain.Staff
		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("mapping staff model to domain failed: %w", err)
		}
		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate staff failed: %w", err)
	}

	return results, total, nil
}
