package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type cartRepositoryImpl struct{}

func NewCartRepositoryImpl() repository.CartRepository {
	return &cartRepositoryImpl{}
}

func (r *cartRepositoryImpl) GetWithItemsByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*domain.Cart, error) {
	query := `
		SELECT 
			c.id,
			c.customer_id,
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
		WHERE c.customer_id = $1
		ORDER BY ci.created_at
	`

	rows, err := exec.Query(ctx, query, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query cart with items by customer id failed: %w", err)
	}
	defer rows.Close()

	var cart *domain.Cart

	for rows.Next() {
		var (
			cID       uuid.UUID
			custID    uuid.UUID
			createdAt time.Time
			updatedAt *time.Time
			itemID    *uuid.UUID
			productID *uuid.UUID
			shopID    *uuid.UUID
			quantity  *int
		)

		err := rows.Scan(
			&cID,
			&custID,
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
				ID:         cID,
				CustomerID: custID,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Items:      []domain.CartItem{},
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

func (r *cartRepositoryImpl) NewCart(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*domain.Cart, error) {
	query := `
		INSERT INTO carts (customer_id)
		VALUES ($1)
		RETURNING id, customer_id, created_at, updated_at
	`

	var cart domain.Cart

	err := exec.QueryRow(ctx, query, customerID).Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert cart failed: %w", err)
	}

	cart.Items = []domain.CartItem{}

	return &cart, nil
}

func (r *cartRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	cart *domain.Cart,
) error {
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
			_, err := exec.Exec(ctx, updateItemQuery,
				cart.ID,
				item.ProductID,
				item.ShopID,
			)
			if err != nil {
				return fmt.Errorf("update cart items failed: %w", err)
			}
			continue
		}

		_, err := exec.Exec(ctx, insertItemQuery,
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

	_, err := exec.Exec(ctx, updateCartQuery, cart.ID)
	if err != nil {
		return fmt.Errorf("update cart failed: %w", err)
	}

	return nil
}

func (r *cartRepositoryImpl) DeleteByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) error {
	query := `
		UPDATE cart_items ci
		SET deleted_at = NOW(), updated_at = NOW()
		FROM carts c
		WHERE ci.cart_id = c.id
			AND c.customer_id = $1
			AND ci.deleted_at IS NULL
	`

	_, err := exec.Exec(ctx, query, customerID)
	if err != nil {
		return fmt.Errorf("delete customer cart items failed: %w", err)
	}

	return nil
}
