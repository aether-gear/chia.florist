package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProductRepository(conn *database.Connection) repository.ProductRepository {
	return &productRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *productRepositoryImpl) FindProducts(params repository.FindProductParams) ([]domain.Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	baseQuery := `
		FROM products p
	`

	selectQuery := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.slug,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at
	`

	conditions = append(conditions, "p.deleted_at IS NULL")

	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf("p.id = $%d", argPos))
		args = append(args, *params.ID)
		argPos++
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", argPos))
		args = append(args, "%"+*params.Name+"%")
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countArgs := append([]any{}, args...)
	countQuery := "SELECT COUNT(DISTINCT p.id) " + baseQuery + whereClause

	var total int
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("query count products failed: %w", err)
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

	limitPos := argPos
	offsetPos := argPos + 1

	query := selectQuery + baseQuery + whereClause +
		fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query products failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Product

	for rows.Next() {
		var item domain.Product

		err := rows.Scan(
			&item.ID,
			&item.SKU,
			&item.Name,
			&item.Slug,
			&item.Description,
			&item.Status,
			&item.Price,
			&item.Weight,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ArchivedAt,
			&item.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("mapping product model to domain failed: %w", err)
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate products failed: %w", err)
	}

	return results, total, nil
}
func (r *productRepositoryImpl) GetByID(id uuid.UUID) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.slug,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at
		FROM products p
		WHERE p.id = $1
		LIMIT 1
	`

	var result domain.Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.SKU,
		&result.Name,
		&result.Slug,
		&result.Description,
		&result.Status,
		&result.Price,
		&result.Weight,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.ArchivedAt,
		&result.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query product by id failed: %w", err)
	}

	return &result, nil
}

func (r *productRepositoryImpl) FindByIDs(ids []uuid.UUID) ([]domain.Product, error) {
	if len(ids) == 0 {
		return []domain.Product{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.slug,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at
		FROM products p
		WHERE p.id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query products by many ids failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Product

	for rows.Next() {
		var item domain.Product

		err := rows.Scan(
			&item.ID,
			&item.SKU,
			&item.Name,
			&item.Slug,
			&item.Description,
			&item.Status,
			&item.Price,
			&item.Weight,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ArchivedAt,
			&item.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("mapping product model to domain failed: %w", err)
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products failed: %w", err)
	}

	return results, nil
}

func (r *productRepositoryImpl) CreateProduct(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (
			id,
			sku,
			name,
			slug,
			description,
			status,
			base_price,
			weight,
			created_at
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.SKU,
		product.Name,
		product.Slug,
		product.Description,
		product.Status,
		product.Price,
		product.Weight,
		product.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert product failed: %w", err)
	}

	return nil
}
