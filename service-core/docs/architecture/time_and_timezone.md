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

### Infrastructure Layer & Database Schema Setup
- **PostgreSQL Session Timezone**: The DB connection pool (`internal/infra/db/connection.go`) executes `SET timezone TO 'Asia/Jakarta'` via `AfterConnect` callback on every `pgxpool` connection.
- **Docker Container Environment**: `docker-compose.yml` configures `PGTZ: Asia/Jakarta` and `TZ: Asia/Jakarta` for PostgreSQL and application containers.
- **Database Schema Standard (`TIMESTAMPTZ`)**: 100% of date/time columns in database migrations (`/migrations`) MUST specify `TIMESTAMPTZ` (Timestamp with time zone) instead of plain `TIMESTAMP`. This guarantees PostgreSQL stores UTC epoch internally while formatting and comparing time relative to `Asia/Jakarta` (`+07:00`).

### Unit Testing & Determinism
- Unit tests can mock the current time using `clock.NewMockClock(fixedTime)` or `clock.SetDefault(mockClock)`.
- `MockClock` allows setting fixed timestamps and manipulating time dynamically using `Advance(duration)`.

---

## Summary of Conventions
| Scope | Strategy | Example / API |
| :--- | :--- | :--- |
| Application Current Time | `clock.Now()` | `now := clock.Now()` |
| DB Connection Pool | Session Timezone | `SET timezone TO 'Asia/Jakarta'` |
| DB Table Columns | `TIMESTAMPTZ` | `created_at TIMESTAMPTZ DEFAULT NOW()` |
| Docker Container | Environment TZ | `PGTZ=Asia/Jakarta`, `TZ=Asia/Jakarta` |
| Incoming External Timestamps | Convert to WIB | `t = clock.InAppLocation(externalTimestamp)` |
| Outbound External Requests | Convert to provider spec | `refundKey := fmt.Sprintf("refund-%s-%d", id, clock.Now().Unix())` |
| Test Environment | Deterministic Mock | `mc := clock.NewMockClock(time.Date(...))` |
