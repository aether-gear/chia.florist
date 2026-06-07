# Better Auth Pipeline Implementation Plan

## Analysis of Current State

I have reviewed the current state of your `auth` module under `internal/modules/auth/`. Here is what you currently have:
- **Core Domain**: You have a basic `Account` model and interfaces for password hashing and token generation.
- **Usecases**: `LoginEmail` and `Register` usecases that handle creating an account and verifying credentials to generate an access token.
- **Token Service**: A `JWTService` that generates a single JWT string containing the `user_id` as a claim.
- **HTTP Handlers**: Registration and Login handlers are set up. There is a `GetByID` handler that expects a `user_id` in the context, but no middleware exists to parse tokens and inject this context.

### What is missing / needs improvement:
1. **Refresh Tokens**: Currently, you only generate a single access token. If you make it short-lived (secure), the user has to log in frequently. If you make it long-lived (convenient), it's a security risk because JWTs are stateless and hard to revoke.
2. **Session Revocation (Logout)**: You cannot easily log a user out across devices because JWTs cannot be invalidated without a stateful check.
3. **Auth Middleware**: To secure other modules, you need a reusable middleware that intercepts requests, validates the JWT, and extracts the `user_id`.

---

## User Review Required

> [!WARNING]
> Adding refresh tokens will require changing the response schema of the login endpoint (adding a `refresh_token` field) and adding a new endpoint `/auth/refresh`. Does this fit well with your frontend/client expectations?

> [!IMPORTANT]
> To handle logout and token revocation, I propose creating a `Session` table in your database. This table will track active refresh tokens. When a user logs out, we delete their session. Is this acceptable, or would you prefer a stateless approach (where we can't easily force logouts)?

## Proposed Changes

### 1. Database & Domain (Sessions)
We need to track refresh tokens to allow revocation.
#### [NEW] `internal/modules/auth/domain/session.go`
- Define a `Session` struct containing `ID`, `AccountID`, `RefreshToken`, `UserAgent`, `ClientIP`, `IsBlocked`, `ExpiresAt`, `CreatedAt`.
- Add a `SessionRepository` interface to handle saving and retrieving sessions.

#### [NEW] Database Migration
- Create a migration to add the `sessions` table in the database.

### 2. Token Service Updates
#### [MODIFY] `internal/modules/auth/domain/token.go`
- Change `TokenService.Generate` to return two tokens: `AccessToken` and `RefreshToken`.

#### [MODIFY] `internal/modules/auth/infra/service/token_service.go`
- Implement the generation of both Access (short-lived, e.g., 15 mins) and Refresh (long-lived, e.g., 7 days) tokens.

### 3. Usecases
#### [MODIFY] `internal/modules/auth/usecase/login_email.go`
- Update to return both access and refresh tokens.
- Save the session in the database (via `SessionRepository`).

#### [NEW] `internal/modules/auth/usecase/refresh_token.go`
- A usecase that takes a refresh token, validates it against the `Session` database table, and issues a new access/refresh token pair.

#### [NEW] `internal/modules/auth/usecase/logout.go`
- A usecase that takes a session ID or refresh token and deletes/invalidates the session in the database.

### 4. HTTP Delivery & Middleware
#### [NEW] `internal/common/middleware/auth.go`
- Create an `AuthMiddleware` that reads the `Authorization: Bearer <token>` header.
- Uses `TokenService.Validate` to check the access token.
- If valid, injects `user_id` into the `http.Request` context. This middleware can be used on **any** handler in your project.

#### [MODIFY] `internal/modules/auth/delivery/http/handler.go`
- Update `SignInEmail` response to include `refresh_token`.
- Add `RefreshToken` endpoint handler.
- Add `Logout` endpoint handler.

---

## What you should avoid

1. **Long-lived Access Tokens**: Avoid relying on a single JWT that lasts for hours or days.
2. **Skipping Token Revocation**: Don't implement auth without a way to revoke access (e.g. if an account is compromised).
3. **Putting too much data in the JWT Payload**: Keep the JWT payload small (just `user_id` and maybe `role`). Don't store email, name, or other sensitive details in the token itself.

## Verification Plan

### Automated/Manual Verification
- Create a user and log in to receive `access_token` and `refresh_token`.
- Make an authenticated request using the `access_token` and the new `AuthMiddleware` to the `GetByID` endpoint.
- Wait for the access token to expire (or simulate it).
- Use the `/auth/refresh` endpoint with the `refresh_token` to get a new pair.
- Call the `/auth/logout` endpoint and verify the session is deleted.
- Attempt to refresh again and verify it is rejected.
