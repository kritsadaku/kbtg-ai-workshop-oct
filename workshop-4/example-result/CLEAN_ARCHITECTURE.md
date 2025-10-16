# Clean Architecture Project Structure

```
temp-backend/
├── main.go                          # Application entry point
├── go.mod                          # Go module file
├── go.sum                          # Go dependencies checksum
├── users.db                        # SQLite database
├── test_api.sh                     # API testing script
├── README.md                       # Project documentation
│
├── internal/                       # Private application code
│   ├── domain/                     # Enterprise business rules
│   │   └── user.go                # User entity and interfaces
│   │
│   ├── usecase/                    # Application business rules
│   │   └── user_usecase.go        # User business logic
│   │
│   ├── repository/                 # Interface adapters - Data layer
│   │   └── user_repository.go     # User data access implementation
│   │
│   └── handler/                    # Interface adapters - Presentation layer
│       └── user_handler.go        # HTTP handlers for user endpoints
│
└── pkg/                           # Public/shared code
    └── database/                  # Database configuration
        └── database.go            # Database connection and setup
```

## Clean Architecture Layers

### 🔵 Domain Layer (Core)
- **Location**: `internal/domain/`
- **Purpose**: Contains enterprise business rules
- **Dependencies**: None (innermost layer)
- **Contents**:
  - `User` entity
  - `UserRepository` interface
  - `UserUsecase` interface
  - Request/Response DTOs

### 🟢 Usecase Layer (Application Business Rules)
- **Location**: `internal/usecase/`
- **Purpose**: Contains application-specific business rules
- **Dependencies**: Only depends on Domain layer
- **Contents**:
  - `UserUsecase` implementation
  - Business logic validation
  - Orchestration of domain entities

### 🟡 Repository Layer (Interface Adapters)
- **Location**: `internal/repository/`
- **Purpose**: Data access implementations
- **Dependencies**: Domain layer + external libraries (database)
- **Contents**:
  - `UserRepository` implementation
  - Database queries and operations
  - Data mapping

### 🔴 Handler Layer (Interface Adapters)
- **Location**: `internal/handler/`
- **Purpose**: HTTP presentation layer
- **Dependencies**: Domain layer + web framework
- **Contents**:
  - HTTP handlers
  - Request/Response parsing
  - HTTP status codes

### 📦 Infrastructure (Frameworks & Drivers)
- **Location**: `pkg/`, `main.go`
- **Purpose**: External tools and frameworks
- **Contents**:
  - Database connection
  - Web server setup
  - Configuration

## Benefits of This Structure

### ✅ **Dependency Rule**
- Inner layers don't depend on outer layers
- Business logic is independent of frameworks
- Easy to test and maintain

### ✅ **Separation of Concerns**
- Each layer has a single responsibility
- Business logic is isolated from infrastructure
- Easy to modify without affecting other layers

### ✅ **Testability**
- Each layer can be tested independently
- Easy to mock dependencies
- Business logic can be tested without database/HTTP

### ✅ **Maintainability**
- Clear structure makes code easy to navigate
- Changes in one layer don't affect others
- Easy to add new features

### ✅ **Scalability**
- Easy to add new entities/features
- Can swap implementations (e.g., database, web framework)
- Supports multiple interfaces (HTTP, gRPC, CLI)

## API Endpoints (Same as before)

### Health Check
```bash
curl http://localhost:3000/
# Response: {"message":"User Management API is running","version":"2.0.0","architecture":"Clean Architecture"}
```

### User Operations
```bash
# Get all users
curl http://localhost:3000/users

# Get user by ID
curl http://localhost:3000/users/1

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"สมชาย","last_name":"ใจดี","phone":"081-234-5678","email":"somchai@example.com","membership_level":"Gold","points":15420}'

# Update user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"points":20000,"membership_level":"Platinum"}'

# Delete user
curl -X DELETE http://localhost:3000/users/1
```

## How to Run

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Build the application**:
   ```bash
   go build .
   ```

3. **Run the server**:
   ```bash
   go run main.go
   ```

4. **Test APIs**:
   ```bash
   ./test_api.sh
   ```

## Key Improvements in Clean Architecture Version

1. **Better Organization**: Code is organized into clear layers
2. **Dependency Injection**: Components are injected rather than created directly
3. **Interface-Based Design**: Business logic depends on interfaces, not concrete implementations
4. **Better Error Handling**: Centralized error handling with proper HTTP status codes
5. **Easier Testing**: Each layer can be tested independently
6. **Future-Proof**: Easy to add new features or change implementations

## Example: Adding a New Feature

To add a new feature (e.g., User Search):

1. **Domain**: Add method to `UserRepository` interface
2. **Repository**: Implement the search query
3. **Usecase**: Add search business logic
4. **Handler**: Add HTTP endpoint
5. **Main**: Wire everything together

The dependency rule ensures changes flow from outer to inner layers, maintaining system integrity.