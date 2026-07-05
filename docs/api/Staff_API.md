# Staff API Documentation

This document covers the staff and admin-facing APIs for the Chia Florist service.
Endpoints are organized by access level: **Public**, **Staff**, and **Admin**.

## TODO

- [x] Public API
  - [x] Products
    - [x] Find Products
    - [x] Get Product Detail
  - [x] Shops
    - [x] Get Shop Addresses
    - [x] Get Shop Couriers
    - [x] Get Shop Products
- [x] Staff API
  - [x] Authentication
    - [x] Staff Sign In
    - [x] Me
    - [x] Log Out
  - [x] Shops
    - [x] Save Shop
  - [x] Inventory
    - [x] Add Inventory
  - [x] Profile
    - [x] Get Current User
    - [x] Update User
- [x] Staff Admin API
  - [x] Staff Management
    - [x] Create Staff
    - [x] Add Staff Account
    - [x] Find Staff
  - [x] Customer Management
    - [x] Find Customers
  - [x] Payment
    - [x] List Payment Method
    - [x] List Payment Account
    - [x] Create Payment Account
  - [x] Orders
    - [x] List Orders
    - [x] Get Order
  - [x] Audit Logs
    - [x] Find Audit Logs
    - [x] Get Audit Log

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

# Staff API

These endpoints require a valid staff session (any staff role).
Authentication is handled via a session cookie set at sign-in.

## Authentication

### Staff Sign In

- **Method**: `POST`
- **Endpoint**: `/auth/staff/signin`
- **Description**: Authenticate a staff account using email and password. Sets a session cookie on success.
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
- **Description**: Return the identity and roles of the currently authenticated staff account.
- **Authentication**: Staff (any role)
- **Request Body**: None

#### Response `200 OK`

```json
{
  "account_id": "51e20db6-5bdb-4f6a-b2b7-8d40c0db857d",
  "account_type": "staff",
  "is_authenticated": true,
  "roles": [
    { "code": "staff", "name": "Staff" }
  ],
  "permissions": [
    { "code": "staff" }
  ]
}
```

### Log Out

- **Method**: `POST`
- **Endpoint**: `/auth/logout`
- **Description**: Invalidate the current session and clear the refresh token cookie.
- **Authentication**: Staff (any role)
- **Request Body**: None

#### Response `200 OK`

```json
{ "message": "logout success" }
```

## Shops

### Save Shop

- **Method**: `POST`
- **Endpoint**: `/shops`
- **Description**: Create a new shop or update an existing shop owned by the authenticated staff. Omit `id` to create; supply `id` to update.
- **Authentication**: Staff (any role)
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
> The `is_active` flag is **only applied** when the actor holds the **staff admin** role.
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
| `403 Forbidden`    | Authenticated user does not hold a staff role. |

## Inventory

### Add Inventory

- **Method**: `POST`
- **Endpoint**: `/inventory`
- **Description**: Create a new inventory for an existing shop or existing product.
- **Authentication**: Staff (any role)
- **Request Body**:

  ```json
  {
      "product_id": "string - uuid (required)",
      "shop_id": "string - uuid (required)",
      "stock": "integer (required)"
  }
  ```

#### Response `200 OK`

