## Product API Documentation

### Find Products

Retrieves a paginated list of products from the catalog. Supports filtering by name and ID.

#### Endpoint

`GET {link}/product`

#### Query Parameters



| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| `page` | `integer` | `1` | The page number to retrieve. Defaults to 1 if ≤ 0.|
| `limit` | `integer` | `10` | The number of items per page. Defaults to 10 if ≤ 0.|
| `name` | `string` | - | Filter products by name.|
| `id` | `string` | - | Filter by a specific product ID.|

---

#### Response Structure

`Success Response (200 OK)`

**Top-Level Fields:**

* `products`: An array of product objects.


* `page`: Current page number.


* `limit`: Items per page limit.


* `total`: Total count of products matching the criteria.

**Product Object Details:**

* `id`: Unique identifier of the product.


* `sku`: Stock Keeping Unit.


* `name`: Product name.


* `slug`: URL-friendly version of the name.


* `price`: Unit price.


* `total_stock`: Current total stock quantity.


* `image.thumbnail`: URL string to the product thumbnail image.


* `is_available`: Boolean indicating if the product is active and has available stock (Total - Reserved > 0).

#### Example Request

`GET {link}/product?page=1&limit=2&name=keyboard`

#### Example Response

```json
{
    "products": [
        {
            "id": "e0686de0-b1ce-4459-999c-ac1c69ada522",
            "sku": "EVT-WED-008",
            "name": "Eternal Union Centerpiece",
            "slug": "eternal-union-centerpiece",
            "is_available": false,
            "price": 120000,
            "stock": 0,
            "images": {}
        },
        {
            "id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
            "sku": "EVT-GOP-007",
            "name": "Prosperity Grand Opening Stand",
            "slug": "prosperity-grand-opening-stand",
            "is_available": true,
            "price": 150000,
            "stock": 837,
            "images": {
                "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/2ceea56c-352f-4a48-a262-f60e9ee85b1c/321e530c-b31e-49e8-8054-c6c85984d386.jpg"
            }
        }
    ],
    "page": 1,
    "limit": 2,
    "total": 1
}
```

#### Error Handling
*   **500 Internal Server Error**: Returned if there is a database failure or internal logic error during execution.[cite: 1]

```