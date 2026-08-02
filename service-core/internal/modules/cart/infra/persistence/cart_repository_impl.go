package persistence

import (
	"context"
	"encoding/json"
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
			ci.product_variant_type,
			ci.product_id,
			ci.shop_id,
			ci.quantity,
			ci.custom_design
		FROM carts c
		LEFT JOIN
			cart_items ci ON ci.cart_id = c.id
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
			cID                uuid.UUID
			custID             uuid.UUID
			createdAt          time.Time
			updatedAt          *time.Time
			itemID             *uuid.UUID
			ProductVariantType *string
			productID          *uuid.UUID
			shopID             *uuid.UUID
			quantity           *int
			customDesign       []byte
		)

		err := rows.Scan(
			&cID,
			&custID,
			&createdAt,
			&updatedAt,
			&itemID,
			&ProductVariantType,
			&productID,
			&shopID,
			&quantity,
			&customDesign,
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
			it := domain.ProductVariantTypeStandard
			if ProductVariantType != nil && *ProductVariantType == string(domain.ProductVariantTypeCustom) {
				it = domain.ProductVariantTypeCustom
			}
			cart.Items = append(cart.Items, domain.CartItem{
				ID:                 *itemID,
				ProductVariantType: it,
				ProductID:          productID,
				ShopID:             *shopID,
				Quantity:           *quantity,
				CustomDesign:       json.RawMessage(customDesign),
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
		INSERT INTO carts (
			customer_id
		)
		VALUES ($1)
		RETURNING
			id,
			customer_id,
			created_at,
			updated_at
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
	const insertStandardQuery = `
		INSERT INTO cart_items (
			cart_id,
			product_variant_type,
			product_id,
			shop_id,
			quantity
		)
		VALUES ($1,'standard',$2,$3,$4)
		ON CONFLICT (cart_id, product_id)
		WHERE deleted_at IS NULL
			AND product_variant_type = 'standard'
		DO UPDATE SET
			quantity   = EXCLUDED.quantity,
			updated_at = NOW()
	`

	const insertCustomQuery = `
		INSERT INTO cart_items (
			cart_id,
			product_variant_type,
			shop_id,
			quantity,
			custom_design
		)
		VALUES ($1,'custom',$2,$3,$4)
	`

	const softDeleteByProductQuery = `
		UPDATE cart_items
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE cart_id = $1
			AND product_id = $2
			AND shop_id = $3
			AND deleted_at IS NULL
	`

	const softDeleteByIDQuery = `
		UPDATE cart_items
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	for _, item := range cart.Items {
		if item.DeletedAt != nil {
			if item.ProductVariantType == domain.ProductVariantTypeCustom {
				if _, err := exec.Exec(ctx, softDeleteByIDQuery, item.ID); err != nil {
					return fmt.Errorf("soft-delete custom cart item failed: %w", err)
				}
			} else {
				if _, err := exec.Exec(ctx, softDeleteByProductQuery,
					cart.ID, item.ProductID, item.ShopID,
				); err != nil {
					return fmt.Errorf("soft-delete standard cart item failed: %w", err)
				}
			}
			continue
		}

		if item.ProductVariantType == domain.ProductVariantTypeCustom {
			if _, err := exec.Exec(ctx, insertCustomQuery,
				cart.ID, item.ShopID, item.Quantity, []byte(item.CustomDesign),
			); err != nil {
				return fmt.Errorf("insert custom cart item failed: %w", err)
			}
		} else {
			if _, err := exec.Exec(ctx, insertStandardQuery,
				cart.ID, item.ProductID, item.ShopID, item.Quantity,
			); err != nil {
				return fmt.Errorf("insert standard cart item failed: %w", err)
			}
		}
	}

	const updateCartQuery = `
		UPDATE carts
		SET
			updated_at = NOW()
		WHERE id = $1
	`

	if _, err := exec.Exec(ctx, updateCartQuery, cart.ID); err != nil {
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
		SET
			deleted_at = NOW(),
			updated_at = NOW()
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
