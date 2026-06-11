# Public and Customer API Documentation

This document outlines the public (core) and customer-facing APIs based on the current routing configuration. Each endpoint includes its method, path, request details, and a placeholder for the response.

## Public API (Core)

These endpoints are accessible to anyone (no authentication required) or require core/public access logic.

### Authentication

#### Sign Up
- **Method**: `POST`
- **Endpoint**: `/auth/signup`
- **Description**: Register a new user account.
- **Request Body**:
  ```json
  {
    "name": "string (required)",
    "username": "string (required)",
    "email": "string (required)",
    "password": "string (required)",
    "phone": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "verification code sent",
    "challenge_id": "0379abef-4745-40b6-bf7b-04cc750e8f25"
  }
  ```

#### Verify Account
- **Method**: `POST`
- **Endpoint**: `/auth/verify`
- **Description**: Verify a newly created account using an OTP challenge.
- **Request Body**:
  ```json
  {
    "user_agent": "string",
    "ip_address": "string",
    "challenge_id": "string (required)",
    "otp": "integer (required)"
  }
  ```
- **Response**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Set-Cookie | chast="value" |

  - **Body**
    ```json
    {
        "message": "verify success"
    }
    ```

#### Sign In
- **Method**: `POST`
- **Endpoint**: `/auth/signin`
- **Description**: Authenticate a user via email.
- **Request Body**:
  ```json
  {
    "email": "string (required)",
    "password": "string (required)",
    "user_agent": "string",
    "ip_address": "string"
  }
  ```
- **Response**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Set-Cookie | chast="value" |

  - **Body**
    ```json
    {
      "message": "login success"
    }
    ```

#### Me
- **Method**: `GET`
- **Endpoint**: `/auth/me`
- **Description**: Get the information of current user.
- **Response**:
  ```json
  {
    "account_id": "8b477309-1492-48b0-b9ac-f493aea10cc6",
    "account_type": "customer",
    "is_authenticated": true
  }
  ```

#### Log Out
- **Method**: `POST`
- **Endpoint**: `/auth/logout`
- **Description**: For log out. It will clear session and refresh token of the user.
- **Response**:
  ```json
  {
    "message": "logout success"
  }
  ```

### Products

#### Find Products
- **Method**: `GET`
- **Endpoint**: `/products/`
- **Description**: Retrieve a list of products.
- **Request Query**:
  - `page`: index page
  - `limit`: how many product in each page
- **Request Body**: None (Uses Query Parameters)
- **Response**:
  ```json
  {
    "limit": 10,
    "page": 1,
    "total": 2,
    "products": [
      {
        "id": "ac788034-bd4c-476a-97d1-9c1def714042",
        "sku": "EVT-BNR-002",
        "name": "Grand Celebration Vinyl Banner",
        "slug": "grand-celebration-vinyl-banner",
        "is_available": false,
        "price": 45000,
        "stock": 0,
        "images": {
          "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/ac788034-bd4c-476a-97d1-9c1def714042/ee820e14-4bd9-46a3-8482-24ec3151a4ac.jpg"
        }
      },
      {
          "id": "b40dcc46-8328-4fcd-af77-42ecc9511606",
          "sku": "EVT-BDY-003",
          "name": "Confetti & Blooms Basket",
          "slug": "confetti-blooms-basket",
          "is_available": true,
          "price": 60000,
          "stock": 561,
          "images": {
              "thumbnail": ""
        }
      }
    ]
  }
  ```

#### Get Product Detail
- **Method**: `GET`
- **Endpoint**: `/products/{id}`
- **Description**: Retrieve details of a specific product.
- **Request Parameters**:
  - `id`: Product UUID in path
- **Request Body**: None
- **Response**:
  ```json
  {
    "id": "71be3ee1-17b4-4bb8-8f80-eae6ad93a844",
    "sku": "EVT-GRAD-006",
    "name": "The Scholar’s Cap Arrangement",
    "slug": "the-scholar-s-cap-arrangement",
    "description": "A modern, bold mix of blue irises and yellow tulips, complete with a commemorative graduation pick.",
    "is_available": true,
    "price": 55000,
    "weight": 1000,
    "stock": 72,
    "updated_at": null,
    "images": [
      {
        "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/4c39bfb4-330b-47a5-8ae2-2066c841c43c.jpg",
        "preview": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/4c39bfb4-330b-47a5-8ae2-2066c841c43c.jpg",
        "detail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/4c39bfb4-330b-47a5-8ae2-2066c841c43c.jpg"
      },
      {
        "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/ecd9dd18-6008-4af8-91d6-c13029fb15bf.png",
        "preview": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/ecd9dd18-6008-4af8-91d6-c13029fb15bf.png",
        "detail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/ecd9dd18-6008-4af8-91d6-c13029fb15bf.png"
      }
    ]
  }
  ```

