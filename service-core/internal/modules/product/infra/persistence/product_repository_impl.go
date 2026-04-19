package persistence

import (
	"context"
	"errors"
	"fmt"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"strings"
	"time"

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

func (r *productRepositoryImpl) FindProducts(params repository.FindProductParams) ([]repository.ProductWithInventory, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		conditions []string
		args       []any
		argPos     = 1
	)

	baseQuery := `
		FROM products p
		LEFT JOIN inventory i ON i.product_id = p.id
	`

	selectQuery := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at,
			COALESCE(i.stock, 0) AS stock,
			COALESCE(i.reserved_stock, 0) AS reserved_stock
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

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
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

	query := selectQuery + baseQuery + whereClause +
		fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []repository.ProductWithInventory

	for rows.Next() {
		var item repository.ProductWithInventory

		err := rows.Scan(
			&item.Product.ID,
			&item.Product.SKU,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Status,
			&item.Product.Price,
			&item.Product.Weight,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
			&item.Product.ArchivedAt,
			&item.Product.DeletedAt,
			&item.Inventory.Stock,
			&item.Inventory.ReservedStock,
		)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, item)
	}

	return results, total, nil
}

func (r *productRepositoryImpl) GetByID(id uuid.UUID) (*repository.ProductWithInventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at,
			COALESCE(i.stock, 0) AS stock,
			COALESCE(i.reserved_stock, 0) AS reserved_stock
		FROM products p
		LEFT JOIN inventory i on i.product_id = p.id
		WHERE p.id = $1
		LIMIT 1
	`

	var result repository.ProductWithInventory

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.Product.ID,
		&result.Product.SKU,
		&result.Product.Name,
		&result.Product.Description,
		&result.Product.Status,
		&result.Product.Price,
		&result.Product.Weight,
		&result.Product.CreatedAt,
		&result.Product.UpdatedAt,
		&result.Product.ArchivedAt,
		&result.Product.DeletedAt,
		&result.Inventory.Stock,
		&result.Inventory.ReservedStock,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *productRepositoryImpl) FindByIDs(ids []uuid.UUID) ([]repository.ProductWithInventory, error) {
	if len(ids) == 0 {
		return []repository.ProductWithInventory{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.description,
			p.status,
			p.base_price,
			p.weight,
			p.created_at,
			p.updated_at,
			p.archived_at,
			p.deleted_at,
			COALESCE(i.stock, 0) AS stock,
			COALESCE(i.reserved_stock, 0) AS reserved_stock
		FROM products p
		LEFT JOIN inventory i on i.product_id = p.id
		WHERE p.id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []repository.ProductWithInventory

	for rows.Next() {
		var item repository.ProductWithInventory

		err := rows.Scan(
			&item.Product.ID,
			&item.Product.SKU,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Status,
			&item.Product.Price,
			&item.Product.Weight,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
			&item.Product.ArchivedAt,
			&item.Product.DeletedAt,
			&item.Inventory.Stock,
			&item.Inventory.ReservedStock,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, item)
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
			description,
			status,
			base_price,
			weight,
			created_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Query(ctx, query,
		product.ID,
		product.SKU,
		product.Name,
		product.Description,
		product.Status,
		product.Price,
		product.Weight,
		product.CreatedAt,
	)

	return err
}

func (r *productRepositoryImpl) CreateInventory(inventory *domain.Inventory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO inventory (
			id,
			product_id,
			stock,
			reserved_stock,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Query(ctx, query,
		inventory.ID,
		inventory.ProductID,
		inventory.Stock,
		inventory.ReservedStock,
		inventory.CreatedAt,
	)

	return err
}
