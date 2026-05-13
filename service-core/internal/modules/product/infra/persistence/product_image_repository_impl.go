package persistence

import (
	"context"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"service-core/internal/shared/image"

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

func (r *productImageRepositoryImpl) ListByProductIDs(
	productIDs []uuid.UUID,
) (map[uuid.UUID][]domain.ProductImage, error) {
	result := make(map[uuid.UUID][]domain.ProductImage)
	if len(productIDs) == 0 {
		return result, nil
	}

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
		WHERE product_id = ANY($1::uuid[])
		ORDER BY display_order ASC
	`

	productIDStrings := make([]string, len(productIDs))
	for i, id := range productIDs {
		productIDStrings[i] = id.String()
	}

	rows, err := r.db.Query(ctx, query, productIDStrings)
	if err != nil {
		return nil, fmt.Errorf("query product images failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r productImageRow

		err := rows.Scan(
			&r.ID,
			&r.ProductID,
			&r.ThumbURL,
			&r.PreviewURL,
			&r.DetailURL,
			&r.ThumbKey,
			&r.PreviewKey,
			&r.DetailKey,
			&r.IsPrimary,
			&r.DisplayOrder,
			&r.CreatedAt,
			&r.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan product image: %w", err)
		}

		img := domain.ProductImage{
			ID:        r.ID,
			ProductID: r.ProductID,

			Variants: map[image.ResolutionType]domain.ImageVariant{
				domain.ResolutionThumbnail: {
					Type: domain.ResolutionThumbnail,
					Key:  r.ThumbKey,
				},
				domain.ResolutionPreview: {
					Type: domain.ResolutionPreview,
					Key:  r.PreviewKey,
				},
				domain.ResolutionDetail: {
					Type: domain.ResolutionDetail,
					Key:  r.DetailKey,
				},
			},

			IsPrimary:    r.IsPrimary,
			DisplayOrder: r.DisplayOrder,
			Metadata:     domain.ProductImageMetadata{},
			CreatedAt:    r.CreatedAt,
		}

		result[img.ProductID] = append(result[img.ProductID], img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image failed: %w", err)
	}

	return result, nil
}

func (r *productImageRepositoryImpl) ListByProductID(productID uuid.UUID) ([]domain.ProductImage, error) {
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

	var rowsData []productImageRow
	for rows.Next() {
		var r productImageRow

		err := rows.Scan(
			&r.ID,
			&r.ProductID,
			&r.ThumbURL,
			&r.ThumbKey,
			&r.PreviewURL,
			&r.PreviewKey,
			&r.DetailURL,
			&r.DetailKey,
			&r.IsPrimary,
			&r.DisplayOrder,
			&r.CreatedAt,
			&r.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan product image: %w", err)
		}

		rowsData = append(rowsData, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image failed: %w", err)
	}

	images := make([]domain.ProductImage, 0, len(rowsData))

	for _, r := range rowsData {
		img := domain.ProductImage{
			ID:        r.ID,
			ProductID: r.ProductID,

			Variants: map[image.ResolutionType]domain.ImageVariant{
				domain.ResolutionThumbnail: {
					Type: domain.ResolutionThumbnail,
					Key:  r.ThumbKey,
				},
				domain.ResolutionPreview: {
					Type: domain.ResolutionPreview,
					Key:  r.PreviewKey,
				},
				domain.ResolutionDetail: {
					Type: domain.ResolutionDetail,
					Key:  r.DetailKey,
				},
			},

			IsPrimary:    r.IsPrimary,
			DisplayOrder: r.DisplayOrder,

			Metadata:  domain.ProductImageMetadata{},
			CreatedAt: r.CreatedAt,
		}

		images = append(images, img)
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
			image.Variants[domain.ResolutionThumbnail].Key,
			image.Variants[domain.ResolutionPreview].Key,
			image.Variants[domain.ResolutionDetail].Key,
			image.Variants[domain.ResolutionThumbnail].Key,
			image.Variants[domain.ResolutionPreview].Key,
			image.Variants[domain.ResolutionDetail].Key,
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