### Locations

#### List Provinces
- **Method**: `GET`
- **Endpoint**: `/provinces/`
- **Description**: Retrieve a list of all provinces.
- **Request Body**: None
- **Response**:
  ```json
  {
    "provinces": [
      {
        "id": "1",
        "name": "NUSA TENGGARA BARAT (NTB)"
      },
      {
        "id": "2",
        "name": "MALUKU"
      },
      {
        "id": "3",
        "name": "KALIMANTAN SELATAN"
      },
      {
        "id": "4",
        "name": "KALIMANTAN TENGAH"
      },
      {
        "id": "5",
        "name": "JAWA BARAT"
      }
    ]
  }
  ```

#### List Cities in Province
- **Method**: `GET`
- **Endpoint**: `/provinces/{id}/cities`
- **Description**: Retrieve a list of cities for a specific province.
- **Request Parameters**:
  - `id`: Province ID in path
- **Request Body**: None
- **Response**:
  ```json
  {
    "cities": [
      {
        "id": "135",
        "province_id": "",
        "name": "JAKARTA BARAT"
      },
      {
        "id": "136",
        "province_id": "",
        "name": "JAKARTA SELATAN"
      },
      {
        "id": "137",
        "province_id": "",
        "name": "JAKARTA PUSAT"
      },
      {
        "id": "138",
        "province_id": "",
        "name": "JAKARTA UTARA"
      },
      {
        "id": "139",
        "province_id": "",
        "name": "JAKARTA TIMUR"
      },
      {
        "id": "141",
        "province_id": "",
        "name": "KEPULAUAN SERIBU"
      }
    ]
  }
  ```

#### List Districts in City
- **Method**: `GET`
- **Endpoint**: `/cities/{id}/districts`
- **Description**: Retrieve a list of districts for a specific city.
- **Request Parameters**:
  - `id`: City ID in path
- **Request Body**: None
- **Response**:
  ```json
  {
    "districts": [
      {
        "id": "1354",
        "city_id": "",
        "name": "CAKUNG"
      },
      {
        "id": "1355",
        "city_id": "",
        "name": "CIPAYUNG"
      },
      {
        "id": "1356",
        "city_id": "",
        "name": "CIRACAS"
      },
      {
        "id": "1357",
        "city_id": "",
        "name": "DUREN SAWIT"
      },
      {
        "id": "1358",
        "city_id": "",
        "name": "JATINEGARA"
      },
      {
        "id": "1359",
        "city_id": "",
        "name": "KRAMAT JATI"
      },
      {
        "id": "1360",
        "city_id": "",
        "name": "MAKASAR"
      },
      {
        "id": "1361",
        "city_id": "",
        "name": "MATRAMAN"
      },
      {
        "id": "1362",
        "city_id": "",
        "name": "PASAR REBO"
      },
      {
        "id": "1363",
        "city_id": "",
        "name": "PULO GADUNG"
      }
    ]
  }
  ```

#### List Villages in District
- **Method**: `GET`
- **Endpoint**: `/districts/{id}/villages`
- **Description**: Retrieve a list of villages for a specific district.
- **Request Parameters**:
  - `id`: District ID in path
- **Request Body**: None
- **Response**:
  ```json
  {
    "villages": [
      {
        "id": "17700",
        "district_id": "",
        "name": "BALI MESTER"
      },
      {
        "id": "17701",
        "district_id": "",
        "name": "BIDARACINA"
      },
      {
        "id": "17702",
        "district_id": "",
        "name": "CIPINANG BESAR SELATAN"
      },
      {
        "id": "17703",
        "district_id": "",
        "name": "CIPINANG BESAR UTARA"
      },
      {
        "id": "17704",
        "district_id": "",
        "name": "CIPINANG CEMPEDAK"
      },
      {
        "id": "17705",
        "district_id": "",
        "name": "CIPINANG MUARA"
      },
      {
        "id": "17706",
        "district_id": "",
        "name": "KAMPUNG MELAYU"
      },
      {
        "id": "17707",
        "district_id": "",
        "name": "RAWA BUNGA"
      }
    ]
  }
  ```

