## Implementation

### Objective

Implement all concrete storage providers and service wiring based on Part 1 contracts.

### Files to Create

* infra/storage/local.go
* infra/storage/supabase.go
* modules/product/infra/service/image_service.go

### Files to Update

* bootstrap/config.go
* bootstrap/container.go

### Implementation Tasks

#### Local Provider

Implement:

* file persistence
* path normalization
* secure deletion

#### Supabase Provider

Implement:

* authenticated upload
* delete object
* public/signed URL generation
* bucket validation

Must consume config only through injected config object.

Never hardcode:

* project URL
* API key
* bucket

* (Keys are on /internal/shared/config/supabase_storage_config.go) 

#### Provider Selection

Container must support:

```go
switch cfg.Storage.Provider {
case "local":
case "supabase":
}
```

#### Product Image Service

Responsibilities:

* deterministic object key generation
* metadata validation
* MIME validation
* provider delegation

Path:

```text
products/{product_id}/{uuid}.jpg
```

### Validation Rules

Accept:

* jpg
* png
* webp

Reject:

* unsupported MIME
* oversize payloads
* malformed filenames

### Success Criteria

Implementation succeeds if:

* Local upload works
* Supabase upload works
* URL generation works
* Delete works
* Swapping providers requires container change only
