package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productImageRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProductImageRepository(conn *database.Connection) repository.ProductImageRepository {
	return &productImageRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *productImageRepositoryImpl) FindByProductID(productID uuid.UUID) ([]domain.ProductImage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			product_id,
			thumbnail_url,
			preview_url,
			detail_url,
			thumbnail_key,
			preview_key,
			detail_key,
			is_primary,
			display_order,
			created_at,
			deleted_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product images failed: %w", err)
	}
	defer rows.Close()

	var images []domain.ProductImage
	for rows.Next() {
		var img domain.ProductImage
		err := rows.Scan(
			&img.ID,
			&img.ProductID,
			&img.Thumbnail.URL,
			&img.Preview.URL,
			&img.Detail.URL,
			&img.Thumbnail.Key,
			&img.Preview.Key,
			&img.Detail.Key,
			&img.IsPrimary,
			&img.DisplayOrder,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping product image model to domain failed: %w", err)
		}

		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image failed: %w", err)
	}

	return images, nil
}

func (r *productImageRepositoryImpl) Create(images []domain.ProductImage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx create product image failed: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO product_images (
			id,
			product_id,
			thumbnail_url,
			preview_url,
			detail_url,
			thumbnail_key,
			preview_key,
			detail_key,
			is_primary,
			display_order,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	for _, image := range images {
		_, err := r.db.Exec(ctx, query,
			image.ID,
			image.ProductID,
			image.Thumbnail.URL,
			image.Preview.URL,
			image.Detail.URL,
			image.Thumbnail.Key,
			image.Preview.Key,
			image.Detail.Key,
			image.IsPrimary,
			image.DisplayOrder,
			image.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert product image failed: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx insert product image failed: %w", err)
	}

	return nil
}

func (r *productImageRepositoryImpl) SoftDeleteByProductID(productID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE
			product_images 
		SET
			deleted_at = NOW()
		WHERE
			product_id = $1
	`

	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("delete product images failed: %w", err)
	}

	return nil
}
