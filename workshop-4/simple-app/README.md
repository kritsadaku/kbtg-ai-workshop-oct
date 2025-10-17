# Simple App - Clean Architecture API with Rust & Axum

A REST API built with Rust using Axum framework, following Clean Architecture principles with SQLite database integration.

## 🏗️ Architecture Overview

This project implements Clean Architecture with the following layers:

```
src/
├── domain/          # Business Logic Layer (Entities & Repository Interfaces)
├── application/     # Use Cases Layer (Business Rules)
├── infrastructure/  # Data Access Layer (Database Implementation)
├── presentation/    # Interface Layer (HTTP Handlers & Routes)
└── main.rs         # Application Entry Point
```

### Architecture Layers

#### 1. **Domain Layer** (`src/domain/`)
- **Entities**: Core business objects (`User`)
- **Repository Interfaces**: Abstract contracts for data access
- **Business Rules**: Entity validation and business logic
- **Dependencies**: None (pure business logic)

#### 2. **Application Layer** (`src/application/`)
- **Use Cases**: Application-specific business rules (`UserService`)
- **Orchestration**: Coordinates between domain and infrastructure
- **Dependencies**: Domain layer only

#### 3. **Infrastructure Layer** (`src/infrastructure/`)
- **Database Implementation**: SQLite repository implementation
- **External Services**: Database connections and queries
- **Dependencies**: Domain layer for interfaces

#### 4. **Presentation Layer** (`src/presentation/`)
- **HTTP Handlers**: API endpoint implementations
- **Request/Response Models**: Data transfer objects
- **Routing**: URL path definitions
- **Dependencies**: Application layer for use cases

## 🚀 Features

- **RESTful API** with full CRUD operations for Users
- **Clean Architecture** with proper separation of concerns
- **SQLite Database** with automatic schema creation
- **Swagger Documentation** with interactive API explorer
- **Input Validation** with proper error handling
- **Email Uniqueness** constraint enforcement
- **Automatic Database Setup** creates DB file if not exists

## 📊 Database Schema

Based on the LBK Points System entities specification:

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    member_since TEXT NOT NULL,
    membership_level TEXT NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
```

## 🔗 API Endpoints

### Core Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `GET` | `/` | Health check / Hello World | - |
| `GET` | `/users` | List all users (with pagination) | - |
| `GET` | `/users/{id}` | Get user by ID | - |
| `POST` | `/users` | Create new user | `CreateUserRequest` |
| `PUT` | `/users/{id}` | Update existing user | `UpdateUserRequest` |
| `DELETE` | `/users/{id}` | Delete user | - |

### Documentation
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/swagger-ui` | Interactive API documentation |
| `GET` | `/api-docs/openapi.json` | OpenAPI specification |

### Query Parameters

**List Users** (`GET /users`):
- `limit` (optional): Number of users to return (default: 100)
- `offset` (optional): Number of users to skip (default: 0)

Example: `GET /users?limit=10&offset=20`

## 📝 Request/Response Models

### CreateUserRequest
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+66812345678",
  "email": "john.doe@example.com",
  "membership_level": "Gold"
}
```

### UpdateUserRequest
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "phone": "+66812345678",
  "email": "john.smith@example.com",
  "membership_level": "Platinum"
}
```

### User Response
```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+66812345678",
  "email": "john.doe@example.com",
  "member_since": "2025-10-17T10:30:00Z",
  "membership_level": "Gold",
  "points": 1500,
  "created_at": "2025-10-17T10:30:00Z",
  "updated_at": "2025-10-17T10:30:00Z"
}
```

## 🛠️ Setup & Installation

### Prerequisites
- **Rust** (latest stable version)
- **Cargo** (comes with Rust)

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd simple-app

# Install dependencies
cargo build

# Run the application
cargo run
```

### Development
```bash
# Check code without building
cargo check

# Run with auto-reload (install cargo-watch first)
cargo install cargo-watch
cargo watch -x run

# Run tests
cargo test
```

## 🔧 Configuration

### Database
- **Type**: SQLite
- **File**: `users.db` (created automatically)
- **Location**: Same directory as the executable

### Server
- **Host**: `0.0.0.0`
- **Port**: `3000`
- **Environment**: Development

## 📡 Usage Examples

### Create a User
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "phone": "+66812345678",
    "email": "alice@example.com",
    "membership_level": "Silver"
  }'
```

### Get All Users
```bash
curl http://localhost:3000/users?limit=5&offset=0
```

### Get User by ID
```bash
curl http://localhost:3000/users/1
```

### Update User
```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Smith",
    "membership_level": "Gold"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:3000/users/1
```

## 🔍 Testing with Swagger UI

1. Start the server: `cargo run`
2. Open browser: http://localhost:3000/swagger-ui
3. Explore and test all API endpoints interactively

## 🏛️ Clean Architecture Benefits

### 1. **Separation of Concerns**
- Each layer has a single responsibility
- Business logic is isolated from external concerns
- Easy to understand and maintain

### 2. **Testability**
- Business logic can be tested independently
- Mock implementations for external dependencies
- Clear interfaces for testing

### 3. **Flexibility**
- Easy to swap database implementations
- Framework-independent business logic
- Adaptable to changing requirements

### 4. **Maintainability**
- Clear module boundaries
- Reduced coupling between components
- Easier debugging and troubleshooting

## 📁 Project Structure Details

```
src/
├── domain/
│   ├── mod.rs              # Module declarations
│   ├── user.rs             # User entity with business logic
│   └── repository.rs       # Repository trait definition
├── application/
│   ├── mod.rs              # Module declarations
│   └── user_service.rs     # User use cases and business rules
├── infrastructure/
│   ├── mod.rs              # Module declarations
│   └── repository.rs       # SQLite implementation
├── presentation/
│   ├── mod.rs              # Module declarations
│   ├── handlers.rs         # HTTP request handlers
│   └── routes.rs           # Route definitions
└── main.rs                 # Application bootstrap
```

## 🔒 Error Handling

The API provides consistent error responses:

```json
{
  "error": "Detailed error message"
}
```

### Common HTTP Status Codes
- `200 OK` - Successful GET/PUT requests
- `201 Created` - Successful POST requests
- `204 No Content` - Successful DELETE requests
- `400 Bad Request` - Invalid input data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

## 🚀 Deployment

### Production Build
```bash
cargo build --release
```

### Docker (Optional)
```dockerfile
FROM rust:1.70 as builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/target/release/simple-app .
EXPOSE 3000
CMD ["./simple-app"]
```

## 📚 Dependencies

### Core Dependencies
- **axum**: Web application framework
- **tokio**: Async runtime
- **sqlx**: Async SQL toolkit
- **serde**: Serialization/deserialization
- **chrono**: Date and time handling

### Development Dependencies
- **utoipa**: OpenAPI documentation
- **utoipa-swagger-ui**: Swagger UI integration

## 🤝 Contributing

1. Follow Clean Architecture principles
2. Write tests for new features
3. Update documentation as needed
4. Ensure all tests pass before submitting

## 📄 License

This project is licensed under the MIT License.

---

**Happy Coding! 🦀**