# Implementation Plan: Backend Custom Product Support (`service-core`)

Reference document: [`D:\__Projects\kage\chia.florist\docs\CUSTOM_PRODUCT.md`](file:///d:/__Projects/kage/chia.florist/docs/CUSTOM_PRODUCT.md)

This plan is structured in two milestones. **Milestone 1** makes the database and cart infrastructure ready to accept and store custom items. **Milestone 2** wires custom items through the order creation pipeline. Each milestone is independently deployable and testable.

---

## Background

The client (`e-commerce`) sends a unified `POST /orders` body with items discriminated by `product_variant_type: "standard" | "custom"`. Custom items carry a `custom_design` object (v1.0.0 schema) instead of a `product_id`. The backend currently rejects this because:

1. `cart_items.product_id` is `NOT NULL` with a FK to `products(id)`.
2. `order_items.product_id` is `NOT NULL` with a FK to `products(id)`.
3. `AddItemUsecase` unconditionally calls `inventoryRepo.Reserve`.
4. `CreateOrderUsecase` unconditionally calls `pricingService.Calculate` on every item's `ProductID`.
5. No table exists to store the `custom_design` JSONB snapshot.

---

## Milestone 1 — Product & Cart: Schema + Custom Cart Item Support

> **Goal**: The database accepts custom items in `cart_items`. `POST /carts/items` with `product_variant_type: "custom"` stores the design snapshot and returns it in `GET /carts`. Existing standard item paths are unaffected.

---

### M1 · Step 1 — Database Schema

#### `migrations/0008_create_cart_items_table.up.sql` (Updated)

```sql
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    cart_id UUID NOT NULL,
    product_variant_type VARCHAR(32) NOT NULL DEFAULT 'standard', -- "standard" | "custom"
    product_id UUID,                                     -- nullable (NULL when product_variant_type = 'custom')
    shop_id UUID NOT NULL,

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    custom_design JSONB,                                 -- full v1.0.0 payload for custom items

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_cart
        FOREIGN KEY(cart_id)
        REFERENCES carts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_cart_item_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE RESTRICT,

    CONSTRAINT check_cart_items_type
        CHECK (
            (product_variant_type = 'standard' AND product_id IS NOT NULL)
            OR (product_variant_type = 'custom' AND product_id IS NULL AND custom_design IS NOT NULL)
        )
);

CREATE UNIQUE INDEX unique_standard_product_per_cart
ON cart_items(cart_id, product_id)
WHERE deleted_at IS NULL AND product_variant_type = 'standard';

CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_shop_id ON cart_items(shop_id);
```

#### `migrations/0032_create_order_items_table.up.sql` (Updated)

```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL,
    product_variant_type VARCHAR(32) NOT NULL DEFAULT 'standard', -- "standard" | "custom"

    shop_id UUID NOT NULL,
    shop_name VARCHAR(255) NOT NULL,

    product_id UUID,                                     -- nullable (NULL when product_variant_type = 'custom')
    product_name VARCHAR(255) NOT NULL,

    quantity INTEGER NOT NULL,
    unit_price BIGINT NOT NULL,
    subtotal BIGINT NOT NULL,

    courier_code VARCHAR(100),
    courier_service VARCHAR(100),
    shipping_fee_total BIGINT NOT NULL,

    CONSTRAINT fk_order_items_order_id
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_shop_id
        FOREIGN KEY (shop_id)
        REFERENCES shops(id),

    CONSTRAINT fk_order_items_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT check_order_items_type
        CHECK (
            (product_variant_type = 'standard' AND product_id IS NOT NULL)
            OR (product_variant_type = 'custom' AND product_id IS NULL)
        ),

    CONSTRAINT order_items_quantity_check CHECK (quantity > 0),
    CONSTRAINT order_items_unit_price_check CHECK (unit_price >= 0),
    CONSTRAINT order_items_subtotal_check CHECK (subtotal >= 0)
);
```

#### `migrations/0103_create_order_item_custom_designs_table.up.sql` (New)

```sql
CREATE TABLE order_item_custom_designs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_item_id UUID NOT NULL UNIQUE,
    version VARCHAR(16) NOT NULL DEFAULT '1.0.0',
    physical_size_id VARCHAR(64) NOT NULL,
    preview_url TEXT,

    header_text_upper VARCHAR(255),
    body_text_upper TEXT,
    header_text_lower VARCHAR(255),
    body_text_lower TEXT,

    design_snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_custom_design_order_item
        FOREIGN KEY (order_item_id)
        REFERENCES order_items(id)
        ON DELETE CASCADE
);
```

> [!NOTE]
> The rollback file (`0103_custom_product_support.down.sql`) must reverse every `ALTER TABLE` and drop the new index/table in reverse order.

---

### M1 · Step 2 — Domain Layer Updates

#### [MODIFY] [`cart/domain/cart_item.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/domain/cart_item.go)

Add `ProductVariantType`, make `ProductID` a pointer, add `CustomDesign` raw JSON, and add a `CartItemID` string for deduplication of custom items:

```go
package domain

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type ProductVariantType string

const (
    ProductVariantTypeStandard ProductVariantType = "standard"
    ProductVariantTypeCustom   ProductVariantType = "custom"
)

type CartItem struct {
    ID           uuid.UUID
    ProductVariantType     ProductVariantType
    ProductID    *uuid.UUID      // nil when ProductVariantType == ProductVariantTypeCustom
    ShopID       uuid.UUID
    Quantity     int
    CustomDesign json.RawMessage // populated when ProductVariantType == ProductVariantTypeCustom
    DeletedAt    *time.Time
}
```

#### [MODIFY] [`cart/domain/cart.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/domain/cart.go)

Add `AddCustomItem` and update existing helpers to be pointer-aware on `ProductID`:

```go
// AddCustomItem appends a new custom design item to the cart.
// Custom items are never deduplicated — each canvas design is independent.
func (c *Cart) AddCustomItem(shopID uuid.UUID, qty int, design json.RawMessage) error {
    if qty <= 0 {
        return ErrInvalidQuantity
    }
    c.Items = append(c.Items, CartItem{
        ID:           uuid.New(),
        ProductVariantType:     ProductVariantTypeCustom,
        ProductID:    nil,
        ShopID:       shopID,
        Quantity:     qty,
        CustomDesign: design,
    })
    return nil
}
```

Also update `RemoveItem`, `HasItem`, `FindItem`, and `HasProductInAnotherShop` to guard against nil `ProductID` (skip custom items in those checks).

---

### M1 · Step 3 — New Usecase: `AddCustomItemUsecase`

#### [NEW] [`cart/usecase/add_custom_item.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/usecase/add_custom_item.go)

This is a **dedicated usecase** to keep custom item logic cleanly separated from the existing `AddItemUsecase` (which must not be modified beyond minimal pointer-safety fixes):

```go
package usecase

import (
    "context"
    "encoding/json"
    "fmt"

    apperrors "service-core/internal/common/errors"
    "service-core/internal/modules/cart/domain"
    "service-core/internal/modules/cart/repository"
    transaction "service-core/internal/shared/transaction"

    "github.com/google/uuid"
)

type AddCustomItemInput struct {
    CustomerID   uuid.UUID
    ShopID       uuid.UUID
    Quantity     int
    ProductName  string          // Human-readable label, e.g. "Custom Board — Selamat & Sukses"
    PhysicalSizeID string        // e.g. "medium" — used for server-side price lookup in M2
    CustomDesign json.RawMessage // Full v1.0.0 payload from client
}

type AddCustomItemUsecase struct {
    executor   transaction.Executor
    transactor transaction.Transactor
    cartRepo   repository.CartRepository
}

func NewAddCustomItemUsecase(
    executor transaction.Executor,
    transactor transaction.Transactor,
    cartRepo repository.CartRepository,
) *AddCustomItemUsecase {
    return &AddCustomItemUsecase{executor: executor, transactor: transactor, cartRepo: cartRepo}
}

func (u *AddCustomItemUsecase) Execute(ctx context.Context, input AddCustomItemInput) error {
    if input.ShopID == uuid.Nil {
        return apperrors.NewInvalidInput(domain.ErrInvalidShopID.Error())
    }
    if input.Quantity <= 0 {
        return apperrors.NewInvalidInput(domain.ErrInvalidQuantity.Error())
    }
    if len(input.CustomDesign) == 0 {
        return apperrors.NewInvalidInput("custom_design is required")
    }

    cart, err := u.cartRepo.GetWithItemsByCustomerID(ctx, u.executor, input.CustomerID)
    if err != nil {
        return fmt.Errorf("failed to load cart: %w", err)
    }
    if cart == nil {
        cart, err = u.cartRepo.NewCart(ctx, u.executor, input.CustomerID)
        if err != nil {
            return fmt.Errorf("failed to create cart: %w", err)
        }
    }

    if err := cart.AddCustomItem(input.ShopID, input.Quantity, input.CustomDesign); err != nil {
        return apperrors.NewInvalidInput(err.Error())
    }

    return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
        if err := u.cartRepo.Save(ctx, exec, cart); err != nil {
            return fmt.Errorf("failed to save cart with custom item: %w", err)
        }
        return nil
    })
}
```

---

### M1 · Step 4 — Repository & Persistence Updates

#### [MODIFY] [`cart/infra/persistence/cart_repository_impl.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/infra/persistence/cart_repository_impl.go)

Three queries need updating:

**`GetWithItemsByCustomerID`** — extend SELECT to include `product_variant_type` and `custom_design`:
```sql
SELECT
    c.id, c.customer_id, c.created_at, c.updated_at,
    ci.id, ci.product_variant_type, ci.product_id, ci.shop_id, ci.quantity, ci.custom_design
FROM carts c
LEFT JOIN cart_items ci
    ON ci.cart_id = c.id AND ci.deleted_at IS NULL
WHERE c.customer_id = $1
ORDER BY ci.created_at
```

**`Save`** — split INSERT into two paths:

- *Standard path* (unchanged conflict resolution on `(cart_id, product_id)`).
- *Custom path* (always INSERT; never upsert — each design is unique):
```sql
-- For custom items:
INSERT INTO cart_items (cart_id, product_variant_type, shop_id, quantity, custom_design)
VALUES ($1, 'custom', $2, $3, $4)
```

**`RemoveItem`** — must accept `*uuid.UUID` for product_id; custom items should be removed by `cart_item.id` not `product_id`.

---

### M1 · Step 5 — HTTP Delivery Layer Updates

#### [MODIFY] [`cart/delivery/http/dto.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/dto.go)

```go
// Replaces addItemRequest
type addItemRequest struct {
    ProductVariantType     string          `json:"product_variant_type"` // "standard" | "custom"
    ProductID    *string         `json:"product_id,omitempty"`
    ShopID       string          `json:"shop_id"`
    Quantity     int             `json:"quantity"`
    ProductName  string          `json:"product_name,omitempty"`
    PhysicalSizeID string        `json:"physical_size_id,omitempty"`
    CustomDesign json.RawMessage `json:"custom_design,omitempty"`
}

// cartItemView extended with optional custom fields
type cartItemView struct {
    CartItemID   uuid.UUID            `json:"cart_item_id"`
    ProductVariantType     string               `json:"product_variant_type"`
    ProductID    *uuid.UUID           `json:"product_id,omitempty"`
    ShopID       uuid.UUID            `json:"shop_id"`
    Name         string               `json:"name"`
    Price        int64                `json:"price"`
    Subtotal     int64                `json:"subtotal"`
    Quantity     int                  `json:"quantity"`
    Image        productImageResponse `json:"images"`
    CustomDesign json.RawMessage      `json:"custom_design,omitempty"`
}
```

#### [MODIFY] [`cart/delivery/http/handler.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/handler.go)

Branch the `AddItem` handler on `product_variant_type`:
```go
func (h *Handler) AddItem(c echo.Context) error {
    var req addItemRequest
    if err := c.Bind(&req); err != nil { ... }

    switch req.ProductVariantType {
    case "custom":
        return h.addCustomItem(c, req, authCtx)
    default:
        return h.addStandardItem(c, req, authCtx)
    }
}
```

---

### M1 · Verification

```powershell
cd d:\__Projects\kage\chia.florist\service-core
go build ./...
go test ./internal/modules/cart/...
```

**Manual smoke test**:
1. `POST /carts/items` with `product_variant_type: "custom"` → `200 OK`.
2. `GET /carts` → response includes the custom item with `custom_design` JSONB echoed back.
3. `POST /carts/items` with `product_variant_type: "standard"` → unchanged behavior, no regression.

---
---

## Milestone 2 — Order: Custom Items Through Checkout & Creation

> **Goal**: `POST /checkout`, `POST /checkout/calculate`, and `POST /orders` accept shops that contain a mix of standard and custom items. Custom items bypass inventory checks and use server-resolved pricing from `physical_size_id`. The `order_item_custom_designs` table is populated atomically with the order.

**Pre-requisite**: Milestone 1 must be deployed and migrated.

---

### M2 · Step 1 — Shared Domain Type: `OrderItemCustomDesign`

#### [NEW] [`order/domain/order_item_custom_design.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/domain/order_item_custom_design.go)

```go
package domain

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

// OrderItemCustomDesign stores the full v1.0.0 design snapshot
// produced by the e-commerce canvas simulator.
// One record per custom order item; cascade-deleted with the parent order_item.
type OrderItemCustomDesign struct {
    ID              uuid.UUID
    OrderItemID     uuid.UUID
    Version         string
    PhysicalSizeID  string
    PreviewURL      *string         // CDN URL from Supabase Storage; nil if not uploaded
    HeaderTextUpper *string
    BodyTextUpper   *string
    HeaderTextLower *string
    BodyTextLower   *string
    DesignSnapshot  json.RawMessage // Full v1.0.0 payload
    CreatedAt       time.Time
}
```

#### [MODIFY] [`order/domain/order_item.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/domain/order_item.go)

Make `ProductID` a pointer and add `ProductVariantType`:

```go
type OrderItem struct {
    ID      uuid.UUID
    OrderID uuid.UUID

    ProductVariantType string  // "standard" | "custom"

    ShopID   uuid.UUID
    ShopName string

    ProductID   *uuid.UUID // nil when ProductVariantType == "custom"
    ProductName string

    Quantity  int
    UnitPrice int64
    Subtotal  int64

    CourierCode    *string
    CourierService *string
    ShippingFee    int64
}
```

---

### M2 · Step 2 — Pricing Service: Custom Item Price Resolution

#### [MODIFY] [`order/repository/query.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/repository/query.go)

Add a custom item pricing input type. The server resolves price from a **price table keyed by `physical_size_id`** (not from the client's `unit_price` — never trust client pricing):

```go
// CustomItemPricingInput is sent alongside standard PricingItemInput
// for items where product_variant_type == "custom".
type CustomItemPricingInput struct {
    PhysicalSizeID string
    ProductName    string
    Quantity       int
    CustomDesign   json.RawMessage
}

// Extend PricingShopInput
type PricingShopInput struct {
    ShopID         uuid.UUID
    CourierCode    *string
    CourierService *string
    Items          []PricingItemInput       // standard items
    CustomItems    []CustomItemPricingInput // custom items
}
```

The `PricingService.Calculate` implementation must:
1. Look up a `custom_product_prices` table (or a hard-coded config map for Phase 1) keyed by `physical_size_id` → `unit_price`.
2. Build `PricingItemResult` entries for custom items using the server-resolved price.
3. Include custom item subtotals in shop `Subtotal` and overall `Subtotal`.

> [!IMPORTANT]
> **Never use the `unit_price` sent by the client for custom items.** The price must be resolved server-side from `physical_size_id`.

**Custom price table (Phase 1 — static seed data)**:

| `physical_size_id` | Label | Price (IDR) |
| :--- | :--- | :--- |
| `small` | 1.2m × 0.8m | 350,000 |
| `medium` | 1.8m × 1.2m | 650,000 |
| `large` | 2.2m × 1.5m | 850,000 |
| `xlarge` | 2.5m × 2.0m | 1,200,000 |

#### [NEW] `migrations/0104_create_custom_product_sizes.up.sql`

```sql
CREATE TABLE custom_product_sizes (
    id          VARCHAR(64)  PRIMARY KEY,
    label       VARCHAR(128) NOT NULL,
    price       BIGINT       NOT NULL CHECK (price > 0),
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO custom_product_sizes (id, label, price) VALUES
    ('small',  '1.2m × 0.8m',  350000),
    ('medium', '1.8m × 1.2m',  650000),
    ('large',  '2.2m × 1.5m',  850000),
    ('xlarge', '2.5m × 2.0m', 1200000);
```

---

### M2 · Step 3 — Checkout Usecase Updates

#### [MODIFY] [`cart/usecase/checkout.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/usecase/checkout.go)

Extend `CheckoutItemInput` and `CheckoutShopInput`:

```go
type CheckoutCustomItemInput struct {
    PhysicalSizeID string
    ProductName    string
    Quantity       int
    CustomDesign   json.RawMessage
}

type CheckoutShopInput struct {
    ShopID      uuid.UUID
    Items       []CheckoutItemInput       // standard
    CustomItems []CheckoutCustomItemInput // custom
    Courier     *SelectedCourierInput
}
```

In `Execute`, forward `CustomItems` into `PricingShopInput.CustomItems` so `PricingService.Calculate` resolves them server-side.

Also extend `CheckoutItemResult` to optionally carry `CustomDesign` back to the client (for checkout preview rendering):

```go
type CheckoutItemResult struct {
    ProductVariantType    string          // "standard" | "custom"
    ProductID   *uuid.UUID
    ShopID      uuid.UUID
    Name        string
    Price       int64
    Quantity    int
    Subtotal    int64
    TotalWeight int
    CustomDesign json.RawMessage `json:",omitempty"`
}
```

---

### M2 · Step 4 — HTTP Checkout DTO Updates

#### [MODIFY] [`cart/delivery/http/dto.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/dto.go)

```go
type checkoutCustomItemRequest struct {
    PhysicalSizeID string          `json:"physical_size_id"`
    ProductName    string          `json:"product_name"`
    Quantity       int             `json:"quantity"`
    CustomDesign   json.RawMessage `json:"custom_design"`
}

// Extend checkoutShopRequest
type checkoutShopRequest struct {
    ShopID      string                      `json:"shop_id"`
    Items       []checkoutItemRequest       `json:"items"`
    CustomItems []checkoutCustomItemRequest `json:"custom_items,omitempty"`
    Courier     *selectedCourierRequest     `json:"courier"`
}

// Extend checkoutItemResponse
type checkoutItemResponse struct {
    ProductVariantType     string          `json:"product_variant_type"`
    ProductID    *uuid.UUID      `json:"product_id,omitempty"`
    ShopID       uuid.UUID       `json:"shop_id"`
    Name         string          `json:"name"`
    Price        int64           `json:"price"`
    Quantity     int             `json:"quantity"`
    Subtotal     int64           `json:"subtotal"`
    CustomDesign json.RawMessage `json:"custom_design,omitempty"`
}
```

---

### M2 · Step 5 — Order Creation Usecase Updates

#### [MODIFY] [`order/usecase/create_order.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/usecase/create_order.go)

**Input struct extension**:
```go
type OrderCustomItemInput struct {
    PhysicalSizeID string
    ProductName    string
    Quantity       int
    CustomDesign   json.RawMessage
}

type OrderShopInput struct {
    ShopID      uuid.UUID
    ShopName    string
    Courier     *OrderCourierInput
    Items       []OrderItemInput       // standard
    CustomItems []OrderCustomItemInput // custom
}
```

**In `Execute()`**:

1. **Pricing**: pass `CustomItems` into `PricingShopInput.CustomItems` alongside standard `Items`.

2. **Order item construction**: for each resolved custom item in `pricingResult.Shops[].Items`, set `ProductVariantType = "custom"`, `ProductID = nil`. Attach the original `CustomDesign` from the matching input.

3. **Inventory skip**: the existing `inventoryRepo.Reserve` loop must check `orderItem.ProductVariantType`:
   ```go
   for _, item := range orderItems {
       if item.ProductVariantType == "custom" {
           continue // no inventory to reserve
       }
       if err := u.inventoryRepo.Reserve(ctx, exec, *item.ProductID, item.ShopID, item.Quantity); err != nil {
           return fmt.Errorf("failed to reserve inventory: %w", err)
       }
   }
   ```

4. **Cart cleanup**: `cart.RemoveItem` currently matches on `ProductID`. Add a parallel `RemoveCustomItemByID` path on `CartItem.ID` so custom items are properly removed from the cart after order placement.

5. **Custom design persistence**: inside the DB transaction, after `orderItemRepo.SaveBulk`, save design records:
   ```go
   for _, item := range orderItems {
       if item.ProductVariantType == "custom" && item.CustomDesign != nil {
           design := domain.OrderItemCustomDesign{
               ID:             uuid.New(),
               OrderItemID:    item.ID,
               Version:        "1.0.0",
               PhysicalSizeID: item.PhysicalSizeID,
               PreviewURL:     extractPreviewURL(item.CustomDesign),
               DesignSnapshot: item.CustomDesign,
               CreatedAt:      now,
           }
           // extract header/body text from snapshot for quick staff access
           design.HeaderTextUpper = extractText(item.CustomDesign, "sections.upper.header.text")
           design.BodyTextUpper   = extractText(item.CustomDesign, "sections.upper.body.text")
           design.HeaderTextLower = extractText(item.CustomDesign, "sections.lower.header.text")
           design.BodyTextLower   = extractText(item.CustomDesign, "sections.lower.body.text")

           if err := u.customDesignRepo.Save(ctx, exec, design); err != nil {
               return fmt.Errorf("failed to save custom design: %w", err)
           }
       }
   }
   ```

---

### M2 · Step 6 — New Repository Interface & Implementation

#### [MODIFY] [`order/repository/interface.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/repository/interface.go)

```go
type OrderItemCustomDesignRepository interface {
    Save(ctx context.Context, exec transaction.Executor, design domain.OrderItemCustomDesign) error
    GetByOrderItemID(ctx context.Context, exec transaction.Executor, orderItemID uuid.UUID) (*domain.OrderItemCustomDesign, error)
    ListByOrderID(ctx context.Context, exec transaction.Executor, orderID uuid.UUID) ([]domain.OrderItemCustomDesign, error)
}
```

#### [NEW] `order/infra/persistence/order_item_custom_design_repository_impl.go`

Implement `Save` (INSERT into `order_item_custom_designs`) and `GetByOrderItemID`/`ListByOrderID` (SELECT with JOIN on `order_items.order_id`).

---

### M2 · Step 7 — Order Item Repository & GET Order Response

#### [MODIFY] `order/infra/persistence/order_item_repository_impl.go`

Update `SaveBulk` INSERT to write `product_variant_type` and handle nullable `product_id`:
```sql
INSERT INTO order_items (
    id, order_id, product_variant_type, shop_id, shop_name,
    product_id, product_name, quantity, unit_price, subtotal,
    courier_code, courier_service, shipping_fee_total
)
VALUES (...)
```

Update `ListByOrderID` SELECT to include `product_variant_type`.

#### [MODIFY] `order/delivery/http/` (handler & DTOs)

The `GET /orders/{id}` response items must expose `product_variant_type` and, for custom items, `custom_design` (or a pointer to `preview_url` for lightweight rendering):

```json
{
  "id": "...",
  "items": [
    {
      "product_variant_type": "standard",
      "product_id": "e7b0c950-...",
      "product_name": "Papan Bunga Grand Opening",
      ...
    },
    {
      "product_variant_type": "custom",
      "product_id": null,
      "product_name": "Custom Board — Selamat & Sukses",
      "unit_price": 650000,
      "custom_design": {
        "preview_url": "https://chia-florist.supabase.co/...",
        "physical_size_id": "medium",
        "sections": { "upper": { "header": { "text": "Selamat & Sukses" } } }
      }
    }
  ]
}
```

---

### M2 · Verification

```powershell
cd d:\__Projects\kage\chia.florist\service-core
go build ./...
go test ./internal/modules/order/... ./internal/modules/cart/...
```

**Manual integration flow**:
1. `POST /carts/items` (`product_variant_type: "custom"`) → cart contains custom item ✓
2. `POST /checkout` with `custom_items` array → response returns server-resolved price for `physical_size_id: "medium"` = IDR 650,000 ✓
3. `POST /orders` with mixed standard + custom items → HTTP 201; DB contains `order_items` row with `product_variant_type = 'custom'` and `product_id = NULL`; `order_item_custom_designs` row contains full `design_snapshot` ✓
4. `GET /orders/{id}` → custom item visible with `preview_url` and `product_name` ✓
5. `POST /orders` with standard items only → no regression ✓

---

## Cross-Milestone Notes

> [!IMPORTANT]
> **Dependency wiring**: `CreateOrderUsecase` and `CheckoutUsecase` constructors in `cmd/` (or the app bootstrap) must receive the new `customDesignRepo` and the updated `PricingService` implementation. Update DI wiring after each milestone.

> [!NOTE]
> **`previewBase64` in the order payload**: The client spec says `previewBase64` may be present in `custom_design.assets`. The backend should **ignore** `previewBase64` entirely — it is only used client-side for canvas rendering. The canonical image reference is `assets.previewUrl` (Supabase CDN). The backend stores only `previewUrl` in `order_item_custom_designs.preview_url`.

> [!WARNING]
> **`order/domain/invoice_item.go`** and its repository also need `product_variant_type` + nullable `product_id` parity with `order_items`. Extend these in M2 Step 5 alongside `OrderItem` to avoid invoice rendering issues on the control panel.
