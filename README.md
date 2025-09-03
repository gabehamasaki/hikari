# Hikari 🌅

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

**Hikari** (光 - "light" in Japanese) is a lightweight, fast, and elegant HTTP web framework for Go. It provides a minimalistic yet powerful foundation for building modern web applications and APIs with built-in logging, recovery, and graceful shutdown capabilities.

## ✨ Features

- 🚀 **Lightweight and Fast** - Minimal overhead with maximum performance
- 🛡️ **Built-in Recovery** - Automatic panic recovery to prevent crashes
- 📝 **Structured Logging** - Beautiful colored logs with Uber's Zap logger
- 🔗 **Route Parameters** - Support for dynamic route parameters (`:param`)
- 🧩 **Middleware Support** - Extensible middleware system
- 🎯 **Context-based** - Rich context with JSON binding, query params, and more
- 🛑 **Graceful Shutdown** - Proper server shutdown handling with signals
- 📊 **Request Logging** - Automatic request/response logging with timing

## 🚀 Quick Start

### Installation

```bash
go mod init your-project
go get github.com/gabehamasaki/hikari-go
```

### Basic Example

```go
package main

import (
    "net/http"
    "github.com/gabehamasaki/hikari-go/pkg/hikari"
)

func main() {
    app := hikari.New(":8080")

    app.GET("/hello/:name", func(c *hikari.Context) {
        c.JSON(http.StatusOK, map[string]string{
            "message": "Hello, " + c.Param("name") + "!",
            "status":  "success",
        })
    })

    app.ListenAndServe()
}
```

Run your application:
```bash
go run main.go
```

Visit `http://localhost:8080/hello/world` to see your app in action!

## 📚 Documentation

### Creating an App

```go
app := hikari.New(":8080")
```

### HTTP Methods

Hikari supports all standard HTTP methods:

```go
app.GET("/users", getUsersHandler)
app.POST("/users", createUserHandler)
app.PUT("/users/:id", updateUserHandler)
app.PATCH("/users/:id", patchUserHandler)
app.DELETE("/users/:id", deleteUserHandler)
```

### Route Parameters

Extract parameters from URLs using the `:param` syntax:

```go
app.GET("/users/:id", func(c *hikari.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, map[string]string{"user_id": id})
})

app.GET("/posts/:category/:id", func(c *hikari.Context) {
    category := c.Param("category")
    id := c.Param("id")
    c.JSON(http.StatusOK, map[string]string{
        "category": category,
        "post_id": id,
    })
})
```

### Context Methods

The `Context` provides various methods to handle requests and responses:

#### Response Methods
```go
// JSON response
c.JSON(http.StatusOK, map[string]interface{}{
    "message": "Success",
    "data": data,
})

// Set status code
c.Status(http.StatusCreated)

// Set headers
c.SetHeader("X-Custom-Header", "value")
```

#### Request Methods
```go
// Get route parameter
name := c.Param("name")

// Get query parameter
page := c.Query("page")

// Get form value
email := c.FormValue("email")

// Bind JSON request body to struct
var user User
if err := c.Bind(&user); err != nil {
    c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
    return
}

// Get request method and path
method := c.Method()
path := c.Path()
```

### Middleware

Create and use custom middleware:

```go
// CORS middleware example
func CORSMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Access-Control-Allow-Origin", "*")
            c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

            if c.Method() == "OPTIONS" {
                c.Status(http.StatusOK)
                return
            }

            next(c)
        }
    }
}

// Use middleware
app.Use(CORSMiddleware())
```

### Built-in Features

Hikari comes with several built-in features:

#### 🛡️ Recovery Middleware
Automatically recovers from panics and logs the error:

```go
// This is built-in and always enabled
// No need to add recovery middleware manually
```

#### 📝 Request Logging
Beautiful structured logging with request details:

```
2024-09-03 15:04:05  INFO  Request started  {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321"}
2024-09-03 15:04:05  INFO  Request completed {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "status": 200, "duration": "2.5ms"}
```

#### 🛑 Graceful Shutdown
Handles shutdown signals gracefully:

```go
// Built-in - handles SIGINT/SIGTERM automatically
app.ListenAndServe()
```

## 🏗️ Project Structure

```
your-project/
├── main.go
├── go.mod
├── go.sum
└── internal/
    └── handlers/
        ├── users.go
        └── posts.go
```

## 📝 Example: RESTful API

```go
package main

import (
    "net/http"
    "strconv"
    "github.com/gabehamasaki/hikari-go/pkg/hikari"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "John Doe", Email: "john@example.com"},
    {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {
    app := hikari.New(":8080")

    // Middleware
    app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Content-Type", "application/json")
            next(c)
        }
    })

    // Routes
    app.GET("/users", getUsers)
    app.GET("/users/:id", getUser)
    app.POST("/users", createUser)

    app.ListenAndServe()
}

func getUsers(c *hikari.Context) {
    c.JSON(http.StatusOK, map[string]interface{}{
        "data": users,
        "count": len(users),
    })
}

func getUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
        return
    }

    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, map[string]interface{}{"data": user})
            return
        }
    }

    c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
}

func createUser(c *hikari.Context) {
    var newUser User
    if err := c.Bind(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
        return
    }

    newUser.ID = len(users) + 1
    users = append(users, newUser)

    c.JSON(http.StatusCreated, map[string]interface{}{"data": newUser})
}
```

## 🛠️ Requirements

- Go 1.21 or higher
- Dependencies:
  - `go.uber.org/zap` - Structured logging
  - `go.uber.org/multierr` - Error handling

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by popular web frameworks like Gin and Echo
- Built with ❤️ and Go
- Named after the Japanese word for "light" (光)

---

**Hikari** - Fast, lightweight, and beautiful web framework for Go 🌅
