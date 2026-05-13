```                                                                
 ▄▄▄▄ ▄▄▄▄▄ ▄▄▄▄  ▄▄ ▄▄ ▄▄  ▄▄▄▄ ▄▄▄▄▄      ▄▄▄▄  ▄▄▄  ▄▄▄▄  ▄▄▄▄▄ 
███▄▄ ██▄▄  ██▄█▄ ██▄██ ██ ██▀▀▀ ██▄▄  ▄▄▄ ██▀▀▀ ██▀██ ██▄█▄ ██▄▄  
▄▄██▀ ██▄▄▄ ██ ██  ▀█▀  ██ ▀████ ██▄▄▄     ▀████ ▀███▀ ██ ██ ██▄▄▄ 

A backend service designed to support the e-commerce operations of 
    Chia Florist.
Built with Go, this project emphasizes simplicity, clarity, and 
    serves as a practical portfolio for developers.
```

### Objective
- Build a modular and scalable backend service to support core e-commerce operations such as product management, transactions, and order handling
- Demonstrate clean architecture principles with a feature-based structure for maintainability and clarity
- Provide a practical learning and reference project for developers exploring backend development with Go
- Serve as a portfolio project showcasing real-world backend design and implementation

### How to Run
#### 1. Prequisites

Make sure the system for development already installed:

- Go (Version 1.20 or higher recommended)

Check Go installation with

```bash
go version
```

#### 2. Install Dependencies
```bash
go mod tidy
```

#### 3. Run The Application

```bash
go run cmd/server/main.go 
```

    Verify server by seeing active log from the application.

### Project Structure

```
├── cmd/
│   └── (migrate, seed, server)
│       └── main.go

├── internal/

│   ├── bootstrap/
│   │   ├── app.go
│   │   ├── config.go
│   │   ├── container.go
│   │   ├── infra.go
│   │   └── router.go

│   ├── common/
│   │   ├── errors/
│   │   │   ├── constructor.go
│   │   │   ├── error.go
│   │   │   └── resolver.go
│   │   ├── http/
│   │   │   ├── multipart/
│   │   │   │   ├── parser.go
│   │   │   │   ├── types.go
│   │   │   │   └── validator.go
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
│   │   ├── db/
│   │   │   ├── connection.go
│   │   │   ├── migrate.go
│   │   │   └── seed.go
│   │   ├── middleware/
│   │   └── storage/
│   │       ├── supabase/
│   │       │   ├── bucket.go
│   │       │   ├── object.go
│   │       │   ├── provider.go
│   │       │   └── resolver.go
│   │       ├── migrate.go
│   │       └── provider.go

│   ├── modules/
│   │   └── (address, audit, auth, cart, courier, inventory
│   │       │       location, payment, product, security_policy,
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
│       ├── image/
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

### Dependency Flow

* Dependencies always point **inward**
* Domain is the **center and most stable layer**

```
delivery → usecase → domain (interface)
                       ↑
                  repository (implementation)
```

<br><br>In any case,
<br>`Happy coding (￣▽￣)╭ `