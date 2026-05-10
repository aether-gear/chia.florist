package persistence

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	database "service-core/internal/infra/db"
// 	"service-core/internal/modules/product/domain"
// 	"service-core/internal/modules/product/repository"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type productImageRepositoryImpl struct {
// 	db *pgxpool.Pool
// }

// func NewProductImageRepository(conn *database.Connection) repository.ProductImageRepository {
// 	return &productImageRepositoryImpl{
// 		db: conn.Pool,
// 	}
// }

// func (r *productImageRepositoryImpl) CreateProductImage(image *domain.ProductImage) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := `
// 		INSERT INTO product_images (
// 			id,
// 			product_id,
// 			catalog_url,
// 			cart_url,
// 			is_primary,
// 			display_order,
// 			created_at
// 		) VALUES ($1, $2, $3, $4, $5, $6, $7)
// 	`

// 	_, err := r.db.Exec(ctx, query,
// 		image.ID,
// 		image.ProductID,
// 		image.CatalogURL,
// 		image.CartURL,
// 		image.IsPrimary,
// 		image.DisplayOrder,
// 		image.CreatedAt,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("insert product image failed: %w", err)
// 	}

// 	return nil
// }

// func (r *productImageRepositoryImpl) CreateDetailImages(images []domain.ProductDetailImage) error {
// 	if len(images) == 0 {
// 		return nil
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Use batch insert or simple loop. Since length is small, a simple loop or transaction works.
// 	// For simplicity, let's use a transaction
// 	tx, err := r.db.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("begin tx for detail images failed: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	query := `
// 		INSERT INTO product_detail_images (
// 			id,
// 			product_image_id,
// 			image_url,
// 			display_order
// 		) VALUES ($1, $2, $3, $4)
// 	`

// 	for _, img := range images {
// 		_, err := tx.Exec(ctx, query, img.ID, img.ProductImageID, img.URL, img.DisplayOrder)
// 		if err != nil {
// 			return fmt.Errorf("insert detail image failed: %w", err)
// 		}
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("commit tx for detail images failed: %w", err)
// 	}

// 	return nil
// }

// func (r *productImageRepositoryImpl) GetByProductID(productID uuid.UUID) ([]domain.ProductImage, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Get primary images first
// 	query := `
// 		SELECT
// 			id, product_id, catalog_url, cart_url, is_primary, display_order, created_at
// 		FROM product_images
// 		WHERE product_id = $1
// 		ORDER BY display_order ASC
// 	`

// 	rows, err := r.db.Query(ctx, query, productID)
// 	if err != nil {
// 		return nil, fmt.Errorf("query product images failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var images []domain.ProductImage
// 	var imageIDs []uuid.UUID
// 	imageMap := make(map[uuid.UUID]*domain.ProductImage)

// 	for rows.Next() {
// 		var img domain.ProductImage
// 		err := rows.Scan(
// 			&img.ID,
// 			&img.ProductID,
// 			&img.CatalogURL,
// 			&img.CartURL,
// 			&img.IsPrimary,
// 			&img.DisplayOrder,
// 			&img.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan product image failed: %w", err)
// 		}
// 		img.DetailImages = []domain.ProductDetailImage{}
// 		images = append(images, img)
// 		imageIDs = append(imageIDs, img.ID)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate product images failed: %w", err)
// 	}

// 	if len(images) == 0 {
// 		return images, nil
// 	}

// 	for i := range images {
// 		imageMap[images[i].ID] = &images[i]
// 	}

// 	// Get detail images
// 	detailQuery := `
// 		SELECT
// 			id, product_image_id, image_url, display_order
// 		FROM product_detail_images
// 		WHERE product_image_id = ANY($1)
// 		ORDER BY display_order ASC
// 	`

// 	detailRows, err := r.db.Query(ctx, detailQuery, imageIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("query product detail images failed: %w", err)
// 	}
// 	defer detailRows.Close()

// 	for detailRows.Next() {
// 		var dImg domain.ProductDetailImage
// 		err := detailRows.Scan(
// 			&dImg.ID,
// 			&dImg.ProductImageID,
// 			&dImg.URL,
// 			&dImg.DisplayOrder,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan detail image failed: %w", err)
// 		}

// 		if parent, ok := imageMap[dImg.ProductImageID]; ok {
// 			parent.DetailImages = append(parent.DetailImages, dImg)
// 		}
// 	}

// 	if err := detailRows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate detail images failed: %w", err)
// 	}

// 	// Copy back from map to slice since we modified the pointer contents
// 	var result []domain.ProductImage
// 	for _, img := range images {
// 		result = append(result, *imageMap[img.ID])
// 	}

// 	return result, nil
// }

// func (r *productImageRepositoryImpl) DeleteByProductID(productID uuid.UUID) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// ON DELETE CASCADE will handle detail_images deletion
// 	query := `DELETE FROM product_images WHERE product_id = $1`

// 	_, err := r.db.Exec(ctx, query, productID)
// 	if err != nil {
// 		return fmt.Errorf("delete product images failed: %w", err)
// 	}

// 	return nil
// }
