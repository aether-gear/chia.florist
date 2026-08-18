# Foreign Key Constraint Fix (`fk_order_items_product_id`) - Legacy Code
**06/08/2026 10:29:00**

Module: _**service-core**_

## Symptom & Log
```json
2026-08-06 10:19:48 ERROR request failed {
  "request_id": "d6906334-8a1f-4bb3-ad65-470584289614",
  "layer": "middleware",
  "method": "POST",
  "path": "/order",
  "status_code": 500,
  "error": "failed to save order items: query to save order item: ERROR: insert or update on table \"order_items\" violates foreign key constraint \"fk_order_items_product_id\" (SQLSTATE 23503)",
  "type": "INTERNAL_ERROR"
}
```

## Root Cause Summary
When placing a custom product order (`POST /order`), `order_items` inserts failed with:
`ERROR: insert or update on table "order_items" violates foreign key constraint "fk_order_items_product_id" (SQLSTATE 23503)`

- **Legacy Go Code**: `domain.OrderItem.ProductID` was typed as `uuid.UUID` (non-pointer). Custom items without a catalog product ID evaluated to `uuid.Nil` (`"00000000-0000-0000-0000-000000000000"`).
- **SQL Persistence**: `order_item_repository_impl.go` omitted `product_variant_type` from the SQL INSERT statement and sent `"00000000-0000-0000-0000-000000000000"` as `product_id`. PostgreSQL attempted to look up `"00000000-..."` in the `products` table and threw SQLSTATE 23503.

## SQL Migration & Database ALTER Queries

### 1. Migration File Updates
Updated [0034_create_invoice_items_table.up.sql](file:///D:/__Projects/kage/chia.florist/service-core/migrations/0034_create_invoice_items_table.up.sql) to allow nullable `product_id` and include `product_variant_type`:
```sql
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    product_variant_type product_variant_type NOT NULL DEFAULT 'standard',
    shop_id UUID NOT NULL,
    shop_name VARCHAR(255) NOT NULL,
    product_id UUID, -- Nullable for custom items
    product_name VARCHAR(255) NOT NULL,
    ...
```

### 2. ALTER Query for Existing Database Environments
Run the following SQL statements on existing database environments:

```sql
-- 1. Ensure product_id is nullable for custom invoice items
ALTER TABLE invoice_items ALTER COLUMN product_id DROP NOT NULL;

-- 2. Add product_variant_type column to invoice_items if missing
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'invoice_items' AND column_name = 'product_variant_type'
    ) THEN
        ALTER TABLE invoice_items ADD COLUMN product_variant_type product_variant_type NOT NULL DEFAULT 'standard';
    END IF;
END $$;
```

## Code Changes

### Domain Layer
- **[order_item.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/domain/order_item.go)** & **[invoice_item.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/domain/invoice_item.go)**:
  - Updated `ProductID` to `*uuid.UUID` (nullable pointer).
  - Added `ProductVariantType cartDomain.ProductVariantType`.

### Persistence Layer
- **[order_item_repository_impl.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/infra/persistence/order_item_repository_impl.go)** & **[invoice_item_repository_impl.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/infra/persistence/invoice_item_repository_impl.go)**:
  - Included `product_variant_type` in `INSERT`, `ON CONFLICT DO UPDATE`, and `SELECT` SQL queries.
  - Passes SQL `NULL` for `product_id` when `item.ProductID == nil`.

### Usecase Layer
- **[create_order.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/usecase/create_order.go)**:
  - Populates `ProductVariantType` (`'custom'` or `'standard'`).
  - Skips catalog inventory reservations and cart removal for custom products where `ProductID == nil`.
- **[process_payment_webhook.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/payment/usecase/process_payment_webhook.go)**, **[sync_pending_payments.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/payment/usecase/sync_pending_payments.go)**, **[expire_past_due_payments.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/payment/usecase/expire_past_due_payments.go)**, & **[update_order_status.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/usecase/update_order_status.go)**:
  - Added `if item.ProductID == nil { continue }` checks for inventory commit/release operations.

## Verification
- Executed full test suite: `go test ./...` -> **PASS**
