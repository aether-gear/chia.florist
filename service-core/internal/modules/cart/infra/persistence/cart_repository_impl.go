package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "service-core/internal/infra/db"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type cartRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCartRepositoryImpl(conn *database.Connection) repository.CartRepository {
	return &cartRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *cartRepositoryImpl) GetWithItemsByUserID(userID uuid.UUID) (*domain.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.id,
			c.user_id,
			c.created_at,
			c.updated_at,
			ci.id,
			ci.product_id,
			ci.shop_id,
			ci.quantity
		FROM carts c
		LEFT JOIN cart_items ci 
			ON ci.cart_id = c.id 
			AND ci.deleted_at IS NULL
		WHERE c.user_id = $1
		ORDER BY ci.created_at
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query cart with items by user id failed: %w", err)
	}
	defer rows.Close()

	var cart *domain.Cart

	for rows.Next() {
		var (
			cID       uuid.UUID
			uID       uuid.UUID
			createdAt time.Time
			updatedAt *time.Time
			itemID    *uuid.UUID
			productID *uuid.UUID
			shopID    *uuid.UUID
			quantity  *int
		)

		err := rows.Scan(
			&cID,
			&uID,
			&createdAt,
			&updatedAt,
			&itemID,
			&productID,
			&shopID,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("mapping cart with items model to domain failed: %w", err)
		}

		if cart == nil {
			cart = &domain.Cart{
				ID:        cID,
				UserID:    uID,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Items:     []domain.CartItem{},
			}
		}

		if itemID != nil {
			cart.Items = append(cart.Items, domain.CartItem{
				ID:        *itemID,
				ProductID: *productID,
				ShopID:    *shopID,
				Quantity:  *quantity,
			})
		}
	}

	if cart == nil {
		return nil, nil
	}

	return cart, nil
}

func (r *cartRepositoryImpl) NewCart(userID uuid.UUID) (*domain.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO carts (user_id)
		VALUES ($1)
		RETURNING id, user_id, created_at, updated_at
	`

	var cart domain.Cart

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert cart failed: %w", err)
	}

	cart.Items = []domain.CartItem{}

	return &cart, nil
}

func (r *cartRepositoryImpl) Save(cart *domain.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx save cart failed: %w", err)
	}
	defer tx.Rollback(ctx)

	updateItemQuery := `
		UPDATE cart_items
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE cart_id = $1 AND product_id = $2 AND shop_id = $3 AND deleted_at IS NULL
	`

	insertItemQuery := `
		INSERT INTO cart_items (
			cart_id,
			product_id,
			shop_id,
			quantity
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (
			cart_id,
			product_id
		) 
		WHERE deleted_at IS NULL
		DO UPDATE 
		SET 
			quantity = EXCLUDED.quantity,
			updated_at = NOW()
	`

	for _, item := range cart.Items {
		if item.DeletedAt != nil {
			_, err := tx.Exec(ctx, updateItemQuery,
				cart.ID,
				item.ProductID,
				item.ShopID,
			)
			if err != nil {
				return fmt.Errorf("update cart items failed: %w", err)
			}
			continue
		}

		_, err := tx.Exec(ctx, insertItemQuery,
			cart.ID,
			item.ProductID,
			item.ShopID,
			item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("insert cart failed: %w", err)
		}
	}

	updateCartQuery := `
		UPDATE carts
		SET updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateCartQuery, cart.ID)
	if err != nil {
		return fmt.Errorf("update cart failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx save cart failed: %w", err)
	}

	return nil
}
