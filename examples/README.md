# Hikari-Go Examples

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

This folder contains practical examples demonstrating the features of the Hikari-Go framework. Each example is a complete and functional application that you can run and study.

## 📋 Examples List

### 1. [Todo App](./todo-app/)
**Port:** `:8080`

A complete REST API for task management demonstrating:
- Basic CRUD (Create, Read, Update, Delete)
- Dynamic route parameters
- Query parameters for filtering
- Custom middleware (CORS)
- Data validation
- Organized JSON structures

**How to run:**
```bash
cd examples/todo-app
go run main.go
```

### 2. [User Management](./user-management/)
**Port:** `:8081`

Complete user management system with authentication and authorization:
- Registration and login system
- Token-based authentication
- Authentication middleware
- Role-based access control (user/admin)
- Password hashing
- Data validation
- Protected endpoints

**How to run:**
```bash
cd examples/user-management
go run main.go
```

**Default users:**
- Admin: `admin/admin123`
- User: `john/password123`

### 3. [File Upload](./file-upload/)
**Port:** `:8082`

File upload and management system:
- Single file upload
- Multiple file upload
- File download
- Static file serving
- File type and size validation
- File listing and removal
- Health check

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
- Todo App: http://localhost:8080
- User Management: http://localhost:8081
- File Upload: http://localhost:8082

## 📊 Demonstrated Features

### Routing & HTTP Methods
- `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- Route parameters (`:id`, `:name`)
- Query parameters (`?status=completed`)
- Dynamic paths (`/static/*`)

### Middlewares
- Global middleware (`app.Use()`)
- Custom middleware (CORS)
- Authentication middleware
- Authorization middleware
- Middleware chaining

### Request/Response Handling
- JSON binding (`c.Bind()`)
- JSON responses (`c.JSON()`)
- Form data handling
- File uploads (multipart/form-data)
- Custom headers
- Custom status codes

### Validation & Security
- Input validation
- Password hashing
- Token authentication
- Role-based authorization
- File type validation
- Directory traversal prevention

### Error Handling
- HTTP error handling
- Structured error responses
- Built-in recovery middleware
- Contextual logging

## 🛠️ Structure of Each Example

Each example application contains:

```
example-name/
├── main.go          # Main application code
├── README.md        # Specific documentation
└── ...             # Additional files when needed
```

## 📚 How to Study the Examples

1. **Start with Todo App** - It's the simplest and shows basic concepts
2. **Move to User Management** - Adds authentication and authorization
3. **Finish with File Upload** - Demonstrates file manipulation and uploads

For each example:
1. Read the specific README
2. Examine the `main.go` code
3. Run the application
4. Test the endpoints with curl or web interface
5. Try modifying the code

## 🔧 Requirements

- Go 1.24.4 or higher
- Dependencies specified in `go.mod`:
  - `go.uber.org/zap` (logging)

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
