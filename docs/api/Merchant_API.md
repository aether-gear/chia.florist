# Merchant API Documentation

This document covers the merchant staff and admin-facing APIs for the Chia Florist service.
Endpoints are organized by access level: **Public**, **Merchant Staff**, and **Merchant Admin**.

## TODO

- [ ] Public API
  - [ ] Products
    - [ ] Find Products
    - [ ] Get Product Detail
  - [ ] Shops
    - [ ] Get Shop Addresses
    - [ ] Get Shop Couriers
    - [ ] Get Shop Products
- [ ] Merchant Staff API
  - [ ] Authentication
    - [ ] Merchant Sign In
    - [ ] Me
    - [ ] Log Out
  - [ ] Shops
    - [ ] Save Shop
- [ ] Merchant Admin API
  - [ ] Merchant Management
    - [ ] Create Merchant
    - [ ] Add Merchant Account
    - [ ] Find Merchants
  - [ ] Customer Management
    - [ ] Find Customers

# Public API

No authentication is required for these endpoints.

## Products

### Find Products

- **Method**: `GET`
- **Endpoint**: `/products`
- **Description**: Retrieve a paginated list of products with optional filtering and sorting.
- **Authentication**: None (public)
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter products by product ID. |
| `name`    | string | No       | Search products by name. |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. |
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field      | Example              | Description                       |
|------------|----------------------|-----------------------------------|
| `date`     | `sort=date:desc`     | Sort by creation date.            |
| `name`     | `sort=name:asc`      | Sort by product name.             |
| `price`    | `sort=price:asc`     | Sort by product price.            |
| `weight`   | `sort=weight:desc`   | Sort by product weight.           |
| `status`   | `sort=status:asc`    | Sort by product status.           |
| `modified` | `sort=modified:desc` | Sort by last modified date.       |
| `archived` | `sort=archived:asc`  | Sort by archived status.          |
| `stock`    | `sort=stock:desc`    | Sort by available stock quantity. |

**Examples**:
- `GET /products?page=1&limit=10`
- `GET /products?name=coffee`
- `GET /products?sort=price:asc`
- `GET /products?name=coffee&sort=stock:desc&page=1&limit=20`

#### Response `200 OK`

```json
{
  "limit": 10,
  "page": 1,
  "total": 2,
  "products": [
    {
      "id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "sku": "EVT-ANV-001",
      "name": "Anniversary",
      "slug": "anniversary",
      "status": "active",
      "is_available": true,
      "price": 85000,
      "stock": 1291,
      "banner": {
        "thumbnail": "https://example.com/thumbnail.jpg",
        "preview": null,
        "detail": null
      },
      "availability": [
        { "slug": "Chia Medan Satria", "name": "chia-medan-satria", "stock": 981 },
        { "slug": "Chia Cikarang",     "name": "chia-cikarang",     "stock": 310 }
      ]
    }
  ]
}
```

### Get Product Detail

- **Method**: `GET`
- **Endpoint**: `/products/{slug}`
- **Description**: Retrieve details of a specific product by its slug.
- **Authentication**: None (public)
- **Request Body**: None

#### Path Parameters

| Parameter | Type   | Description       |
|-----------|--------|-------------------|
| `slug`    | string | The product slug. |

#### Response `200 OK`

```json
{
  "id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
  "sku": "EVT-GOP-007",
  "name": "Grand Opening",
  "slug": "grand-opening",
  "status": "active",
  "is_available": true,
  "price": 150000,
  "stock": 837,
  "description": "A towering two-tier floral stand featuring red anthuriums.",
  "weight": 4500,
  "updated_at": null,
  "banner": {
    "thumbnail": "https://example.com/thumbnail.jpg",
    "preview": "https://example.com/preview.jpg",
    "detail": "https://example.com/detail.jpg"
  },
  "gallery": [
    {
      "thumbnail": "https://example.com/thumbnail.jpg",
      "preview": "https://example.com/preview.jpg",
      "detail": "https://example.com/detail.jpg"
    }
  ],
  "availability": [
    { "slug": "Chia Cipinang",     "name": "chia-cipinang",     "stock": 109 },
    { "slug": "Chia Medan Satria", "name": "chia-medan-satria", "stock": 728 }
  ]
}
```

