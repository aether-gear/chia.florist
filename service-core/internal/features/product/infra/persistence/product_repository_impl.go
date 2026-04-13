package persistence

import (
	"context"
	"fmt"
	domain "service-core/internal/features/product/domain"
	"service-core/internal/features/product/repository"
	database "service-core/internal/infra/db"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProductRepository(conn *database.Connection) *productRepositoryImpl {
	return &productRepositoryImpl{db: conn.Pool}
}

func (r *productRepositoryImpl) FindProducts(params repository.FindProductParams) ([]domain.Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	query := `
		SELECT id, sku, name, description, status,
		       base_price, weight,
		       created_at, updated_at, archived_at, deleted_at
		FROM products
	`

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Name+"%")
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM products"
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
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []domain.Product

	for rows.Next() {
		var m ProductModel

		err := rows.Scan(
			&m.ID,
			&m.SKU,
			&m.Name,
			&m.Description,
			&m.Status,
			&m.BasePrice,
			&m.Weight,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.ArchivedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		d, err := m.ToDomain()
		if err != nil {
			return nil, 0, err
		}

		results = append(results, *d)
	}

	return results, total, nil
}

func (r *productRepositoryImpl) GetById(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, sku, name, description, status,
		       base_price, weight,
		       created_at, updated_at, archived_at, deleted_at
		FROM products
		WHERE id = $1
		LIMIT 1
	`

	var m ProductModel

	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.SKU,
		&m.Name,
		&m.Description,
		&m.Status,
		&m.BasePrice,
		&m.Weight,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.ArchivedAt,
		&m.DeletedAt,
	)

	if err != nil {
		return &domain.Product{}, err
	}

	d, err := m.ToDomain()
	if err != nil {
		return &domain.Product{}, err
	}

	return d, nil
}
