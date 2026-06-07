## Final Chosen Architecture

### Core Philosophy

```text id="v7m2qx"
User is the central identity.
```

Everything branches from `users`.

The system separates:

* identity
* authentication
* business domains
* authorization

instead of mixing everything into roles.

## Core Structure

```text id="m4q8tv"
users
└── accounts (1:1 authentication)

users
├── consumers
└── merchants

merchants
└── merchant_users
```

## Entity Responsibilities

### users

Represents:

```text id="z1v7pk"
human identity
```

Contains:

* name
* username
* profile
* personal information

### accounts

Represents:

```text id="x5m2rv"
authentication layer
```

Contains:

* email
* password hash
* sessions
* OTP
* MFA
* account status

Important:

```text id="w8q4ny"
1 user = 1 account
```

### consumers

Represents:

```text id="t3v9mx"
customer/business-consumer capability
```

This separates customer domain logic from authentication.

Examples:

* loyalty
* addresses
* purchase history
* customer preferences

A user becomes customer-capable if:

* consumer row exists

### merchants

Represents:

```text id="p6q1tw"
merchant/business capability
```

This is NOT just a role.

Merchant is a business domain.

Examples:

* shop settings
* KYC
* operational data
* merchant onboarding

A user becomes merchant-capable if:

* merchant row exists

### merchant_users

Represents:

```text id="r2m8vx"
merchant team membership
```

Used later for:

* merchant staff
* merchant admins
* operator permissions

Example:

| merchant | user | role  |
| -- | -- | -- |
| shop-a   | john | owner |
| shop-a   | jane | staff |

## Why This Architecture Was Chosen

Because:

```text id="n9v3pk"
merchant and consumer are business domains
NOT merely roles
```

This creates cleaner boundaries.

## Authentication Flow

```text id="u5q7mx"
login
→ validate account
→ resolve user
→ resolve capabilities
    → consumer?
    → merchant?
→ create session
```

Authentication stays centralized.

## Authorization Strategy

RBAC becomes:

```text id="g4m1rv"
authorization layer
```

NOT identity modeling.

## Example Roles

```text id="x7v2tw"
customer
merchant_owner
merchant_staff
admin
```

## Policies Still Matter

Real authorization becomes:

```text id="m3q8nx"
validate ownership
+
validate merchant membership
+
validate permissions
```

NOT just:

```text id="j1v5pk"
role == admin
```

## Benefits

### Unified Identity

One human identity across system.

### Clean Domain Separation

Consumer and merchant logic stay isolated.

### Scalable

Supports:

* customer-only users
* merchant-only users
* users that are both

### Easier UX

No duplicate logins required.

### Future-Proof

Can later support:

* merchant teams
* merchant staff
* advanced RBAC
* scoped permissions
* policy engines

## Final Mental Model

```text id="h8m4qx"
User = who the human is
Account = how they authenticate
Consumer = customer capability
Merchant = business capability
RBAC = what they are allowed to do
```
