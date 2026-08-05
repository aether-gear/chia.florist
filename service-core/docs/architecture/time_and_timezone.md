# Time Management & Timezone Strategy

## Overview
This document outlines the application's canonical time and timezone management strategy. To ensure consistency across business rules, timestamp comparisons, order and payment expirations, audit logging, and third-party integrations, the application enforces a single canonical business timezone and a centralized time provider.

---

## Canonical Business Timezone
- **Location**: `Asia/Jakarta` (WIB, Western Indonesia Time, UTC+7)
- **Fallback**: Fixed Zone `WIB` (+25200 seconds / +07:00 offset) in environments without IANA tzdata installed.
- **Embedded Tzdata**: The application embeds `time/tzdata` to guarantee consistent IANA timezone lookups across platforms (Linux, macOS, Windows, Docker containers).

---

## Centralized Clock (`internal/common/clock`)

All components in domain, usecase, delivery, and infrastructure layers MUST retrieve current time through the centralized clock package rather than calling `time.Now()` or `time.Now().UTC()` directly.

### Key API Methods
- **`clock.Now() time.Time`**: Returns current wall-clock time with its location set to `Asia/Jakarta`.
- **`clock.Location() *time.Location`**: Returns the `Asia/Jakarta` location reference.
- **`clock.InAppLocation(t time.Time) time.Time`**: Converts any `time.Time` struct to `Asia/Jakarta` location.
- **`clock.ToUTC(t time.Time) time.Time`**: Converts a `time.Time` struct to UTC (for external provider requests requiring UTC).
- **`clock.BeginningOfDay(t time.Time)` / `clock.EndOfDay(t time.Time)`**: Helper functions to obtain exact day boundaries in WIB.

---

## Architecture Boundaries & Adapter Rules

### Domain & Usecase Layers
- Domain models and usecase business logic MUST use `clock.Now()` for stamping timestamps (`CreatedAt`, `UpdatedAt`, `ConfirmedAt`, `ExpiresAt`, `HandlingExpiresAt`).
- No business logic or domain code may call `time.Now()` or `time.Now().UTC()`.

### Infrastructure Adapters & External Gateways
- **Inbound Data**: Webhooks, HTTP responses, or database drivers returning external timestamps MUST normalize those timestamps into the canonical `Asia/Jakarta` timezone using `clock.InAppLocation(t)` before handing data to usecases or domain entities.
- **Outbound Data**: Adapters communicating with external APIs (e.g. Midtrans Payment Gateway, logistics providers) are responsible for converting canonical timestamps into provider-required formats (e.g. UTC, ISO-8601 strings, Unix Epoch seconds) at the infrastructure boundary.

### Unit Testing & Determinism
- Unit tests can mock the current time using `clock.NewMockClock(fixedTime)` or `clock.SetDefault(mockClock)`.
- `MockClock` allows setting fixed timestamps and manipulating time dynamically using `Advance(duration)`.

---

## Summary of Conventions
| Scope | Strategy | Example / API |
| :--- | :--- | :--- |
| Application Current Time | `clock.Now()` | `now := clock.Now()` |
| Incoming External Timestamps | Convert to WIB | `t = clock.InAppLocation(externalTimestamp)` |
| Outbound External Requests | Convert to provider spec | `refundKey := fmt.Sprintf("refund-%s-%d", id, clock.Now().Unix())` |
| Test Environment | Deterministic Mock | `mc := clock.NewMockClock(time.Date(...))` |
