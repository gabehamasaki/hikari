# Todo App Example

A simple REST API for task management using Hikari-Go.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Features

- Complete CRUD for tasks
- Filtering by status (completed/pending)
- Custom CORS middleware
- Data validation
- Organized JSON structure

## How to run

```bash
cd examples/todo-app
go run main.go
```

The server will start at `http://localhost:8080`

## Endpoints

### GET /
Returns information about the API and list of available endpoints.

### GET /todos
Lists all tasks.

**Query Parameters:**
- `status`: `completed` or `pending` to filter tasks

**Example:**
```bash
curl http://localhost:8080/todos
curl http://localhost:8080/todos?status=completed
```

### GET /todos/:id
Returns a specific task by ID.

**Example:**
```bash
curl http://localhost:8080/todos/1
```

### POST /todos
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
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Study Go","content":"Learn about goroutines"}'
```

### PUT /todos/:id
Updates an existing task.

**Body:**
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
