# Bug Investigation & Resolution Report: `/carts/checkout/calculate` 400 Bad Request
**06/08/2026 10:15:00**

Module: _**service-core**_

## Symptom & Log
```json
2026-08-06 10:03:28 WARN request warning {
  "request_id": "b50fc1da-de3e-47f2-a34a-491796ac1ab5",
  "layer": "middleware",
  "method": "POST",
  "path": "/carts/checkout/calculate",
  "status_code": 400,
  "error": "invalid product id",
  "type": "BAD_REQUEST"
}
```

## Technical Root Cause Analysis

The issue was caused by a **2-part desynchronization** between the database schema, HTTP DTO response contract, and Pricing Service hydration:

```
[User Action: Add Custom Item to Cart]
                │
                ▼
[Database: cart_items table]
  • product_variant_type = 'custom'
  • product_id = NULL
  • custom_design = {"layout": {...}}
                │
                ▼
[Endpoint: POST /carts/checkout]
  • Returned DTO `checkoutItemResponse`:
    { "is_custom": true, "cart_item_id": "...", "product_id": null }
  • ❌ MISSING: `product_variant_type: "custom"` was omitted from JSON response DTO!
                │
                ▼
[Frontend Action: Calculate Checkout]
  • Re-sent item to `POST /carts/checkout/calculate`:
    { "cart_item_id": "...", "quantity": 1 }  (omitted `is_custom` and `product_variant_type`)
                │
                ▼
[Backend: cart/delivery/http/handler.go `parseCheckoutInput`]
  • Handlers checked `itemReq.ProductVariantType == "custom"`.
  • Since string was absent/defaulted, `isCustom` evaluated to `FALSE`.
  • Item had `product_id == null`.
  • Handler triggered `!isCustom && productID == nil` and threw:
    `400 BAD_REQUEST: "invalid product id"`
```

## Step-by-Step Fix

### 1. DTO Contract Fix (`product_variant_type` Sync)
Added `product_variant_type` (`"custom"` or `"standard"`) to `checkoutItemResponse` in [dto.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/dto.go) and populated it in [handler.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/handler.go):
```go
type checkoutItemResponse struct {
    ProductID          *uuid.UUID `json:"product_id,omitempty"`
    CartItemID         *uuid.UUID `json:"cart_item_id,omitempty"`
    ProductVariantType string     `json:"product_variant_type"` // "custom" | "standard"
    IsCustom           bool       `json:"is_custom"`
    ...
}
```

### 2. Pricing Service Database Hydration
Injected `CartRepository` into `PricingService` ([pricing.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/infra/service/pricing.go) & [container.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/bootstrap/container.go)).
When `PricingService.Calculate()` receives an item with `cart_item_id`, it queries PostgreSQL to fetch the saved cart item and automatically hydrates `is_custom = true`, `product_variant_type = "custom"`, and `custom_design` JSON directly from the database.

### 3. Handler Parsing Fallback
Updated `parseCheckoutInput` ([cart/delivery/http/handler.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/cart/delivery/http/handler.go)):
- Items with `cart_item_id != nil` or `product_id == nil` are recognized as custom items.
- Only items that are standard catalog items without any product ID or cart item ID return `400 BAD_REQUEST: "invalid product id"`.

## Verification
- Unit test suite passed cleanly (`go test ./...`).
- `/carts/checkout/calculate` now returns `200 OK` with full server-calculated shipping, custom item pricing, and payment options.
