# Customer API Documentation

This document covers the customer-facing APIs for the Chia Florist service.
Endpoints are organized by access level: **Public** and **Authenticated Customer**.

## TODO

- [ ] Public API
  - [ ] Authentication
    - [ ] Sign Up
    - [ ] Verify Account
    - [ ] Sign In
  - [ ] Products
    - [ ] Find Products
    - [ ] Get Product Detail
  - [ ] Shops
    - [ ] Find Shops
    - [ ] Get Shop
    - [ ] Get Shop Addresses
    - [ ] Get Shop Couriers
    - [ ] Get Shop Products
  - [ ] Locations
    - [ ] List Provinces
    - [ ] List Cities by Province
    - [ ] List Districts by City
    - [ ] List Villages by District
- [ ] Authenticated Customer API
  - [ ] Authentication
    - [ ] Me
    - [ ] Log Out
  - [ ] Profile
    - [ ] Get Current User
  - [ ] Addresses
    - [ ] List My Addresses
    - [ ] Save My Address
    - [ ] Delete My Address
  - [ ] Cart
    - [ ] Get Cart
    - [ ] Add Item
    - [ ] Update Item
    - [ ] Remove Item
  - [ ] Checkout
    - [ ] Estimate Checkout
    - [ ] Checkout

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

- **Set-Cookie**: `<access_token_cookie>=<value>`
- **Body**:
  ```json
  { "message": "login success" }
  ```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `email` or `password` is missing or empty. |
| `401 Unauthorized` | Invalid email or password. |

## Products

### Find Products

- **Method**: `GET`
- **Endpoint**: `/products`
- **Description**: Retrieve a paginated list of products with optional filtering and sorting.
- **Authentication**: None
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
- `GET /products?name=anniversary`
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
- **Endpoint**: `/users/me`
- **Description**: Retrieve the profile of the currently authenticated customer.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "me": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Jane Doe",
    "username": "janedoe",
    "phone": "+6281234567890",
    "last_login_at": "2026-06-17T10:00:00Z"
  }
}
```

#### Error Responses

| Status             | Condition                   |
|--------------------|-----------------------------|
| `401 Unauthorized` | Missing or invalid session. |
| `404 Not Found`    | User profile not found.     |

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
- **Description**: Retrieve the current cart of the authenticated customer, including all items with subtotals and a grand total.
- **Authentication**: Customer
- **Request Body**: None

#### Response `200 OK`

```json
{
  "cart_id": "f1e2d3c4-b5a6-7890-fedc-ba0987654321",
  "total": 170000,
  "items": [
    {
      "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      "shop_id":    "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "name":       "Anniversary",
      "price":      85000,
      "quantity":   2,
      "subtotal":   170000,
      "images": {
        "thumbnail": "https://example.com/thumbnail.jpg"
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
- **Description**: Add a product to the cart. If the item already exists, quantity is incremented.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "product_id": "string (UUID, required)",
    "shop_id":    "string (UUID, required)",
    "quantity":   1
  }
  ```

#### Response `200 OK`

```json
{ "message": "item added" }
```

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | `product_id` or `shop_id` is not a valid UUID, or `quantity` is `<= 0`. |
| `401 Unauthorized` | Missing or invalid session. |

### Update Item

- **Method**: `PUT`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Update the quantity of a specific item in the cart. Setting quantity to `0` is treated as removal depending on the usecase implementation.
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

### Remove Item

- **Method**: `DELETE`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Remove a specific item from the cart entirely.
- **Authentication**: Customer
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

## Checkout

### Estimate Checkout

- **Method**: `POST`
- **Endpoint**: `/carts/checkout/calculate`
- **Description**: Calculate shipping costs and order totals for a given set of items and shops without placing an order. A courier must be provided for each shop.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "address_id": "string (UUID, optional)",
    "shops": [
      {
        "shop_id": "string (UUID, required)",
        "items": [
          {
            "product_id": "string (UUID, required)",
            "quantity":   1
          }
        ],
        "courier": {
          "code":    "jne",
          "service": "REG"
        }
      }
    ]
  }
  ```

#### Important Notes

> A `courier` object is **required** for every shop entry in this endpoint.
> Use `GET /shops/{shopID}/couriers` to retrieve available courier codes for each shop.

#### Response `200 OK`

```json
{
  "address": {
    "id":             "d4e5f6a7-b8c9-0123-defa-234567890123",
    "recipient_name": "Jane Doe",
    "phone":          "+6281234567890",
    "full_address":   "Jl. Bunga Indah No. 10, Bekasi"
  },
  "shops": [
    {
      "shop_id":  "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "subtotal": 85000,
      "total":    93000,
      "selected_courier": {
        "code":    "jne",
        "service": "REG",
        "fee":     8000
      },
      "items": [
        {
          "product_id": "9886edf6-087b-48e7-b00a-d79dd092e8d4",
          "shop_id":    "c3d4e5f6-a7b8-9012-cdef-123456789012",
          "name":       "Anniversary",
          "price":      85000,
          "quantity":   1,
          "subtotal":   85000
        }
      ],
      "cost_couriers": [
        {
          "code":    "jne",
          "name":    "Jalur Nugraha Ekakurir",
          "service": "REG",
          "etd":     "2-3",
          "fee":     8000
        }
      ]
    }
  ],
  "subtotal":       85000,
  "total_shipping": 8000,
  "total":          93000
}
```

#### Checkout Response Fields

| Field            | Type   | Description |
|------------------|--------|-------------|
| `address`        | object | Delivery address details. |
| `shops`          | array  | Per-shop breakdown of items, couriers, and totals. |
| `subtotal`       | int64  | Sum of all item prices before shipping. |
| `total_shipping` | int64  | Total shipping fee across all shops. |
| `total`          | int64  | Grand total including shipping (only present when courier is selected). |

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing or invalid UUIDs, `quantity <= 0`, or courier missing for any shop. |
| `401 Unauthorized` | Missing or invalid session. |

### Checkout

- **Method**: `POST`
- **Endpoint**: `/carts/checkout`
- **Description**: Place an order for the selected items. Courier selection per shop is optional at this stage.
- **Authentication**: Customer
- **Request Body**:
  ```json
  {
    "address_id": "string (UUID, optional)",
    "shops": [
      {
        "shop_id": "string (UUID, required)",
        "items": [
          {
            "product_id": "string (UUID, required)",
            "quantity":   1
          }
        ],
        "courier": {
          "code":    "jne",
          "service": "REG"
        }
      }
    ]
  }
  ```

#### Important Notes

> Unlike Estimate Checkout, `courier` is **optional** per shop in this endpoint.
> The response shape is identical to Estimate Checkout.

#### Response `200 OK`

Same structure as the [Estimate Checkout response](#estimate-checkout).

#### Error Responses

| Status             | Condition |
|--------------------|-----------|
| `400 Bad Request`  | Missing or invalid UUIDs, or `quantity <= 0`. |
| `401 Unauthorized` | Missing or invalid session. |
