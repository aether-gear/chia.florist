# Product Image Feature Implementation Task

## Objective

Implement product image persistence and image resolution handling for product module.

Scope is **domain + repository + infra implementation only**.

**Do NOT create usecases, handlers, DTOs, or API routes.**

## Requirements

Product image variants required:

### Catalog

* single image
* medium/small optimized
* intended for product listing

### Detail

* multiple images
* higher resolution
* intended for product detail page/gallery

### Cart

* single image
* small optimized
* intended for cart display

## 1. Migration

Create migration file under:

`service-core/migrations`

Create table for product images.

Expected fields:

* id
* product_id (FK)
* catalog_url
* cart_url
* is_primary
* display_order
* created_at

For detail images:

Choose one of these approaches (pick the cleaner one based on existing product structure):

### Option A (preferred)

Separate table:

`product_detail_images`

Fields:

* id
* product_image_id
* image_url
* display_order

### Option B

Single table with image_type

## 2. Domain

Create product image domain models.

Expected:

### ProductImage

Contains:

* ID
* ProductID
* CatalogURL
* CartURL
* IsPrimary
* DisplayOrder
* DetailImages

### ProductDetailImage

Contains:

* ID
* ProductImageID
* URL
* DisplayOrder

## 3. Repository Interfaces

Create **2 repositories**

### A. Product Image Repository

Responsible for DB persistence.

Methods should support:

* create product image
* create detail images
* get by product id
* delete by product id

Implementation required.

### B. Product Image Upload Repository

Responsible for storage upload.

Methods should support:

* upload catalog image
* upload cart image
* upload detail image
* delete uploaded images

Implementation required.

Use existing storage/image service.

## 4. Existing Image Service Review

Inspect and modify if needed:

`service-core/internal/modules/product/infra/service/image_service.go`

Tasks:

* adapt to support catalog/cart/detail variants
* support image resizing
* support multiple detail image generation if needed
* keep backward compatibility where possible

## 5. Global Resolution Handling

Resolution logic must be reusable globally.

Find proper placement under:

`service-core/internal`

Create shared image processing/resolution package there.

Expected responsibility:

* resize image
* generate variants
* resolution presets

Suggested structure:

`internal/shared/image`

or

`internal/infra/image`

Choose whichever matches existing architecture better.

## Resolution Targets

Use reasonable defaults:

### Catalog

~400-600px

### Cart

~96-150px

### Detail

~1200px

Keep aspect ratio.

Optimize file size.

## Deliverables

Must finish:

* migration
* domain models
* repository interfaces
* repository implementations
* upload repository implementation
* global image resolution utility
* modified image_service.go

Must NOT add:

* usecases
* handlers
* controller/routes
* tests
* extra features

## Notes - Folder Structure

Below is the folder structure of the project. Hope this can give better visualization
```
├── cmd/
│   └── (migrate, seed, server)
│       └── main.go

├── internal/

│   ├── bootstrap/
│   │   ├── app.go
│   │   ├── config.go
│   │   ├── container.go
│   │   └── router.go

│   ├── common/
│   │   ├── errors/
│   │   │   ├── constructor.go
│   │   │   ├── error.go
│   │   │   └── resolver.go
│   │   ├── http/
│   │   │   ├── handler.go
│   │   │   ├── method.go
│   │   │   └── response.go
│   │   ├── logger/
│   │   │   ├── logger.go
│   │   │   ├── slog.go
│   │   │   └── zap.go
│   │   └── middleware/
│   │       ├── chain.go
│   │       ├── logging.go
│   │       ├── middleware.go
│   │       ├── recovery.go
│   │       └── response.go

│   ├── infra/
│   │   ├── cache/
│   │   ├── db/
│   │   │   ├── config.go
│   │   │   ├── connection.go
│   │   │   ├── migrate.go
│   │   │   └── seed.go
│   │   ├── storage/
│   │   ├── queue/
│   │   └── waf/

│   ├── modules/
│   │   └── (address, audit, auth, cart, courier, inventory
│   │       │       location, payment, product, security_policy, inventory
│   │       │       shipment, shop, threat_intel, user)/
│   │       ├── delivery/
│   │       │   └── http/
│   │       │       ├── dto.go
│   │       │       └── handler.go
│   │       ├── domain/
│   │       ├── infra/
│   │       │   ├── service/
│   │       │   └── persistence/
│   │       │       └── <interface>_repository_impl.go
│   │       ├── repository/
│   │       │   ├── query.go
│   │       │   └── interface.go
│   │       └── usecase/

│   ├── seed/
│   │   └── (courier, location)/
│   │       ├── source/
│   │       └── seeder.go

│   └── shared/
│       ├── config/
│       ├── conversion/
│       ├── loader/
│       └── slug/

├── migrations/

├── test/

├── .env.example
├── .gitignore
├── README.md
├── go.mod
└── go.sum
```