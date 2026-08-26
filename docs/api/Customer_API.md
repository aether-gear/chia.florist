# Customer API Documentation

This document covers the customer-facing APIs for the Chia Florist service.
Endpoints are organized by access level: **Public** and **Authenticated Customer**.

## TODO

- [ ] Public API
  - [X] Authentication
    - [X] Sign Up
    - [X] Verify Account
    - [X] Sign In
    - [ ] Google Sign In
    - [ ] Forgot Password
    - [ ] Verify Forgot Password
    - [ ] Reset Password
  - [X] Products
    - [X] Find Products
    - [X] Get Product Detail
  - [ ] Shops
    - [ ] Find Shops
    - [ ] Get Shop
    - [ ] Get Shop Addresses
    - [ ] Get Shop Couriers
    - [ ] Get Shop Products
  - [X] Locations
    - [X] List Provinces
    - [X] List Cities by Province
    - [X] List Districts by City
    - [X] List Villages by District

- [ ] Authenticated Customer API
  - [X] Authentication
    - [X] Me
    - [X] Log Out
  - [ ] Profile
    - [ ] Get Current User
    - [ ] Update User
    - [ ] List Orders
    - [ ] Delete Account
  - [X] Addresses
    - [X] List My Addresses
    - [X] Save My Address
    - [X] Delete My Address
  - [X] Cart
    - [X] Get Cart
    - [X] Add Item
    - [X] Update Item
    - [X] Remove Item
    - [X] Remove Custom Item
  - [ ] Checkout
    - [ ] Estimate Checkout
    - [ ] Calculate Checkout
  - [ ] Orders
    - [ ] Create Order
    - [ ] Find My Orders
    - [ ] Get Order
    - [ ] Get Order Payment Details
    - [X] Check Order Payment Status
    - [X] Get Order Tracking Timeline


# Public API

No authentication is required for these endpoints.

## Authentication

### Sign Up

- **Method**: `POST`
- **Endpoint**: `/auth/signup`
- **Description**: Register a new customer account. A verification OTP is sent to the provided email. The returned `challenge_id` is required to complete registration via the Verify Account endpoint.
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "name":     "string (optional)",
    "username": "string (required)",
    "email":    "string (required)",
    "password": "string (required)",
    "phone":    "string (optional)"
  }
  ```

#### Response `201 Created`

```json
{
  "message":      "verification code sent",
  "challenge_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### Error Responses

| Status            | Condition |
|-------------------|-----------|
| `400 Bad Request` | `email`, `password`, or `username` is missing or empty. |

### Verify Account

- **Method**: `POST`
- **Endpoint**: `/auth/verify`
- **Description**: Complete account registration by submitting the OTP received via email. Sets a session cookie on success.
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "challenge_id": "string (UUID, required)",
    "otp":          123456,
    "user_agent":   "string (optional)",
    "ip_address":   "string (optional)"
  }
  ```

#### Important Notes

> `otp` must be exactly **6 digits** as an integer (e.g. `123456`).
> `challenge_id` is the UUID returned from the Sign Up response.

#### Response `201 Created`

- **Set-Cookie**: `<access_token_cookie>=<value>`
- **Body**:
  ```json
  { "message": "verify success" }
  ```

#### Error Responses

| Status            | Condition |
|-------------------|-----------|
| `400 Bad Request` | `challenge_id` is missing, not a valid UUID, or `otp` is not a 6-digit number. |
| `404 Not Found`   | No pending verification challenge found for the given `challenge_id`. |

### Sign In

- **Method**: `POST`
- **Endpoint**: `/auth/signin`
- **Description**: Authenticate an existing customer account using email and password. Sets a session cookie on success.
- **Authentication**: None
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

- **Set-Cookie**: `chast=<value>`
- **Body**:
  ```json
  { "message": "login success" }
  ```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `email` or `password` is missing or empty. |
| `401 Unauthorized` | Invalid email or password. |

### Google Sign In

- **Method**: `GET`
- **Endpoint**: `/auth/google/login`
- **Description**: Authenticate a user with _Google_ account, for both sign up and sign in.
- **Authentication**: None
- **Request Body**: None

#### Response `200 OK`

- **Set-Cookie**: `chast=<value>`
- **Body**:
  ```json
  { "message": "login success" }
  ```

### Forgot Password

- **Method**: `POST`
- **Endpoint**: `/forgot-password`
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

## Products

### Find Products

- **Method**: `GET`
- **Endpoint**: `/products`
- **Description**: Retrieve a paginated list of products with optional filtering and sorting.
- **Authentication**: None
- **Request Body**: None

#### Query Parameters

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `id`        | UUID   | No       | Filter products by product ID. |
| `name`      | string | No       | Search products by name. |
| `shop_id`   | UUID   | No       | Filter products stocked/available at a specific shop ID. |
| `shop_slug` | string | No       | Filter products stocked/available at a specific shop slug. |
| `page`      | int    | No       | Page number. Defaults to `1`. |
| `limit`     | int    | No       | Number of results per page. |
| `sort`      | string | No       | Comma-separated sort expressions. Format: `<field>:<direction>`. |

> **Store Selection & Header Context**: You can pass `X-Shop-ID: <UUID>` header as an alternative to specify store context for catalog browsing. If neither query parameters nor the header are provided, products across all shops will be returned.

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
- `GET /products?name=anniversary`
- `GET /products?shop_id=c3d4e5f6-a7b8-9012-cdef-123456789012`
- `GET /products?shop_slug=chia-medan-satria`
- `GET /products?sort=price:asc`

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
- **Description**: Retrieve full details of a specific product by its slug.
- **Authentication**: None
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

#### Error Responses

| Status          | Condition                          |
|-----------------|------------------------------------|
| `404 Not Found` | No product found for the given slug. |

## Shops

### Find Shops

- **Method**: `GET`
- **Endpoint**: `/shops`
- **Description**: Retrieve a paginated list of shops.
- **Authentication**: None
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `id`      | UUID   | No       | Filter by exact shop ID. |
| `name`    | string | No       | Filter by shop name. |
| `active`  | bool   | No       | Filter shops by active status (`true` for store picker). |
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`    | string | No       | Comma-separated sort expressions. |

