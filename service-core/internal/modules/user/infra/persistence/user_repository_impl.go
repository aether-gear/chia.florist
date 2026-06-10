package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepositoryImpl struct{}

func NewUserRepositoryImpl() repository.UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) FindUsers(
	ctx context.Context,
	exec transaction.Executor,
	params repository.FindUserParams,
) ([]domain.User, int, error) {
	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	query := `
		SELECT 
			u.id, 
			u.name, 
			u.username, 
			u.phone,
			u.created_at, 
			u.updated_at, 
			u.deleted_at,
			a.last_login_at
		FROM users u
		LEFT JOIN accounts a ON a.user_id = u.id
	`

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Username != nil {
		conditions = append(conditions, fmt.Sprintf("username ILIKE $%d", argPos))
		args = append(args, "%"+*params.Username+"%")
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
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

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	var total int
	err := exec.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count users failed: %w", err)
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users failed: %w", err)
	}
	defer rows.Close()

	var results []domain.User

	for rows.Next() {
		var m domain.User

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

func (r *userRepositoryImpl) GetByID(
	ctx context.Context,
	exec transaction.Executor,
	id uuid.UUID,
) (*domain.User, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.username,
			u.phone,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			a.last_login_at
		FROM users u
		LEFT JOIN accounts a ON a.user_id = u.id
		WHERE u.id = $1
		LIMIT 1
	`

	var m domain.User

	err := exec.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return &domain.User{}, fmt.Errorf("query user by id failed: %w", err)
	}

	return &m, nil
}

func (r *userRepositoryImpl) GetByUsername(
	ctx context.Context,
	exec transaction.Executor,
	username string,
) (*domain.User, error) {
	query := `
		SELECT 
			u.id,
			u.name,
			u.username,
			u.phone,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			a.last_login_at
		FROM users u
		LEFT JOIN accounts a ON a.user_id = u.id
		WHERE u.username = $1
		LIMIT 1
	`

	var m domain.User
	err := exec.QueryRow(ctx, query, username).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query user by username failed: %w", err)
	}

	return &m, nil
}

func (r *userRepositoryImpl) CreateUser(
	ctx context.Context,
	exec transaction.Executor,
	props repository.CreateUserProps,
) error {
	query := `
		INSERT INTO users (
			id,
			name,
			username,
			phone,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := exec.Exec(ctx, query,
		props.ID,
		props.Name,
		props.Username,
		props.Phone,
		props.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user failed: %w", err)
	}

	return nil
}
