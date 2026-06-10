# Task: Build Merchant Profile Settings View (Backend Migration Ready)

## Objective

Implement the Merchant Profile Settings page in the Merchant SPA.

Backend support for this feature is not yet available in `service-core`.

The purpose of this task is to build the complete UI, state management, validation, and page structure using temporary hardcoded data while keeping the implementation ready for future API migration.

This is a frontend implementation task only.

## Branch Preparation

Before starting development:

```bash
git checkout develop
git pull origin main
```

Create your working branch from the updated `develop` branch.

> Keep local `main` clean and do not perform development directly on `main`.

## Important Implementation Requirements

### Current State

Profile Settings APIs do not exist yet in `service-core`.

### Expected Approach

Use a temporary data source such as:

```ts
const mockMerchantProfile = {
  ...
}
```

The page architecture should be designed so that the temporary data can later be replaced with API calls with minimal refactoring.

Example:

```ts
const loadMerchantProfile = async () => {
  // TODO: Replace with API call when service-core endpoint becomes available
  return mockMerchantProfile
}
```

### Do Not

* Create backend endpoints
* Modify service-core
* Create mock servers
* Add temporary local storage persistence
* Add fake API layers that do not reflect future integration needs

## Page Structure

The Profile Settings page should contain the following sections.

Tabs are recommended.

### 1. Identity

Fields:

* Account ID
* Merchant ID
* Merchant Name / Display Name
* Merchant Slug
* Profile Photo
* Cover Banner
* Description / Bio

### 2. Contact

Fields:

* Email
* Phone Number
* WhatsApp Number
* Customer Service Contact
* Address
* Country
* Province / State
* City
* District
* Postal Code
* Full Address
* Latitude
* Longitude

### 3. Settings

Fields:

* Preferred Language
* Preferred Currency
* Timezone

### 4. Operational

Fields:

* Opening Hours
* Closing Hours
* Business Days
* Delivery Radius
* Pickup Available

### 5. Financial

Fields:

* Bank Account Name
* Bank Name
* Bank Account Number
* E-Wallet Information
* Tax Number (Optional)

## Architecture Requirements

Design the page as if the backend already exists.

Suggested separation:

```text
pages/
└── merchant-profile/

components/
├── ProfileIdentitySection
├── ProfileContactSection
├── ProfileSettingsSection
├── ProfileOperationalSection
└── ProfileFinancialSection

viewmodels/
└── useMerchantProfileViewModel

models/
└── MerchantProfile
```

The goal is to minimize future migration effort when service-core APIs become available.

## Acceptance Criteria

* Profile Settings page is implemented.
* All required sections are present.
* Tabs navigation is functional.
* Temporary hardcoded data is displayed.
* Form state management is implemented.
* Validation is implemented where applicable.
* Architecture is prepared for future API integration.
* No backend development performed.
* No service-core modifications performed.

## Pull Request Requirements

Target Branch: `develop`

Assignee: you

Reviewer: Deuterrr

Include:

* Notes regarding future API integration points.

## Non-Goals

The following are OUT OF SCOPE:

- Creating duplicate backend implementations
- Developing directly on main branch
- Merging directly to main branch