## Shops

### Find Shops

- **Method**: `GET`
- **Endpoint**: `/shops`
- **Description**: Retrieve a paginated list of shops with optional filtering and sorting.
- **Authentication**: None (public)
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter shops by shop ID. |
| `name`    | string | No       | Search shops by name. |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. |
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field      | Example              | Description                       |
|------------|----------------------|-----------------------------------|
| `name`     | `sort=name:asc`      | Sort by shop name.                |
| `active`   | `sort=active:desc`   | Sort by shop status.              |
| `date`     | `sort=date:desc`     | Sort by creation date.            |
| `modified` | `sort=modified:desc` | Sort by last modified date.       |

**Examples**:
- `GET /shops?page=1&limit=10`
- `GET /shops?name=coffee`
- `GET /shops?sort=date:asc`
- `GET /shops?name=coffee&sort=date:desc&page=1&limit=20`

#### Response `200 OK`

```json
{
  "limit": 10,
  "page": 1,
  "total": 2,
  "shops": [
    {
      "id": "427db07b-dbad-43ee-9199-49a2849a4e30",
      "name": "Chia Bogor",
      "slug": "chia-bogor",
      "description": "Toko cabang berlokasi di Bogor, Jawa Barat",
      "is_active": true,
      "created_at": "2026-06-17T14:50:22.544511Z",
      "updated_at": null
    },
    {
      "id": "06011c23-73c1-4d97-a1e4-54cb6c98f398",
      "name": "Chia Kranji",
      "slug": "chia-kranji",
      "description": "Toko cabang berlokasi di Kranji, Kota Bekasi",
      "is_active": true,
      "created_at": "2026-06-17T14:49:36.708634Z",
      "updated_at": null
    }
  ]
}
```

### Get Shop Addresses

- **Method**: `GET`
- **Endpoint**: `/shops/{shopID}/addresses`
- **Description**: Retrieve all addresses registered for a specific shop.
- **Authentication**: None (public)
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                |
|-----------|---------------|----------------------------|
| `shopID`  | UUID (string) | The unique ID of the shop. |

#### Response `200 OK`

```json
{
  "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "addresses": [
    {
      "id": "d4e5f6a7-b8c9-0123-defa-234567890123",
      "label": "Main Branch",
      "phone": "+6281234567890",
      "is_active": true,
      "province_id": "32",
      "city_id": "3204",
      "district_id": "320401",
      "village_id": "3204010001",
      "full_address": "Jl. Bunga Indah No. 10, Bekasi",
      "postal_code": "17520",
      "created_at": "2026-01-15T09:00:00Z",
      "updated_at": "2026-05-01T12:00:00Z"
    }
  ]
}
```

#### Error Responses

| Status            | Condition                                 |
|-------------------|-------------------------------------------|
| `400 Bad Request` | `shopID` in the path is not a valid UUID. |

### Get Shop Couriers

- **Method**: `GET`
- **Endpoint**: `/shops/{shopID}/couriers`
- **Description**: Retrieve the list of courier services configured for a specific shop.
- **Authentication**: None (public)
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                |
|-----------|---------------|----------------------------|
| `shopID`  | UUID (string) | The unique ID of the shop. |

#### Response `200 OK`

```json
{
  "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "couriers": [
    { "code": "jne",     "active": true  },
    { "code": "sicepat", "active": false }
  ]
}
```

#### Error Responses

| Status            | Condition                                 |
|-------------------|-------------------------------------------|
| `400 Bad Request` | `shopID` in the path is not a valid UUID. |

### Get Shop Products

- **Method**: `GET`
- **Endpoint**: `/shops/{shopID}/products`
- **Description**: Retrieve all products listed in a specific shop along with their current inventory levels.
- **Authentication**: None (public)
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                |
|-----------|---------------|----------------------------|
| `shopID`  | UUID (string) | The unique ID of the shop. |

#### Response `200 OK`

