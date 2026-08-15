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
    - [x] Forgot Password
    - [x] Verify Forgot Password
    - [x] Reset Password
  - [x] Shops
    - [x] Save Shop
    - [x] Delete Shop
    - [x] Create Shop Address
    - [x] Update Shop Address
    - [x] Delete Shop Address
  - [x] Inventory
    - [x] Add Inventory
    - [x] Update Inventory
    - [x] Remove Inventory
  - [x] Profile
    - [x] Get Current User
    - [x] Update User
  - [x] Shipments
    - [x] Update Shipment Status
    - [x] Update Shipment
- [x] Staff Admin API
  - [x] Products Management
      - [x] Save Product
      - [x] Get Product Stats
      - [x] Delete Product
  - [x] Staff Management
    - [x] Create Staff
    - [x] Find Staff
    - [x] Update Staff
    - [x] Delete Staff
    - [x] List Staff Accounts
    - [x] Add Staff Account
    - [x] Remove Staff Account
  - [x] Customer Management
    - [x] Find Customers
  - [x] Payment
    - [x] List Payment Method
    - [x] Save Payment Method
    - [x] Save Payment Instruction
    - [x] List Payment Account
    - [x] Create Payment Account
  - [X] Orders
    - [X] List Orders
    - [X] Get Order
    - [X] Get Order Tracking Timeline (Customer Endpoint)

  - [X] Audit Logs
    - [X] Find Audit Logs
    - [X] Get Audit Log
  - [ ] Analytics
    - [ ] Get Order Metrics
    - [ ] Get Payment Metrics
    - [ ] Get Shipment Metrics
    - [ ] Get Inventory Metrics
    - [ ] Get Product Metrics
  - [ ] WAF Security Policy
    - [ ] WAF Rules
    - [ ] List Rules
    - [ ] Create Rule
    - [ ] Toggle Rule
    - [ ] Delete Rule
    - [ ] IP Access Control
    - [ ] List IP Config
    - [ ] Update IP Action
    - [ ] Filters
    - [ ] Get Filters
    - [ ] Update Filter
    - [ ] Threat Intelligence
    - [ ] Analyze IP
    - [ ] Get Geolocation

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

### Forgot Password

- **Method**: `POST`
- **Endpoint**: `/staff/forgot-password`
- **Description**: Create session for forgot password on account based on input email.
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "email": "string (required)",
  }
  ```

#### Response `200 OK`

- **Body**:
  ```json
  {
    "message": "if the email is registered you will receive a reset code shortly",
    "challenge_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890" // can be null
  }
  ```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `email` is missing or empty. |
| `403 Forbidden`    | The email is not registered or type account missmatch. |

### Verify Forgot Password

- **Method**: `POST`
- **Endpoint**: `/forgot-password/verify`
- **Description**: Verify the forgot password challenge based on the challenge ID and reset code.
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "challenge_id": "string (required)",
    "otp":          "string (required)"
  }
  ```

#### Response `200 OK`

- **Body**:
  ```json
  {
    "message": "otp verified",
    "challenge_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
  ```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `challenge_id` or `otp` is missing or empty. |
| `404 Not Found`    | The challenge ID is not registered or type account missmatch. |
| `409 Conflict`    | The session purpose is not reset the password, already consumed, already verified, or already expired. |

### Reset Password

