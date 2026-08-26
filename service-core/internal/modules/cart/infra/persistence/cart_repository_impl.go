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
			ci.custom_design,
			ci.item_options
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
			rawOptions         []byte
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
			&rawOptions,
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

			var itemOptions domain.ItemOptions
			if len(rawOptions) > 0 {
				_ = json.Unmarshal(rawOptions, &itemOptions)
			}

			cart.Items = append(cart.Items, domain.CartItem{
				ID:                 *itemID,
				ProductVariantType: it,
				ProductID:          productID,
				ShopID:             *shopID,
				Quantity:           *quantity,
				ItemOptions:        itemOptions.Normalized(),
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
			id,
			cart_id,
			product_variant_type,
			product_id,
			shop_id,
			quantity,
			item_options
		)
		VALUES ($1,$2,'standard',$3,$4,$5,$6::jsonb)
		ON CONFLICT (id)
		DO UPDATE SET
			shop_id      = EXCLUDED.shop_id,
			quantity     = EXCLUDED.quantity,
			item_options = EXCLUDED.item_options,
			updated_at   = NOW()
	`

	const insertCustomQuery = `
	    INSERT INTO cart_items (
	        id,
	        cart_id,
	        product_variant_type,
	        shop_id,
	        quantity,
	        custom_design
	    )
	    VALUES ($1,$2,'custom',$3,$4,$5::jsonb)
	    ON CONFLICT (id)
	    DO UPDATE SET
	        shop_id       = EXCLUDED.shop_id,
	        quantity      = EXCLUDED.quantity,
	        custom_design = EXCLUDED.custom_design,
	        updated_at    = NOW()
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
			if _, err := exec.Exec(ctx, softDeleteByIDQuery, item.ID); err != nil {
				return fmt.Errorf("soft-delete cart item failed: %w", err)
			}
			continue
		}

		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}

		var (
			query     string
			args      []any
			errFormat string
		)

		switch item.ProductVariantType {
		case domain.ProductVariantTypeCustom:
			query = insertCustomQuery
			var customDesignArg any
			if len(item.CustomDesign) > 0 {
				customDesignArg = string(item.CustomDesign)
			}
			args = []any{
				item.ID,
				cart.ID,
				item.ShopID,
				item.Quantity,
				customDesignArg,
			}
			errFormat = "insert custom cart item failed: %w"

		case domain.ProductVariantTypeStandard:
			query = insertStandardQuery
			optBytes, _ := json.Marshal(item.ItemOptions.Normalized())
			args = []any{
				item.ID,
				cart.ID,
				item.ProductID,
				item.ShopID,
				item.Quantity,
				string(optBytes),
			}
			errFormat = "insert standard cart item failed: %w"

		default:
			return fmt.Errorf("unsupported product variant type: %q", item.ProductVariantType)
		}

		if _, err := exec.Exec(ctx, query, args...); err != nil {
			return fmt.Errorf(errFormat, err)
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
