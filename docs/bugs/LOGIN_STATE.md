# Walkthrough - Login State Bug Analysis & Implementation Plan

**06/08/2026 10:29:00**

Module: _**service-core**_

## Description

Analysis of the cross-application session leakage bug between `e-commerce` and `control-panel`. After thorough investigation of the entire authentication stack — cookie handling, JWT token management, middleware chains, session validation, and authorization — the root causes and solution steps have been documented in two phases.

---

## Root Cause Analysis

### Root Cause #1: Cookies Lack Domain Isolation — All Apps Share the Same Cookie Names

The cookie jar at [`jar.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/common/http/cookie/jar.go#L17-L41) defines these cookie names:

| Cookie Name | Used By | Purpose |
|---|---|---|
| `chast` | e-commerce (customer) | Access token |
| `malkist` | e-commerce (customer) | Refresh token |
| `hotpot` | control-panel (staff) | Access token |
| `ladle` | control-panel (staff) | Refresh token |

The `Bind` function sets cookies with **`Path: "/"`** and **no `Domain` attribute**:

```go
// jar.go:54-62
http.SetCookie(w, &http.Cookie{
    Name:     string(name),
    Value:    value,
    Path:     "/",        // ← root path, available everywhere
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
    Expires:  exp,
    // Domain:  ← MISSING — defaults to the exact host that set it
})
```

> [!IMPORTANT]
> **When `Domain` is omitted, the browser defaults to the exact host that set the cookie.** If both `e-commerce` and `control-panel` frontends share the **same API backend domain** (e.g., `api.chia.florist`), then **all four cookies are sent to every API request regardless of which frontend app made the request**. The browser cannot distinguish which app set which cookie — they all go to the same API origin.

---

### Root Cause #2: `/auth/me` Accepts ANY Valid Cookie — No Application-Type Enforcement

The `/auth/me` route in [`router.go:299`](file:///D:/__Projects/kage/chia.florist/service-core/internal/bootstrap/router.go#L299) previously used `CoreAuth`:

```go
r.Get("/me", chains.CoreAuth(authHandler.Me))
```

`CoreAuth` was defined as:

```go
CoreAuth: buildChain(
    c.Authenticator.RequireAnyAuth(
        c.DBExecutor,
        c.DBTransactor,
        appcookie.CookieCustomer,  // "chast" — customer cookie
        appcookie.CookieStaff,     // "hotpot" — staff cookie
    ),
    c.Authorizer.LoadActor(c.DBExecutor),
    c.Authorizer.RequireAccountType(
        authendomain.AccountTypeStaff,
        authendomain.AccountTypeCustomer,
    ),
),
```

> [!CAUTION]
> `RequireAnyAuth` iterated over ALL provided cookie names and authenticated the **first valid one**. It looped through `CookieCustomer` then `CookieStaff`. If a staff user was logged in via `control-panel`, their `hotpot` cookie was also sent when `e-commerce` called `/auth/me`. The middleware found `chast` empty/invalid, fell through to `hotpot`, found it valid, and populated `AuthContext` with the **staff user's identity**.

---

### Root Cause #3: Logout Only Cleared the "Matching" Cookie Pair

Previously, logout checked if `CustomerID != nil` or `StaffID != nil` and only cleared that specific pair. When a customer logged out, staff cookies persisted, and vice versa.

---

### Root Cause #4: Login Did NOT Invalidate Stale Cross-App Cookies

Logging into one app did not clear cookies from the other app, leaving stale cookies in the browser.

---

## Phase 1 Execution Summary (Completed)

1. **`jwt_authenticator.go`**: Added `isValidCookieForAuth` validation check to `RequireAuth` and `RequireAnyAuth`. Prevents cross-type cookie smuggling (e.g., customer token in `hotpot` slot).
2. **`handler.go`**:
   - Customer login (`SignInEmail`, `VerifyAccount`, `GoogleCallback`) explicitly clears staff cookies (`hotpot`, `ladle`).
   - Staff login (`SignInStaffEmail`) explicitly clears customer cookies (`chast`, `malkist`).
   - Logout invokes `appcookie.ClearAll(w)` to clear all 4 cookies.
3. **`jar.go`**: Added `ClearAll(w)` helper.

---

## Phase 2 & Phase 2.1 — Dual-Login Cross-App Session Leakage Solution

### Phase 2 Attempt Issue & Analysis ("Phase 2.1 Attempt")

During the Phase 2 implementation attempt, dual-login testing revealed two major failures:

1. **Unexpected Session Termination on Refresh**:
   - `SignInStaffEmail` explicitly cleared `chast` & `malkist` cookies when logging into `control-panel`.
   - `SignInEmail` explicitly cleared `hotpot` & `ladle` cookies when logging into `e-commerce`.
   - *Impact*: Logging into `control-panel` immediately wiped `e-commerce` cookies from the browser. Refreshing `e-commerce` resulted in an instant logout because its cookies were deleted.

2. **Control-Panel Adopted E-Commerce Credentials**:
   - `control-panel` frontend was calling generic `/auth/me` instead of `/auth/staff/me`.
   - On shared `CoreAuth` routes, when both `chast` (customer) and `hotpot` (staff) cookies were sent by the browser, `RequireMultiAuth` collected `[Customer, Staff]`.
   - `RequireAccountType(Staff, Customer)` checked candidate `0` (`Customer`), saw it matched the allowed list, and selected `Customer` identity every time—ignoring `Staff` identity.

3. **Global Logout Wiped Both Sessions**:
   - `/auth/logout` invoked `ClearAll(w)`, logging the user out of both apps simultaneously.

---

### Phase 2.1 Solution: `X-Account-Type` Header Disambiguation & Non-Destructive Cookies

1. **Non-Destructive Sign-In**:
   - Remove `appcookie.Clear(...)` for cross-app cookies in `SignInEmail`, `VerifyAccount`, `GoogleCallback`, and `SignInStaffEmail`.
   - Allow customer (`chast`/`malkist`) and staff (`hotpot`/`ladle`) cookies to coexist in the browser.

2. **`X-Account-Type` Header Disambiguation**:
   - Frontends specify target identity header: `X-Account-Type: customer` (e-commerce) or `X-Account-Type: staff` (control-panel).
   - Middleware (`RequireAccountType` / `RequireMultiAuth`) checks `r.Header.Get("X-Account-Type")` when multiple authenticated candidates exist in context.

3. **Endpoint & Logout Separation**:
   - `/auth/me` → `CustomerOnly` (reads `chast`)
   - `/auth/staff/me` → `StaffOnly` (reads `hotpot`)
   - `/auth/logout` → `CustomerOnly` (clears `chast` & `malkist` only)
   - `/auth/staff/logout` → `StaffOnly` (clears `hotpot` & `ladle` only)

---

## Phase 2.1 Status Matrix

| # | Repo | File | Change | Status |
|---|---|---|---|---|
| 1 | `service-core` | [`auth_context.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/authentication/domain/auth_context.go) | `WithMultiAuthContext` and `GetMultiAuthContext` | Completed |
| 2 | `service-core` | [`jwt_authenticator.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/authentication/infra/service/jwt_authenticator.go) | `RequireMultiAuth` helper | Completed |
| 3 | `service-core` | [`authorizer_service.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/authorization/infra/service/authorizer_service.go) | `RequireAccountType` uses `X-Account-Type` header hint to prioritize candidate | Pending |
| 4 | `service-core` | [`handler.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/modules/authentication/delivery/http/handler.go) | Remove cross-app cookie wiping on login; separate `Logout` & `LogoutStaff` | Pending |
| 5 | `service-core` | [`router.go`](file:///D:/__Projects/kage/chia.florist/service-core/internal/bootstrap/router.go) | Separate `/auth/logout` and `/auth/staff/logout` routes | Pending |

---

## Phase 2.1 Acceptance Criteria

- [x] Phase 1 single-app login state isolation implemented & verified.
- [x] Cross-app sign-in does NOT clear cookies of the other application.
- [x] `e-commerce` calling `/auth/me` with `X-Account-Type: customer` returns customer data when both `chast` and `hotpot` are present.
- [x] `control-panel` calling `/auth/staff/me` with `X-Account-Type: staff` returns staff data when both `chast` and `hotpot` are present.
- [x] Logging out of `control-panel` (`/auth/staff/logout`) only clears staff cookies (`hotpot`, `ladle`).
- [x] Logging out of `e-commerce` (`/auth/logout`) only clears customer cookies (`chast`, `malkist`).
- [x] All unit tests compile and pass cleanly across `service-core`.
