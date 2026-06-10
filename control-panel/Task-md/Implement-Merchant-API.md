# Task: Migrate Merchant Authentication & Management API from Service-Core to Merchant SPA

## Objective

Implement the existing Merchant Authentication and Merchant Management APIs from `service-core` into the Merchant SPA frontend.

This is a **migration task**, not a backend development task. The APIs already exist in `service-core`; your responsibility is to integrate and consume them within the Merchant SPA application.

## Branch Preparation

Before starting development:

1. Checkout your local `develop` branch.
2. Pull the latest changes from remote:

   ```bash
   git checkout develop
   git pull origin main
   ```
3. Create your working branch from the updated `develop` branch.

> Important:
> Development must be performed on your local `develop` branch (or a feature branch created from it) to ensure your local `main` branch remains clean and aligned with production-ready code.

## Backend Setup

Before implementation, read the `service-core` documentation and run the backend locally.

Expected outcomes:

* Update file `.env` based on example and make sure git is not associate it.
* Verify backend `service-core` can run in locally.
* Understand available Merchant APIs.
* Verify request/response payloads.
* Test API integration locally before submitting.

## Scope

### Merchant Authentication

#### Merchant Sign In

* Endpoint: `POST /auth/merchant/signin`
* Implement merchant login flow.
* Handle authentication cookie returned by backend.
* Handle success and error states.

Request:

```json
{
  "email": "string",
  "password": "string",
  "user_agent": "string",
  "ip_address": "string"
}
```

Response:

```json
{
  "message": "login success"
}
```

### Merchant Management (Admin)

#### Create Merchant

* Endpoint: `POST /merchants`
* Create merchant form and submission flow.

Request:

```json
{
  "name": "string",
  "description": "string",
  "logo_url": "string",
  "banner_url": "string"
}
```

Response:

```json
{
  "message": "merchant successfully created"
}
```

#### Add Merchant Account

* Endpoint: `POST /merchants/{merchantID}/accounts`
* Create merchant staff account creation flow.
* Ensure merchant ID is correctly supplied from the selected merchant context.

Request:

```json
{
  "email": "string",
  "name": "string",
  "username": "string",
  "password": "string",
  "phone": "string"
}
```

Response:

```json
{
  "message": "verify success"
}
```

## Acceptance Criteria

* `service-core` can run locally without any changes needed (only credentials).
* Merchant login is integrated with backend API.
* Authentication cookie is handled correctly.
* Merchant creation flow is functional.
* Merchant account creation flow is functional.
* API error handling is implemented.
* UI follows existing Merchant SPA patterns and architecture.
* Backend tested against local `service-core`.
* No backend development is required.
* No duplicate API implementation should be created.

## Pull Request Requirements

* Target Branch: `develop`
* Assignee: <Friend Name>
* Reviewer: Deuterrr

Include:

* Any API integration issues encountered during implementation.

## Non-Goals

The following are OUT OF SCOPE:

- Creating new backend endpoints
- Modifying `service-core` APIs
- Refactoring `service-core`
- Creating mock APIs
- Creating duplicate backend implementations
- Developing directly on main branch
- Merging directly to main branch

Consume existing APIs from service-core only.