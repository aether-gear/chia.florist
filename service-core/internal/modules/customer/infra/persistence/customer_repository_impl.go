package persistence

import (
	"context"
	"fmt"
	"strings"

	"service-core/internal/modules/customer/domain"
	"service-core/internal/modules/customer/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

type customerRepositoryImpl struct{}

func NewCustomerRepositoryImpl() repository.CustomerRepository {
	return &customerRepositoryImpl{}
}

func (r *customerRepositoryImpl) FindCustomers(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindCustomerParams,
) ([]domain.Customer, int, error) {
	baseQuery := `
		FROM users u
		LEFT JOIN accounts a ON a.user_id = u.id
	`

	selectQuery := `
		SELECT 
			u.id, 
			u.name, 
			u.username, 
			u.phone,
			u.created_at, 
			u.updated_at, 
			u.deleted_at,
			a.last_login_at
	`

	// Build filters
	// Apply search criteria and soft-delete constraints
	whereClause := ""
	notDeletedCondition := "u.deleted_at IS NULL"
	onlyCustomerCondition := "a.type = 'customer'"

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	conditions = append(conditions,
		notDeletedCondition,
		onlyCustomerCondition,
	)

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("u.id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("u.name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Name+"%")
		argPos++
	}

	if params.Username != nil {
		conditions = append(conditions, fmt.Sprintf("u.username ILIKE $%d", argPos))
		args = append(args, "%"+*params.Username+"%")
		argPos++
	}

	if params.Email != nil {
		conditions = append(conditions, fmt.Sprintf("a.email ILIKE $%d", argPos))
		args = append(args, "%"+*params.Email+"%")
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
		return nil, 0, fmt.Errorf("query count users failed: %w", err)
	}

	// Build sorting expressions
	// Convert requested sort keys into SQL ORDER BY clauses
	var userSortKeys = map[query.SortKey]string{
		repository.CustomerSortLatest:    "u.created_at",
		repository.CustomerSortName:      "u.name",
		repository.CustomerSortUsername:  "u.username",
		repository.CustomerSortPhone:     "u.phone",
		repository.CustomerSortModify:    "u.updated_at",
		repository.CustomerSortLastLogin: "a.last_login_at",
	}

	var sortClauses []string
	for _, sort := range params.Sorts {
		colName, exists := userSortKeys[sort.By]
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

	orderBy := "ORDER BY u.created_at DESC"
	if len(sortClauses) > 0 {
		orderBy = "ORDER BY " + strings.Join(sortClauses, ", ")
	}

	// Apply pagination
	// Calculate limit and offset values
	limitPos := argPos
	offsetPos := argPos + 1

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

	// The execution
	finalQuery := selectQuery + baseQuery + whereClause + " " + orderBy +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	rows, err := exec.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Customer
	for rows.Next() {
		var m domain.Customer
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Username,
			&m.Phone,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
			&m.LastLoginAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("mapping user model to domain failed: %w", err)
		}

		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users failed: %w", err)
	}

	return results, total, nil
}
