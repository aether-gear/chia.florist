# Auth Session Isolation — Bug Report

**Date:** 2026-09-05  
**Reporter:** Antigravity (AGY)  
**Status:** Phase 1 complete · Phase 2.1 in progress  
**Severity:** High  
**Scope:** Staff (control-panel) vs. Customer (e-commerce) authentication isolation

---

## Overview

Two interrelated bugs affect the authentication system when a user is logged into both the e-commerce frontend and the control-panel simultaneously:

1. **Cross-logout** — logging out of one app can expire or disrupt the session of the other.
2. **Cookie sharing** — both frontends read each other's cookies when they share the same API origin.

> [!NOTE]
> See also [`docs/bugs/LOGIN_STATE.md`](./LOGIN_STATE.md) for the prior Phase 1 / Phase 2.1 analysis.
> This document captures the full picture discovered during the 2026-09-05 re-investigation.

---

## Bug 1 — Cross-Logout (Session Bleed on Logout)

### Symptoms

- User logs into e-commerce (`chast` cookie set) and control-panel (`hotpot` cookie set) in the same browser.
- Logging out of the e-commerce frontend also degrades or invalidates the control-panel session (or vice versa).
- Calls to `CoreAuth` routes after one logout silently pick up the *other* app's session.

### Root Cause Trace

#### `LogoutUsecase.Execute` — [`logout.go`](../../internal/modules/authentication/usecase/logout.go)

Revokes by `SessionID` only. This is correct for a single-app model but produces unpredictable effects when the logout handler clears the wrong set of cookies or when `CoreAuth` falls back to a still-live session.

```go
// Revokes session S1 when customer logs out
u.refreshTknRepo.RevokeBySessionID(ctx, exec, authCtx.SessionID)
u.sessionRepo.RevokeByID(ctx, exec, authCtx.SessionID)
```

#### `RequireMultiAuth` — [`jwt_authenticator.go`](../../internal/modules/authentication/infra/service/jwt_authenticator.go)

Used by `CoreAuth`. Collects all valid auth contexts from all provided cookies and sets `collected[0]` as the primary identity:

```go
ctx = domain.WithAuthContext(ctx, collected[0])  // first valid wins
```

After a customer logout (`chast` cleared), any subsequent `CoreAuth` request still carries `hotpot`. The middleware finds `chast` absent, falls through to `hotpot`, and authenticates the user as **staff** — even though they intended to be logged out of the customer session.

#### `Logout` / `LogoutStaff` handlers — [`handler.go`](../../internal/modules/authentication/delivery/http/handler.go)

Each handler only clears its own-app cookies:

```go
// Logout (customer) — does NOT clear hotpot / ladle
appcookie.Clear(w, appcookie.CookieCustomer)
appcookie.Clear(w, appcookie.CookieCustomerRefresh)

// LogoutStaff — does NOT clear chast / malkist
appcookie.Clear(w, appcookie.CookieStaff)
appcookie.Clear(w, appcookie.CookieStaffRefresh)
```

`ClearAll()` exists in `jar.go` but was removed from the logout handlers after Phase 2.1 analysis showed it caused unintended full-logout across both apps.

---

## Bug 2 — Cookie Sharing Between Frontends

### Symptoms

- E-commerce frontend receives the `hotpot`/`ladle` (staff) cookies set by the control-panel.
- Control-panel receives the `chast`/`malkist` (customer) cookies set by e-commerce.
- A fresh login to one app can be authenticated using the stale session of the other.

### Root Cause Trace

#### `Bind()` — [`jar.go`](../../internal/common/http/cookie/jar.go)

```go
http.SetCookie(w, &http.Cookie{
    Name:     string(name),
    Value:    value,
    Path:     "/",          // scoped to entire server root
    HttpOnly: true,
    Secure:   isSecureCookie(),
    SameSite: http.SameSiteLaxMode,
    Expires:  exp,
    // Domain: ""           // NOT SET — defaults to exact request host
})
```

When `Domain` is omitted, the cookie is **host-only** — it belongs to the exact host that set it. If both frontends call the same API host (e.g., `api.chia.florist`), the browser sends **all four cookies** with every request, regardless of which frontend initiated it.

#### `Session` struct — [`session.go`](../../internal/modules/authentication/domain/session.go)

```go
type Session struct {
    ID     uuid.UUID
    UserID uuid.UUID
    // No AccountType, no AppScope
}
```

Sessions carry no `AccountType` field. A staff session and a customer session for the same `UserID` are structurally identical in the DB — they can only be told apart by the JWT claims (`StaffID` vs `CustomerID`), not the session record itself.

#### `isValidCookieForAuth` — [`jwt_authenticator.go`](../../internal/modules/authentication/infra/service/jwt_authenticator.go)

```go
func isValidCookieForAuth(cookie appcookie.CookieName, authCtx *domain.AuthContext) bool {
    switch cookie {
    case appcookie.CookieCustomer:
        return authCtx.CustomerID != nil
    case appcookie.CookieStaff:
        return authCtx.StaffID != nil
    }
}
```

This check exists and is correct — it prevents a customer JWT from validating in a staff-cookie slot. However, it only helps when the right cookie name is passed to the middleware. On `Core`/`CoreAuth` routes, both cookie names are always passed, so this guard is bypassed.

---

## Cookie Map

| Cookie Name | App       | Purpose         |
|-------------|-----------|-----------------|
| `chast`     | e-commerce (customer) | Access token  |
| `malkist`   | e-commerce (customer) | Refresh token |
| `hotpot`    | control-panel (staff) | Access token  |
| `ladle`     | control-panel (staff) | Refresh token |

All four cookies are set with `Path: "/"` and no `Domain`. In a shared-origin deployment, all four are sent by the browser on every request to the API.

---

## Is This Intentional?

**No.** The distinct cookie names (`chast`, `hotpot`, etc.) imply the intention of app-level isolation. The current behavior is an emergent side-effect of a monolith API serving two separate frontends without cookie domain or subdomain separation.

---

## Prior Work

| Phase     | Status    | Summary |
|-----------|-----------|---------|
| Phase 1   | ✅ Done   | `isValidCookieForAuth` guard added; `ClearAll()` helper added to `jar.go`; login endpoints revised |
| Phase 2   | ⚠️ Partial | `RequireMultiAuth` introduced; `/auth/me` and `/auth/staff/me` separated; `CoreAuth` updated |
| Phase 2.1 | 🔄 Pending | `X-Account-Type` header disambiguation; `RequireAccountType` update in authorizer; route + handler cleanup |

See [`LOGIN_STATE.md`](./LOGIN_STATE.md) for detailed Phase 2.1 acceptance criteria and status matrix.

---

## Recommended Fix Strategy

See [`docs/implementation_plan.md`](../implementation_plan.md) for the full fix plan.

### Summary

1. **`X-Account-Type` header disambiguation** (Phase 2.1) — frontend sends `X-Account-Type: customer` or `staff`; middleware uses it to pick the right `AuthContext` from `MultiAuthContext` instead of defaulting to `collected[0]`.
2. **`RequireAccountType` update** — read `X-Account-Type` header to select the correct candidate from `MultiAuthContext`.
3. **Non-destructive login** — login endpoints must NOT clear cross-app cookies.
4. **Strict per-app logout** — each logout handler clears only its own cookies; `ClearAll` is removed from logout flows.
5. **Long-term: network isolation** — deploy staff API on a separate subdomain to achieve browser-level cookie isolation with zero application code overhead.