#### Response `200 OK`

```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "shops": [
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "name": "Chia Medan Satria",
      "slug": "chia-medan-satria",
      "description": "Our Medan Satria branch.",
      "is_active": true,
      "created_at": "2026-01-01T08:00:00Z",
      "updated_at": "2026-05-10T09:00:00Z"
    }
  ]
}
```

### Get Shop

- **Method**: `GET`
- **Endpoint**: `/shops/{shopID}`
- **Description**: Retrieve details of a specific shop by its ID.
- **Authentication**: None
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                |
|-----------|---------------|----------------------------|
| `shopID`  | UUID (string) | The unique ID of the shop. |

#### Response `200 OK`

```json
{
  "shop": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "name": "Chia Medan Satria",
    "slug": "chia-medan-satria",
    "description": "Our Medan Satria branch.",
    "is_active": true,
    "created_at": "2026-01-01T08:00:00Z",
    "updated_at": "2026-05-10T09:00:00Z"
  }
}
```

#### Error Responses

| Status            | Condition                                 |
|-------------------|-------------------------------------------|
| `400 Bad Request` | `shopID` in the path is not a valid UUID. |
| `404 Not Found`   | No shop found for the given ID.           |

### Get Shop Addresses

- **Method**: `GET`
- **Endpoint**: `/shops/{shopID}/addresses`
- **Description**: Retrieve all addresses registered for a specific shop.
- **Authentication**: None
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
- **Authentication**: None
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
- **Authentication**: None
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

## Locations

> All location endpoints follow a cascading hierarchy: **Province → City → District → Village**.
> Use the `id` from each response as the path parameter for the next level.

### List Provinces

- **Method**: `GET`
- **Endpoint**: `/provinces`
- **Description**: Retrieve all available provinces.
- **Authentication**: None
- **Request Body**: None

#### Response `200 OK`

```json
{
  "provinces": [
    { "id": "32", "name": "Jawa Barat" },
    { "id": "31", "name": "DKI Jakarta" }
  ]
}
```

### List Cities by Province

- **Method**: `GET`
- **Endpoint**: `/provinces/{id}/cities`
- **Description**: Retrieve all cities within a given province.
- **Authentication**: None
- **Request Body**: None

#### Path Parameters

| Parameter | Type   | Description                          |
|-----------|--------|--------------------------------------|
| `id`      | string | Province ID from List Provinces. |

#### Response `200 OK`

```json
{
  "cities": [
    { "id": "3204", "province_id": "32", "name": "Kabupaten Bandung" }
  ]
}
```

### List Districts by City

- **Method**: `GET`
- **Endpoint**: `/cities/{id}/districts`
- **Description**: Retrieve all districts within a given city.
- **Authentication**: None
- **Request Body**: None

#### Path Parameters

| Parameter | Type   | Description                       |
|-----------|--------|-----------------------------------|
| `id`      | string | City ID from List Cities. |

#### Response `200 OK`

```json
{
  "districts": [
    { "id": "320401", "city_id": "3204", "name": "Ciwidey" }
  ]
}
```

### List Villages by District

- **Method**: `GET`
- **Endpoint**: `/districts/{id}/villages`
- **Description**: Retrieve all villages within a given district.
- **Authentication**: None
- **Request Body**: None

#### Path Parameters

| Parameter | Type   | Description                             |
|-----------|--------|-----------------------------------------|
| `id`      | string | District ID from List Districts. |

#### Response `200 OK`

```json
{
  "villages": [
    { "id": "3204010001", "district_id": "320401", "name": "Panyocokan" }
  ]
}
```

# Authenticated Customer API

These endpoints require a valid customer session set via the Sign In or Verify Account cookie.

## Authentication

### Me

- **Method**: `GET`
- **Endpoint**: `/auth/me`
- **Description**: Return the identity and roles of the currently authenticated account.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "account_id": "51e20db6-5bdb-4f6a-b2b7-8d40c0db857d",
  "account_type": "customer",
  "is_authenticated": true,
  "roles": [
    { "code": "customer", "name": "Customer" }
  ],
  "permissions": [
    { "code": "customer" }
  ]
}
```

#### Error Responses

| Status             | Condition                       |
|--------------------|---------------------------------|
| `401 Unauthorized` | Missing or invalid session.     |

### Log Out

- **Method**: `POST`
- **Endpoint**: `/auth/logout`
- **Description**: Invalidate the current session and clear the access token cookie.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{ "message": "logout success" }
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |

## Profile

### Get Current User

- **Method**: `GET`
- **Endpoint**: `/profile`
- **Description**: Retrieve the profile of the currently authenticated customer.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "profile": {
    "customer_id": "d0118ea9-ab28-4338-89dd-3d3ca5f89880",
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
- **Description**: Update the profile of the currently authenticated customer.
- **Authentication**: Customer
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
    "customer_id": "d0118ea9-ab28-4338-89dd-3d3ca5f89880",
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

### Delete Account

- **Method**: `DELETE`
- **Endpoint**: `/profile`
- **Description**: Delete the account of the currently authenticated customer.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "message": "account deleted successfully"
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Account must have type customer.     |

## Addresses

### List My Addresses

- **Method**: `GET`
- **Endpoint**: `/users/me/addresses`
- **Description**: Retrieve all saved delivery addresses for the authenticated customer.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "addresses": [
    {
      "address_id": "d4e5f6a7-b8c9-0123-defa-234567890123",
      "user_id":    "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "receiver_name": "Jane Doe",
      "phone":      "+6281234567890",
      "is_default": true,
      "province_id": "32",
      "city_id":     "3204",
      "district_id": "320401",
      "village_id":  "3204010001",
      "full_address": "Jl. Bunga Indah No. 10",
      "postal_code": "17520",
      "created_at":  "2026-01-15T09:00:00Z",
      "updated_at":  "2026-05-01T12:00:00Z"
    }
  ]
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |

### Save My Address

- **Method**: `POST`
- **Endpoint**: `/users/me/addresses`
- **Description**: Create a new delivery address or update an existing one. Omit `address_id` to create; supply it to update.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "address_id":    "string (UUID, optional — omit to create, supply to update)",
    "receiver_name": "string (optional)",
    "phone":         "string (optional)",
    "is_default":    "string (optional, e.g. \"true\" or \"false\")",
    "province_id":   "string (required)",
    "city_id":       "string (required)",
    "district_id":   "string (required)",
    "village_id":    "string (required)",
    "full_address":  "string (required)",
    "postal_code":   "string (required)"
  }
  ```

