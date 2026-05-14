package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
			deleted_at
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
			deleted_at
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

func (r *userRepositoryImpl) GetByUsername(username string) (*domain.User, error) {
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
			deleted_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	var m domain.User

	err := r.db.QueryRow(ctx, query, username).Scan(
		&m.ID,
		&m.Name,
		&m.Username,
		&m.Phone,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query user by username failed: %w", err)
	}

	return &m, nil
}

func (r *userRepositoryImpl) CreateUser(props repository.CreateUserProps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	_, err := r.db.Exec(ctx, query,
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