- **Method**: `POST`
- **Endpoint**: `/forgot-password/reset`
- **Description**: Reset the password based on the challenge ID.
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "challenge_id": "string (required)",
    "new_password": "string (required)"
  }
  ```

#### Response `200 OK`

- **Body**:
  ```json
  {
    "message": "password reset successful",
  }
  ```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `challenge_id` or `new_password` is missing or empty. |
| `403 Forbidden`    | The challenge is not verified or the session purpose is not reset the password. |
| `404 Not Found`    | The challenge ID is not registered or the challenge purpose is invalid. |
| `409 Conflict`     | The challenge is already consumed, already verified, or already expired. |

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

### Delete Shop

- **Method**: `DELETE`
- **Endpoint**: `/shops/{shopID}`
- **Description**: Soft delete a shop by ID. The shop record is marked as deleted (`deleted_at` timestamp set) and will no longer be accessible in shop queries.
- **Authentication**: Staff Admin (`RoleStaffAdmin`)
- **Request Body**: None

#### Path Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `shopID`  | string (UUID) | Yes      | The unique ID of the shop to delete. |

#### Response `200 OK`

```json
{
  "message": "shop successfully deleted"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `shopID` path parameter is missing or not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold the staff admin role (`RoleStaffAdmin`). |
| `404 Not Found`    | Shop with the given `shopID` does not exist or is already soft-deleted. |

### Create Shop Address

- **Method**: `POST`
- **Endpoint**: `/shops/{shopID}/addresses`
- **Description**: Add a new address to a specific shop.
- **Authentication**: Staff (any role)
- **Request Body**:

  ```json
  {
    "label":        "string (required, e.g. \"Warehouse\")",
    "phone":        "string (optional)",
    "is_active":    "string (required, e.g. \"true\" or \"false\")",
    "province_id":  "string (required)",
    "city_id":      "string (required)",
    "district_id":  "string (required)",
    "village_id":   "string (required)",
    "full_address": "string (required)",
    "postal_code":  "string (required)"
  }
  ```

#### Response `200 OK`

```json
{
  "message": "address successfully created"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid request body or missing required location fields. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |

### Update Shop Address

- **Method**: `PUT`
- **Endpoint**: `/shops/{shopID}/addresses/{addressID}`
- **Description**: Update an existing shop address by ID. Setting `is_active = true` will automatically unset active status on other addresses for the shop.
- **Authentication**: Staff (any role)
- **Request Body**:

  ```json
  {
    "label":        "string (required)",
    "phone":        "string (optional)",
    "is_active":    "string (required, e.g. \"true\" or \"false\")",
    "province_id":  "string (required)",
    "city_id":      "string (required)",
    "district_id":  "string (required)",
    "village_id":   "string (required)",
    "full_address": "string (required)",
    "postal_code":  "string (required)"
  }
  ```

#### Path Parameters

| Parameter   | Type          | Required | Description |
|-------------|---------------|----------|-------------|
| `shopID`    | string (UUID) | Yes      | The unique ID of the shop. |
| `addressID` | string (UUID) | Yes      | The unique ID of the shop address to update. |

#### Response `200 OK`

```json
{
  "message": "address successfully updated"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `shopID` / `addressID` path parameter is not a valid UUID or missing required fields. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | Shop address with the given `addressID` does not exist for the shop or is soft-deleted. |

### Delete Shop Address

- **Method**: `DELETE`
- **Endpoint**: `/shops/{shopID}/addresses/{addressID}`
- **Description**: Soft delete a shop address by ID. Active (`is_active = true`) shop addresses cannot be deleted directly without making another address active first.
- **Authentication**: Staff (any role)
- **Request Body**: None

#### Path Parameters

| Parameter   | Type          | Required | Description |
|-------------|---------------|----------|-------------|
| `shopID`    | string (UUID) | Yes      | The unique ID of the shop. |
| `addressID` | string (UUID) | Yes      | The unique ID of the shop address to delete. |

#### Response `200 OK`

```json
{
  "message": "address successfully deleted"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `shopID` or `addressID` path parameter is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | Shop address with the given `addressID` does not exist for the shop or is soft-deleted. |
| `409 Conflict`     | Cannot delete active shop address (`is_active = true`). |

## Inventory

### Add Inventory

- **Method**: `POST`
- **Endpoint**: `/shops/{shopID}/products/{productID}/inventories`
- **Description**: Add inventory stock for a product in a shop.
- **Authentication**: Staff (any role)
- **Path Parameters**:

| Parameter   | Type | Required | Description |
|-------------|------|----------|-------------|
| `shopID`    | UUID | Yes      | The unique ID of the shop. |
| `productID` | UUID | Yes      | The unique ID of the product. |

- **Request Body**:

  ```json
  {
      "stock": 10
  }
  ```

#### Response `200 OK`

```json
{ "message": "inventory successfully added" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `stock` is less than 0, or either `shopID` or `productID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | Either product or shop is missing. |
| `409 Conflict`     | Inventory already exists for this product and shop. |

### Update Inventory

- **Method**: `PUT`
- **Endpoint**: `/shops/{shopID}/products/{productID}/inventories`
- **Description**: Update total stock level for a product's inventory in a shop.
- **Authentication**: Staff (any role)
- **Path Parameters**:

| Parameter   | Type | Required | Description |
|-------------|------|----------|-------------|
| `shopID`    | UUID | Yes      | The unique ID of the shop. |
| `productID` | UUID | Yes      | The unique ID of the product. |

- **Request Body**:

  ```json
  {
      "stock": 50
  }
  ```

#### Response `200 OK`

```json
{ "message": "inventory successfully updated" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `stock` is less than 0, either ID is invalid, or the new stock is less than the current reserved stock. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | Inventory record not found. |

### Remove Inventory

- **Method**: `DELETE`
- **Endpoint**: `/shops/{shopID}/products/{productID}/inventories`
- **Description**: Delete the inventory record for a product in a shop.
- **Authentication**: Staff (any role)
- **Path Parameters**:

| Parameter   | Type | Required | Description |
|-------------|------|----------|-------------|
| `shopID`    | UUID | Yes      | The unique ID of the shop. |
| `productID` | UUID | Yes      | The unique ID of the product. |

- **Request Body**: None

#### Response `200 OK`

```json
{ "message": "inventory successfully removed" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Either `shopID` or `productID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | Inventory record not found. |
| `409 Conflict`     | Cannot delete inventory with active reservations (`ReservedStock > 0`). |

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

## Products Management

### Save Product

- **Method**: `POST`
- **Endpoint**: `/products`
- **Description**: Create a new product or update an existing product owned by the authenticated staff. Omit `id` to create; supply `id` to update.
- **Authentication**: Staff (any role)
- **Request Body**:

  ```json
  {
    "id":                      "string (UUID, optional — omit to create, supply to update)",
    "sku":                     "string (required)",
    "name":                    "string (required)",
    "description":             "string (optional)",
    "is_available":            "bool (required)",
    "status":                  "string (optional): active, inactive or archived",
    "price":                   "int (required)",
    "weight":                  "float (optional)",
    "cost_price":              "int (optional, procurement cost per unit in cents)",
    "supplier_lead_time_days": "int (optional, average replenishment duration in days)"
  }
  ```

#### Response `200 OK`

```json
{ "message": "product successfully saved" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `name` or `sku` is empty, `id` is not a valid UUID, or `status` is invalid. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |

### Add Product Image

- **Method**: `POST`
- **Endpoint**: `/products/id/{id}/images`
- **Description**: Add an image to a product by product ID.
- **Authentication**: Staff (any role)
- **Request Multi-part**: Yes
- **Request Body**:

  ```json
  {
    "image": "file (required)"
  }
  ```

#### Response `200 OK`

```json
{ "message": "product image successfully added" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `image` is not a valid file. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not hold a staff role. |
| `404 Not Found`    | No product with the given `id` exists. |

### Get Product Stats

- **Method**: `GET`
- **Endpoint**: `/products/stats`
- **Description**: Retrieve a paginated list of products with performance statistics and metadata (views, gross margins, and multi-window sales velocity).
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter stats by a specific product ID. |
| `name`    | string | No       | Filter stats by product name (partial match). |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field | Example | Description |
|---|---|---|
| `latest` | `sort=latest:desc` | Sort by product creation date. |
| `name` | `sort=name:asc` | Sort alphabetically by product name. |
| `price` | `sort=price:desc` | Sort by product base price. |
| `view_count` | `sort=view_count:desc` | Sort by detail page view count (most-seen). |
| `sales_30d` | `sort=sales_30d:desc` | Sort by unit sales velocity in the last 30 days. |
| `sales_7d` | `sort=sales_7d:desc` | Sort by unit sales velocity in the last 7 days. |
| `revenue` | `sort=revenue:desc` | Sort by revenue contribution. |
| `gross_margin` | `sort=gross_margin:desc` | Sort by gross margin percentage. |

> Default sort: `latest:desc`.

**Examples**:
- `GET /products/stats?page=1&limit=10`
- `GET /products/stats?sort=view_count:desc`

#### Response `200 OK`

```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "stats": [
    {
      "id": "e2c3d4e5-f6a7-8901-bcde-f12345678901",
      "sku": "PROD-001",
      "name": "Red Roses Bouquet",
      "slug": "red-roses-bouquet",
      "status": "active",
      "price": 10000,
      "cost_price": 5000,
      "supplier_lead_time_days": 3,
      "gross_margin_pct": 50,
      "view_count": 120,
      "stock": 45,
      "sales_velocity_7d": 12,
      "sales_velocity_30d": 50,
      "sales_velocity_90d": 150,
      "conversion_rate": 41.67,
      "revenue_contribution_percentage": 15.5,
      "return_rate": null,
      "average_rating": null,
      "review_count": null,
      "thumbnail": "https://example.com/images/red-roses-thumb.jpg"
    }
  ]
}
```

#### Error Responses

| Status | Condition |
|---|---|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden` | Authenticated user does not have the staff admin role. |

#### Delete Product

- **Method**: `DELETE`
- **Endpoint**: `/products/id/{id}`
- **Description**: Permanently remove a product by ID.
- **Authentication**: Staff Admin
- **Request Body**: None

##### Path Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | string | Yes      | The product ID to delete. |

##### Response `204 No Content`

Empty body.

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` path parameter is missing. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No product with the given `id` exists. |

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
    "username":    "string (required)",
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
| `400 Bad Request`  | `name` or `username` is empty, or optional string fields are present but empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `409 Conflict`     | A user with this username already exists. |

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

### Update Staff

- **Method**: `PUT`
- **Endpoint**: `/staff/{staffID}`
- **Description**: Updates the metadata (name, description, logo, banner) of a staff entity and its associated root user record.
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

#### Path Parameters

| Parameter | Type          | Required | Description                 |
|-----------|---------------|----------|-----------------------------|
| `staffID` | UUID (string) | Yes      | The ID of the target staff. |

#### Response `200 OK`

```json
{
  "message": "staff successfully updated"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` is not a valid UUID, or `name` is empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No staff entity found with the given `staffID`. |

### Delete Staff

- **Method**: `DELETE`
- **Endpoint**: `/staff/{staffID}`
- **Description**: Soft-deletes a staff entity (`deleted_at = NOW()`) and cascade removes all associated staff memberships.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Required | Description                 |
|-----------|---------------|----------|-----------------------------|
| `staffID` | UUID (string) | Yes      | The ID of the target staff. |

#### Response `200 OK`

```json
{
  "message": "staff successfully deleted"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No staff entity found with the given `staffID`. |

### List Staff Accounts

- **Method**: `GET`
- **Endpoint**: `/staff/{staffID}/accounts`
- **Description**: Retrieve all user accounts bound to the specified staff entity through staff memberships.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Required | Description                 |
|-----------|---------------|----------|-----------------------------|
| `staffID` | UUID (string) | Yes      | The ID of the target staff. |

#### Response `200 OK`

```json
{
  "staff_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
  "total": 2,
  "accounts": [
    {
      "account_id": "a1b2c3d4-087b-48e7-b00a-d79dd092e8d4",
      "user_id": "56d88b08-ad99-4c91-9571-15b5bae95591",
      "email": "jane@chia.florist",
      "name": "Jane Doe",
      "username": "janedoe",
      "phone": "+628123456789",
      "avatar_url": "https://example.com/avatar.png",
      "role": {
        "id": "e0b8529e-2495-46aa-b2b7-a36c92fcbe09",
        "code": "staff_admin",
        "name": "Staff Admin"
      },
      "last_login_at": "2026-08-15T15:10:00Z",
      "created_at": "2026-06-26T19:47:28Z"
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No staff entity found with the given `staffID`. |

### Add Staff Account

- **Method**: `POST`
- **Endpoint**: `/staff/{staffID}/accounts`
- **Description**: Register and assign a new account to a staff. The actor must be an admin of the target staff.
- **Authentication**: Staff Admin
- **Request Body**:

  ```json
  {
    "email":    "string (required)",
    "password": "string (required)"
  }
  ```

#### Path Parameters

| Parameter | Type          | Required | Description                 |
|-----------|---------------|----------|-----------------------------|
| `staffID` | UUID (string) | Yes      | The ID of the target staff. |

#### Response `201 Created`

```json
{ "message": "staff account successfully created" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` is not a valid UUID, or `email` or `password` is empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No staff entity found with the given `staffID`. |
| `409 Conflict`     | An account with this email already exists. |

### Remove Staff Account

- **Method**: `DELETE`
- **Endpoint**: `/staff/{staffID}/accounts/{accountID}`
- **Description**: Unbind and remove a specific user account from a staff unit. Prevents self-removal and primary staff owner removal.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Path Parameters

| Parameter   | Type          | Required | Description                    |
|-------------|---------------|----------|--------------------------------|
| `staffID`   | UUID (string) | Yes      | The ID of the target staff.    |
| `accountID` | UUID (string) | Yes      | The ID of the account to unbind. |

#### Response `200 OK`

```json
{
  "message": "staff account successfully removed"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `staffID` or `accountID` is not a valid UUID, or actor is attempting self-removal. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role, or attempting to remove primary staff owner. |
| `404 Not Found`    | No staff entity or account membership found. |


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
- **Description**: Retrieve a list of payment methods with optional custom sorting. Includes the payment instruction configuration for each payment method if present.
- **Authentication**: Staff
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `sort`    | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

#### Sort Fields

| Field      | Example            | Description                       |
|------------|--------------------|-----------------------------------|
| `latest`   | `sort=latest:desc` | Sort by creation date.            |
| `name`     | `sort=name:asc`    | Sort alphabetically by name.      |
| `code`     | `sort=code:desc`   | Sort by payment method code.      |
| `type`     | `sort=type:asc`    | Sort by payment method type.      |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=name:asc,latest:desc`.

#### Response `200 OK`

```json
{
    "methods": [
        {
            "id": "0137d751-5188-447a-b630-1bf858f4f866",
            "name": "QRIS",
            "code": "qris",
            "provider": "midtrans",
            "type": "qr_code",
            "is_active": true,
            "description": "QRIS payment via Midtrans",
            "fee_type": "",
            "fee_fixed": 0,
            "fee_percentage": 0,
            "instruction": null
        },
        {
            "id": "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
            "name": "Mandiri",
            "code": "mandiri",
            "provider": "manual",
            "type": "bank_transfer",
            "is_active": true,
            "description": "Manual Bank Transfer to Mandiri",
            "fee_type": "flat",
            "fee_fixed": 2500,
            "fee_percentage": 0,
            "instruction": {
                "id": "21213fdf-7cf2-4354-bf31-a288a6706c41",
                "content": "Please transfer exactly **Rp {{amount}}** to Mandiri Account **{{va_number}}**.",
                "created_at": "2026-07-10T14:50:00Z"
            }
        }
    ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have a staff role. |

### Save Payment Method

- **Method**: `POST`
- **Endpoint**: `/payments/methods`
- **Description**: Save payment method (creates a new payment method or updates an existing one if the `id` is provided).
- **Authentication**: Staff Admin

- **Request Body**:

  ```json
  {
    "id": "string (optional, UUID)",
    "name": "string (required)",
    "code": "string (required)",
    "provider": "string (required)",
    "type": "string (required)",
    "is_active": "string (required)",
    "description": "string (required)",
    "fee_type": "string (required)",
    "fee_amount": "string (optional)",
    "fee_percentage": "string (optional)"
  }
  ```

#### Important Notes

> Account number is optional but required for creating account that based on `bank_transfer` payment method.

#### Response `200 OK`

```json
{
    "message": "payment method successfully created"
}
```

```json
{
    "message": "payment method successfully updated"
}
```

### Save Payment Instruction

- **Method**: `POST`
- **Endpoint**: `/payments/methods/{methodID}/instruction`
- **Description**: Create or update the payment instruction details for a specific payment method. If an instruction already exists for the given payment method, calling this endpoint will update its content.
- **Authentication**: Staff Admin
- **Path Parameters**:

| Parameter  | Type | Required | Description |
|------------|------|----------|-------------|
| `methodID` | UUID | Yes      | The unique ID of the payment method. |

- **Request Body**:

  ```json
  {
    "content": "string (required, markdown content)"
  }
  ```

#### Response `200 OK`

```json
{
    "message": "payment instruction successfully saved"
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `methodID` is not a valid UUID, request body is invalid, or `content` is empty. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | Payment method with the given ID does not exist. |

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
      "customer_id": "8ce91a56-deea-46ac-9330-de65d64daa32",
      "address_id": "48956fd0-bcea-44a2-b598-af999d7abc7a",
      "status": "pending",
      "subtotal": 150000,
      "shipping_fee": 15000,
      "total": 165000,
      "created_at": "2026-06-21T08:45:00Z",
      "updated_at": null,
      "address": {
        "id": "48956fd0-bcea-44a2-b598-af999d7abc7a",
        "customer_id": "8ce91a56-deea-46ac-9330-de65d64daa32",
        "receiver_name": "Jane Doe",
        "phone": "+628123456789",
        "is_default": true,
        "province_id": "31",
        "city_id": "3173",
        "district_id": "317305",
        "village_id": "3173051003",
        "full_address": "Jl. Mawar No. 12",
        "postal_code": "11530"
      },
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
- **Endpoint**: `/orders/{orderID}`
- **Description**: Retrieve the order by its ID. This API is meant for staff administrators to view order details, including payment and shipment tracking history.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                   |
|-----------|---------------|-------------------------------|
| `orderID` | UUID (string) | The ID of the order to query. |

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
 },
 "shipment": {
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "in_transit",
  "fulfillment_method": "courier",
  "courier": "jne",
  "service": "REG",
  "tracking_number": "JNE001928374",
  "cost": 770000,
  "created_at": "2026-06-27T10:09:05.904647Z",
  "events": [
   {
    "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "status": "packed",
    "description": "Order packed",
    "location": "Jakarta Store",
    "timestamp": "2026-06-27T11:00:00Z"
   },
   {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "status": "picked_up",
    "description": "Picked up by courier",
    "location": "Jakarta Store",
    "timestamp": "2026-06-27T13:30:00Z"
   }
  ]
 }
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `orderID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | Order not found. |

### Get Order Tracking Timeline (Customer Endpoint)

- **Method**: `GET`
- **Endpoint**: `/users/me/orders/{orderID}/tracking`
- **Description**: Retrieve the chronological tracking timeline of a shipment, merging internal shop events (e.g. `packed`, `picked_up`) with external courier updates (e.g. transit manifests from Komerce) when available. (Customer-facing endpoint documented here for completeness).
- **Authentication**: Customer
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                     |
|-----------|---------------|---------------------------------|
| `orderID` | UUID (string) | The ID of the order to track.   |

#### Response `200 OK`

```json
{
  "order_id": "f1e2d3c4-b5a6-7890-fedc-ba0987654321",
  "shipment_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "courier": "jne",
  "tracking_number": "JNE001928374",
  "timeline": [
    {
      "status": "packed",
      "description": "Order packed and ready for courier pickup",
      "location": "Jakarta Store",
      "timestamp": "2026-07-14T02:00:00Z"
    },
    {
      "status": "manifested",
      "description": "Shipment booked",
      "location": "Jakarta Hub",
      "timestamp": "2026-07-14T03:00:00Z"
    },
    {
      "status": "transit",
      "description": "Departed JNE Hub",
      "location": "Bekasi Office",
      "timestamp": "2026-07-14T04:30:00Z"
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `orderID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Order or shipment not found, or order does not belong to the authenticated customer. |

## Shipment Management

### Update Shipment Status

- **Method**: `PATCH`
- **Endpoint**: `/shipments/{shipmentID}/status`
- **Description**: Update the status of a shipment. Transitions the shipment status and logs a shipment event. If the shipment status transitions to `delivered`, the parent order's status is also automatically updated to `delivered`.
- **Authentication**: Staff
- **Path Parameters**:

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `shipmentID` | `string (UUID)` | Yes | The ID of the shipment to update. |

- **Request Body**:

```json
{
  "status": "string (required, enum: created, packed, labelled, picked_up, in_transit, out_for_delivery, delivered, failed, returned, cancelled)",
  "description": "string (optional)",
  "location": "string (optional)"
}
```

#### Response `200 OK`

```json
{
  "id": "7ca19532-6bb0-47b8-936a-2ee3d6790b9b",
  "order_id": "e4a31771-4638-4e89-a292-624e723927d1",
  "status": "packed",
  "fulfillment_method": "courier",
  "tracking_number": "JNE123456789",
  "courier": "jne",
  "service": "REG",
  "cost": 15000,
  "weight": 1000,
  "created_at": "2026-07-13T18:34:00Z"
}
```

#### Error Responses

| Status | Condition |
| :--- | :--- |
| `400 Bad Request` | Missing `status` in body, invalid `shipmentID` format, or invalid request payload. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden` | Authenticated user does not have the staff role. |
| `404 Not Found` | Shipment not found. |
| `422 Unprocessable Entity` | Invalid status transition (e.g., trying to transition from a terminal state or moving backwards). |

### Update Shipment

- **Method**: `PATCH`
- **Endpoint**: `/shipments/{shipmentID}`
- **Description**: Update a shipment's metadata such as the tracking number, courier code, or service name. Typically used in manual logistics mode when tracking numbers or waybill information becomes available later.
- **Authentication**: Staff
- **Path Parameters**:

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `shipmentID` | `string (UUID)` | Yes | The ID of the shipment to update. |

- **Request Body**:

```json
{
  "tracking_number": "string (optional)",
  "courier": "string (optional)",
  "service": "string (optional)"
}
```

#### Response `200 OK`

```json
{
  "id": "7ca19532-6bb0-47b8-936a-2ee3d6790b9b",
  "order_id": "e4a31771-4638-4e89-a292-624e723927d1",
  "status": "created",
  "fulfillment_method": "courier",
  "tracking_number": "JNE999888777",
  "courier": "jne",
  "service": "YES",
  "cost": 15000,
  "weight": 1000,
  "created_at": "2026-07-13T18:34:00Z"
}
```

#### Error Responses

| Status | Condition |
| :--- | :--- |
| `400 Bad Request` | Invalid `shipmentID` format, or invalid request payload (e.g. invalid fields). |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden` | Authenticated user does not have the staff role. |
| `404 Not Found` | Shipment not found. |

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

## WAF Security Policy

All WAF management endpoints are scoped to the `/api/` prefix and are protected by **Staff Admin** authentication.
They are intentionally bypassed by the WAF middleware itself — the WAF only inspects traffic outside `/api/`.

### WAF Rules

#### List Rules

- **Method**: `GET`
- **Endpoint**: `/api/rules`
- **Description**: Retrieve all configured WAF detection rules.
- **Authentication**: Staff Admin
- **Request Body**: None

##### Response `200 OK`

```json
{
  "rules": [
    {
      "id": "1000",
      "description": "SQL Injection Detection (Basic)",
      "pattern": "(?i)(union\\s+select|select\\s+.*\\s+from|delete\\s+from|drop\\s+table)",
      "tags": ["sqli", "owasp"],
      "impact": "high",
      "enabled": true,
      "created_at": "2026-07-05T10:00:00Z",
      "updated_at": "2026-07-05T10:00:00Z"
    }
  ]
}
```

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

#### Create Rule

- **Method**: `POST`
- **Endpoint**: `/api/rules`
- **Description**: Add a new WAF detection rule. The `id` is auto-generated. The rule is enabled by default.
- **Authentication**: Staff Admin
- **Request Body**:

  ```json
  {
    "description": "string (required)",
    "pattern":     "string (required) — must be a valid Go regular expression",
    "tags":        ["string"],
    "impact":      "string (optional) — e.g. \"high\", \"medium\", \"low\""
  }
  ```

##### Response `201 Created`

```json
{
  "id": "1001",
  "description": "XSS Detection",
  "pattern": "(?i)<script",
  "tags": ["xss"],
  "impact": "medium",
  "enabled": true,
  "created_at": "2026-07-05T12:00:00Z",
  "updated_at": "2026-07-05T12:00:00Z"
}
```

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `description` or `pattern` is empty, or `pattern` is not a valid regular expression. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

#### Toggle Rule

- **Method**: `PUT`
- **Endpoint**: `/api/rules/{id}`
- **Description**: Enable or disable an existing WAF rule by ID.
- **Authentication**: Staff Admin
- **Request Body**:

  ```json
  {
    "enabled": true
  }
  ```

##### Path Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | string | Yes      | The rule ID (e.g. `"1001"`). |

##### Response `204 No Content`

Empty body.

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` path parameter is missing. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No rule with the given `id` exists. |

#### Delete Rule

- **Method**: `DELETE`
- **Endpoint**: `/api/rules/{id}`
- **Description**: Permanently remove a WAF detection rule by ID.
- **Authentication**: Staff Admin
- **Request Body**: None

##### Path Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | string | Yes      | The rule ID to delete. |

##### Response `204 No Content`

Empty body.

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `id` path parameter is missing. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `404 Not Found`    | No rule with the given `id` exists. |

### IP Access Control

#### List IP Config

- **Method**: `GET`
- **Endpoint**: `/api/ip`
- **Description**: Retrieve all IP addresses that are currently banned, whitelisted, or ignored by the WAF.
- **Authentication**: Staff Admin
- **Request Body**: None

##### Response `200 OK`

```json
{
  "entries": [
    { "ip": "1.2.3.4",   "status": "banned",      "reason": "Auto-Banned: SQL Injection" },
    { "ip": "10.0.0.1",  "status": "whitelisted",  "reason": "" },
    { "ip": "192.168.1.5", "status": "ignored",    "reason": "" }
  ]
}
```

##### IP Status Values

| Status        | Effect |
|---------------|--------|
| `banned`      | All requests from this IP are blocked with `403 Forbidden`. |
| `whitelisted` | All WAF rule and keyword checks are skipped for this IP. |
| `ignored`     | Requests pass through; audit logging is suppressed for this IP. |

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

#### Update IP Action

- **Method**: `POST`
- **Endpoint**: `/api/ip`
- **Description**: Set or clear the access control status for an IP address. This is idempotent — submitting the same action twice is safe.
- **Authentication**: Staff Admin
- **Request Body**:

  ```json
  {
    "ip":     "string (required) — the IP address to act on",
    "action": "string (required) — one of: \"ban\", \"whitelist\", \"ignore\", \"reset\"",
    "reason": "string (optional) — only meaningful when action is \"ban\""
  }
  ```

##### Action Values

| Action        | Effect |
|---------------|--------|
| `ban`         | Blocks the IP. Sets status to `banned`. |
| `whitelist`   | Bypasses WAF checks for the IP. Sets status to `whitelisted`. |
| `ignore`      | Suppresses audit logging for the IP. Sets status to `ignored`. |
| `reset`       | Removes the IP from all access control lists entirely. |

> **Note**: An IP can only hold one status at a time. Submitting a new action always replaces the previous one.

##### Response `204 No Content`

Empty body.

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `ip` or `action` is missing, or `action` is not one of the valid values. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Filters

#### Get Filters

- **Method**: `GET`
- **Endpoint**: `/api/filters`
- **Description**: Retrieve the current keyword blocklist and URL whitelist used by the WAF payload inspector.
- **Authentication**: Staff Admin
- **Request Body**: None

##### Response `200 OK`

```json
{
  "keywords":         ["xp_cmdshell", "eval(", "<script"],
  "whitelisted_urls": ["/health", "/public/"]
}
```

##### Field Descriptions

| Field             | Description |
|-------------------|-------------|
| `keywords`        | Strings whose presence anywhere in the request payload (path, query, headers, body) causes an immediate block. Case-insensitive. |
| `whitelisted_urls` | URL path prefixes that bypass all WAF rule and keyword evaluation. Matched with `strings.HasPrefix`. |

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

#### Update Filter

- **Method**: `POST`
- **Endpoint**: `/api/filters`
- **Description**: Add or remove a single keyword or URL whitelist entry.
- **Authentication**: Staff Admin
- **Request Body**:

  ```json
  {
    "type":   "string (required) — \"keyword\" or \"url\"",
    "action": "string (required) — \"add\" or \"remove\"",
    "value":  "string (required) — the keyword or URL prefix to add/remove"
  }
  ```

##### Examples

```json
// Block a new keyword
{ "type": "keyword", "action": "add",    "value": "xp_cmdshell" }

// Remove an existing keyword
{ "type": "keyword", "action": "remove", "value": "xp_cmdshell" }

// Whitelist a URL prefix (bypasses WAF checks)
{ "type": "url",     "action": "add",    "value": "/webhooks/" }

// Remove a URL whitelist entry
{ "type": "url",     "action": "remove", "value": "/webhooks/" }
```

> **Note**: Adding a `keyword` entry is idempotent (duplicate inserts are silently ignored).

##### Response `204 No Content`

Empty body.

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `value` is empty, `type` is not `keyword` or `url`, or `action` is not `add` or `remove`. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Threat Intelligence

#### Analyze IP

- **Method**: `GET`
- **Endpoint**: `/api/analyze/{ip}`
- **Description**: Proxy query to VirusTotal v3 API to fetch the security reputation analysis of a given IP address.
- **Authentication**: Staff Admin
- **Headers**:
  - `X-VT-Key`: `string (optional) — optional custom VirusTotal API key to overwrite the server-configured API key`
- **Request Body**: None

##### Response `200 OK`
Returns the raw JSON response payload from VirusTotal's `/v3/ip_addresses/{ip}` endpoint. Example response snippet:

```json
{
  "data": {
    "id": "8.8.8.8",
    "type": "ip_address",
    "attributes": {
      "as_owner": "Google LLC",
      "asn": 15169,
      "last_analysis_stats": {
        "harmless": 76,
        "malicious": 0,
        "suspicious": 0,
        "undetected": 11
      }
    }
  }
}
```

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid IP address format, or VirusTotal API key is not configured/supplied. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `500 Internal Server Error` | Threat intelligence provider connection failed. |

#### Get Geolocation

- **Method**: `GET`
- **Endpoint**: `/api/geo/{ip}`
- **Description**: Proxy query to resolve the geolocation details of a given IP address. Uses `ip2location.io` if configured with a key, falling back automatically to a free `ip-api.com` lookup (mapping response fields for UI compatibility).
- **Authentication**: Staff Admin
- **Request Body**: None

##### Response `200 OK`
Returns geolocation parameters mapped to standard IP2Location fields.

```json
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "country_name": "United States of America",
  "region_name": "California",
  "city_name": "Mountain View",
  "latitude": 37.405992,
  "longitude": -122.078515,
  "zip_code": "94043",
  "time_zone": "-07:00",
  "asn": "15169",
  "as": "Google LLC"
}
```

##### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid IP address format. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
| `500 Internal Server Error` | Geolocation provider connection failed. |

## Analytics

### Get Order Metrics

- **Method**: `GET`
- **Endpoint**: `/analytics/orders`
- **Description**: Retrieve order performance summary KPI card data, time-series breakdown, top products, and top shops over a specified date range.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter     | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `from`        | string | No       | Start date in RFC3339 (`2026-01-01T00:00:00Z`) or `YYYY-MM-DD` (`2026-01-01`). Defaults to 30 days ago. |
| `to`          | string | No       | End date in RFC3339 (`2026-07-31T23:59:59Z`) or `YYYY-MM-DD` (`2026-07-31`). Defaults to current time. |
| `granularity` | string | No       | Time-series grouping. Allowed values: `daily` (default), `weekly`, `monthly`. |
| `shop_id`     | UUID   | No       | Filter metrics by specific shop ID. |
| `top_n`       | int    | No       | Number of top products and shops to return. Defaults to `10`. |

#### Response `200 OK`

```json
{
  "summary": {
    "total_orders": 150,
    "total_gmv": 45000000,
    "total_revenue": 42000000,
    "total_shipping_fee": 3000000,
    "aov": 300000,
    "cancellation_rate": 0.0667,
    "pending_count": 10,
    "confirmed_count": 20,
    "processing_count": 30,
    "shipped_count": 40,
    "delivered_count": 40,
    "cancelled_count": 10
  },
  "time_series": [
    {
      "date": "2026-07-01T00:00:00Z",
      "order_count": 5,
      "gmv": 1500000,
      "aov": 300000
    }
  ],
  "top_products": [
    {
      "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "product_name": "Anniversary Bouquet",
      "quantity": 25,
      "revenue": 7500000
    }
  ],
  "top_shops": [
    {
      "shop_id": "427db07b-dbad-43ee-9199-49a2849a4e30",
      "shop_name": "Chia Bogor",
      "revenue": 20000000,
      "orders": 65
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid date format or invalid `shop_id` UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Get Payment Metrics

- **Method**: `GET`
- **Endpoint**: `/analytics/payments`
- **Description**: Retrieve payment summary statistics (total paid, pending, expired, refunded, payment success rate, avg time to pay) and breakdown by payment method.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `from`    | string | No       | Start date in RFC3339 or `YYYY-MM-DD`. Defaults to 30 days ago. |
| `to`      | string | No       | End date in RFC3339 or `YYYY-MM-DD`. Defaults to current time. |

#### Response `200 OK`

```json
{
  "summary": {
    "total_paid": 42000000,
    "total_pending": 2000000,
    "total_expired": 1000000,
    "total_refunded": 500000,
    "payment_success_rate": 0.9333,
    "avg_time_to_pay": 145.5
  },
  "breakdown": [
    {
      "method_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "method_name": "BCA Virtual Account",
      "method_type": "bank_transfer",
      "count": 100,
      "amount": 30000000,
      "success_rate": 0.95
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid date format. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Get Shipment Metrics

- **Method**: `GET`
- **Endpoint**: `/analytics/shipments`
- **Description**: Retrieve shipment status summary (delivered, failed, returned, cancelled), overall delivery rate, average fulfillment duration in seconds, and breakdown by courier/service.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `from`    | string | No       | Start date in RFC3339 or `YYYY-MM-DD`. Defaults to 30 days ago. |
| `to`      | string | No       | End date in RFC3339 or `YYYY-MM-DD`. Defaults to current time. |
| `top_n`   | int    | No       | Number of top couriers to return. Defaults to `10`. |

#### Response `200 OK`

```json
{
  "summary": {
    "total": 120,
    "delivered": 110,
    "failed": 3,
    "returned": 2,
    "cancelled": 5,
    "delivery_rate": 0.9167,
    "avg_fulfillment_sec": 7200.0
  },
  "couriers": [
    {
      "courier": "JNE",
      "service": "REG",
      "count": 80,
      "delivery_rate": 0.9375,
      "avg_cost": 25000
    }
  ]
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid date format. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Get Inventory Metrics

- **Method**: `GET`
- **Endpoint**: `/analytics/inventory`
- **Description**: Retrieve inventory performance summary across all products or a specific shop, including total stock, reserved, available, stockout count, and low-stock count.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter             | Type   | Required | Description |
|-----------------------|--------|----------|-------------|
| `shop_id`             | UUID   | No       | Filter inventory metrics by specific shop ID. |
| `low_stock_threshold` | int    | No       | Threshold for low stock warning count. Defaults to `5`. |

#### Response `200 OK`

```json
{
  "total_products": 45,
  "total_stock": 2500,
  "total_reserved": 150,
  "total_available": 2350,
  "stockout_count": 2,
  "low_stock_count": 4
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid `shop_id` UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |

### Get Product Metrics

- **Method**: `GET`
- **Endpoint**: `/analytics/products`
- **Description**: Retrieve top products ranked by revenue and units sold, along with gross margin %, 7-day and 30-day sales velocity, and overall invoice void rate.
- **Authentication**: Staff Admin
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `from`    | string | No       | Start date in RFC3339 or `YYYY-MM-DD`. Defaults to 30 days ago. |
| `to`      | string | No       | End date in RFC3339 or `YYYY-MM-DD`. Defaults to current time. |
| `top_n`   | int    | No       | Number of top products to return. Defaults to `10`. |

#### Response `200 OK`

```json
{
  "top_by_revenue": [
    {
      "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "product_name": "Grand Opening Stand",
      "revenue": 15000000,
      "units_sold": 100,
      "conversion_rate": 0.0,
      "return_rate": null,
      "gross_margin_pct": 42.5,
      "sales_velocity_7d": 25,
      "sales_velocity_30d": 100
    }
  ],
  "top_by_volume": [
    {
      "product_id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
      "product_name": "Single Red Rose",
      "revenue": 5000000,
      "units_sold": 200,
      "conversion_rate": 0.0,
      "return_rate": null,
      "gross_margin_pct": 50.0,
      "sales_velocity_7d": 50,
      "sales_velocity_30d": 200
    }
  ],
  "avg_conversion": 0.0,
  "avg_return_rate": 0.0,
  "invoice_void_rate": 0.015
}
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid date format. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Authenticated user does not have the staff admin role. |
