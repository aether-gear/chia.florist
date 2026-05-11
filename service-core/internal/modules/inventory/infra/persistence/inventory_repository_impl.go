package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/inventory/domain"
	"service-core/internal/modules/inventory/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type inventoryRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(conn *database.Connection) repository.InventoryRepository {
	return &inventoryRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *inventoryRepositoryImpl) GetByProductAndShop(productID uuid.UUID, shopID uuid.UUID) (*domain.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			product_id,
			shop_id,
			stock,
			reserved_stock,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1 AND shop_id = $2
		LIMIT 1
	`

	var inventory domain.Inventory
	err := r.db.QueryRow(ctx, query, productID, shopID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.ShopID,
		&inventory.TotalStock,
		&inventory.ReservedStock,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query inventory by product and shop failed: %w", err)
	}

	return &inventory, nil
}

func (r *inventoryRepositoryImpl) ListByProduct(productID uuid.UUID) ([]domain.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			product_id,
			shop_id,
			stock,
			reserved_stock,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query inventory by product failed: %w", err)
	}
	defer rows.Close()

	var inventories []domain.Inventory

	for rows.Next() {
		var inventory domain.Inventory
		if err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.ShopID,
			&inventory.TotalStock,
			&inventory.ReservedStock,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("mapping inventory model to domain failed: %w", err)
		}

		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inventory failed: %w", err)
	}

	return inventories, nil
}

func (r *inventoryRepositoryImpl) ListByProducts(productIDs []uuid.UUID) (map[uuid.UUID][]domain.Inventory, error) {
	result := make(map[uuid.UUID][]domain.Inventory)
	if len(productIDs) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			product_id,
			shop_id,
			stock,
			reserved_stock,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = ANY($1::uuid[])
		ORDER BY created_at ASC
	`

	productIDStrings := make([]string, len(productIDs))
	for i, id := range productIDs {
		productIDStrings[i] = id.String()
	}

	rows, err := r.db.Query(ctx, query, productIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query inventory by products failed: %w", err)
	}
	defer rows.Close()

	var inventories []domain.Inventory

	for rows.Next() {
		var inventory domain.Inventory
		if err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.ShopID,
			&inventory.TotalStock,
			&inventory.ReservedStock,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("mapping inventory model to domain failed: %w", err)
		}

		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inventory failed: %w", err)
	}

	for _, inventory := range inventories {
		result[inventory.ProductID] = append(result[inventory.ProductID], inventory)
	}

	return result, nil
}

func (r *inventoryRepositoryImpl) ListByShop(shopID uuid.UUID) ([]domain.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			product_id,
			shop_id,
			stock,
			reserved_stock,
			created_at,
			updated_at
		FROM inventory
		WHERE shop_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("query inventory by shop failed: %w", err)
	}
	defer rows.Close()

	var inventories []domain.Inventory

	for rows.Next() {
		var inventory domain.Inventory
		if err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.ShopID,
			&inventory.TotalStock,
			&inventory.ReservedStock,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("mapping inventory model to domain failed: %w", err)
		}

		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inventory failed: %w", err)
	}

	return inventories, nil
}

func (r *inventoryRepositoryImpl) Create(inventory *domain.Inventory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO inventory (
			id,
			product_id,
			shop_id,
			stock,
			reserved_stock,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		inventory.ID,
		inventory.ProductID,
		inventory.ShopID,
		inventory.TotalStock,
		inventory.ReservedStock,
		inventory.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert inventory failed: %w", err)
	}

	return nil
}
