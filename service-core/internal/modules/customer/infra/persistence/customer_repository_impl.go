package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/customer/domain"
	"service-core/internal/modules/customer/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type customerRepositoryImpl struct{}

func NewCustomerRepositoryImpl() repository.CustomerRepository {
	return &customerRepositoryImpl{}
}

func (r *customerRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	customer domain.Customer,
) error {
	query := `
		INSERT INTO customers (
			id,
			user_id,
			created_at
		) VALUES ($1, $2, $3)
	`

	_, err := exec.Exec(ctx, query,
		customer.ID,
		customer.UserID,
		customer.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert customer failed: %w", err)
	}
	return nil
}

func (r *customerRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.Customer, error) {
	query := `
		SELECT
			id,
			user_id,
			created_at,
			updated_at,
			deleted_at
		FROM customers
		WHERE id = $1
		LIMIT 1
	`

	var m domain.Customer
	err := exec.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query customer by id failed: %w", err)
	}
	return &m, nil
}

func (r *customerRepositoryImpl) GetByUserID(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
) (*domain.Customer, error) {
	query := `
		SELECT
			id,
			user_id,
			created_at,
			updated_at,
			deleted_at
		FROM customers
		WHERE user_id = $1
		LIMIT 1
	`

	var m domain.Customer
	err := exec.QueryRow(ctx, query, userID).Scan(
		&m.ID,
		&m.UserID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query customer by user_id failed: %w", err)
	}
	return &m, nil
}

func (r *customerRepositoryImpl) GetProfileByUserID(
	ctx context.Context,
	exec transaction.Executor,
	userID uuid.UUID,
) (*domain.CustomerProfile, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			u.name,
			u.username,
			u.phone,
			u.avatar_url,
			c.created_at,
			c.updated_at
		FROM customers c
		INNER JOIN users u
			ON u.id = c.user_id
		WHERE c.user_id = $1
	`

	var profile domain.CustomerProfile
	err := exec.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Name,
		&profile.Username,
		&profile.Phone,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound("customer not found")
		}

		return nil, fmt.Errorf("get staff profile by user id failed: %w", err)
	}

	return &profile, nil
}

func (r *customerRepositoryImpl) FindCustomers(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindCustomerParams,
) ([]domain.CustomerProfile, int, error) {
	baseQuery := `
		FROM customers m
		INNER JOIN users u ON u.id = m.user_id
		LEFT JOIN accounts a ON a.user_id = u.id
	`

	selectQuery := `
		SELECT 
			m.id, 
			m.user_id,
			u.name, 
			u.username, 
			u.phone,
			u.avatar_url,
			m.created_at, 
			m.updated_at, 
			a.last_login_at
	`

	// Build filters
	// Apply search criteria and soft-delete constraints
	whereClause := ""
	notDeletedCondition := "m.deleted_at IS NULL AND u.deleted_at IS NULL"
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
		conditions = append(conditions, fmt.Sprintf("m.id = $%d", argPos))
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
		repository.CustomerSortLatest:    "m.created_at",
		repository.CustomerSortName:      "u.name",
		repository.CustomerSortUsername:  "u.username",
		repository.CustomerSortPhone:     "u.phone",
		repository.CustomerSortModify:    "m.updated_at",
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

	orderBy := "ORDER BY m.created_at DESC"
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

	var results []domain.CustomerProfile
	for rows.Next() {
		var m domain.CustomerProfile
		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Name,
			&m.Username,
			&m.Phone,
			&m.AvatarURL,
			&m.CreatedAt,
			&m.UpdatedAt,
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
