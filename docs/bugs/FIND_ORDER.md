# Walkthrough - Order HTTP Handler Nil Pointer Dereference Fix
**06/08/2026 10:40:00**

Module: _**service-core**_

## Root Cause Summary
When listing or fetching orders containing custom product items (e.g. `GET /users/me/orders`), the HTTP handler panicked with:
```
panic: runtime error: invalid memory address or nil pointer dereference
  at order/delivery/http.buildOrderResponse(handler.go:49)
```

In `buildOrderResponse` ([handler.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go#L49)), `item.ProductID.String()` directly dereferenced `item.ProductID`. Because custom product items have `item.ProductID == nil`, calling `.String()` on a `nil` pointer threw a runtime panic.

## Code Fixes Implemented

### 1. HTTP Response DTO ([dto.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/dto.go))
Updated `orderItemResponse`:
- Changed `ProductID` from string to a pointer string (`*string`) with `json:"product_id,omitempty"`.
- Added `ProductVariantType string json:"product_variant_type"` and `IsCustom bool json:"is_custom"`.

### 2. HTTP Handler Response Builder ([handler.go](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go))
Updated `buildOrderResponse` to safely handle pointer dereferencing:
```go
var productIDStr *string
if item.ProductID != nil {
    s := item.ProductID.String()
    productIDStr = &s
}

variantType := string(item.ProductVariantType)
if variantType == "" {
    if item.ProductID == nil {
        variantType = "custom"
    } else {
        variantType = "standard"
    }
}

items[j] = orderItemResponse{
    ID:                 item.ID.String(),
    ProductID:          productIDStr,
    ProductVariantType: variantType,
    IsCustom:           item.ProductID == nil || variantType == "custom",
    ProductName:        item.ProductName,
    Quantity:           item.Quantity,
    UnitPrice:          item.UnitPrice,
    Subtotal:           item.Subtotal,
    ShopID:             item.ShopID.String(),
    ShopName:           item.ShopName,
    CourierCode:        item.CourierCode,
    CourierService:     item.CourierService,
    ShippingFeeTotal:   item.ShippingFee,
}
```
