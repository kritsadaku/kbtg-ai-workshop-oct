# GitHub Copilot Instructions for Temp-Backend Project

## Project Overview
This is a Go-based REST API backend built with **Clean Architecture** principles. The system manages users and point transfers with a comprehensive ledger system. 

### Core Technologies
- **Language**: Go 1.23.1
- **Framework**: Fiber v2 (web framework)
- **Database**: SQLite3 with raw SQL queries
- **Architecture**: Clean Architecture (Domain, Repository, UseCase, Handler layers)

## Architecture Patterns

### Clean Architecture Structure
```
internal/
├── domain/          # Business entities and interfaces
├── handler/         # HTTP request handlers (Controllers)
├── repository/      # Data access layer
└── usecase/         # Business logic layer
pkg/
└── database/        # Database connection and migrations
```

### Layer Responsibilities
1. **Domain Layer**: Contains business entities, DTOs, and interface definitions
2. **Handler Layer**: HTTP request/response handling and routing
3. **UseCase Layer**: Business logic implementation and validation
4. **Repository Layer**: Database operations and data persistence

## Code Generation Guidelines

### 1. Domain Entities and Interfaces
- All business entities should be defined in `internal/domain/`
- Use struct tags for JSON serialization: `json:"field_name"`
- Create separate request/response DTOs for API operations
- Define repository and usecase interfaces in domain package
- Use pointer fields for optional updates: `*string`, `*float64`

**Example Pattern:**
```go
// Entity
type EntityName struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Request DTOs
type CreateEntityRequest struct {
    Name string `json:"name" validate:"required"`
}

type UpdateEntityRequest struct {
    Name *string `json:"name,omitempty"`
}

// Repository Interface
type EntityRepository interface {
    GetAll() ([]*Entity, error)
    GetByID(id int) (*Entity, error)
    Create(entity *Entity) (*Entity, error)
    Update(id int, updates map[string]interface{}) (*Entity, error)
    Delete(id int) error
    Exists(id int) (bool, error)
}
```

### 2. Handler Layer (HTTP Controllers)
- Use Fiber framework for HTTP handling
- Place all handlers in `internal/handler/`
- Constructor pattern: `NewXxxHandler(usecase domain.XxxUsecase) *XxxHandler`
- Standard HTTP status codes and error responses
- Parse URL parameters with `strconv.Atoi()` for IDs
- Use `c.BodyParser()` for request body parsing

**Handler Pattern:**
```go
type EntityHandler struct {
    entityUsecase domain.EntityUsecase
}

func NewEntityHandler(entityUsecase domain.EntityUsecase) *EntityHandler {
    return &EntityHandler{entityUsecase: entityUsecase}
}

func (h *EntityHandler) GetEntity(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
    }
    
    entity, err := h.entityUsecase.GetEntityByID(id)
    if err != nil {
        if err.Error() == "entity not found" {
            return c.Status(404).JSON(fiber.Map{"error": "Entity not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch entity"})
    }
    
    return c.JSON(entity)
}
```

### 3. UseCase Layer (Business Logic)
- Implement business logic and validation
- Constructor pattern: `NewXxxUsecase(repo domain.XxxRepository) domain.XxxUsecase`
- Comprehensive input validation before database operations
- Return domain-specific errors
- Handle edge cases and business rules

**UseCase Pattern:**
```go
type entityUsecase struct {
    entityRepo domain.EntityRepository
}

func NewEntityUsecase(entityRepo domain.EntityRepository) domain.EntityUsecase {
    return &entityUsecase{entityRepo: entityRepo}
}

func (u *entityUsecase) CreateEntity(req *domain.CreateEntityRequest) (*domain.Entity, error) {
    // Validation
    if err := u.validateCreateRequest(req); err != nil {
        return nil, err
    }
    
    // Business logic
    entity := &domain.Entity{
        Name: req.Name,
    }
    
    return u.entityRepo.Create(entity)
}
```

### 4. Repository Layer (Data Access)
- Use raw SQL queries with `database/sql` package
- Constructor pattern: `NewXxxRepository(db *sql.DB) domain.XxxRepository`
- Handle SQL errors appropriately (`sql.ErrNoRows` returns nil)
- Use parameterized queries to prevent SQL injection
- Format timestamps consistently: `"2006-01-02 15:04:05"`

