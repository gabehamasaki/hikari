# Hikari 🌅

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

**Hikari** (光 - "light" in Japanese) is a lightweight, fast, and elegant HTTP web framework for Go. It provides a minimalistic yet powerful foundation for building modern web applications and APIs with built-in logging, recovery, and graceful shutdown capabilities.

## ✨ Features

- 🚀 **Lightweight and Fast** - Minimal overhead with maximum performance
- 🛡️ **Built-in Recovery** - Automatic panic recovery to prevent crashes
- 📝 **Structured Logging** - Beautiful colored logs with Uber's Zap logger
- 🔗 **Route Parameters** - Support for dynamic route parameters (`:param`) and wildcards (`*`)
- 🧩 **Middleware Support** - Extensible middleware system (global and per-route)
- 🎯 **Context-based** - Rich context with JSON binding, query params, and more
- 🛑 **Graceful Shutdown** - Proper server shutdown handling with signals
- 📊 **Request Logging** - Automatic contextual logging with timing and User-Agent
- 📁 **File Server** - Serve static files easily
- ⚙️ **Configured Timeouts** - Pre-configured read/write timeouts (5s)

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
        c.JSON(http.StatusOK, hikari.H{
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

Hikari supports all standard HTTP methods with optional per-route middleware:

```go
// Without specific middleware
app.GET("/users", getUsersHandler)
app.POST("/users", createUserHandler)

// With route-specific middleware
app.PUT("/users/:id", updateUserHandler, authMiddleware, validationMiddleware)
app.PATCH("/users/:id", patchUserHandler, authMiddleware)
app.DELETE("/users/:id", deleteUserHandler, authMiddleware, adminMiddleware)
```

### Route Parameters

Extract parameters from URLs using the `:param` syntax and wildcards `*`:

```go
// Simple parameters
app.GET("/users/:id", func(c *hikari.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, hikari.H{"user_id": id})
})

// Multiple parameters
app.GET("/posts/:category/:id", func(c *hikari.Context) {
    category := c.Param("category")
    id := c.Param("id")
    c.JSON(http.StatusOK, hikari.H{
        "category": category,
        "post_id": id,
    })
})

// Wildcard - captures multiple path segments
app.GET("/files/*", func(c *hikari.Context) {
    filepath := c.Wildcard() // Ex: "docs/api/v1/users.md"
    c.JSON(http.StatusOK, hikari.H{"file": filepath})
})

// Combining parameters and wildcard
app.GET("/api/:version/*", func(c *hikari.Context) {
    version := c.Param("version")
    endpoint := c.Wildcard()
    c.JSON(http.StatusOK, hikari.H{
        "version": version,
        "endpoint": endpoint,
    })
})
```

### Context Methods

The `Context` provides various methods to handle requests and responses:

### `hikari.H` Alias

For easier JSON response creation, Hikari provides the `hikari.H` alias:

```go
// Instead of using map[string]any or map[string]interface{}
c.JSON(http.StatusOK, map[string]interface{}{
    "message": "success",
    "data": userData,
})

// Use the cleaner hikari.H alias
c.JSON(http.StatusOK, hikari.H{
    "message": "success",
    "data": userData,
})
```

#### Response Methods
```go
// JSON response
c.JSON(http.StatusOK, hikari.H{
    "message": "Success",
    "data": data,
})

// Plain text response
c.String(http.StatusOK, "Hello, %s!", name)

// Set status code
c.Status(http.StatusCreated)

// Serve static file
c.File("/path/to/file.pdf")

// Set headers
c.SetHeader("X-Custom-Header", "value")

// Get current response status
status := c.GetStatus()

// Get response header
contentType := c.GetHeader("Content-Type")
```

#### Request Methods
```go
// Get route parameter
name := c.Param("name")

// Get wildcard parameter
filepath := c.Wildcard()

// Get query parameter
page := c.Query("page")

// Get form value
email := c.FormValue("email")

// Bind JSON request body to struct
var user User
if err := c.Bind(&user); err != nil {
    c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid JSON"})
    return
}

// Get request method and path
method := c.Method()
path := c.Path()
```

### Middleware

Create and use custom middleware - applicable globally or per specific route:

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

// Authentication middleware
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token == "" {
                c.JSON(http.StatusUnauthorized, hikari.H{"error": "Token required"})
                return
            }
            next(c)
        }
    }
}

