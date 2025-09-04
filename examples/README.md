# Hikari-Go Examples

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

This folder contains practical examples demonstrating the advanced features of the Hikari-Go framework. Each example is a complete and functional application showcasing modern API development patterns with route groups, middleware, and organized structure.

## 🆕 Framework Features Demonstrated

- **🏗️ Route Groups:** Hierarchical route organization with shared prefixes
- **🔧 Middleware Stack:** Global and group-specific middleware application
- **📋 API Versioning:** Professional API structure with `/api/v1` pattern
- **🩺 Health Checks:** Monitoring endpoints for production readiness
- **🔄 Backward Compatibility:** Smooth migration paths
- **🎯 Pattern Normalization:** Automatic route pattern standardization

## 📋 Examples List

### 1. [Todo App](./todo-app/)
**Port:** `:8080` | **API:** `/api/v1`

A modern REST API for task management demonstrating:
- ✅ Complete CRUD operations with route groups
- 🎯 Dynamic route parameters (`:id`)
- 🔍 Query parameters for filtering (`?status=completed`)
- 🌐 Global CORS middleware
- 📝 Data validation and error handling
- 🏗️ Organized route groups (`/api/v1/todos`)
- 🩺 Health check endpoint
- 📊 JSON response standardization

**Route Structure:**
```
/api/v1/
├── /todos/     → Complete todo management
├── /health     → Service health check
└── /           → API information
```

**How to run:**
```bash
cd examples/todo-app
go run main.go
```

### 2. [User Management](./user-management/)
**Port:** `:8081` | **API:** `/api/v1`

Advanced user management system with hierarchical access control:
- 🔐 Complete authentication system (register/login/logout)
- 🎫 JWT-like token-based authentication
- 🛡️ Cascading middleware (auth → admin)
- 👑 Role-based access control (user/admin)
- 🔒 Password hashing and security
- 📊 Admin statistics and monitoring
- 🏗️ Hierarchical route groups with inherited middleware
- 🩺 Health check and system monitoring

**Route Structure:**
```
/api/v1/
├── /auth/      → Public authentication
├── /profile/   → [AUTH] Profile management
├── /users/     → [AUTH] User operations
└── /admin/     → [ADMIN] Admin-only operations
    ├── /stats  → System statistics
    └── /users/ → Advanced user management
```

**How to run:**
```bash
cd examples/user-management
go run main.go
```

**Default users:**
- Admin: `admin/admin123`
- User: `john/password123`

### 3. [File Upload](./file-upload/)
**Port:** `:8082` | **API:** `/api/v1`

Complete file upload and management system demonstrating:
- 📤 Single and multiple file upload with validation
- 📋 File metadata management and listing
- 📥 Secure file download with proper headers
- 🌐 Static file serving with security checks
- ✅ File type and size validation (10MB limit)
- 🗑️ File deletion and cleanup
- 📊 Upload statistics and system information
- 🏗️ Organized route groups (`/files`, `/upload`)
- 🩺 Health check with storage validation

**Route Structure:**
```
/api/v1/
├── /files/     → File management operations
├── /upload/    → Upload operations (single/multiple)
├── /download/  → Secure file download
├── /health     → Health check with storage status
└── /info       → System statistics
/static/*       → Direct static file serving
```

**How to run:**
```bash
cd examples/file-upload
go run main.go
```

**Test interface:** Open `test.html` in your browser after starting the server.

## 🚀 Quick Start

To test all examples quickly:

```bash
# Terminal 1 - Todo App
cd examples/todo-app && go run main.go

# Terminal 2 - User Management
cd examples/user-management && go run main.go

# Terminal 3 - File Upload
cd examples/file-upload && go run main.go
```

**Access URLs:**
- Todo App: http://localhost:8080/api/v1
- User Management: http://localhost:8081/api/v1
- File Upload: http://localhost:8082/api/v1

## 📊 Advanced Framework Features Demonstrated

### 🏗️ Route Groups & Organization
```go
// Hierarchical route groups with shared middleware
v1Group := app.Group("/api/v1")
{
    authGroup := v1Group.Group("/auth")           // Public routes
    protectedGroup := v1Group.Group("/users", authMiddleware)  // Auth required
    adminGroup := v1Group.Group("/admin", authMiddleware, adminMiddleware) // Admin only
}
```

### 🔧 Middleware Stack Patterns
- **Global Middleware:** Applied to all routes (`app.Use()`)
- **Group Middleware:** Applied to route groups with inheritance
- **Route Middleware:** Applied to specific routes
- **Cascading Middleware:** Multiple middleware layers (auth → admin)

### 🎯 Modern API Patterns
- **API Versioning:** `/api/v1` structure for evolution
- **Resource Grouping:** Logical endpoint organization
- **Health Checks:** Production-ready monitoring endpoints
- **Error Standardization:** Consistent error response formats

### 🔄 Pattern Normalization
- Automatic route pattern cleanup (`//` → `/`)
- Consistent trailing slash handling
- Parameter validation and security

### Request/Response Handling
- JSON binding (`c.Bind()`) with validation
- Structured JSON responses (`c.JSON()`)
- File uploads (multipart/form-data) with security
- Custom headers and status codes
- Wildcard parameters (`*`) handling

### Security & Validation
- Input validation and sanitization
- Password hashing (SHA-256)
- Token-based authentication
- Role-based authorization with hierarchy
- File type and size validation
- Path traversal protection

## 🛠️ Structure of Each Example

Each example application contains:

```
example-name/
├── main.go                    # Main application with route groups
├── README.md                  # Comprehensive documentation
├── README.pt-BR.md           # Portuguese documentation
└── requests/
    └── test-sequence.http     # Complete API testing sequences
```

## 📚 Learning Path

### 1. **Start with Todo App** (Basic Concepts)
- Route groups and REST API patterns
- Global middleware application
- Basic CRUD with validation
- Health monitoring

### 2. **User Management** (Advanced Security)
- Hierarchical route groups
- Cascading middleware (auth → admin)
- Token authentication
- Role-based access control
- Protected endpoints

### 3. **File Upload** (File Operations)
- File handling and validation
- Multiple upload strategies
- Static file serving with security
- System monitoring and stats

### 📝 Testing the Examples

Each example includes comprehensive HTTP test files:
```bash
# Open in VS Code and execute requests sequentially
examples/todo-app/requests/test-sequence.http
examples/user-management/requests/test-sequence.http
examples/file-upload/requests/test-sequence.http
```

## 🔧 Requirements

- Go 1.21+ or higher
- Dependencies specified in `go.mod`:
  - `go.uber.org/zap` (structured logging)

## 📝 Development Notes

These examples were created for educational purposes and demonstrate:

- **Best practices** for web development in Go
- **REST patterns** for APIs
- **Structuring** web applications
- **Proper error handling**
- **Clear API documentation**

For production use, consider implementing:
- Persistent database
- Configuration via environment variables
- Structured logging
- Automated testing
- Docker containers
- CI/CD pipelines

## 🤝 Contributing

Want to add more examples? Consider creating:
- WebSocket chat
- GraphQL API
- Microservices
- Background jobs
- Template engine integration
- Database integration (PostgreSQL, MongoDB)

Each new example should follow the same structure and include complete documentation.
