# Todo App Example

A modern REST API for task management using Hikari-Go with grouped routes and advanced features.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Features

- ✅ Complete CRUD for tasks with REST API
- 🔍 Filtering by status (completed/pending)
- 🌐 Global CORS middleware
- 📝 Data validation and error handling
- 🏗️ Organized route groups structure
- 📊 Health check endpoint
- 🔄 Backward compatibility
- 🎯 Clean JSON responses

## How to run

```bash
cd examples/todo-app
go run main.go
```

The server will start at `http://localhost:8080`

## API Structure

The API uses a versioned group structure for better organization:

```
/                    → API information (backward compatibility)
/api/v1/
├── /                → API v1 information
├── /health          → Health check
└── /todos/
    ├── GET /        → List todos
    ├── POST /       → Create todo
    ├── GET /:id     → Get specific todo
    ├── PUT /:id     → Update todo
    ├── DELETE /:id  → Delete todo
    └── PATCH /:id/toggle → Toggle completion status
```

## Endpoints

### GET /
Returns general API information with backward compatibility.

### GET /api/v1/
Returns API v1 information and available endpoints.

### GET /api/v1/health
Health check endpoint for monitoring.

**Example:**
```bash
curl http://localhost:8080/api/v1/health
```

### GET /api/v1/todos
Lists all tasks with optional filtering.

**Query Parameters:**
- `status`: `completed` or `pending` to filter tasks

**Examples:**
```bash
curl http://localhost:8080/api/v1/todos
curl http://localhost:8080/api/v1/todos?status=completed
curl http://localhost:8080/api/v1/todos?status=pending
```

### GET /api/v1/todos/:id
Returns a specific task by ID.

**Example:**
```bash
curl http://localhost:8080/api/v1/todos/1
```

### POST /api/v1/todos
Creates a new task.

**Body:**
```json
{
  "title": "New task",
  "content": "Task description"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Study Go","content":"Learn about goroutines"}'
```

### PUT /api/v1/todos/:id
Updates an existing task.

**Body:**
```json
{
  "title": "Updated title",
  "content": "Updated content",
  "completed": true
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Study Go - Advanced","completed":true}'
```

### DELETE /api/v1/todos/:id
Deletes a task.

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

### PATCH /api/v1/todos/:id/toggle
Toggles the completion status of a task.

**Example:**
```bash
curl -X PATCH http://localhost:8080/api/v1/todos/1/toggle
```

## Code Structure

### Route Groups
The application uses Hikari-Go's group feature for better organization:

```go
// API v1 routes group
v1Group := app.Group("/api/v1")
{
    // Home page
    v1Group.GET("/", homePage)

    // Todos group
    todosGroup := v1Group.Group("/todos")
    {
        todosGroup.GET("/", getTodos)
        todosGroup.POST("/", createTodo)
        todosGroup.GET("/:id", getTodo)
        todosGroup.PUT("/:id", updateTodo)
        todosGroup.DELETE("/:id", deleteTodo)
        todosGroup.PATCH("/:id/toggle", toggleTodo)
    }

    // Health check endpoint
    v1Group.GET("/health", healthCheck)
}
```

### Middlewares
- **CORS Middleware**: Handles cross-origin requests
- **JSON Middleware**: Sets JSON content-type for all responses

## Testing

Use the provided HTTP test file:
```
examples/todo-app/requests/test-sequence.http
```

Open this file in VS Code and execute the requests sequentially to test all functionality.

## Sample Data

The application starts with 3 sample todos to demonstrate the functionality:
1. "Learn Hikari-Go" (pending)
2. "Build Todo API" (completed)
3. "Add tests" (pending)
```json
{
  "title": "Updated title",
  "content": "Updated content",
  "completed": true
}
```

### DELETE /todos/:id
Removes a task.

**Example:**
```bash
curl -X DELETE http://localhost:8080/todos/1
```

### PATCH /todos/:id/toggle
Toggles the completion status of a task.

**Example:**
```bash
curl -X PATCH http://localhost:8080/todos/1/toggle
```

## Demonstrated Features

- **Routing**: Different HTTP methods and route parameters
- **JSON Binding**: Automatic JSON deserialization
- **Query Parameters**: Filtering using query strings
- **Custom Middleware**: Custom CORS middleware
- **Error Handling**: Validation and error handling
- **Response Formatting**: Structured JSON responses