// Use middleware globally (applies to all routes)
app.Use(CORSMiddleware())
app.Use(AuthMiddleware())

// Use route-specific middleware
app.GET("/public", publicHandler) // No middleware
app.GET("/protected", protectedHandler, AuthMiddleware()) // Only auth
app.POST("/admin", adminHandler, AuthMiddleware(), AdminMiddleware()) // Multiple middlewares
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
Beautiful contextual structured logging with detailed request information:

```
2024-09-04 15:04:05  INFO  Request started  {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "user_agent": "Mozilla/5.0..."}
2024-09-04 15:04:05  INFO  Request completed {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "user_agent": "Mozilla/5.0...", "status": 200, "duration": "2.5ms"}
```

The logger is automatically enriched with contextual information and available in handlers:

```go
app.GET("/debug", func(c *hikari.Context) {
    c.Logger.Info("Processing debug request",
        zap.String("user_id", userID))
    // ... handler logic
})
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

## 📝 Example: Complete RESTful API

```go
package main

import (
    "net/http"
    "strconv"
    "github.com/gabehamasaki/hikari-go/pkg/hikari"
    "go.uber.org/zap"
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

// Simple authentication middleware
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token != "Bearer valid-token" {
                c.JSON(http.StatusUnauthorized, hikari.H{
                    "error": "Invalid or missing token"})
                return
            }
            next(c)
        }
    }
}

func main() {
    app := hikari.New(":8080")

    // Global middleware
    app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Content-Type", "application/json")
            next(c)
        }
    })

    // Public routes
    app.GET("/", func(c *hikari.Context) {
        c.JSON(http.StatusOK, hikari.H{
            "message": "Hikari API is running!",
            "version": "1.0.0",
        })
    })

    app.GET("/users", getUsers)
    app.GET("/users/:id", getUser)

    // Protected routes (with specific middleware)
    app.POST("/users", createUser, AuthMiddleware())
    app.PUT("/users/:id", updateUser, AuthMiddleware())
    app.DELETE("/users/:id", deleteUser, AuthMiddleware())

    // Wildcard route for serving files
    app.GET("/files/*", func(c *hikari.Context) {
        filepath := c.Wildcard()
        c.Logger.Info("Serving file", zap.String("file", filepath))
        c.File("./static/" + filepath)
    })

    // Text response route
    app.GET("/health", func(c *hikari.Context) {
        c.String(http.StatusOK, "OK - Server is running perfectly!")
    })

    app.ListenAndServe()
}

func getUsers(c *hikari.Context) {
    c.JSON(http.StatusOK, hikari.H{
        "data": users,
        "count": len(users),
    })
}

func getUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid user ID"})
        return
    }

    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, hikari.H{"data": user})
            return
        }
    }

    c.JSON(http.StatusNotFound, hikari.H{"error": "User not found"})
}

func createUser(c *hikari.Context) {
    var newUser User
    if err := c.Bind(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid JSON"})
        return
    }

    newUser.ID = len(users) + 1
    users = append(users, newUser)

    c.Logger.Info("New user created",
        zap.Int("user_id", newUser.ID),
        zap.String("user_name", newUser.Name))

    c.JSON(http.StatusCreated, hikari.H{"data": newUser})
}

func updateUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid user ID"})
        return
    }

    var updatedUser User
    if err := c.Bind(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid JSON"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            updatedUser.ID = id
            users[i] = updatedUser
            c.JSON(http.StatusOK, hikari.H{"data": updatedUser})
            return
        }
    }

    c.JSON(http.StatusNotFound, hikari.H{"error": "User not found"})
}

func deleteUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, hikari.H{"error": "Invalid user ID"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(http.StatusOK, hikari.H{"message": "User deleted successfully"})
            return
        }
    }

    c.JSON(http.StatusNotFound, hikari.H{"error": "User not found"})
}
```

## 🛠️ Requirements

- Go 1.24 or higher
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
