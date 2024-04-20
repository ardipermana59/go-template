# Go Gin GORM API with Consistent Error Handling & Tests

A production-ready RESTful API built with Go, Gin, GORM, featuring JWT authentication, role-based authorization, consistent error responses, and comprehensive tests.

## ✨ Features

- ✅ **Consistent Error Responses**: All errors follow the same format
- ✅ **User Authentication**: Register & Login with JWT
- ✅ **Password Security**: bcrypt hashing
- ✅ **Authorization**: Role-based (User, Admin) & Owner-based
- ✅ **Middleware**: Auth & Role middleware
- ✅ **Validation**: Readable error messages
- ✅ **Integration Tests**: Comprehensive test coverage
- ✅ **Clean Architecture**: Repository → Service → Handler pattern
- ✅ **Database**: MySQL with GORM ORM

## 📁 Project Structure

```
go-gin-gorm-api/
├── cmd/api/
│   └── main.go                     # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── internal/
│   ├── auth/
│   │   └── jwt.go                  # JWT service
│   ├── middleware/
│   │   └── auth.go                 # Auth & Role middleware
│   ├── common/
│   │   ├── apperror/
│   │   │   └── errors.go           # Consistent error definitions
│   │   └── response/
│   │       └── response.go         # Standard API responses
│   ├── user/
│   │   ├── model.go                # User models & DTOs
│   │   ├── repository.go           # User database operations
│   │   ├── service.go              # User business logic
│   │   └── handler.go              # User HTTP handlers
│   └── post/
│       ├── model.go                # Post models & DTOs
│       ├── repository.go           # Post database operations
│       ├── service.go              # Post business logic
│       └── handler.go              # Post HTTP handlers
├── pkg/
│   ├── database/
│   │   └── database.go             # Database connection
│   └── validator/
│       └── validator.go            # Custom validation formatter
├── tests/
│   └── integration/
│       └── user_test.go            # Integration tests
├── .env                            # Environment variables
├── .env.test                       # Test environment variables
├── docker-compose.yml              # MySQL container
├── Makefile                        # Build commands
└── test.http                       # API test requests
```

## 🚀 Quick Start


### 2. Start MySQL Database

```bash
cd go-gin-gorm-api
make docker-up
```

### 3. Run Application

```bash
make run
```

Server will start on: **http://localhost:8080**

### 4. Run Tests

```bash
make test
```

## 🎯 Consistent Error Response Format

All errors now follow the same format as validation errors:

### Success Response
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Validation Error Response
```json
{
  "success": false,
  "message": "Validation failed",
  "error": [
    {
      "field": "name",
      "message": "The name must be at least 3 characters"
    },
    {
      "field": "email",
      "message": "The email must be a valid email address"
    },
    {
      "field": "password",
      "message": "The password must be at least 6 characters"
    }
  ]
}
```

### Business Logic Error Response
```json
{
  "success": false,
  "message": "Failed to register user",
  "error": [
    {
      "field": "email",
      "message": "The email has already been taken"
    }
  ]
}
```

### Authentication Error Response
```json
{
  "success": false,
  "message": "Login failed",
  "error": [
    {
      "field": "credentials",
      "message": "The provided credentials are invalid"
    }
  ]
}
```

### Authorization Error Response
```json
{
  "success": false,
  "message": "Failed to update post",
  "error": [
    {
      "field": "ownership",
      "message": "You don't have permission to modify this resource"
    }
  ]
}
```

### Not Found Error Response
```json
{
  "success": false,
  "message": "Post not found",
  "error": [
    {
      "field": "post",
      "message": "The post could not be found"
    }
  ]
}
```

## 📝 Error Types Defined

The system includes predefined errors for common scenarios:

- `EmailAlreadyExists()` - Email is already registered
- `InvalidCredentials()` - Wrong email or password
- `OldPasswordIncorrect()` - Old password doesn't match
- `PostNotFound()` - Post doesn't exist
- `UserNotFound()` - User doesn't exist
- `Unauthorized()` - Not authorized
- `OwnershipRequired()` - Not the resource owner
- `InvalidID()` - Invalid ID format
- `DatabaseError()` - Database operation failed

## 🧪 Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Test Coverage Includes
- ✅ User Registration (success & duplicate email)
- ✅ User Login (success & invalid credentials)
- ✅ Get Profile (with & without token)
- ✅ Validation Errors
- ✅ Authorization Errors
- ✅ Consistent Error Formats

## 📚 API Endpoints

### Public Endpoints
```http
POST   /api/v1/auth/register        # Register new user
POST   /api/v1/auth/login           # Login user
GET    /api/v1/posts                # Get all posts
GET    /api/v1/posts/:id            # Get post by ID
GET    /api/v1/users/:user_id/posts # Get user's posts
```

### Protected Endpoints (Requires Authentication)
```http
GET    /api/v1/profile              # Get user profile
PUT    /api/v1/profile              # Update profile
PUT    /api/v1/change-password      # Change password
POST   /api/v1/posts                # Create post
GET    /api/v1/posts/my             # Get my posts
PUT    /api/v1/posts/:id            # Update own post
DELETE /api/v1/posts/:id            # Delete own post
```

### Admin Only Endpoints
```http
GET    /api/v1/admin/users          # Get all users
GET    /api/v1/admin/users/:id      # Get user by ID
PUT    /api/v1/admin/users/:id      # Update any user
DELETE /api/v1/admin/users/:id      # Delete user
```

## 🛠️ Make Commands

```bash
make run            # Run the application
make build          # Build binary to bin/api
make test           # Run integration tests
make test-coverage  # Run tests with coverage report
make deps           # Download and tidy dependencies
make docker-up      # Start MySQL container
make docker-down    # Stop MySQL container
make docker-logs    # View MySQL logs
make clean          # Clean build artifacts
```

## 🔒 Security Features

- ✅ Password hashing with bcrypt (cost 10)
- ✅ JWT token with expiration
- ✅ Authorization header validation
- ✅ Role-based access control
- ✅ Owner-based resource protection
- ✅ Input validation
- ✅ SQL injection prevention (GORM ORM)
- ✅ Consistent error messages (no sensitive info leakage)

## 📦 Dependencies

```go
github.com/gin-gonic/gin                  // HTTP framework
github.com/golang-jwt/jwt/v5              // JWT authentication
golang.org/x/crypto                       // bcrypt password hashing
github.com/go-playground/validator/v10    // Input validation
gorm.io/gorm                              // ORM
gorm.io/driver/mysql                      // MySQL driver
github.com/joho/godotenv                  // Environment variables
github.com/stretchr/testify               // Testing framework
```

## 🎓 Testing Best Practices

The project includes comprehensive integration tests that:
- Setup and teardown test database
- Test all major user flows
- Verify consistent error formats
- Test authentication and authorization
- Validate business logic

Example test structure:
```go
func TestRegisterUser(t *testing.T) {
    t.Run("Success - Register new user", func(t *testing.T) {
        // Test successful registration
    })
    
    t.Run("Fail - Email already exists", func(t *testing.T) {
        // Test duplicate email with consistent error format
    })
    
    t.Run("Fail - Validation errors", func(t *testing.T) {
        // Test validation errors
    })
}
```

## 🐛 Troubleshooting

### Database Connection Failed
```bash
make docker-logs    # Check MySQL logs
make docker-down    # Stop MySQL
make docker-up      # Start MySQL
```

### Tests Failing
```bash
# Ensure test database is clean
mysql -u root -p
DROP DATABASE IF EXISTS testdb_test;
CREATE DATABASE testdb_test;
```