#### Important Notes

> `is_default` is a string-encoded boolean (`"true"` / `"false"`).
> Parsed with `strconv.ParseBool`, so `"1"`, `"0"`, `"TRUE"`, `"FALSE"` are also accepted.
>
> Use the Locations endpoints to look up valid `province_id`, `city_id`, `district_id`, and `village_id` values.

#### Response `200 OK`

```json
{ "message": "address saved successfully" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Any required location field is empty, `address_id` is not a valid UUID, or `is_default` cannot be parsed. |
| `401 Unauthorized` | Missing or invalid session. |

### Delete My Address

- **Method**: `DELETE`
- **Endpoint**: `/users/me/addresses/{addressID}`
- **Description**: Delete a saved delivery address belonging to the authenticated customer.
- **Authentication**: Customer
- **Request Body**: None

#### Path Parameters

| Parameter   | Type          | Description                       |
|-------------|---------------|-----------------------------------|
| `addressID` | UUID (string) | The ID of the address to delete.  |

#### Response `200 OK`

```json
{ "message": "address deleted successfully" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `addressID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Address not found. |

## Cart

### Get Cart

- **Method**: `GET`
- **Endpoint**: `/carts`
- **Description**: Retrieve the current cart of the authenticated customer, including standard and custom product items with subtotals and grand total.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "cart_id": "f1e2d3c4-b5a6-7890-fedc-ba0987654321",
  "total": 220000,
  "items": [
    {
      "cart_item_id": "d1a2b3c4-e5f6-7890-abcd-ef1234567890",
      "product_variant_type": "standard",
      "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "shop_id":    "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "shop_name":  "Chia Medan Satria",
      "shop_slug":  "chia-medan-satria",
      "name":       "Anniversary Flower Stand",
      "price":      110000,
      "quantity":   2,
      "subtotal":   220000,
      "item_options": {
        "size": "medium",
        "jambul": "top"
      },
      "images": {
        "thumbnail": "https://example.com/thumbnail.jpg"
      }
    },
    {
      "cart_item_id": "e5f6a7b8-c9d0-1234-ef56-789012345678",
      "product_variant_type": "custom",
      "product_id": null,
      "shop_id":    "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "shop_name":  "Chia Medan Satria",
      "shop_slug":  "chia-medan-satria",
      "name":       "(Custom Board)",
      "price":      0,
      "quantity":   1,
      "subtotal":   0,
      "images": {},
      "custom_design": {
        "metadata": {
          "version": "1.0.0",
          "editorVersion": "1.0.0",
          "platform": "web",
          "checksum": "a8f3b9e1"
        },
        "layout": {
          "physicalSizeId": "medium",
          "upperHeightRatio": 0.58
        },
        "sections": {
          "upper": {
            "header": { "text": "Selamat & Sukses", "fontId": "playfair" },
            "body": { "text": "Jane Doe", "fontId": "inter" }
          },
          "lower": {
            "body": { "text": "PT. Tech Nusantara", "fontId": "inter" }
          }
        },
        "assets": {
          "previewUrl": "https://chia-florist.supabase.co/storage/v1/object/public/custom-previews/designs/2026/08/order-board-a8f3b9e1.png"
        }
      }
    }
  ]
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Cart not found.             |

### Add Item

- **Method**: `POST`
- **Endpoint**: `/carts/items`
- **Description**: Add a standard catalog product or custom-designed flower board item to the customer's cart. Standard items support pick-style option selections (`size` and `jambul`) using a strictly non-negative additive pricing model.
- **Authentication**: Customer

#### Standard Product Pricing Rules

Standard board unit pricing is strictly additive and non-negative:
$$\text{UnitPrice} = \max(0, \text{base\_price} + \text{size\_addon} + \text{jambul\_addon})$$

- **Size Options**:
  - `small` (1.5 × 2.0m - Compact): `+ Rp 0` (Weight: 1,500g)
  - `medium` (1.8 × 2.5m - Standard): `+ Rp 50.000` (Weight: 2,500g)
  - `large` (2.0 × 3.0m - Grand): `+ Rp 100.000` (Weight: 4,000g)
- **Floral Crest (Jambul) Options**:
  - `none`: `+ Rp 0`
  - `top`: `+ Rp 25.000`
  - `bottom`: `+ Rp 25.000`
  - `both`: `+ Rp 50.000`

#### Request Body (Standard Product Item)

```json
{
  "product_variant_type": "standard",
  "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
  "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "quantity": 1,
  "item_options": {
    "size": "medium",
    "jambul": "both"
  }
}
```

#### Request Body (Custom Product Item)

```json
{
  "product_variant_type": "custom",
  "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
  "quantity": 1,
  "product_name": "Custom Board — Selamat & Sukses",
  "physical_size_id": "medium",
  "custom_design": {
    "metadata": {
      "version": "1.0.0",
      "editorVersion": "1.0.0",
      "platform": "web",
      "checksum": "a8f3b9e1"
    },
    "layout": {
      "physicalSizeId": "medium",
      "upperHeightRatio": 0.58
    },
    "sections": {
      "upper": {
        "header": { "text": "Selamat & Sukses", "fontId": "playfair" },
        "body": { "text": "Jane Doe", "fontId": "inter" }
      },
      "lower": {
        "body": { "text": "PT. Tech Nusantara", "fontId": "inter" }
      }
    },
    "assets": {
      "previewUrl": "https://chia-florist.supabase.co/storage/v1/object/public/custom-previews/designs/2026/08/order-board-a8f3b9e1.png"
    }
  }
}
```

#### Field Descriptions

| Field                  | Type   | Required | Description |
|------------------------|--------|----------|-------------|
| `product_variant_type` | string | No       | Variant type discriminator (`"standard"` or `"custom"`). Defaults to `"standard"`. |
| `product_id`           | string | Optional | Product UUID. Required when `product_variant_type` is `"standard"`. |
| `shop_id`              | string | Yes      | Shop UUID where the item is purchased. |
| `quantity`             | int    | Yes      | Quantity of items (`> 0`). |
| `item_options`         | object | Optional | Optional variant attributes (`size`: `"small"` \| `"medium"` \| `"large"`, `jambul`: `"none"` \| `"top"` \| `"bottom"` \| `"both"`). |
| `product_name`         | string | Optional | Human-readable title for custom products. Required when `product_variant_type` is `"custom"`. |
| `physical_size_id`     | string | Optional | Physical size identifier (e.g. `"medium"`). Required when `product_variant_type` is `"custom"`. |
| `custom_design`        | object | Optional | Full v1.0.0 canvas design JSON snapshot payload. Required when `product_variant_type` is `"custom"`. |

#### Response `200 OK`

```json
{ "message": "item added" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing required fields (`product_id` for standard, `custom_design` / `physical_size_id` for custom), invalid UUIDs, or `quantity <= 0`. |
| `401 Unauthorized` | Missing or invalid session. |

### Update Item

- **Method**: `PUT`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Update the quantity of a specific standard item in the cart.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "quantity": 3
  }
  ```

#### Path Parameters

| Parameter   | Type          | Description            |
|-------------|---------------|------------------------|
| `shopID`    | UUID (string) | The shop ID of the item.    |
| `productID` | UUID (string) | The product ID of the item. |

#### Response `200 OK`

```json
{ "message": "item updated" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `shopID` or `productID` is not a valid UUID, or `quantity` is `< 0`. |
| `401 Unauthorized` | Missing or invalid session. |

### Update Item by ID (Quantity and/or Options)

- **Method**: `PUT` or `PATCH`
- **Endpoint**: `/carts/items/{cartItemID}`
- **Description**: Update the quantity and/or `item_options` (size, jambul) of a specific cart item identified by its unique `cartItemID`. Supports both standard and custom items. If `item_options` is omitted, the item retains its existing options. If changing `item_options` produces a combination that already exists in the cart for the same product and shop, the quantities are automatically merged and the target item is soft-deleted.
- **Authentication**: Customer
- **Request Body (Quantity only — preserves current options)**:
  ```json
  {
    "quantity": 3
  }
  ```
- **Request Body (With optional options update)**:
  ```json
  {
    "quantity": 2,
    "item_options": {
      "size": "large",
      "jambul": "both"
    }
  }
  ```

#### Path Parameters

| Parameter    | Type          | Description                             |
|--------------|---------------|-----------------------------------------|
| `cartItemID` | UUID (string) | The unique cart item ID to update.      |

#### Request Fields

| Field                   | Type     | Required | Description                                                    |
|-------------------------|----------|----------|----------------------------------------------------------------|
| `quantity`              | integer  | Yes      | Updated item quantity (must be `> 0`).                         |
| `item_options`          | object   | No       | Optional variant options object. If omitted, existing options are preserved. |
| `item_options.size`     | string   | No       | `"small"` (+Rp 0), `"medium"` (+Rp 50.000), `"large"` (+Rp 100.000). |
| `item_options.jambul`   | string   | No       | `"none"` (+Rp 0), `"top"` (+Rp 25.000), `"bottom"` (+Rp 25.000), `"both"` (+Rp 50.000). |

#### Response `200 OK`

```json
{ "message": "item updated" }
```

#### Error Responses

| Status              | Condition                                                         |
|---------------------|-------------------------------------------------------------------|
| `400 Bad Request`   | Invalid `cartItemID`, missing request body, or `quantity <= 0`.  |
| `401 Unauthorized`  | Missing or invalid session.                                       |
| `404 Not Found`     | Cart or cart item not found.                                      |
| `409 Conflict`      | Insufficient stock at fulfillment shop across all styles.         |

### Change Item Shop

- **Method**: `PATCH`
- **Endpoint**: `/carts/items/{cartItemID}/shop`
- **Description**: Update the fulfillment shop for an existing cart item (standard or custom) prior to checkout. For standard items, stock at the target shop is validated. If the exact same product and style already exists under the target shop in the cart, quantities are merged; distinct styles remain separate items.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "shop_id": "c3d4e5f6-a7b8-9012-cdef-123456789012"
  }
  ```

#### Path Parameters

| Parameter    | Type | Description |
|--------------|------|-------------|
| `cartItemID` | UUID | The cart item ID to transfer to a new shop. |

#### Request Fields

| Field     | Type | Required | Description |
|-----------|------|----------|-------------|
| `shop_id` | UUID | Yes      | Target active shop ID. |

#### Response `200 OK`

```json
{ "message": "item shop updated" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Invalid `cartItemID` or `shop_id` UUID string format. |
| `401 Unauthorized` | Missing or invalid customer authentication session. |
| `404 Not Found`    | Cart item not found or target shop not found/inactive. |
| `409 Conflict`     | Target shop has insufficient stock for standard catalog product. |

### Remove Item by ID

- **Method**: `DELETE`
- **Endpoint**: `/carts/items/{cartItemID}`
- **Description**: Remove a specific item (standard or custom) from the cart by its unique cart item ID.
- **Authentication**: Customer
- **Request Body**: None

#### Path Parameters

| Parameter    | Type          | Description                         |
|--------------|---------------|-------------------------------------|
| `cartItemID` | UUID (string) | The unique ID of the item to remove. |

#### Response `200 OK`

```json
{ "message": "item removed" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `cartItemID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Cart item not found in cart. |

### Remove Item by Product & Shop (Legacy)

- **Method**: `DELETE`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Remove a standard item from the cart by product and shop ID. Optional query parameters `size` and `jambul` can be supplied to target a specific style when multiple styles of the same product exist.
- **Authentication**: Customer
- **Query Parameters**:
  - `size` (optional, string): e.g. `"small"`, `"medium"`, `"large"`.
  - `jambul` (optional, string): e.g. `"none"`, `"top"`, `"bottom"`, `"both"`.
- **Request Body**: None

#### Path Parameters

| Parameter   | Type          | Description                 |
|-------------|---------------|-----------------------------|
| `shopID`    | UUID (string) | The shop ID of the item.    |
| `productID` | UUID (string) | The product ID of the item. |

#### Response `200 OK`

```json
{ "message": "item removed" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `shopID` or `productID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |

### Remove Custom Item (Legacy)

- **Method**: `DELETE`
- **Endpoint**: `/carts/items/custom/{cartItemID}`
- **Description**: Legacy alias for removing a custom-designed item by its `cart_item_id`.
- **Authentication**: Customer
- **Request Body**: None

#### Path Parameters

| Parameter    | Type          | Description |
|--------------|---------------|-------------|
| `cartItemID` | UUID (string) | The unique `cart_item_id` of the custom item to remove. |

#### Response `200 OK`

```json
{ "message": "item removed" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `cartItemID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Custom item not found in cart. |

## Checkout

### Estimate Checkout

- **Method**: `POST`
- **Endpoint**: `/carts/checkout`
- **Description**: Place an order for the selected items. Courier selection per shop is optional at this stage.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "shops": [
      {
        "shop_id": "string (UUID, required)",
        "items": [
          {
            "product_variant_type": "standard",
            "product_id": "string (UUID, required for standard)",
            "quantity": 1,
            "item_options": {
              "size": "medium",
              "jambul": "top"
            }
          },
          {
            "product_variant_type": "custom",
            "product_name": "Custom Board — Selamat & Sukses",
            "physical_size_id": "medium",
            "quantity": 1,
            "custom_design": {
              "metadata": { "version": "1.0.0" },
              "layout": { "physicalSizeId": "medium" }
            }
          }
        ],
        "courier": {
          "code": "jne",
          "service": "REG"
        }
      }
    ]
  }
  ```

#### Important Notes

> The service will **use default address** of the user.
> The service will use the **least** cost of payment method and courier option.
> Only estimation grand total, use checkout calculate api to get real grand total

#### Response `200 OK`

```json
{
  "address": {
    "id": "95120305-02ad-468e-9103-2bb113d41cd7",
    "recipient_name": "My Belle Gweh",
    "phone": "000",
    "full_address": "Blok LA 22B, Jl Sokekarno Saya akan lawan"
  },
  "shops": [
    {
      "shop_id": "8fad2c68-82a2-4578-a550-c625a1691d8a",
      "name": "Chia Medan Satria",
      "slug": "chia-medan-satria",
      "subtotal": 21600000,
      "total": 21850000,
      "items": [
          {
            "product_id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
            "shop_id": "8fad2c68-82a2-4578-a550-c625a1691d8a",
            "name": "Grand Opening",
            "price": 150000,
            "quantity": 144,
            "subtotal": 21600000
          }
      ],
      "cost_couriers": [
        {
          "code": "tiki",
          "name": "Citra Van Titipan Kilat (TIKI)",
          "service": "TRC",
          "etd": "5 day",
          "fee": 50000
        },
        {
          "code": "jne",
          "name": "Jalur Nugraha Ekakurir (JNE)",
          "service": "REG",
          "etd": "1 day",
          "fee": 150000
        },
        {
          "code": "lion",
          "name": "Lion Parcel",
          "service": "REGPACK",
          "etd": "1-2 day",
          "fee": 144450
        },
        {
          "code": "pos",
          "name": "POS Indonesia (POS)",
          "service": "PAKETPOS VALUABLE GOODS",
          "etd": "3 day",
          "fee": 247500
        }
      ]
    },
    {
      "shop_id": "333f6432-a01c-412f-99f4-0f08ca0d8eb1",
      "name": "Chia Cipinang",
      "slug": "chia-cipinang",
      "subtotal": 825000,
      "total": 870000,
      "items": [
          {
            "product_id": "71be3ee1-17b4-4bb8-8f80-eae6ad93a844",
            "shop_id": "333f6432-a01c-412f-99f4-0f08ca0d8eb1",
            "name": "Graduate",
            "price": 55000,
            "quantity": 15,
            "subtotal": 825000
          }
        ],
        "cost_couriers": [
          {
            "code": "tiki",
            "name": "Citra Van Titipan Kilat (TIKI)",
            "service": "TRC",
            "etd": "5 day",
            "fee": 50000
          },
          {
            "code": "jne",
            "name": "Jalur Nugraha Ekakurir (JNE)",
            "service": "REG",
            "etd": "1 day",
            "fee": 150000
          },
          {
            "code": "lion",
            "name": "Lion Parcel",
            "service": "REGPACK",
            "etd": "1-2 day",
            "fee": 144450
          },
          {
            "code": "pos",
            "name": "POS Indonesia (POS)",
            "service": "PAKETPOS VALUABLE GOODS",
            "etd": "3 day",
            "fee": 247500
          }
        ]
    }
  ],
  "subtotal": 22425000,
  "total_shipping": 295000,
  "total": 22720000,
  "payment_methods": [
    {
      "id": "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
      "name": "GoPay",
      "type": "ewallet",
      "description": "GoPay via Midtrans",
      "fee": 0,
      "subtotal": 22720000,
      "total": 22720000
    },
    {
      "id": "074b02e4-e047-4f60-bdb0-cfeb5481d002",
      "name": "DANA",
      "type": "ewallet",
      "description": "DANA via Midtrans",
      "fee": 0,
      "subtotal": 22720000,
      "total": 22720000
    },
    {
      "id": "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
      "name": "Mandiri",
      "type": "bank_transfer",
      "description": "Mandiri Bill Payment via Midtrans",
      "fee": 0,
      "subtotal": 22720000,
      "total": 22720000
    }
  ]
}
```

#### Checkout Response Fields

| Field            | Type   | Description |
|------------------|--------|-------------|
| `address`        | object | Delivery address details. |
| `shops`          | array  | Per-shop breakdown of items, **list of couriers**, and totals. |
| `subtotal`       | int64  | Sum of all item prices before shipping. |
| `total_shipping` | int64  | Total shipping fee across all shops. |
| `payment_methods`| array  | Selection for method of payment. |
| `total`          | int64  | Grand total but only estimation (use `/cart/checkout/calculate` for actual cost). |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing or invalid UUIDs, `quantity <= 0`, or courier missing for any shop. |
| `401 Unauthorized` | Missing or invalid session. |
| `409 Conflict`     | System conflict includes **empty user default address**, **insufficient stock** and **shop has not registered the courier**. |

### Calculate Checkout

- **Method**: `POST`
- **Endpoint**: `/carts/checkout/calculate`
- **Description**: Calculate shipping costs and order totals for a given set of items and shops without placing an order. A courier must be provided for each shop.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "address_id": "string (UUID, required)",
    "payment_method_id": "string (UUID, required)",
    "shops": [
      {
        "shop_id": "string (UUID, required)",
        "items": [
          {
            "product_variant_type": "standard",
            "product_id": "string (UUID, required for standard)",
            "quantity": 1,
            "item_options": {
              "size": "medium",
              "jambul": "both"
            }
          },
          {
            "product_variant_type": "custom",
            "product_name": "Custom Board — Selamat & Sukses",
            "physical_size_id": "medium",
            "quantity": 1,
            "custom_design": {
              "metadata": { "version": "1.0.0" },
              "layout": { "physicalSizeId": "medium" }
            }
          }
        ],
        "courier": {
          "code": "jne",
          "service": "REG"
        }
      }
    ]
  }
  ```

#### Important Notes

> Unlike Estimate Checkout, `address_id`, `payment_method_id` and `courier` per shop are **required** in this endpoint.

#### Response `200 OK`

```json
{
    "address": {
        "id": "cd0ce5ff-d5ed-4b1c-85a4-16fef4d19f01",
        "recipient_name": "Dialyn UwU",
        "phone": "000",
        "full_address": "Blok LA 22B, Jl Sokekarno Saya akan lawan"
    },
    "shops": [
        {
            "shop_id": "8fad2c68-82a2-4578-a550-c625a1691d8a",
            "name": "Chia Medan Satria",
            "slug": "chia-medan-satria",
            "subtotal": 21600000,
            "total": 41040000,
            "selected_courier": {
                "code": "tiki",
                "service": "SDS",
                "fee": 19440000
            },
            "items": [
                {
                    "product_id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
                    "shop_id": "8fad2c68-82a2-4578-a550-c625a1691d8a",
                    "name": "Grand Opening",
                    "price": 150000,
                    "quantity": 144,
                    "subtotal": 21600000
                }
            ]
        },
        {
            "shop_id": "333f6432-a01c-412f-99f4-0f08ca0d8eb1",
            "name": "Chia Cipinang",
            "slug": "chia-cipinang",
            "subtotal": 825000,
            "total": 960000,
            "selected_courier": {
                "code": "tiki",
                "service": "REG",
                "fee": 135000
            },
            "items": [
                {
                    "product_id": "71be3ee1-17b4-4bb8-8f80-eae6ad93a844",
                    "shop_id": "333f6432-a01c-412f-99f4-0f08ca0d8eb1",
                    "name": "Graduate",
                    "price": 55000,
                    "quantity": 15,
                    "subtotal": 825000
                }
            ]
        }
    ],
    "subtotal": 22425000,
    "total_shipping": 19575000,
    "total": 42000000,
    "selected_payment_method": {
        "id": "0137d751-5188-447a-b630-1bf858f4f866",
        "name": "QRIS",
        "type": "qr_code",
        "description": "QRIS payment via Midtrans",
        "fee": 0,
        "subtotal": 42000000,
        "total": 42000000
    }
}
```

#### Checkout Response Fields

| Field            | Type   | Description |
|------------------|--------|-------------|
| `address`        | object | Delivery address details. |
| `shops`          | array  | Per-shop breakdown of items, **selected courier**, and totals. |
| `subtotal`       | int64  | Sum of all item prices before shipping. |
| `total_shipping` | int64  | Total shipping fee across all shops. |
| `selected_payment_method`| object  | Selected method of payment. |
| `total`          | int64  | Grand total including shipping and method of payment. |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing or invalid UUIDs, or `quantity <= 0`. |
| `401 Unauthorized` | Missing or invalid session. |
| `409 Conflict`     | System conflict includes **empty user default address**, **insufficient stock** and **shop has not registered the courier**. |

## Orders

### Find My Orders

- **Method**: `GET`
- **Endpoint**: `/users/me/orders`
- **Description**: Retrieve the user's orders.
- **Authentication**: Customer
- **Request Body**: None

#### Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `page`    | int    | No       | Page number. Defaults to `1`. |
| `limit`   | int    | No       | Number of results per page. Defaults to `10`. |
| `sort`    | string | No       | Comma-separated sort expressions. |
| `status`    | string | No       | Filter by order status. |

#### Sort Fields

| Field      | Example              | Description                       |
|------------|----------------------|-----------------------------------|
| `latest`   | `sort=latest:desc`   | Sort by creation date.            |
| `date`     | `sort=date:asc`      | Sort by order date.               |
| `number`   | `sort=number:asc`    | Sort by order number.             |
| `total`    | `sort=total:desc`    | Sort by total amount.             |
| `status`   | `sort=status:asc`    | Sort by order status.             |
| `modified` | `sort=modified:desc` | Sort by last modified date.       |

> Default sort: `latest:desc`. Multiple fields can be chained, e.g. `sort=latest:asc,number:desc`.

**Examples**:
- `GET /users/me/orders?page=1&limit=10`
- `GET /users/me/orders?status=pending`
- `GET /users/me/orders?sort=total:desc`

#### Response `200 OK`

```json
{
	"limit": 10,
	"orders": [
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
					"product_variant_type": "standard",
					"product_id": "480eec7c-d950-4927-a570-1fc3dc20df67",
					"product_name": "Prosperity Grand Opening Stand",
					"quantity": 3,
					"unit_price": 150000,
					"subtotal": 450000,
					"item_options": {
						"size": "small",
						"jambul": "none"
					},
					"shop_id": "7e5e335a-ec5b-4399-a8f6-1ea7dd8f0974",
					"shop_name": "dayum",
					"courier_code": "tiki",
					"courier_service": "SDS",
					"shipping_fee": 770000
				},
				{
					"id": "f84a9f49-b3b3-4c8e-8fee-dc5d4a0207ac",
					"product_variant_type": "custom",
					"product_id": null,
					"product_name": "Custom Board — Selamat & Sukses",
					"quantity": 1,
					"unit_price": 650000,
					"subtotal": 650000,
					"shop_id": "7e5e335a-ec5b-4399-a8f6-1ea7dd8f0974",
					"shop_name": "dayum",
					"courier_code": "tiki",
					"courier_service": "SDS",
					"shipping_fee": 0,
					"custom_design": {
						"physical_size_id": "medium",
						"preview_url": "https://chia-florist.supabase.co/storage/v1/object/public/custom-previews/designs/2026/08/order-board-a8f3b9e1.png"
					}
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
		},
		{
			"id": "01606ea8-5952-40c8-a4cf-ac206dbbf096",
			"number": "ORD-20260627-56D32D",
			"customer_id": "7466a260-6dd9-45db-8a70-3254bfc2dc98",
			"address_id": "f279f798-2de1-4ebd-a660-568b835b3a52",
			"status": "pending",
			"subtotal": 1100000,
			"shipping_fee": 1100000,
			"total": 2200000,
			"created_at": "2026-06-27T10:05:48.635958Z",
			"items": [
				{
					"id": "97e40886-5a36-47e5-8c1b-7d30e62b05e6",
					"product_id": "1aa696a8-ddf6-4718-9d30-bc510646ad70",
					"product_name": "The Scholar’s Cap Arrangement",
					"quantity": 20,
					"unit_price": 55000,
					"subtotal": 1100000,
					"shop_id": "7e5e335a-ec5b-4399-a8f6-1ea7dd8f0974",
					"shop_name": "dayum",
					"courier_code": "tiki",
					"courier_service": "SDS",
					"shipping_fee": 1100000
				}
			],
			"payment": {
				"id": "58826ccf-a45b-437e-8bc5-23bb37590954",
				"status": "pending",
				"provider": "midtrans",
				"amount": 2200000,
				"expires_at": "2026-06-28T10:05:52Z",
				"created_at": "2026-06-27T10:05:48.635958Z"
			}
		}
	],
	"page": 1,
	"total": 2
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | User profile not found.     |

### Get Order

- **Method**: `GET`
- **Endpoint**: `/users/me/orders/{id}`
- **Description**: Retrieve the order by its ID.
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

### Get Order Payment Details

- **Method**: `GET`
- **Endpoint**: `/users/me/orders/{orderID}/payment`
- **Description**: Retrieve detailed payment recovery information for the order. This includes the status, amount, expiration timestamp, active gateway channel data (such as QRIS string/deep links), manual account info, and the rendered step-by-step payment instructions.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

##### Example (Gateway Payment - QRIS)
```json
{
  "payment_id": "1d3a0355-ce51-4346-a07c-b8bc839e85f1",
  "status": "pending",
  "amount": 1220000,
  "expires_at": "2026-06-28T10:09:09Z",
  "channel_type": "qr_code",
  "display_name": "QRIS",
  "action_url": "data:image/png;base64,qrstring...",
  "instruction": "# QRIS Payment Instructions\n\nPlease scan the QR code before **2026-06-28T10:09:09Z**."
}
```

##### Example (Manual Payment - Bank Transfer)
```json
{
  "payment_id": "2d3a0355-ce51-4346-a07c-b8bc839e85f2",
  "status": "pending",
  "amount": 250000,
  "expires_at": "2026-06-28T10:09:09Z",
  "account_name": "Chia Florist",
  "account_number": "1234567890",
  "instruction": "# Bank Transfer Instructions\n\nPlease transfer **IDR 250,000** to BCA account **1234567890**."
}
```

#### Response Fields

| Field            | Type   | Description |
|------------------|--------|-------------|
| `payment_id`     | string | The unique ID of the payment transaction. |
| `status`         | string | Current payment status (`pending`, `paid`, `failed`, `expired`, `cancelled`, `refunded`). |
| `amount`         | int    | Total payment amount. |
| `expires_at`     | string | RFC3339 formatted expiration timestamp. |
| `channel_type`   | string | (Optional) Gateway payment channel method type. |
| `display_name`   | string | (Optional) Gateway payment channel display name. |
| `action_url`     | string | (Optional) Gateway payment channel QR code, deep link, or VA number. |
| `account_name`   | string | (Optional) Manual payment account holder name. |
| `account_number` | string | (Optional) Manual payment account number. |
| `phone_number`   | string | (Optional) Manual payment phone number. |
| `qr_string`      | string | (Optional) Manual payment QR string. |
| `instruction`    | string | (Optional) Rendered step-by-step markdown payment instructions. |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Order or payment not found, or order does not belong to the authenticated customer. |

### Check Order Payment Status

- **Method**: `POST`
- **Endpoint**: `/users/me/orders/{orderID}/payment/check`
- **Description**: Query the external payment provider (e.g. Midtrans) directly to check and sync the payment status of the order. This is useful when the payment status is stuck as `pending` after successful payment, bypassing the background reconciler.
- **Authentication**: Customer
- **Request Body**: None

#### Path Parameters

| Parameter | Type          | Description                   |
|-----------|---------------|-------------------------------|
| `orderID` | UUID (string) | The ID of the order to check. |

#### Response `200 OK`

```json
{
  "status": "paid",
  "synced": true
}
```

#### Response Fields

| Field    | Type   | Description                                                                                             |
|----------|--------|---------------------------------------------------------------------------------------------------------|
| `status` | string | Current payment status (`pending`, `paid`, `failed`, `expired`, `cancelled`, `refunded`).               |
| `synced` | bool   | Indicates if the payment status was resolved and updated/synchronized in the system during this call.  |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `orderID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Order or payment not found, or order does not belong to the authenticated customer. |

### Get Order Tracking Timeline

- **Method**: `GET`
- **Endpoint**: `/users/me/orders/{orderID}/tracking`
- **Description**: Retrieve the chronological tracking timeline of a shipment, merging internal shop events (e.g. `packed`, `picked_up`) with external courier updates (e.g. transit manifests from Komerce) when available.
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

#### Response Fields

| Field             | Type   | Description |
|-------------------|--------|-------------|
| `order_id`        | string | The unique ID of the order. |
| `shipment_id`     | string | The unique ID of the shipment. |
| `courier`         | string | The courier code (e.g. `jne`, `self_delivery`). |
| `tracking_number` | string | (Optional) The tracking number if shipped via courier. |
| `timeline`        | array  | Chronological list of tracking timeline events. |

#### Timeline Event Fields

| Field         | Type   | Description |
|---------------|--------|-------------|
| `status`      | string | Status/manifest code of the event. |
| `description` | string | Description of the tracking event. |
| `location`    | string | Location where the event occurred. |
| `timestamp`   | string | Date and time when the event occurred (ISO 8601). |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `orderID` is not a valid UUID. |
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | Order or shipment not found, or order does not belong to the authenticated customer. |

### Create Order

- **Method**: `POST`
- **Endpoint**: `/order`
- **Description**: Place a new order and create the associated invoice, order items, invoice items, and payment details. It locks the inventory for standard catalog products (bypassing stock reservation for custom flower board items) and returns the assigned payment account and instructions.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "address_id": "string (UUID, required)",
    "selected_payment": {
      "id": "string (UUID, required)",
      "is_manual": true
    },
    "shops": [
      {
        "shop_id": "string (UUID, required)",
        "name": "Chia Medan Satria",
        "selected_courier": {
          "code": "jne",
          "service": "REG"
        },
        "items": [
          {
            "product_variant_type": "standard",
            "product_id": "e7b0c950-6d33-4f51-b851-93c10a421234",
            "name": "Prosperity Grand Opening Stand",
            "quantity": 1,
            "item_options": {
              "size": "large",
              "jambul": "both"
            }
          },
          {
            "product_variant_type": "custom",
            "name": "Custom Board — Selamat & Sukses",
            "physical_size_id": "medium",
            "quantity": 1,
            "custom_design": {
              "metadata": {
                "version": "1.0.0",
                "checksum": "a8f3b9e1"
              },
              "layout": { "physicalSizeId": "medium", "upperHeightRatio": 0.58 },
              "sections": {
                "upper": { "header": { "text": "Selamat & Sukses", "fontId": "playfair" } },
                "lower": { "body": { "text": "PT. Tech Nusantara", "fontId": "inter" } }
              },
              "assets": {
                "previewUrl": "https://chia-florist.supabase.co/storage/v1/object/public/custom-previews/designs/2026/08/order-board-a8f3b9e1.png"
              }
            }
          }
        ]
      }
    ]
  }
  ```

#### Important Notes

> Only manual payments are currently supported. If `is_manual` is set to `false`, the request will return a `403 Forbidden` response.
>
> `name` in `shops` refers to the shop's name, and `name` in `items` refers to the product's name.
>
> `product_variant_type` defaults to `"standard"` if omitted. When set to `"custom"`, stock reservation is bypassed and the item's price is resolved server-side based on `physical_size_id`.

#### Response `200 OK`

```json
{
  "order_id": "f1e2d3c4-b5a6-7890-fedc-ba0987654321",
  "instruction": "# Payment Instructions\n\nPlease transfer **IDR 22,720,000** to the following account:\n- Bank: BCA\n- Account Name: Chia Florist\n- Account Number: 1234567890\n\nPlease complete the payment within 24 hours.",
  "payment_account": {
    "account_name": "Chia Florist",
    "account_number": "1234567890"
  }
}
```

#### Response Fields

| Field             | Type   | Description |
|-------------------|--------|-------------|
| `order_id`        | string | The unique ID of the created order. |
| `instruction`     | string | Markdown-formatted payment instructions. |
| `payment_account` | object | Details of the assigned payment account. |

#### Payment Account Fields

| Field            | Type   | Description |
|------------------|--------|-------------|
| `account_name`   | string | Name of the account holder. |
| `account_number` | string | (Optional) Bank account number. |
| `phone_number`   | string | (Optional) Phone number for e-wallet accounts. |
| `qr_string`      | string | (Optional) Raw QR code string for QRIS payments. |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing or invalid UUIDs, empty names, or `quantity <= 0`. |
| `401 Unauthorized` | Missing or invalid session. |
| `403 Forbidden`    | Requesting a non-manual payment method (gateway integration is not available yet). |
| `404 Not Found`    | Payment method or address not found. |
| `409 Conflict`     | System conflict including **no available payment account**, **insufficient stock**, or **courier/shipping issues**. |
