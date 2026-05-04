```                                                                
 ▄▄▄▄ ▄▄▄▄▄ ▄▄▄▄  ▄▄ ▄▄ ▄▄  ▄▄▄▄ ▄▄▄▄▄      ▄▄▄▄  ▄▄▄  ▄▄▄▄  ▄▄▄▄▄ 
███▄▄ ██▄▄  ██▄█▄ ██▄██ ██ ██▀▀▀ ██▄▄  ▄▄▄ ██▀▀▀ ██▀██ ██▄█▄ ██▄▄  
▄▄██▀ ██▄▄▄ ██ ██  ▀█▀  ██ ▀████ ██▄▄▄     ▀████ ▀███▀ ██ ██ ██▄▄▄ 

A backend service designed to support the e-commerce operations of 
    Chia Florist.
Built with Go, this project emphasizes simplicity, clarity, and 
    serves as a practical portfolio for developers.
```

## Objective
- Build a modular and scalable backend service to support core e-commerce operations such as product management, transactions, and order handling
- Demonstrate clean architecture principles with a feature-based structure for maintainability and clarity
- Provide a practical learning and reference project for developers exploring backend development with Go
- Serve as a portfolio project showcasing real-world backend design and implementation

## How to Run
### 1. Prequisites

Make sure the system for development already installed:

- Go (Version 1.20 or higher recommended)

Check Go installation with

```bash
go version
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run The Application

```bash
go run cmd/server/main.go 
```

    Verify server by seeing active log from the application.

## Project Structure

```
├── cmd/
│   ├── migrate/
│   │   └── main.go                # Application entry point
│   └── server/
│       └── main.go                # Application entry point

├── internal/
│   ├── bootstrap/
│   │   ├── app.go
│   │   ├── router.go
│   │   └── container.go

│   ├── common/
│   │   ├── errors/
│   │   │   ├── error.go
│   │   │   └── mapper.go
│   │   ├── http/
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── context.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   └── logger/
│   │       └── logger.go

│   ├── modules/
│   │   └── (...)/
│   │       ├── domain/
│   │       │   └── <entity>.go
│   │       ├── repository/
│   │       │   └── <repository>.go
│   │
│   │       ├── usecase/
│   │
│   │       ├── delivery/
│   │       │   ├── grpc/
│   │       │   ├── kafka/
│   │       │   └── http/
│   │       │       ├── handler.go
│   │       │       ├── dto.go       # response shaping
│   │       │       └── router.go    # route registration
│   │
│   │       └── infra/
│   │           ├── service/
│   │           └── persistence/
│   │               ├── <repository_impl>.go
│   │               ├── <model>.go
│   │               └── mapper.go

│   ├── infra/          # Application-wide infrastructure
│   │   ├── db/         # DB connection setup
│   │   │   ├── postgres.go
│   │   │   ├── mysql.go
│   │   │   └── migrate.go
│   │   ├── logger/
│   │   │   └── logger.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       ├── logging.go
│   │       └── recovery.go

│   └── shared/      # pure utilities 

├── go.mod
└── go.sum
```

## Dependency Flow

```
delivery → usecase → domain (interface)
                       ↑
                  repository (implementation)
```

* Dependencies always point **inward**
* Domain is the **center and most stable layer**

```
Happy coding (￣▽￣)╭ 
```