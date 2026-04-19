package persistence

import (
	"context"
	"errors"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	"time"

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
		return nil, err
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
			quantity  *int
		)

		err := rows.Scan(
			&cID, &uID, &createdAt, &updatedAt,
			&itemID, &productID, &quantity,
		)
		if err != nil {
			return nil, err
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

	var cartRow CartModel

	err := r.db.QueryRow(ctx, query, userID).
		Scan(&cartRow.ID, &cartRow.UserID, &cartRow.CreatedAt, &cartRow.UpdatedAt)
	if err != nil {
		return nil, err
	}

	cart := &domain.Cart{
		ID:        cartRow.ID,
		UserID:    cartRow.UserID,
		Items:     []domain.CartItem{},
		CreatedAt: cartRow.CreatedAt,
		UpdatedAt: cartRow.UpdatedAt,
	}

	return cart, nil
}

func (r *cartRepositoryImpl) Save(cart *domain.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range cart.Items {

		if item.DeletedAt != nil {
			_, err := tx.Exec(ctx, `
				UPDATE cart_items
				SET deleted_at = NOW(), updated_at = NOW()
				WHERE cart_id = $1 AND product_id = $2 AND deleted_at IS NULL
			`, cart.ID, item.ProductID)
			if err != nil {
				return err
			}
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO cart_items (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (cart_id, product_id) 
			WHERE deleted_at IS NULL
			DO UPDATE 
			SET 
				quantity = EXCLUDED.quantity,
				updated_at = NOW()
		`, cart.ID, item.ProductID, item.Quantity)

		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE carts SET updated_at = NOW() WHERE id = $1
	`, cart.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
