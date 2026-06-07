## Refactor: Move External Service Initialization to Bootstrap

### Background

Currently, initialization of external services is handled directly inside the container layer.

This includes provider construction and infrastructure readiness setup such as:

* Database connection initialization
* Storage provider initialization
* Future external services (cache, queue, third-party providers)

This makes the container responsible for both dependency wiring and infrastructure lifecycle orchestration.

### Problem

The current responsibility boundary is too broad.

Container should only focus on dependency composition.

Infrastructure initialization belongs to a dedicated bootstrap layer responsible for startup orchestration.

Current flow:

```text
container -> initialize external services -> wire dependencies
```

Target flow:

```text
bootstrap -> initialize external services
container -> wire initialized dependencies
runtime -> execute
```

### Goal

Move all external service initialization into `internal/bootstrap/`.

Create a bootstrap file with an appropriate name (recommended examples):

* `app.go`
* `bootstrap.go`
* `infra.go`

Choose whichever best fits the existing project structure.

### Required Changes

#### 1. Create bootstrap orchestration layer

Under:

```text
internal/bootstrap/
```

Responsible for:

* Initializing database connections
* Initializing storage providers
* Running infra readiness checks
* Returning initialized app dependencies

#### 2. Remove external initialization from container

Container should no longer:

* Open DB connections
* Initialize provider clients
* Perform provider setup checks

It should only receive already-initialized dependencies.

#### 3. Update constructor boundaries

External service constructors should:

* Return errors
* Never terminate process internally (`log.Fatal`, panic, etc.)

Bootstrap decides failure handling.

#### 4. Update command entrypoints

Refactor all entrypoints to use bootstrap:

* `cmd/api`
* `cmd/migrate`
* `cmd/seed`

Pattern:

```go
app, err := bootstrap.New(cfg)
if err != nil {
    return err
}
defer app.Close()
```

### Acceptance Criteria

* All external service initialization moved out of container
* New bootstrap orchestration layer exists under `internal/bootstrap/`
* Container only composes dependencies
* Constructors return errors instead of killing process
* App startup flow is centralized

### Notes

Current architectural direction:

```text
app config -> what is needed
provider config -> how to connect
bootstrap -> implementation selection + initialization
container -> dependency wiring
runtime -> execution
```

This refactor aligns the project with cleaner separation of concerns and makes future provider swaps easier.
