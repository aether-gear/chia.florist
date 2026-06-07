# Simple Lifecycle Flow

```text id="oq3a8m"
Login
 ↓
Validate Session
 ↓
Resolve Actor
 ↓
Resolve Roles
 ↓
Resolve Permissions
 ↓
Inject Actor Into Context
 ↓
Authorization Check
 ↓
Handler
 ↓
Service
```

---

# Simple Authorization Workflow

```text id="b7w1c5"
Account
 ↓
Membership
 ↓
Role
 ↓
Permission
```

---

# Necessary Domains Only

You are short on time, so keep ONLY these.

---

# 1. Authentication Domain

Already exists.

Purpose:

```text id="mqn0n8"
identity + session validation
```

Contains:

```text id="c6a7sz"
account
session
token
auth middleware
```

---

# 2. Authorization Domain

NEW core layer.

Purpose:

```text id="gozv4h"
capability resolution
```

Contains ONLY:

```text id="1j2h3f"
actor
role
permission
membership
```

---

# 3. Merchant Domain

Purpose:

```text id="9sm9k0"
business tenant/store
```

Contains:

```text id="5yn7zr"
merchant
merchant profile
merchant settings
```

NOT auth.

---

# Minimal Database Model

ONLY necessary tables.

---

# accounts

(authentication identity)

```sql id="s28dxr"
accounts
- id
- user_id
- email
- password_hash
- status
```

---

# users

(profile)

```sql id="ajz8eo"
users
- id
- name
- username
```

---

# merchants

(business tenant)

```sql id="m6v0q9"
merchants
- id
- name
- slug
- status
```

---

# merchant_memberships

THIS is the important table.

```sql id="8j3dy2"
merchant_memberships
- id
- merchant_id
- account_id
- role_id
```

This connects:

```text id="6z1zph"
account ↔ merchant ↔ role
```

---

# roles

Keep simple.

```sql id="pf0p4w"
roles
- id
- code
```

Example:

```text id="uv0qwa"
platform_admin
merchant_owner
merchant_staff
customer
```

---

# permissions

```sql id="x5m3m7"
permissions
- id
- code
```

Example:

```text id="a7n9s2"
product.write
order.manage
merchant.manage
```

---

# role_permissions

```sql id="4j0tcm"
role_permissions
- role_id
- permission_id
```

DONE.

That’s enough.

---

# Necessary Models Only

---

# Actor

MOST important new model.

```go id="d1l7fz"
type Actor struct {
    AccountID uuid.UUID
    UserID    uuid.UUID

    MerchantID *uuid.UUID

    Roles       []Role
    Permissions []Permission
}
```

---

# Role

```go id="9f5k0n"
type Role struct {
    ID   uuid.UUID
    Code string
}
```

---

# Permission

```go id="d3n6rk"
type Permission struct {
    ID   uuid.UUID
    Code string
}
```

---

# Membership

```go id="3n0vya"
type MerchantMembership struct {
    MerchantID uuid.UUID
    AccountID  uuid.UUID
    RoleID     uuid.UUID
}
```

---

# What To Refactor First

Now the important part.

---

# STEP 1 — Refactor Auth Identity

Current:

```text id="l9u0mw"
JWT → UserID
```

Refactor into:

```text id="km4d3x"
JWT → AccountID
```

because:

```text id="z7k3sh"
Account = authentication identity
```

This is your FIRST priority.

---

# STEP 2 — Introduce Actor Resolver Middleware

NEW middleware after auth.

Pipeline:

```text id="l2t1so"
RequireAuth
 ↓
ResolveActor
```

Responsibilities:

```text id="h7y5eo"
- load account
- load roles
- load permissions
- build actor
- inject actor into context
```

---

# STEP 3 — Introduce Permission Check Middleware

Example:

```go id="8s4z1j"
RequirePermission(permission.ProductWrite)
```

DO NOT check roles directly in handlers anymore.

---

# STEP 4 — Create Merchant Module

Now merchant becomes clean.

Merchant is:

```text id="5w8f2g"
business tenant
```

NOT identity/auth.

---

# STEP 5 — Add Merchant Memberships

Now:

```text id="7d0k9q"
account can operate merchant
```

cleanly.

---

# Recommended Minimal Folder Structure

Keep it VERY small.

```text id="z8h2jq"
modules/
├── authentication/
├── authorization/
├── merchant/
└── user/
```

---

# authorization module

ONLY:

```text id="v1k9sx"
actor
role
permission
membership
middleware
```

No policy engine yet.

---

# What You Should NOT Build Yet

Skip:

```text id="z4m7fr"
- ABAC
- Casbin
- policy DSL
- dynamic policy graph
- generic authorization framework
```

You do not need them yet.

---

# Final Important Direction

Your architecture is becoming:

```text id="f9j3wd"
Authentication
    ↓
Actor Resolution
    ↓
Authorization
    ↓
Business Domain
```

That is the correct scalable direction for your project.
