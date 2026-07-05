# Logging Standards

> Reference document for all contributors. Covers the structured logging
> conventions used in this project, including required fields, layer
> responsibilities, helper usage, and sensitive field redaction rules.

---

## Architecture Overview

All logs share a single correlation anchor — the `request_id` — generated
once per HTTP request and propagated through context automatically.

```
Request arrives
    │
    ▼
Logging middleware          → layer=middleware  request_id stamped onto ctx
    │
    ▼
Auth middleware             → actor_id stamped onto ctx (when authenticated)
    │
    ▼
Handler → Usecase           → layer=usecase     calls applogger.Business()
    │
    ▼
Infra / External service    → layer=infra       calls applogger.Dependency()
    │
    ▼
Response                    → middleware emits final log with status_code
```

Every log entry at every layer automatically carries `request_id`, `actor_id`,
and `client_ip` from context — no manual attachment needed.

---

## Required Fields

Every log entry must include these fields:

| Field        | Type     | Source                        |
|--------------|----------|-------------------------------|
| `timestamp`  | string   | Logger implementation         |
| `level`      | string   | Logger implementation         |
| `request_id` | string   | Injected by Logging middleware |
| `layer`      | string   | Set by caller or helper       |
| `message`    | string   | `msg` argument to log call    |

---

## Recommended Fields by Layer

### Middleware (`layer=middleware`)

| Field        | Description                          |
|--------------|--------------------------------------|
| `method`     | HTTP method (GET, POST, …)           |
| `path`       | URL path                             |
| `status_code`| HTTP response status code            |
| `duration_ms`| Total request duration in ms         |
| `user_agent` | Client User-Agent header             |
| `client_ip`  | Remote address (auto from context)   |
| `category`   | Always `"system"` for HTTP logs      |

### Usecase (`layer=usecase`)

| Field      | Description                                     |
|------------|-------------------------------------------------|
| `action`   | Business event name, e.g. `"order_created"`     |
| `actor_id` | Authenticated user ID (auto from context)        |
| Any domain-specific fields relevant to the event |

### Infra (`layer=infra`)

| Field        | Description                                   |
|--------------|-----------------------------------------------|
| `provider`   | Dependency name, e.g. `"midtrans"`, `"smtp"`  |
| `operation`  | Operation performed, e.g. `"charge"`          |
| `duration_ms`| Duration of the external call in ms           |
| `status`     | `"success"` or `"failure"`                   |
| `error`      | Error message (only on failure)               |

---

## Helper Functions

Import path: `service-core/internal/common/logger`

### `applogger.Business()` — usecase layer

```go
applogger.Business(ctx, log, "order_created",
    applogger.Field{Key: "order_id", Value: id.String()},
    applogger.Field{Key: "total",    Value: 250000},
)
```

Automatically attaches `layer=usecase` and `action`. Do NOT call
`log.Info(ctx, ...)` directly from usecases — use this helper instead.

### `applogger.Dependency()` — infra layer

```go
start := time.Now()
resp, err := midtransClient.Charge(req)
applogger.Dependency(ctx, log, "midtrans", "charge",
    time.Since(start).Milliseconds(), err)
```

Automatically attaches `layer=infra`, `provider`, `operation`, `duration_ms`,
and `status`. Logs at `Error` level when `err != nil`.

### `applogger.LogError()` — any layer

```go
applogger.LogError(ctx, log, applogger.LayerUsecase, err,
    applogger.Field{Key: "order_id", Value: id.String()},
)
```

Automatically attaches `layer` and `error`. Use this instead of `log.Error`
when you only have an error to report.

### `applogger.Redact()` — sanitising caller-supplied data

```go
safe := applogger.Redact(fields)
log.Info(ctx, "user_updated", safe...)
```

Must be called before logging any user-supplied or request-derived field slice.
See the Forbidden Fields section below.

---

## Sensitive Field Redaction

### Forbidden Fields

The following field keys must **never** appear in any log entry with their real values.
Any field matching these keys (case-insensitive) will be replaced with `[REDACTED]`
by `applogger.Redact()`.

```
password
access_token
refresh_token
authorization_header
authorization
otp
cvv
credit_card_number
card_number
session_cookie
cookie
secret
private_key
api_key
api_secret
```

### Rules

1. **Call `Redact()` before logging any user-supplied field slice.**
2. **Never log raw request bodies.** Parse and select specific fields.
3. **Never use `fmt.Sprintf("%+v", struct{})`.** Structs often contain
   sensitive fields that bypass key-based redaction.
4. **Never log full HTTP headers.** Extract only what is needed
   (e.g. `User-Agent`) and avoid logging `Authorization`.

---

## Log Volume Guidelines

### Do NOT log at the repository level for every query

```
// BAD — creates noise, exposes implementation details
FindUser
FindAddress
CreateOrder
CreatePayment
```

### DO log meaningful events at the usecase and infra levels

```
// GOOD — describes what happened, not how
order_created
inventory_reserved
midtrans_charge
shipping_calculated
```

Logs should answer business and operational questions, not trace SQL calls.

---

## Environment Strategy

### Development

- Format: human-readable text (zap `DevelopmentConfig`, slog `TextHandler`)
- Colour-coded log levels
- Optimised for local debugging

### Production

- Format: structured JSON (zap `ProductionConfig`, slog `JSONHandler`)
- One JSON object per line
- Optimised for log aggregators (Datadog, Loki, CloudWatch, etc.)

---

## Out of Scope

The following decisions are deferred until the logging structure is stable:

- Database log storage
- JSON file storage
- ELK / OpenSearch integration
- Loki / Grafana integration
- Log retention and archival policies
- Alerting systems based on log events

---

## Acceptance Criteria Checklist

- [x] Request ID generated and propagated through all layers
- [x] Middleware logs request lifecycle events (`layer=middleware`)
- [x] Usecase logs business events (`applogger.Business()`)
- [x] Infrastructure logs dependency interactions (`applogger.Dependency()`)
- [x] Centralized logging abstraction (`Logger` interface + helpers)
- [x] Automatic error logging (`applogger.LogError()`)
- [x] Manual business logging supported (`applogger.Business()`)
- [x] Consistent log schema defined (required + recommended fields above)
- [x] Development and production strategies documented
- [x] Sensitive field redaction rules documented
- [x] Logging standards documented for future contributors
- [x] Architecture prepared for future storage backends (`AuditLogger`)
