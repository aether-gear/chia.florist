package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepositoryImpl(conn *database.Connection) repository.UserRepository {
	return &userRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *userRepositoryImpl) FindUsers(params repository.FindUserParams) ([]domain.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	query := `
		SELECT 
			id, 
			name, 
			username, 
			phone,
			created_at, 
			updated_at, 
			deleted_at, 
			last_login_at
		FROM users
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
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count users failed: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
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

func (r *userRepositoryImpl) GetByID(id uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			id,
			name,
			username,
			phone,
			created_at,
			updated_at,
			deleted_at,
			last_login_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var m domain.User

	err := r.db.QueryRow(ctx, query, id).Scan(
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
		return &domain.User{}, fmt.Errorf("query user by id failed: %w", err)
	}

	return &m, nil
}

func (r *userRepositoryImpl) GetUserWithAccount(id uuid.UUID) (*repository.UserWithAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			id,
			name,
			username,
			email,
			phone,
			created_at,
			updated_at,
			deleted_at,
			last_login_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var m repository.UserWithAccount

	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Name,
		&m.Username,
		&m.Email,
		&m.Phone,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
		&m.LastLoginAt,
	)
	if err != nil {
		return &repository.UserWithAccount{}, fmt.Errorf("query user with account by id failed: %w", err)
	}

	return &m, nil
}