```json
{
  "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "products": [
    {
      "id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "sku": "EVT-ANV-001",
      "name": "Anniversary",
      "slug": "anniversary",
      "description": "A beautiful anniversary bouquet.",
      "status": "active",
      "price": 85000,
      "weight": 1.5,
      "inventory": {
        "total_stock": 100,
        "reserved_stock": 19,
        "available": 81
      },
      "created_at": "2026-01-10T08:00:00Z",
      "updated_at": "2026-06-01T10:30:00Z"
    }
  ]
}
```

#### Inventory Fields

| Field            | Type | Description                                        |
|------------------|------|----------------------------------------------------|
| `total_stock`    | int  | Total units in stock.                              |
| `reserved_stock` | int  | Units reserved for pending orders.                 |
| `available`      | int  | Units available for purchase (`total - reserved`). |

**Product `status` values**: `"active"`, `"inactive"`, `"archived"`.

#### Error Responses

| Status            | Condition                                 |
|-------------------|-------------------------------------------|
| `400 Bad Request` | `shopID` in the path is not a valid UUID. |

# Merchant Staff API

These endpoints require a valid merchant session (any merchant role).
Authentication is handled via a session cookie set at sign-in.

## Authentication

### Merchant Sign In

- **Method**: `POST`
- **Endpoint**: `/auth/merchant/signin`
- **Description**: Authenticate a merchant account using email and password. Sets a session cookie on success.
- **Authentication**: None (public endpoint)
- **Request Body**:
  ```json
  {
    "email":      "string (required)",
    "password":   "string (required)",
    "user_agent": "string (optional)",
    "ip_address": "string (optional)"
  }
  ```

#### Response `200 OK`

- **Set-Cookie**: `hotpot=<value>`
- **Body**:
  ```json
  { "message": "login success" }
  ```

### Me

- **Method**: `GET`
- **Endpoint**: `/auth/me`
- **Description**: Return the identity and roles of the currently authenticated merchant account.
- **Authentication**: Merchant (any role)
- **Request Body**: None

#### Response `200 OK`

```json
{
  "account_id": "51e20db6-5bdb-4f6a-b2b7-8d40c0db857d",
  "account_type": "merchant",
  "is_authenticated": true,
  "roles": [
    { "code": "merchant_staff", "name": "Merchant Staff" }
  ],
  "permissions": [
    { "code": "merchant_staff" }
  ]
}
```

### Log Out

- **Method**: `POST`
- **Endpoint**: `/auth/logout`
- **Description**: Invalidate the current session and clear the refresh token cookie.
- **Authentication**: Merchant (any role)
- **Request Body**: None

#### Response `200 OK`

```json
{ "message": "logout success" }
```

## Shops

### Save Shop

- **Method**: `POST`
- **Endpoint**: `/shops`
- **Description**: Create a new shop or update an existing shop owned by the authenticated merchant. Omit `id` to create; supply `id` to update.
- **Authentication**: Merchant (any role)
- **Request Body**:
  ```json
  {
    "id":          "string (UUID, optional — omit to create, supply to update)",
    "name":        "string (required)",
    "description": "string (optional)",
    "is_active":   "string (required, e.g. \"true\" or \"false\")"
  }
  ```

#### Important Notes

> `is_active` is a string-encoded boolean. The server parses it with `strconv.ParseBool`,
> so `"true"`, `"1"`, `"TRUE"` and `"false"`, `"0"`, `"FALSE"` are all accepted.
>
> The `is_active` flag is **only applied** when the actor holds the **merchant admin** role.
> For non-admin staff the field is silently ignored.

#### Response `200 OK`

