# Go E-Wallet Service 🪙💸

A robust, scalable backend service for E-Wallet applications built with Golang (Gin Framework) using Clean Architecture principles. This project implements ACID Database Transactions to ensure data integrity during financial operations.

## 🌟 Key Features

- Clean Architecture (Handler -> Usecase -> Repository) to ensure separation of concerns.
- RESTful API with Gin Framework.
- PostgreSQL Database with raw SQL (pgx driver) for maximum performance.
- ACID Transactions for Money Transfer (Atomic operations).
- JWT Authentication (JSON Web Token).
- Unit Testing with Testify (Mocking & Assertions).
- Middleware for secure route protection.

## 🛠 Tech Stack

- Language: Go (Golang)
- Framework: Gin Gonic
- Database: PostgreSQL
- Driver: pgx/v5
- Testing: Testify & Mockery
- Security: Bcrypt (Hashing) & JWT-Go

## 📂 Project Structure

```text
ewallet-service/
├── cmd/
│   └── api/          # Entry point (main.go)
├── config/           # Database Connection
├── internal/
│   ├── handler/      # HTTP Delivery Layer
│   ├── usecase/      # Business Logic Layer
│   ├── repository/   # Data Access Layer (SQL)
│   ├── middleware/   # Auth Middleware
│   └── model/        # Structs & DTOs
├── .env              # Environment Variables
└── database.sql      # SQL Schema
```

## 🚀 How to Run

### 1. Clone the repository

```bash
git clone https://github.com/OrryFrasetyo/ewallet-service.git
cd ewallet-service
```

### 2. Setup Database

- Create a PostgreSQL database named ewallet_db.

- Run the script in database.sql to create tables.

### 3. Setup Environtment

- Create .env file based on example:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=ewallet_db
JWT_SECRET=your_secret_key
```

### 4. Run the Server

```bash
go run cmd/api/main.go
```

### 5. Run Tests

```bash
go test ./internal/usecase/... -v
```

## 🔌API Endpoints

| **Method** |     **Endpoint**     |   **Description**  | **Auth** |
|:----------:|:--------------------:|:------------------:|:--------:|
|    POST    |   /api/v1/register   |  Register new user |    No    |
|    POST    |     /api/v1/login    |  Login & Get Token |    No    |
|    POST    |     /api/v1/topup    |    Topup Balance   |  **Yes** |
|    POST    |   /api/v1/transfer   |   Transfer Money   |  **Yes** |
|     GET    |    /api/v1/balance   | Get Wallet Balance |  **Yes** |
|     GET    | /api/v1/transactions |     Get History    |  **Yes** |


