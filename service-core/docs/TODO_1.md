## Domain, Types, Interfaces, and Contracts

### Objective

Define all storage abstractions, contracts, configuration types, and dependency boundaries before implementation.

### Files to Create

* infra/storage/provider.go
* infra/storage/types.go
* internal/config/supabase_storage_config.go
* modules/product/domain/image.go
* modules/product/repository/image_storage.go

### Responsibilities

#### Storage Provider Interface

Define provider abstraction:

```go
type Provider interface {
    Upload(key string, file io.Reader) (string, error)
    Delete(key string) error
    Exists(key string) (bool, error)
}
```

Requirements:

* Provider agnostic
* Context aware
* No implementation logic

Notes:

* Notice the interface doesn't explicitely saying about `context`.

#### Config Contracts

Define storage config and Supabase config.

Sensitive values must live only inside:

`internal/config/supabase_storage_config.go`

Includes:

* Project URL
* Service role key
* Bucket name
* Public base URL
* Signed URL expiry

#### Product Image Domain Contracts

Define:

* image metadata structure
* path conventions
* MIME constraints
* max upload size contract

### Success Criteria

Implementation phase should be possible without redefining contracts.
