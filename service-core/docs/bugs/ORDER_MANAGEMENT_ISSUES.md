# Order Management Module - Issue Report

> **Date:** 2026-08-19  
> **Module:** `internal/modules/order` & related authorization/delivery handlers  
> **Status:** Open for Resolution  

---

## Executive Summary

This document details 6 technical issues and bugs identified during the architecture and security audit of the Order Management system in `service-core`.

| Issue ID | Severity | Area / Layer | Short Description | Impact |
|---|---|---|---|---|
| **[BUG-1](#bug-1--redundant-status-field-population)** | 🔴 Low (Latent) | Use Case / Repo (`find_orders.go`) | `Status` field is passed alongside `Statuses` to `FindOrderParams` | Code smell, risk of future filter double-application |
| **[BUG-2](#bug-2--non-deterministic-rule-lookup-in-multi-shop-orders)** | 🔴 Medium | Delivery Handler (`handler.go`) | `UpdateOrderStatus` applies rules from only the first matched shop in slice | Rule bypass potential in multi-shop orders |
| **[BUG-3](#bug-3--missing-orderread-permission-check-on-staff-tracking-endpoint)** | 🟡 Medium | Delivery Handler (`handler.go`) | `GetOrderTrackingForStaff` missing shop-level `order:read` check | Information disclosure across unauthorized shop boundaries |
| **[BUG-4](#bug-4--unused-statuses-slice-in-findordersinput-handler)** | 🟡 Low | Delivery Handler (`handler.go`) | Handler query parsing assigns `?statuses=` to `input.Status` instead of `input.Statuses` | Dead DTO field / API design confusion |
| **[BUG-5](#bug-5--external-logistics-call-executed-outside-db-transaction)** | 🟡 Medium | Use Case (`update_order_status.go`) | Komerce external API call is outside the DB transaction without rollback/compensation | Out-of-sync data on DB failure |
| **[MINOR-1](#minor-1--silent-no-op-in-resolveshopfilter-when-getshop-is-nil)** | 🟢 Low | Delivery Handler (`handler.go`) | `resolveShopFilter` returns `(nil, false, nil)` if `getShop` dependency is missing | Silent failure to filter by shop slug |

---

## Detailed Bug Breakdown

---

### BUG-1 — Redundant `Status` Field Population

* **Severity:** 🔴 Low (Latent Bug)
* **File:** [`internal/modules/order/usecase/find_orders.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/usecase/find_orders.go#L140-L153)
* **Location:** Lines 140–153 (`Execute` method)

#### Description
When building `repository.FindOrderParams`, the usecase sets **both** `Status: input.Status` (raw string pointer) and `Statuses: statuses` (parsed string slice).

```go
params := repository.FindOrderParams{
    ID:         input.ID,
    Number:     input.Number,
    CustomerID: input.CustomerID,
    ShopID:     input.ShopID,
    ShopIDs:    input.ShopIDs,
    Status:     input.Status,   // <--- Raw pointer still set
    Statuses:   statuses,       // <--- Slice also set
    ...
}
```

#### Impact
While `order_repository_impl.go` currently checks `len(params.Statuses) > 0` before checking `params.Status`, leaving both populated creates ambiguity and risk of future bugs if repository layer logic changes.

#### Remediation
Set `Status` to `nil` when `Statuses` is populated, or standardize on using `Statuses` exclusively in `FindOrderParams`.

---

### BUG-2 — Non-Deterministic Rule Lookup in Multi-Shop Orders

* **Severity:** 🔴 Medium
* **File:** [`internal/modules/order/delivery/http/handler.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go#L720-L747)
* **Location:** Lines 720–747 (`UpdateOrderStatus` method)

#### Description
When a staff member updates an order status, the handler checks if they have `order:update_status` permission for any shop associated with the order items. However, it saves only the **first** matched shop ID to `matchedShopID` and evaluates staff rules (`allowed_statuses`, `max_order_amount`) solely against that shop:

```go
var matchedShopID uuid.UUID
for _, item := range existingOrder.Items {
    if actor.HasPermission(item.ShopID, authzDomain.PermissionOrderUpdateStatus) {
        hasAccess = true
        matchedShopID = item.ShopID  // <--- Only grabs the first matched item's shop
        break
    }
}

// Rules evaluated only against matchedShopID
if shopRules, exists := actor.Rules[matchedShopID]; exists && shopRules != nil {
    ...
}
```

#### Impact
If an order contains items from multiple shops (e.g. Shop A and Shop B) and staff member has different rules for each shop, the rule enforced is non-deterministic (dependent on DB slice iteration order). Staff could bypass stricter shop rules if a lenient shop item appears first.

#### Remediation
Iterate through **all** shops for which the staff member has permission in the order, and apply the most restrictive rules (e.g. intersection of `allowed_statuses` and minimum `max_order_amount`).

---

### BUG-3 — Missing `order:read` Permission Check on Staff Tracking Endpoint

* **Severity:** 🟡 Medium
* **File:** [`internal/modules/order/delivery/http/handler.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go#L849-L903)
* **Location:** Lines 849–903 (`GetOrderTrackingForStaff` method)

#### Description
`GetOrderTrackingForStaff` verifies that the caller has a staff account (`actor.Type == AccountTypeStaff`), but fails to verify whether the staff member has `order:read` permission for the specific shop(s) associated with that order.

In contrast, `GetOrder` (L383–394) properly validates:
```go
for _, item := range result.Items {
    if actor.HasPermission(item.ShopID, authzDomain.PermissionOrderRead) {
        hasAccess = true
        break
    }
}
```

#### Impact
Any staff account can query shipment and tracking timelines for orders belonging to shops they do not have authorization to view.

#### Remediation
Add shop-level `order:read` permission validation to `GetOrderTrackingForStaff` before returning tracking details.

---

### BUG-4 — Unused `Statuses` Slice in `FindOrdersInput` Handler

* **Severity:** 🟡 Low
* **File:** [`internal/modules/order/delivery/http/handler.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go#L286-L291)
* **Location:** Lines 286–291 (`FindOrders` method)

#### Description
In `FindOrders`, query parameters `status` and `statuses` are parsed as follows:

```go
if status != "" {
    input.Status = &status
} else if statuses := apphttp.Query(r, "statuses"); statuses != "" {
    input.Status = &statuses   // <--- Assigns to input.Status instead of input.Statuses
}
```

#### Impact
`FindOrdersInput.Statuses []string` is never populated by the handler. Comma-splitting works transitively inside the usecase (`find_orders.go`), but the DTO field remains dead code in the HTTP handler.

#### Remediation
Explicitly parse `statuses` into `input.Statuses []string` when `?statuses=` query parameter is supplied.

---

### BUG-5 — External Logistics Call Executed Outside DB Transaction

* **Severity:** 🟡 Medium
* **File:** [`internal/modules/order/usecase/update_order_status.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/usecase/update_order_status.go#L416-L453)
* **Location:** Lines 416–453 (`Execute` method during status transition to `shipped`)

#### Description
When moving an order to `shipped`, the logistics API (`u.logistics.CreateOrder`) is called to create external shipments before entering the database transaction:

```go
komerceResult, err := u.logistics.CreateOrder(ctx, orderInput)
...
err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
    // Persist shipment and update order status
})
```

#### Impact
If the DB transaction fails after `logistics.CreateOrder` succeeds (e.g. DB network timeout, deadlock), a shipment order remains created at the logistics provider (Komerce), but no local shipment record exists in `service-core`.

#### Remediation
Add a rollback/compensation call (e.g. cancel logistics shipment) in case the DB transaction fails, mirroring the compensation pattern implemented in `create_order.go` (`paymentGateway.CancelTransaction`).

---

### MINOR-1 — Silent No-Op in `resolveShopFilter` When `getShop` is `nil`

* **Severity:** 🟢 Low
* **File:** [`internal/modules/order/delivery/http/handler.go`](file:///d:/__Projects/kage/chia.florist/service-core/internal/modules/order/delivery/http/handler.go#L220-L221)
* **Location:** Lines 220–221 (`resolveShopFilter` method)

#### Description
If `h.getShop` is `nil` (due to missing dependency injection), resolving a shop filter via `shop_slug` returns `(nil, false, nil)` without logging or throwing an error.

#### Impact
Queries filtering by shop slug will silently ignore the shop filter and return unfiltered order results.

#### Remediation
Return an error or log a warning when `h.getShop` is uninitialized during a slug lookup attempt.
