# Bookmark User Service

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Gin](https://img.shields.io/badge/Gin-v1.12-008080?style=flat&logo=gin)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=flat&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-Alpine-DC382D?style=flat&logo=redis)
![Docker](https://img.shields.io/badge/Docker-Multi--stage-2496ED?style=flat&logo=docker)

`bookmark-user-service` is a Go microservice responsible for user authentication and self-profile management in the **Bookmark Management** system.

---

## Features

- **Authentication**: User registration and login issuing **RS256 JWT tokens** (RSA 2048-bit key pair).
- **User Profile Management**: Authenticated endpoints to view and update own profile information (`/v1/self/info`).
- **Security & Resilience**: Redis-backed rate limiting middleware to protect endpoints from abuse.
- **Database Migrations**: Automatic schema migrations powered by `golang-migrate`.
- **Health Check**: `/health-check` endpoint checking service status, Redis connections.
- **API Documentation**: Interactive Swagger UI auto-generated via Swaggo.

---

##  Tech Stack

- **Language**: Go `1.26+`
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL 17 + GORM
- **Cache & Rate Limiting**: Redis Alpine + `go-redis`
- **Security**: RS256 JWT (`golang-jwt/jwt v5`) & Custom Password Hasher
- **API Docs**: Swagger UI (`swaggo/gin-swagger`)
- **Shared Package**: `github.com/TNJKL/bookmark-pkg`

---

## Architecture & Project Structure

The project follows **Clean Architecture** to separate transport, business logic, and data access layers:

`Client Request` ➔ `Gin Router & Middlewares` ➔ `Handlers` ➔ `Services` ➔ `Repositories (GORM / Redis)`

```plain
bookmark-user-service/
├── cmd/api/main.go          # Application entrypoint & Swagger annotations
├── docs/                    # Auto-generated Swagger specifications
├── internal/
│   ├── api/                 # Router setup, middlewares & configuration
│   ├── app/
│   │   ├── handler/         # HTTP handlers & DTO validation
│   │   ├── model/           # Domain models (User, HealthCheck)
│   │   ├── repository/        # Data access layer (PostgreSQL & Redis)
│   │   └── service/         # Core business logic (Auth, Password hashing)
│   └── infrastructure/      # Database, Redis & JWT initialization
├── migrations/              # SQL migration files
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Local PostgreSQL + Redis + App setup
└── Makefile                 # Build & run shortcuts
```

---

##  Quick Start

### 1. Environment Setup

Create a `.env` file in the root directory:

```ini
APP_PORT=8080
SERVICE_NAME=user_service
LOG_LEVEL=info
BASE_PATH=/

REDIS_ADDRESS=YOUR_REDIS_ADDRESS

DB_HOST=YOUR_DB_HOST
DB_PORT=YOUR_DB_PORT
DB_USER=YOUR_DB_USER
DB_PASSWORD=YOUR_DB_PASSWORD
DB_NAME=YOUR_DB_NAME
```

### 2. Generate RSA Key Pair

Generate the RSA 2048-bit key pair required for JWT signing & verification (`private.pem` & `public.pem`):

```bash
make generate-rsa-key
```

### 3. Run Application

**Option A: Local Development** (Requires local PostgreSQL & Redis):
```bash
make dev-run
```

**Option B: Docker Compose** (Runs PostgreSQL, Redis & Service in containers):
```bash
make docker-up
```

---

##  API Endpoints

Once the application is running, access the interactive Swagger documentation at:
**http://localhost:8080/swagger/index.html** (or the configured host and port)

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :---: | :--- |
| `GET` | `/health-check` | No | Health check status & instance info |
| `GET` | `/swagger/*any` | No | Swagger UI documentation |
| `POST` | `/v1/users/register` | No | Register a new user |
| `POST` | `/v1/users/login` | No | Authenticate user & receive RS256 JWT |
| `GET` | `/v1/self/info` | Bearer JWT | Get current user profile info |
| `PUT` | `/v1/self/info` | Bearer JWT | Update current user profile info |

---

##  Testing

Run unit tests and check test coverage against the **80% threshold**:

```bash
make test
```

For isolated testing in Docker:
```bash
make docker-test
```

---

##  Useful Makefile Commands

| Command | Description |
| :--- | :--- |
| `make dev-run` | Regenerate Swagger docs and run API locally |
| `make test` | Run tests with coverage threshold check |
| `make generate-rsa-key` | Generate RSA key pair for JWT |
| `make docker-up` | Start all services via Docker Compose |
| `make docker-down` | Stop Docker Compose containers |
| `make swag` | Regenerate Swagger docs |

---

## License

Just a repo I use for learning, so there’s no license here :D