**Repository Pattern:**
```go
type entityRepository struct {
    db *sql.DB
}

func NewEntityRepository(db *sql.DB) domain.EntityRepository {
    return &entityRepository{db: db}
}

func (r *entityRepository) GetByID(id int) (*domain.Entity, error) {
    entity := &domain.Entity{}
    var createdAt, updatedAt string
    
    err := r.db.QueryRow(`
        SELECT id, name, created_at, updated_at 
        FROM entities WHERE id = ?`, id).Scan(
        &entity.ID, &entity.Name, &createdAt, &updatedAt)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    // Parse timestamps
    if entity.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
        return nil, err
    }
    if entity.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
        return nil, err
    }
    
    return entity, nil
}
```

### 5. Database Migrations
- Add table creation in `pkg/database/database.go`
- Use SQLite-compatible SQL syntax
- Include appropriate constraints and indexes
- Follow existing naming conventions (snake_case for columns)

### 6. Testing Patterns
- Create mock repositories for unit testing
- Test files should end with `_test.go`
- Use table-driven tests for multiple scenarios
- Mock struct should implement the repository interface
- Include edge cases and validation testing

**Mock Repository Pattern:**
```go
type mockEntityRepository struct {
    entities  map[int]*domain.Entity
    nextID    int
    shouldErr bool
}

func (m *mockEntityRepository) GetByID(id int) (*domain.Entity, error) {
    if m.shouldErr {
        return nil, errors.New("mock error")
    }
    entity, exists := m.entities[id]
    if !exists {
        return nil, nil
    }
    return entity, nil
}
```

## Business Rules & Validations

### User Management Rules
1. **Name Length Validation**: First name and last name must not exceed 3 characters (supports Unicode/Thai characters)
2. **Membership Levels**: Only allow "Bronze", "Silver", "Gold", "Platinum"
3. **Default Values**: New users default to "Bronze" membership level
4. **Required Fields**: first_name, last_name, phone, email are mandatory

### Transfer System Rules
1. **Idempotency**: Use idempotency keys to prevent duplicate transfers
2. **Point Validation**: Ensure sender has sufficient points before transfer
3. **Status Management**: Handle transfer states (pending, processing, completed, failed, cancelled, reversed)
4. **Ledger Tracking**: Create ledger entries for all point changes
5. **Transaction Safety**: Use database transactions for point transfers

## API Response Patterns

### Success Responses
```go
// Single entity
return c.JSON(entity)

// List with pagination
return c.JSON(fiber.Map{
    "data":     entities,
    "page":     page,
    "pageSize": pageSize,
    "total":    total,
})

// Creation success
return c.Status(201).JSON(entity)

// Deletion success
return c.JSON(fiber.Map{"message": "Entity deleted successfully"})
```

### Error Responses
```go
// Validation errors
return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})

// Not found
return c.Status(404).JSON(fiber.Map{"error": "Entity not found"})

// Server errors
return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
```

## Database Conventions

### Table Design
- Use snake_case for table and column names
- Always include `id`, `created_at`, `updated_at` fields
- Use TEXT type for string fields in SQLite
- Use INTEGER for numeric fields (points are stored as integers)
- Add appropriate foreign key constraints

### Indexing Strategy
- Create indexes on frequently queried columns
- Use composite indexes for multi-column queries
- Index foreign key columns for join performance

## File Organization

### Naming Conventions
- Files: `snake_case.go`
- Packages: lowercase single words
- Types: PascalCase
- Functions: PascalCase (exported), camelCase (private)
- Constants: UPPER_SNAKE_CASE or PascalCase

### Import Organization
```go
import (
    // Standard library
    "database/sql"
    "errors"
    "time"
    
    // Third-party packages
    "github.com/gofiber/fiber/v2"
    
    // Internal packages
    "temp-backend/internal/domain"
)
```

## Development Guidelines

### Error Handling
- Use domain-specific error messages
- Handle `sql.ErrNoRows` by returning nil (not an error)
- Provide meaningful error messages for API consumers
- Log errors appropriately without exposing sensitive information

### Code Quality
- Keep functions focused and single-purpose
- Use dependency injection for testability
- Validate inputs at the usecase layer
- Use interfaces for loose coupling between layers
- Write comprehensive unit tests with mocks

### Performance Considerations
- Use prepared statements for repeated queries
- Implement proper indexing strategy
- Consider pagination for list endpoints
- Use database transactions for multi-step operations

This guide should help GitHub Copilot generate code that follows the established patterns and maintains consistency with the existing codebase architecture.