```json
{ "message": "inventory successfully added" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `stock` is less than 0, either `product_id` or `shop_id` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `403 Not Found`    | Either product or shop is missing. |
| `409 Conflict`     | System conflict includes **inventory already exists for product and shop**. |

## Profile

### Get Current User

- **Method**: `GET`
- **Endpoint**: `/profile`
- **Description**: Retrieve the profile of the currently authenticated staff.
- **Authentication**: Staff (any role)
- **Request Body**: None

#### Response `200 OK`

```json
{
  "profile": {
    "staff_id": "d0118ea9-ab28-4338-89dd-3d3ca5f89880",
    "user_id": "56d88b08-ad99-4c91-9571-15b5bae95591",
    "Name": "Astra Yao",
    "Username": "friedriceuwu",
    "Phone": "021",
    "AvatarURL": "https://pbs.twimg.com/profile_images/1868113521902489600/SeenT54s_400x400.png",
    "LastLoginAt": null,
    "CreatedAt": "2026-06-26T19:47:28.045173Z",
    "UpdatedAt": null
  }
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | User profile not found.     |

### Update Current User

- **Method**: `PUT`
- **Endpoint**: `/profile`
- **Description**: Update the profile of the currently authenticated staff.
- **Authentication**: Staff (any role)
- **Request Body**: None

  ```json
  {
    "name": "string - optional",
    "phone": "string - optional",
    "avatar_url": "string - optional"
  }
  ```

#### Response `200 OK`

```json
{
  "profile": {
    "staff_id": "d0118ea9-ab28-4338-89dd-3d3ca5f89880",
    "user_id": "56d88b08-ad99-4c91-9571-15b5bae95591",
    "Name": "Astra Yao",
    "Username": "friedriceuwu",
    "Phone": "021",
    "AvatarURL": "https://pbs.twimg.com/profile_images/1868113521902489600/SeenT54s_400x400.png",
    "LastLoginAt": null,
    "CreatedAt": "2026-06-26T19:47:28.045173Z",
    "UpdatedAt": null
  }
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | User profile not found.     |

# Staff Admin API

These endpoints require a valid staff session with the **staff admin** role.

## Staff Management

### Create Staff

- **Method**: `POST`
- **Endpoint**: `/staff`
- **Description**: Create a new staff entity. The authenticated account is automatically associated as the staff owner.
- **Authentication**: Staff Admin
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
{ "message": "staff successfully created" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `name` is empty, or optional string fields are present but empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Add Staff Account

- **Method**: `POST`
- **Endpoint**: `/staff/{staffID}/accounts`
- **Description**: Register and assign a new account to a staff. The actor must be an admin of the target staff.
- **Authentication**: Staff Admin
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
| `staffID` | UUID (string) | The ID of the target staff. |

#### Response `201 Created`

```json
{ "message": "verify success" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` is not a valid UUID, or `email`, `name`, or `username` is empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Find Staff

- **Method**: `GET`
- **Endpoint**: `/staff`
- **Description**: Retrieve a paginated list of staff with optional filtering and sorting.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter by exact staff ID. |
| `name`    | string | No       | Filter by staff name (partial match). |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field    | Example            | Description                 |
|----------|--------------------|-----------------------------|
| `latest` | `sort=latest:desc` | Sort by creation date.      |
| `name`   | `sort=name:asc`    | Sort by staff name.      |
| `modified` | `sort=modified:desc` | Sort by last modified date. |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=name:asc,latest:desc`.

**Examples**:

- `GET /staff?page=1&limit=10`
- `GET /staff?name=chia`
- `GET /staff?sort=name:asc`

#### Response `200 OK`

```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "staff": [
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
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No staff match the given filters. |

## Customer Management

### Find Customers

- **Method**: `GET`
- **Endpoint**: `/customers`
- **Description**: Retrieve a paginated list of registered customer accounts with optional filtering and sorting.
- **Authentication**: Staff Admin
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
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No customers match the given filters. |

## Payment Management

### List Payment Method

- **Method**: `GET`
- **Endpoint**: `/payments/methods`
- **Description**: Retrieve a list of payment methods.
- **Authentication**: Staff
- **Request Body**: None

#### Response `200 OK`

```json
{
    "methods": [
        {
            "id": "0137d751-5188-447a-b630-1bf858f4f866",
            "name": "QRIS",
            "type": "qr_code",
            "is_active": true,
            "description": "QRIS payment via Midtrans",
            "fee_type": "",
            "fee_fixed": 0,
            "fee_percentage": 0
        },
        {
            "id": "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
            "name": "GoPay",
            "type": "ewallet",
            "is_active": true,
            "description": "GoPay via Midtrans",
            "fee_type": "",
            "fee_fixed": 0,
            "fee_percentage": 0
        },
        {
            "id": "074b02e4-e047-4f60-bdb0-cfeb5481d002",
            "name": "DANA",
            "type": "ewallet",
            "is_active": true,
            "description": "DANA via Midtrans",
            "fee_type": "",
            "fee_fixed": 0,
            "fee_percentage": 0
        },
        {
            "id": "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
            "name": "Mandiri",
            "type": "bank_transfer",
            "is_active": true,
            "description": "Mandiri Bill Payment via Midtrans",
            "fee_type": "",
            "fee_fixed": 0,
            "fee_percentage": 0
        }
    ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` is provided but is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### List Payment Accounts

- **Method**: `GET`
- **Endpoint**: `/payments/accounts`
- **Description**: Retrieve a list of payment accounts.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Response `200 OK`

```json
{
    "accounts": [
        {
            "id": "5672b98b-3474-4bbe-94de-ccf784ae90dc",
            "method_id": "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
            "account_name": "Mandiri Reyhan",
            "account_number": "1690002799366",
            "phone_number": "0895326204046",
            "qr_string": null
        },
        {
            "id": "01198989-6b57-4005-b7e6-50e797ccca04",
            "method_id": "074b02e4-e047-4f60-bdb0-cfeb5481d002",
            "account_name": "Dana Ilham",
            "account_number": null,
            "phone_number": "081291302897",
            "qr_string": null
        },
        {
            "id": "d8242cd0-7e80-4bbb-95cd-7b0061230ed6",
            "method_id": "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
            "account_name": "GoPay Ilham",
            "account_number": null,
            "phone_number": "081291302897",
            "qr_string": null
        },
        {
            "id": "e16476f0-f7d1-4959-a3cc-07f95e2d6242",
            "method_id": "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
            "account_name": "GoPay Reyhan",
            "account_number": null,
            "phone_number": "0895326204046",
            "qr_string": null
        }
    ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` is provided but is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Create Payment Account

- **Method**: `POST`
- **Endpoint**: `/payments/accounts`
- **Description**: Create payment account per payment method.
- **Authentication**: Customer

- **Request Body**:

  ```json
  {
    "method_id": "string (UUID, required)",
    "account_number": "string (optional)",
    "account_name": "string (required)",
    "phone_number": "string (required)",
    "is_active": "string (required)"
  }
  ```

#### Important Notes

> Account number is optional but required for creating account that based on `bank_transfer` payment method.

#### Response `200 OK`

```json
{
    "message": "payment account successfully created"
}
```

## Order Management

### Find Orders

- **Method**: `GET`
- **Endpoint**: `/orders`
- **Description**: Retrieve a paginated list of registered customer accounts with optional filtering and sorting.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | Page number for pagination. Must be `>= 1`. |
| `limit` | `integer` | No | `10` | Number of items per page. Must be `>= 1`. |
| `sort` | `string` | No | `latest:desc` | Sorting criteria. Format: `field:direction` (e.g. `latest:asc`, `total:desc`). Multiple sorts can be comma-separated. |
| `id` | `string` | No | - | Filter by unique Order UUID. |
| `number` | `string` | No | - | Filter/Search by order number (case-insensitive substring search). |
| `user_id` | `string` | No | - | Filter by customer User UUID. |
| `status` | `string` | No | - | Filter by order status (`pending`, `confirmed`, `processing`, `shipped`, `delivered`, `cancelled`). |

#### Sort Fields

| Field        | Example                | Description                    |
|--------------|------------------------|--------------------------------|
| `latest` / `date`     | `sort=latest:desc`     | Sort by account creation date. |
| `number`       | `sort=number:asc`        | Sort by display number.          |
| `total`   | `sort=total:asc`    | Sort by total paid.              |
| `status`      | `sort=status:asc`       | Sort by status number.          |
| `modified`     | `sort=modified:desc`     | Sort by last modified date.    |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=latest:asc,number:desc`.

**Examples**:

- `GET /orders?page=1&limit=20`
- `GET /orders?number=011&sort=modified:desc`

#### Response `200 OK`

```json
{
  "orders": [
    {
      "id": "e4a31771-4638-4e89-a292-624e723927d1",
      "number": "ORD-20260621-E4A317",
      "user_id": "8ce91a56-deea-46ac-9330-de65d64daa32",
      "address_id": "48956fd0-bcea-44a2-b598-af999d7abc7a",
      "status": "pending",
      "subtotal": 150000,
      "shipping_fee": 15000,
      "total": 165000,
      "created_at": "2026-06-21T08:45:00Z",
      "updated_at": null,
      "items": [
        {
          "id": "2529f895-5b14-4aa1-a8a1-5bf890160441",
          "product_id": "f55b14a1-a8a1-5bf8-9016-0441d2529f89",
          "product_name": "Premium Red Roses Bouquet",
          "quantity": 1,
          "unit_price": 150000,
          "subtotal": 150000,
          "shop_id": "8c1e82f6-30a2-4cc6-be52-f05aa2d6e9ec",
          "shop_name": "Jakarta Florist Central",
          "courier_code": "jne",
          "courier_service": "REG",
          "shipping_fee": 15000
        }
      ]
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` is provided but is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No customers match the given filters. |

### Get Order

- **Method**: `GET`
- **Endpoint**: `/users/me/orders/{id}`
- **Description**: Retrieve the order by its ID. This API meant for staff administrators to view order details.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
 "id": "00e9dc8e-e96d-4cb7-b9ea-8d5d5c8d930f",
 "number": "ORD-20260627-63E4B3",
 "customer_id": "7466a260-6dd9-45db-8a70-3254bfc2dc98",
 "address_id": "f279f798-2de1-4ebd-a660-568b835b3a52",
 "status": "pending",
 "subtotal": 450000,
 "shipping_fee": 770000,
 "total": 1220000,
 "created_at": "2026-06-27T10:09:05.904647Z",
 "items": [
  {
   "id": "e73f9f49-b3b3-4c8e-8fee-dc5d4a0207ab",
   "product_id": "480eec7c-d950-4927-a570-1fc3dc20df67",
   "product_name": "Prosperity Grand Opening Stand",
   "quantity": 3,
   "unit_price": 150000,
   "subtotal": 450000,
   "shop_id": "7e5e335a-ec5b-4399-a8f6-1ea7dd8f0974",
   "shop_name": "dayum",
   "courier_code": "tiki",
   "courier_service": "SDS",
   "shipping_fee": 770000
  }
 ],
 "payment": {
  "id": "1d3a0355-ce51-4346-a07c-b8bc839e85f1",
  "status": "pending",
  "provider": "midtrans",
  "amount": 1220000,
  "expires_at": "2026-06-28T10:09:09Z",
  "created_at": "2026-06-27T10:09:05.904647Z"
 }
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | User profile not found.     |

## Audit Logs

### Find Audit Logs

- **Method**: `GET`
- **Endpoint**: `/api/stats`
- **Description**: Retrieve a paginated list of audit logs with optional filtering and sorting.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | Page number for pagination. Must be `>= 1`. |
| `limit` | `integer` | No | `10` | Number of items per page. Must be `>= 1`. |
| `sort` | `string` | No | `date:desc` | Sorting criteria. Format: `field:direction` (e.g. `date:asc`, `action:desc`). Multiple sorts can be comma-separated. |
| `action` | `string` | No | - | Filter by exact audit action. |
| `user_id` | `string` | No | - | Filter by the user ID (actor_id) who performed the action. |
| `start_date` | `string` | No | - | Start date of the range to query (RFC3339 format, e.g. `2026-07-05T00:00:00Z` or `YYYY-MM-DD`). |
| `end_date` | `string` | No | - | End date of the range to query (RFC3339 format, e.g. `2026-07-05T23:59:59Z` or `YYYY-MM-DD`). |

#### Sort Fields

| Field | Example | Description |
| :--- | :--- | :--- |
| `date` | `sort=date:desc` | Sort by log creation date. |
| `action` | `sort=action:asc` | Sort alphabetically by action name. |

**Examples**:

- `GET /api/stats?page=1&limit=20`
- `GET /api/stats?action=signin&user_id=8ce91a56-deea-46ac-9330-de65d64daa32`
- `GET /api/stats?start_date=2026-07-01&end_date=2026-07-05T23:59:59Z&sort=date:asc`

#### Response `200 OK`

```json
{
  "audit_logs": [
    {
      "id": "e4a31771-4638-4e89-a292-624e723927d1",
      "category": "user_action",
      "action": "signin",
      "resource": "session",
      "resource_id": "session_12345",
      "actor_id": "8ce91a56-deea-46ac-9330-de65d64daa32",
      "outcome": "success",
      "request_id": "req-01ef82a",
      "client_ip": "127.0.0.1",
      "metadata": {
        "user_agent": "Mozilla/5.0"
      },
      "created_at": "2026-07-05T08:45:00Z"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

#### Error Responses

| Status | Condition |
| :--- | :--- |
| `400 Bad Request` | Invalid UUID or invalid date format. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden` | Authenticated user does not have the staff admin role. |

### Get Audit Log

- **Method**: `GET`
- **Endpoint**: `/api/stats/{id}`
- **Description**: Retrieve the details of a single audit log by its ID.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Path Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `id` | `string (UUID)` | Yes | The ID of the audit log to retrieve. |

#### Response `200 OK`

```json
{
  "id": "e4a31771-4638-4e89-a292-624e723927d1",
  "category": "user_action",
  "action": "signin",
  "resource": "session",
  "resource_id": "session_12345",
  "actor_id": "8ce91a56-deea-46ac-9330-de65d64daa32",
  "outcome": "success",
  "request_id": "req-01ef82a",
  "client_ip": "127.0.0.1",
  "metadata": {
    "user_agent": "Mozilla/5.0"
  },
  "created_at": "2026-07-05T08:45:00Z"
}
```

#### Error Responses

| Status | Condition |
| :--- | :--- |
| `400 Bad Request` | The provided `id` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden` | Authenticated user does not have the staff admin role. |
| `404 Not Found` | Audit log with the given ID does not exist. |