```json
{ "message": "shop successfully saved" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `name` is empty, `id` is not a valid UUID, or `is_active` cannot be parsed as a boolean. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a merchant role. |

# Merchant Admin API

These endpoints require a valid merchant session with the **merchant admin** role.

## Merchant Management

### Create Merchant

- **Method**: `POST`
- **Endpoint**: `/merchants`
- **Description**: Create a new merchant entity. The authenticated account is automatically associated as the merchant owner.
- **Authentication**: Merchant Admin
- **Request Body**:
  ```json
  {
    "name":        "string (required)",
    "description": "string (optional)",
    "logo_url":    "string (optional)",
    "banner_url":  "string (optional)"
  }
  ```

#### Response `200 OK`

```json
{ "message": "merchant successfully created" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `name` is empty, or optional string fields are present but empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the merchant admin role. |

### Add Merchant Account

- **Method**: `POST`
- **Endpoint**: `/merchants/{merchantID}/accounts`
- **Description**: Register and assign a new account to a merchant. The actor must be an admin of the target merchant.
- **Authentication**: Merchant Admin
- **Request Body**:
  ```json
  {
    "email":    "string (required)",
    "name":     "string (required)",
    "username": "string (required)",
    "password": "string (required)",
    "phone":    "string (optional)"
  }
  ```

#### Path Parameters

| Parameter    | Type          | Description                    |
|--------------|---------------|--------------------------------|
| `merchantID` | UUID (string) | The ID of the target merchant. |

#### Response `201 Created`

```json
{ "message": "verify success" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `merchantID` is not a valid UUID, or `email`, `name`, or `username` is empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the merchant admin role. |

### Find Merchants

- **Method**: `GET`
- **Endpoint**: `/merchants`
- **Description**: Retrieve a paginated list of merchants with optional filtering and sorting.
- **Authentication**: Merchant Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter by exact merchant ID. |
| `name`    | string | No       | Filter by merchant name (partial match). |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field    | Example            | Description                 |
|----------|--------------------|-----------------------------|
| `latest` | `sort=latest:desc` | Sort by creation date.      |
| `name`   | `sort=name:asc`    | Sort by merchant name.      |
| `modify` | `sort=modify:desc` | Sort by last modified date. |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=name:asc,latest:desc`.

**Examples**:
- `GET /merchants?page=1&limit=10`
- `GET /merchants?name=chia`
- `GET /merchants?sort=name:asc`

#### Response `200 OK`

```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "merchants": [
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Chia Florist",
      "description": "Fresh flowers delivered to your door.",
      "logo_url": "https://example.com/logo.png",
      "banner_url": "https://example.com/banner.png",
      "created_at": "2026-01-01T08:00:00Z"
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` is provided but is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the merchant admin role. |
| `404 Not Found`    | No merchants match the given filters. |

## Customer Management

### Find Customers

- **Method**: `GET`
- **Endpoint**: `/customers`
- **Description**: Retrieve a paginated list of registered customer accounts with optional filtering and sorting.
- **Authentication**: Merchant Admin
- **Request Body**: None

#### Query Parameters

| Parameter  | Type   | Required | Description |
|------------|--------|----------|-------------|
| `id`       | UUID   | No       | Filter by exact customer ID. |
| `name`     | string | No       | Filter by customer display name (partial match). |
| `username` | string | No       | Filter by customer username (partial match). |
| `email`    | string | No       | Filter by customer email address (partial match). |
| `page`     | int    | No       | Page number. Defaults to `1`. |
| `limit`    | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`     | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field        | Example                | Description                    |
|--------------|------------------------|--------------------------------|
| `latest`     | `sort=latest:desc`     | Sort by account creation date. |
| `name`       | `sort=name:asc`        | Sort by display name.          |
| `username`   | `sort=username:asc`    | Sort by username.              |
| `phone`      | `sort=phone:asc`       | Sort by phone number.          |
| `modify`     | `sort=modify:desc`     | Sort by last modified date.    |
| `last_login` | `sort=last_login:desc` | Sort by last login timestamp.  |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=name:asc,last_login:desc`.

**Examples**:
- `GET /customers?page=1&limit=20`
- `GET /customers?name=jane`
- `GET /customers?email=jane@example.com&sort=last_login:desc`

#### Response `200 OK`

```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "users": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Jane Doe",
      "username": "janedoe",
      "phone": "+6281234567890",
      "last_login_at": "2026-06-17T10:00:00Z"
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` is provided but is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the merchant admin role. |
| `404 Not Found`    | No customers match the given filters. |
