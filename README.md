# Golang Fiber Orders Service

A RESTful service built with Go Fiber framework for managing orders and products. This service provides APIs for order management and product operations with PostgreSQL as the database backend.

## Features

- RESTful API endpoints for orders and products
- Swagger documentation
- PostgreSQL database integration using GORM
- Request validation
- Middleware for logging and request validation
- Modular architecture with separation of concerns
- Unit tests

## Tech Stack

| Technology | Version | Description | Link |
|------------|---------|-------------|------|
| Go | 1.24.0 | Programming Language | [golang.org](https://golang.org) |
| Fiber | v2 | Web Framework | [github.com/gofiber/fiber](https://github.com/gofiber/fiber/v2) |
| GORM | v1.30.3 | ORM Library | [gorm.io](https://gorm.io) |
| PostgreSQL | - | Database | [postgresql.org](https://www.postgresql.org) |
| Swagger | v1.16.6 | API Documentation | [github.com/gofiber/swagger](https://github.com/gofiber/swagger) |
| Validator | v10.27.0 | Request Validation | [github.com/go-playground/validator](https://github.com/go-playground/validator) |

## Author

**Srinivas Vadikari**
- 👨‍💻 GitHub: [@svadikari](https://github.com/svadikari)
- 🌐 Project: [golang_fiber_orders](https://github.com/svadikari/golang_fiber_orders)
- 💼 Role: Software Engineer
- 📫 How to reach me:
  - Create an issue in this repository
  - Open a pull request
  - Contact through GitHub discussions

## Project Structure

```
src/
├── main.go                 # Application entry point
├── database/
│   └── database.go        # Database configuration
├── docs/                  # Swagger documentation
├── middleware/
│   ├── logger.go         # Request logging
│   ├── request-validator.go
│   └── users-client.go
├── orders/
│   ├── controllers/      # Order HTTP handlers
│   ├── models/          # Order database models
│   ├── routers/         # Order routes
│   └── schemas/         # Order request/response schemas
└── products/
    ├── controllers/     # Product HTTP handlers
    ├── models/         # Product database models
    ├── repository/     # Product data access layer
    ├── routers/        # Product routes
    ├── schemas/        # Product request/response schemas
    └── services/       # Product business logic
```

## Prerequisites

- Go 1.24.0 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/svadikari/golang_fiber_orders.git
   cd golang_fiber_orders
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up your environment variables for database connection

4. Run the application:
   ```bash
   go run src/main.go
   ```

## API Documentation

The API documentation is available through Swagger UI. After starting the application, visit:
```
http://localhost:8080/swagger/
```

## Testing

To run the tests:
```bash
go test -v ./...
```

## Project Dependencies

Main dependencies:
- github.com/gofiber/fiber/v2 - Web framework
- github.com/gofiber/swagger - Swagger integration
- gorm.io/gorm - ORM
- gorm.io/driver/postgres - PostgreSQL driver
- github.com/go-playground/validator/v10 - Request validation
- github.com/stretchr/testify - Testing framework
- github.com/swaggo/swag - Swagger documentation generator

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Swagger Documentation Generation

Before running the application, generate the Swagger documentation:

1. Install swag CLI:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate Swagger documentation:
   ```bash
   cd src
   swag fmt    # Format swagger comments
   swag init   # Generate swagger documentation
   ```

## Build and Run

1. Build the project:
   ```bash
   # From project root
   go build -o bin/server src/main.go
   ```

2. Run the built binary:
   ```bash
   ./bin/server
   ```

   Or using environment variables:
   ```bash
   PORT=8080 DB_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable" ./bin/server
   ```

## Docker Build (Optional)

1. Build Docker image:
   ```bash
   docker build -t golang-fiber-orders .
   ```

2. Run container:
   ```bash
   docker run -p 8080:8080 \
     -e DB_URL="postgres://user:password@host.docker.internal:5432/dbname?sslmode=disable" \
     golang-fiber-orders
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
