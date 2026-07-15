# 🔗 URL Shortener (Clean Architecture & System Design)

A highly scalable, production-ready URL Shortener built with **Go 1.25**, applying **Clean Architecture** principles and robust system design patterns. This project is crafted as a learning and reference implementation for system design interviews and software engineering best practices.

## ✨ Features

- **Blazing Fast Redirection**: Generates unique, true random 7-character short codes (e.g. `Xk2a0B1`) optimized for read-heavy operations.
- **Collision Handling**: Implements a retry mechanism at the usecase layer to guarantee uniqueness against database collisions during concurrent short code generation.
- **Click Analytics**: Tracks total clicks for every shortened URL asynchronously.
- **Clean Architecture**: Strictly separated layers (Domain, Usecase, Repository, Delivery) ensuring maximum testability and maintainability.
- **Server-Side Rendered (SSR) UI**: Includes a lightweight, vanilla HTML/JS frontend served directly by the Go backend—no separate build step required.
- **Testing & CI**: Comprehensive Unit Tests (using `testify/mock`) and continuous integration pipeline configured with `golangci-lint` via GitHub Actions.
- **Health Checks**: Built-in `/api/v1/health` endpoint for infrastructure monitoring.

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Router**: [go-chi/chi](https://github.com/go-chi/chi) (Lightweight, idiomatic `net/http` router)
- **Database**: PostgreSQL 15
- **DB Driver**: [jackc/pgx/v5](https://github.com/jackc/pgx) (High-performance, pure Go Postgres driver)
- **Frontend**: HTML5, TailwindCSS (via CDN), Vanilla JavaScript
- **Testing**: [testify](https://github.com/stretchr/testify)
- **Linter**: golangci-lint

## 🏗️ Architecture Design

### Clean Architecture
This project adheres to Uncle Bob's Clean Architecture. Dependencies always point **inwards** towards the Domain layer.

1. **Domain (Entities)**: The innermost layer. Contains pure Go structs (`URL`) and interface contracts (`URLRepository`, `URLUsecase`). It has zero knowledge of databases or web frameworks.
2. **Usecase (Business Logic)**: The brain of the application. Orchestrates the flow: validates inputs, generates a random 7-character string, handles potential database constraint collisions with retries, and persists the entity.
3. **Repository (Interface Adapters)**: Database layer. Implements the `URLRepository` interface to execute raw SQL queries against PostgreSQL.
4. **Delivery (Interface Adapters)**: HTTP transport layer. Parses incoming HTTP requests, invokes usecases, and returns JSON/HTML responses.

### 🔄 Request Lifecycle (Code Flow)
To understand how the code executes, let's trace a typical request (e.g., Shortening a URL):
1. **Entry Point (`cmd/api/main.go`)**: The application starts here. It connects to the database, initializes the `Repository` and `Usecase` layers, injects them into the HTTP `Handler` (Delivery), and starts the server.
2. **Delivery Layer (`internal/delivery/http/url_handler.go`)**: The router catches the `POST /api/v1/urls` request. It parses the JSON body to get the original URL, and passes it to the Usecase layer.
3. **Usecase Layer (`internal/usecase/url_usecase.go`)**: The business logic executes.
   - It validates the URL format.
   - It generates a random 7-character short code using the `random` utility.
   - It bundles this data into a `domain.URL` entity and calls `repo.Store()` to save it.
   - If a collision occurs (database unique constraint error), it retries automatically.
4. **Repository Layer (`internal/repository/postgres/url_repo.go`)**: Executes the raw `INSERT INTO urls ...` SQL query to save the data into PostgreSQL.
5. **Response**: The result flows backward. The Usecase returns the `URL` entity to the Handler, and the Handler formats it as a `201 Created` JSON response to the user.

### Folder Structure
```text
url-shortener/
├── .github/
│   └── workflows/
│       └── golangci-lint.yml # CI pipeline configuration
├── cmd/
│   └── api/
│       └── main.go           # Application entrypoint & dependency injection wiring
├── internal/
│   ├── delivery/
│   │   └── http/             # HTTP Handlers (API endpoints & UI rendering)
│   ├── domain/               # Core business models & interface definitions
│   ├── repository/
│   │   └── postgres/         # PostgreSQL database implementation
│   └── usecase/              # Application business rules & Unit Tests
├── migrations/               # Database initialization SQL scripts
├── pkg/
│   └── random/               # Reusable utility for true random string generation
├── web/                      # HTML templates for Frontend SSR
├── .golangci.yml             # Linter rules
├── docker-compose.yml        # Docker configuration for PostgreSQL
├── go.mod                    # Go module dependencies (v1.25)
└── README.md
```

## 🗄️ Database Design

### Why PostgreSQL?
While URL Shorteners seem like a great fit for NoSQL (Key-Value), we opted for an RDBMS (PostgreSQL) for **data integrity**. The `short_code` requires absolute uniqueness. PostgreSQL handles `UNIQUE` constraints flawlessly and provides strict type safety, which prevents anomalies. The standard industry best practice for massive scale is typically an RDBMS as the *Source of Truth*, fronted by a caching layer (like Redis) for fast reads.

### Table: `urls`
| Column Name  | Type | Constraints | Description |
| ------------ | ---- | ----------- | ----------- |
| `id`         | `BIGSERIAL` | `PRIMARY KEY` | Auto-incrementing numeric internal ID |
| `short_code` | `VARCHAR(7)`| `UNIQUE, NOT NULL`| The random 7-character string (e.g., `Xk2a0B1`) |
| `original_url`| `VARCHAR(2048)`| `NOT NULL`| The target long URL |
| `click_count`| `BIGINT` | `DEFAULT 0` | Total clicks tracker |
| `created_at` | `TIMESTAMP` | `DEFAULT CURRENT_TIMESTAMP` | Creation timestamp |

*(Note: An Index is automatically created on `short_code` in our migrations to ensure `O(log N)` lookup times for fast redirection).*

## 🚀 Getting Started

### Prerequisites
- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://www.docker.com/products/docker-desktop) (for the database)

### Running Locally

1. **Start the Database**
   Fire up the PostgreSQL instance. The initialization script in `./migrations` will automatically create the `urls` table.
   ```bash
   docker-compose up -d
   ```

2. **Download Dependencies**
   Ensure dependencies are resolved.
   ```bash
   go mod tidy
   ```

3. **Run the Tests & Linter**
   ```bash
   go test ./...
   golangci-lint run
   ```

4. **Run the Application**
   ```bash
   go run cmd/api/main.go
   ```
   *The server will start at `http://localhost:8080`.*

## 📚 API Endpoints

### 1. Shorten a URL
- **POST** `/api/v1/urls`
- **Body:**
  ```json
  {
    "original_url": "https://en.wikipedia.org/wiki/System_design"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": 1,
    "short_code": "Xk2a0B1",
    "original_url": "https://en.wikipedia.org/wiki/System_design",
    "click_count": 0,
    "created_at": "2026-07-15T10:00:00Z"
  }
  ```

### 2. Redirect to Original URL
- **GET** `/{shortCode}`
- **Response:** `301 Moved Permanently` (Redirects user to original URL)

### 3. Get URL Statistics
- **GET** `/api/v1/urls/{shortCode}/stats`
- **Response (200 OK):**
  ```json
  {
    "short_code": "Xk2a0B1",
    "click_count": 42
  }
  ```

### 4. Health Check
- **GET** `/api/v1/health`
- **Response (200 OK):**
  ```json
  {
    "status": "UP"
  }
  ```