### Payments

#### List Payment Methods
- **Method**: `GET`
- **Endpoint**: `/payments/methods/`
- **Description**: Retrieve available payment methods.
- **Request Body**: None
- **Response**:
  ```json
  // On Progress
  ```

---

## Customer API

These endpoints require authentication and are restricted to users with the `Customer` account type.

### Users

#### Get Current User
- **Method**: `GET`
- **Endpoint**: `/users/me/`
- **Description**: Retrieve profile of the currently logged-in customer.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |

- **Request Body**: None
- **Response**:
  ```json
  {
    "me": {
      "id": "73fe7af7-1ad9-4a7c-aa73-091272af2856",
      "name": "Deu",
      "username": "deu3321",
      "phone": "021",
      "last_login_at": null
    }
  }
  ```

#### List User Addresses
- **Method**: `GET`
- **Endpoint**: `/users/me/addresses/`
- **Description**: Retrieve addresses for the current user.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**: None
- **Response**:
  ```json
  {
    "addresses": [
      {
        "user_id": "28be4cae-6ad9-4ca4-a5b4-7c1d924139c0",
        "receiver_name": "Iam",
        "phone": "000",
        "is_default": false,
        "province_id": "5",
        "city_id": "63",
        "district_id": "611",
        "village_id": "6555",
        "full_address": "Blok LA 22B, Jl Sokekarno Saya akan lawan",
        "postal_code": "17131",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": null
      },
      {
        "user_id": "28be4cae-6ad9-4ca4-a5b4-7c1d924139c0",
        "receiver_name": "Iam",
        "phone": "000",
        "is_default": true,
        "province_id": "5",
        "city_id": "63",
        "district_id": "611",
        "village_id": "6555",
        "full_address": "Blok LA 22B, Jl Sokekarno Saya akan lawan",
        "postal_code": "17131",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": null
      }
    ]
  }
  ```

#### Create User Address
- **Method**: `POST`
- **Endpoint**: `/users/me/addresses/`
- **Description**: Add a new address for the current user.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**:
  ```json
  {
    "receiver_name": "string (required)",
    "phone": "string",
    "is_default": "string",
    "province_id": "string (required)",
    "city_id": "string (required)",
    "district_id": "string (required)",
    "village_id": "string (required)",
    "full_address": "string (required)",
    "postal_code": "string (required)"
  }
  ```
- **Response**:
  ```json
  {
    "message": "address successfully created"
  }
  ```

### Carts

#### Get Cart
- **Method**: `GET`
- **Endpoint**: `/carts/`
- **Description**: Retrieve the current user's shopping cart.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**: None
- **Response**:
  ```json
  {
    "cart_id": "c49227fa-4e52-4545-a179-e94c83eb7edf",
    "items": [
      {
        "product_id": "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
        "shop_id": "8fad2c68-82a2-4578-a550-c625a1691d8a",
        "name": "Prosperity Grand Opening Stand",
        "price": 150000,
        "subtotal": 21600000,
        "quantity": 144,
        "images": {
            "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/2ceea56c-352f-4a48-a262-f60e9ee85b1c/321e530c-b31e-49e8-8054-c6c85984d386.jpg"
        }
      },
      {
        "product_id": "71be3ee1-17b4-4bb8-8f80-eae6ad93a844",
        "shop_id": "333f6432-a01c-412f-99f4-0f08ca0d8eb1",
        "name": "The Scholar’s Cap Arrangement",
        "price": 55000,
        "subtotal": 825000,
        "quantity": 15,
        "images": {
            "thumbnail": "https://mqolpawlannysqjokzoq.supabase.co/storage/v1/object/public/public-assets/products/71be3ee1-17b4-4bb8-8f80-eae6ad93a844/4c39bfb4-330b-47a5-8ae2-2066c841c43c.jpg"
        }
      }
    ],
    "total": 22425000
  }
  ```

