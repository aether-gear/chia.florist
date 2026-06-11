# Merchant Staff and Admin API Documentation

This document outlines the merchant staff and admin-facing APIs based on the current routing configuration. Each endpoint includes its method, path, request details, and a placeholder for the response.

## Authentication

### Merchant Sign In
- **Method**: `POST`
- **Endpoint**: `/auth/merchant/signin`
- **Description**: Authenticate a merchant account using email and password.
- **Authentication**: None (public)
- **Request Body**:
  ```json
  {
    "email": "string (required)",
    "password": "string (required)",
    "user_agent": "string (optional)",
    "ip_address": "string (optional)"
  }
  ```
- **Response**:
  - **Cookie**: 

    | Key | Value |
    | --- | --- |
    | Set-Cookie | hotpot="value" |

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
    "account_id": "51e20db6-5bdb-4f6a-b2b7-8d40c0db857d",
    "account_type": "merchant",
    "is_authenticated": true,
    "roles": [
      {
        "code": "merchant_staff",
        "name": "Merchant Staff"
      }
    ],
    "permissions": [
      {
        "code": "merchant_staff"
      }
    ]
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

---

## Merchant Management (Admin Only)

> These endpoints require merchant authentication via the merchant access token cookie and the **admin** merchant role.

### Create Merchant
- **Method**: `POST`
- **Endpoint**: `/merchants`
- **Description**: Create a new merchant entity. The authenticated account is automatically associated as the merchant owner.
- **Authentication**: Merchant Admin
- **Request Body**:
  ```json
  {
    "name": "string (required)",
    "description": "string (optional)",
    "logo_url": "string (optional)",
    "banner_url": "string (optional)"
  }
  ```
- **Response**:
  ```json
    {
      "message": "merchant successfully created"
    }
  ```

---

### Add Merchant Account
- **Method**: `POST`
- **Endpoint**: `/merchants/{merchantID}/accounts`
- **Description**: Register and assign a new account to a merchant. The actor must be an admin of the target merchant.
- **Authentication**: Merchant Admin
- **Path Parameters**:
  | Parameter    | Type   | Description                        |
  |--------------|--------|------------------------------------|
  | `merchantID` | string (UUID) | The ID of the target merchant |
- **Request Body**:
  ```json
  {
    "email": "string (required)",
    "name": "string (required)",
    "username": "string (required)",
    "password": "string (required)",
    "phone": "string (optional)"
  }
  ```
- **Response**:
  ```json
  {
    "message": "verify success"
  }
  ```