#### Add Item to Cart
- **Method**: `POST`
- **Endpoint**: `/carts/items/`
- **Description**: Add a product to the cart.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**:
  ```json
  {
    "product_id": "string (required)",
    "shop_id": "string (required)",
    "quantity": "integer (required)"
  }
  ```
- **Response**:
  ```json
  {
    "message": "item added"
  }
  ```

#### Update Cart Item
- **Method**: `PUT`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Update the quantity of a cart item.
- **Request Parameters**:
  - `shopID`: Shop ID in path
  - `productID`: Product ID in path
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**:
  ```json
  {
    "quantity": "integer (required)"
  }
  ```
- **Response**:
  ```json
  {
    "message": "item updated"
  }
  ```

#### Remove Cart Item
- **Method**: `DELETE`
- **Endpoint**: `/carts/items/{shopID}/{productID}`
- **Description**: Remove an item from the cart.
- **Request Parameters**:
  - `shopID`: Shop ID in path
  - `productID`: Product ID in path
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**: None
- **Response**:
  ```json
  {
    "message": "item removed"
  }
  ```

### Payments

#### List Payment Accounts
- **Method**: `GET`
- **Endpoint**: `/payments/accounts/`
- **Description**: Retrieve the user's saved payment accounts.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |
    
- **Request Body**: None
- **Response**:
  ```json
  // On Progress
  ```

### Shipping

#### Estimate Shipping Options
- **Method**: `POST`
- **Endpoint**: `/shipping/cost`
- **Description**: Calculate estimated shipping costs based on weight and destination.
- **Request Header**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Cookie | chast="value" |

- **Request Body**:
  ```json
  {
    "origin": "integer (required)",
    "destination": "integer (required)",
    "weight": "integer (required)",
    "couriers": ["string"] (required),
    "price_filter": "string"
  }
  ```
  **Example**:
  ```json
  {
    "origin": 1358,
    "destination": 1119,
    "weight": 10000,
    "couriers": ["tiki", "ncs", "sentral", "wahana", "wahana"],
    "price_filter": "highest"
  }
  ```
  **Courier Availability**
    | Courier            | Code      | Checking Domestic Cost | Checking International Cost | Checking AWB |
    | ------------------ | --------- | ---------------------- | --------------------------- | ------------ |
    | JNE                | `jne`     | ✅                      | ✅                           | ✅            |
    | SiCepat            | `sicepat` | ✅                      | ❌                           | ✅            |
    | IDExpress          | `ide`     | ✅                      | ❌                           | ❌            |
    | SAP Express        | `sap`     | ✅                      | ❌                           | ✅            |
    | Ninja              | `ninja`   | ✅                      | ❌                           | ✅            |
    | J&T Express        | `jnt`     | ✅                      | ❌                           | ✅            |
    | TIKI               | `tiki`    | ✅                      | ✅                           | ✅            |
    | Wahana Express     | `wahana`  | ✅                      | ❌                           | ✅            |
    | POS Indonesia      | `pos`     | ✅                      | ✅                           | ✅            |
    | Sentral Cargo      | `sentral` | ✅                      | ❌                           | ❌            |
    | Lion Parcel        | `lion`    | ✅                      | ❌                           | ✅            |
    | Royal Express Asia | `rex`     | ✅                      | ❌                           | ❌            |



- **Response**:
  ```json
  {
    "couriers": [
      {
        "name": "Nusantara Card Semesta",
        "code": "ncs",
        "service": "NRS",
        "description": "Regular Service",
        "cost": 805000,
        "etd": "5-6 day"
      },
      {
        "name": "Citra Van Titipan Kilat (TIKI)",
        "code": "tiki",
        "service": "REG",
        "description": "Reguler Service",
        "cost": 680000,
        "etd": "5 day"
      },
      {
        "name": "Wahana",
        "code": "wahana",
        "service": "Kargo",
        "description": "Layanan Pengiriman Dengan Minimal Berat 10 Kg",
        "cost": 600000,
        "etd": "11 day"
      },
      {
        "name": "Sentral Cargo",
        "code": "sentral",
        "service": "DARAT NON ELEKTRONIK",
        "description": "Darat Non Elektronik",
        "cost": 200000,
        "etd": "10-12 day"
      }
    ]
  }
  ```
