# Hikari 🌅

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

**Hikari** (光 - "light" in Japanese) is a lightweight, fast, and elegant HTTP web framework for Go. It provides a minimalistic yet powerful foundation for building modern web applications and APIs.

## 📖 Documentation

**🌐 [hikari-go.dev](https://gabehamasaki.github.io/hikari-docs/)**

Visit our comprehensive documentation for detailed guides, examples, and API reference.

## ✨ Features

- 🚀 **Lightweight and Fast** - Minimal overhead with maximum performance
- 🛡️ **Built-in Recovery** - Automatic panic recovery to prevent crashes
- 📝 **Structured Logging** - Beautiful colored logs with Uber's Zap logger
- 🏗️ **Route Groups** - Organize routes with shared prefixes and middleware
- 🧩 **Middleware Support** - Extensible middleware system (global, group, and per-route)
- 🎯 **Context-based** - Rich context with JSON binding, query params, and storage
- 🛑 **Graceful Shutdown** - Proper server shutdown handling with signals

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

    // API v1 group
    v1Group := app.Group("/api/v1")
    {
        v1Group.GET("/hello/:name", func(c *hikari.Context) {
            c.JSON(http.StatusOK, hikari.H{
                "message": "Hello, " + c.Param("name") + "!",
                "status":  "success",
            })
        })

        // Health check
        v1Group.GET("/health", func(c *hikari.Context) {
            c.JSON(http.StatusOK, hikari.H{
                "status": "healthy",
                "service": "my-api",
            })
        })
    }

    app.ListenAndServe()
}
```

Run your application:
```bash
go run main.go
```

Visit `http://localhost:8080/api/v1/hello/world` to see your app in action!

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## �️ Requirements

- Go 1.24 or higher
- Dependencies:
  - `go.uber.org/zap` - Structured logging
  - `go.uber.org/multierr` - Error handling

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Hikari** - Fast, lightweight, and beautiful web framework for Go 🌅

### Creating an App

```go
app := hikari.New(":8080")

// Configure request timeout (default: 30 seconds)
app.SetRequestTimeout(60 * time.Second)
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

### Route Groups

Organize your routes with shared prefixes and middleware using groups. This feature enables clean API structure and hierarchical middleware application:

```go
// Basic route group
apiGroup := app.Group("/api")
{
    apiGroup.GET("/health", healthHandler)
    apiGroup.GET("/version", versionHandler)
}

// Versioned API group
v1Group := app.Group("/api/v1")
{
    // Users resource group
    usersGroup := v1Group.Group("/users")
    {
        usersGroup.GET("/", listUsers)
        usersGroup.POST("/", createUser)
        usersGroup.GET("/:id", getUser)
        usersGroup.PUT("/:id", updateUser)
        usersGroup.DELETE("/:id", deleteUser)
    }

    // Posts resource group
    postsGroup := v1Group.Group("/posts")
    {
        postsGroup.GET("/", listPosts)
        postsGroup.POST("/", createPost)
        postsGroup.GET("/:id", getPost)
    }
}
```

#### Groups with Middleware

Apply middleware to entire groups for shared authentication, logging, or other concerns:

```go
// Global middleware
app.Use(corsMiddleware)

// API v1 with rate limiting
v1Group := app.Group("/api/v1", rateLimitMiddleware)
{
    // Public endpoints (no additional middleware)
    v1Group.GET("/health", healthHandler)

    // Auth group - requires authentication
    authGroup := v1Group.Group("/auth")
    {
        authGroup.POST("/login", loginHandler)
        authGroup.POST("/register", registerHandler)
        authGroup.POST("/logout", logoutHandler, authMiddleware)
    }

    // Protected group - requires authentication
    protectedGroup := v1Group.Group("/protected", authMiddleware)
    {
        protectedGroup.GET("/profile", getProfile)
        protectedGroup.PUT("/profile", updateProfile)

        // Admin group - requires authentication + admin role
        adminGroup := protectedGroup.Group("/admin", adminMiddleware)
        {
            adminGroup.GET("/users", adminListUsers)
            adminGroup.DELETE("/users/:id", adminDeleteUser)
            adminGroup.GET("/stats", adminGetStats)
        }
    }
}
```

#### Nested Groups with Inherited Middleware

Groups automatically inherit middleware from their parent groups:

```go
// Parent group with auth middleware
apiGroup := app.Group("/api", authMiddleware)
{
    // Child group inherits auth + adds logging
    v1Group := apiGroup.Group("/v1", loggingMiddleware)
    {
        // Grandchild inherits auth + logging + adds admin check
        adminGroup := v1Group.Group("/admin", adminMiddleware)
        {
            // This endpoint has all 3 middlewares: auth -> logging -> admin
            adminGroup.GET("/users", getUsersHandler)
        }
    }
}
```

Results in route structure:
- `/api/v1/admin/users` → authMiddleware → loggingMiddleware → adminMiddleware → getUsersHandler

## 🛡️ Middleware System

Hikari provides a powerful and flexible middleware system that supports multiple levels of application:

### Global Middleware
Apply to all routes across your entire application:

```go
app := hikari.New(":8080")

// Global CORS middleware
app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
    return func(c *hikari.Context) {
        c.SetHeader("Access-Control-Allow-Origin", "*")
        c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
        next(c)
    }
})
```

### Group Middleware
Apply middleware to specific route groups:

```go
// API v1 with rate limiting
v1Group := app.Group("/api/v1", RateLimitMiddleware())

// Auth endpoints with security headers
authGroup := v1Group.Group("/auth", SecurityHeadersMiddleware())

// Admin routes with authentication + authorization
adminGroup := v1Group.Group("/admin", AuthMiddleware(), AdminMiddleware())
```

### Route-Specific Middleware
Apply middleware to individual routes:

```go
// Single middleware
app.POST("/users", createUser, AuthMiddleware())

// Multiple middleware (executed in order)
app.DELETE("/users/:id", deleteUser,
    AuthMiddleware(),
    AdminMiddleware(),
    AuditMiddleware())
```

### Middleware Execution Order

Middleware executes in a predictable order:
1. **Global middleware** (in registration order)
2. **Parent group middleware** (outer to inner groups)
3. **Current group middleware**
4. **Route-specific middleware** (in parameter order)
5. **Handler function**

```go
app.Use(GlobalMiddleware())                          // 1. Global

apiGroup := app.Group("/api", APIMiddleware())       // 2. API group
{
    v1Group := apiGroup.Group("/v1", V1Middleware()) // 3. V1 group (inherits API)
    {
        // 4. Route-specific middleware
        v1Group.POST("/users", createUser, AuthMiddleware(), ValidateMiddleware())
    }
}

// Execution order: Global → API → V1 → Auth → Validate → createUser
```

### Common Middleware Examples

#### Authentication Middleware
```go
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token == "" {
                c.JSON(401, hikari.H{"error": "Missing authorization token"})
                return
            }

            // Validate token
            if user, valid := validateJWT(token); valid {
                c.Set("user", user)
                c.Set("authenticated", true)
                next(c)
            } else {
                c.JSON(401, hikari.H{"error": "Invalid token"})
            }
        }
    }
}
```

#### Rate Limiting Middleware
```go
func RateLimitMiddleware(requestsPerMinute int) hikari.Middleware {
    limiter := make(map[string][]time.Time)
    mutex := sync.RWMutex{}

    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            ip := c.ClientIP()
            now := time.Now()

            mutex.Lock()
            defer mutex.Unlock()

            // Clean old requests
            requests := limiter[ip]
            var validRequests []time.Time
            for _, reqTime := range requests {
                if now.Sub(reqTime) < time.Minute {
                    validRequests = append(validRequests, reqTime)
                }
            }

            if len(validRequests) >= requestsPerMinute {
                c.JSON(429, hikari.H{
                    "error": "Rate limit exceeded",
                    "retry_after": 60,
                })
                return
            }

            limiter[ip] = append(validRequests, now)
            next(c)
        }
    }
}
```

#### Request ID Middleware
```go
func RequestIDMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            requestID := c.Request.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = generateUUID() // Your UUID generation function
            }

            c.Set("request_id", requestID)
            c.SetHeader("X-Request-ID", requestID)

            // Add to logger context
            c.Logger = c.Logger.With(zap.String("request_id", requestID))

            next(c)
        }
    }
}
```

#### Recovery Middleware
```go
func RecoveryMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            defer func() {
                if err := recover(); err != nil {
                    c.Logger.Error("Panic recovered",
                        zap.Any("error", err),
                        zap.String("path", c.Path()),
                        zap.String("method", c.Method()))

                    c.JSON(500, hikari.H{
                        "error": "Internal server error",
                        "request_id": c.GetString("request_id"),
                    })
                }
            }()
            next(c)
        }
    }
}
```

### Conditional Middleware

Apply middleware based on conditions:

```go
func ConditionalAuth() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            // Skip auth for health checks
            if c.Path() == "/health" || c.Path() == "/ping" {
                next(c)
                return
            }

            // Apply authentication for other routes
            AuthMiddleware()(next)(c)
        }
    }
}
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

#### Context Storage
```go
// Store values in context (thread-safe)
c.Set("user_id", 123)
c.Set("username", "john_doe")

// Retrieve values from context
userID, exists := c.Get("user_id")
if exists {
    // Use the value
}

// Retrieve with type assertion helpers
userID := c.GetInt("user_id")     // Returns 0 if not found or wrong type
username := c.GetString("username") // Returns "" if not found or wrong type
isActive := c.GetBool("is_active")  // Returns false if not found or wrong type

// Must get (returns nil and logs error if not found)
userID := c.MustGet("user_id")

// Get all stored keys
keys := c.Keys()
```

#### Context Interface (Go's context.Context)
```go
// Create context with timeout
ctx, cancel := c.WithTimeout(5 * time.Second)
defer cancel()

// Create context with cancellation
ctx, cancel := c.WithCancel()
defer cancel()

// Create context with value
ctx := c.WithValue("trace_id", "abc123")

// Access context values
traceID := c.Value("trace_id")

// Check if context is done or has error
select {
case <-c.Done():
    if err := c.Err(); err != nil {
        c.Logger.Error("Context cancelled", zap.Error(err))
        return
    }
default:
    // Continue processing
}
```

### Middleware

Create and use custom middleware - applicable globally, to route groups, or per specific route:

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

            // Store user info for later use
            c.Set("authenticated", true)
            c.Set("user_id", extractUserID(token))
            next(c)
        }
    }
}

// Rate limiting middleware
func RateLimitMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            // Rate limiting logic here
            c.Logger.Info("Rate limit check passed")
            next(c)
        }
    }
}
```

#### Middleware Application Levels

```go
// 1. Global middleware (applies to ALL routes)
app.Use(CORSMiddleware())
app.Use(loggingMiddleware)

// 2. Group middleware (applies to all routes in the group)
apiGroup := app.Group("/api", RateLimitMiddleware())
{
    // All routes here have rate limiting

    protectedGroup := apiGroup.Group("/protected", AuthMiddleware())
    {
        // All routes here have rate limiting + authentication

        adminGroup := protectedGroup.Group("/admin", AdminMiddleware())
        {
            // All routes here have: rate limiting + auth + admin check
            adminGroup.GET("/users", getUsersHandler)
        }
    }
}

// 3. Route-specific middleware (applies only to specific route)
app.GET("/public", publicHandler) // No middleware
app.GET("/auth-only", protectedHandler, AuthMiddleware()) // Only auth
app.POST("/admin-only", adminHandler, AuthMiddleware(), AdminMiddleware()) // Multiple middlewares
```

#### Middleware Execution Order

Middleware executes in the order they are applied:

```go
// Global middleware first
app.Use(middleware1) // Executes 1st
app.Use(middleware2) // Executes 2nd

// Then group middleware (outer to inner)
groupA := app.Group("/api", middleware3) // Executes 3rd
{
    groupB := groupA.Group("/v1", middleware4) // Executes 4th
    {
        // Finally route-specific middleware
        groupB.GET("/users", handler, middleware5) // Executes 5th, then handler
    }
}

// Execution order: middleware1 → middleware2 → middleware3 → middleware4 → middleware5 → handler
```

#### Middleware with Context Storage
You can use the context storage system in middleware to pass data between middlewares and handlers:

```go
// User extraction middleware
func UserMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token != "" {
                // Extract user from token (pseudo code)
                user := extractUserFromToken(token)
                c.Set("user", user)
                c.Set("user_id", user.ID)
                c.Set("is_authenticated", true)
            }
            next(c)
        }
    }
}

// Using stored values in handlers
app.GET("/profile", func(c *hikari.Context) {
    if !c.GetBool("is_authenticated") {
        c.JSON(http.StatusUnauthorized, hikari.H{"error": "Not authenticated"})
        return
    }

    user := c.MustGet("user")
    userID := c.GetInt("user_id")

    c.JSON(http.StatusOK, hikari.H{
        "user": user,
        "user_id": userID,
    })
}, UserMiddleware())
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

## 📝 Example: Complete RESTful API with Route Groups

```go
package main

import (
    "net/http"
    "strconv"
    "time"
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

// Middleware functions
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token != "Bearer valid-token" {
                c.JSON(http.StatusUnauthorized, hikari.H{
                    "error": "Invalid or missing token"})
                return
            }

            // Store auth info in context
            c.Set("authenticated", true)
            c.Set("user_role", "user")
            next(c)
        }
    }
}

func AdminMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            role := c.GetString("user_role")
            if role != "admin" {
                c.JSON(http.StatusForbidden, hikari.H{
                    "error": "Admin access required"})
                return
            }
            next(c)
        }
    }
}

func LoggingMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            start := time.Now()
            next(c)
            duration := time.Since(start)

            c.Logger.Info("Request processed",
                zap.String("method", c.Method()),
                zap.String("path", c.Path()),
                zap.Duration("duration", duration),
                zap.Int("status", c.GetStatus()))
        }
    }
}

func main() {
    app := hikari.New(":8080")

    // Configure request timeout
    app.SetRequestTimeout(60 * time.Second)

    // Global middleware
    app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Content-Type", "application/json")
            c.SetHeader("Access-Control-Allow-Origin", "*")
            next(c)
        }
    })

    // Root endpoint
    app.GET("/", func(c *hikari.Context) {
        c.JSON(http.StatusOK, hikari.H{
            "message": "Hikari API with Route Groups",
            "version": "1.0.0",
            "endpoints": hikari.H{
                "health":  "/health",
                "api_v1":  "/api/v1",
                "admin":   "/api/v1/admin",
            },
        })
    })

    // Health check (public)
    app.GET("/health", func(c *hikari.Context) {
        c.JSON(http.StatusOK, hikari.H{
            "status": "healthy",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    })

    // API v1 group with logging
    v1Group := app.Group("/api/v1", LoggingMiddleware())
    {
        // API info
        v1Group.GET("/", func(c *hikari.Context) {
            c.JSON(http.StatusOK, hikari.H{
                "message": "Welcome to API v1",
                "version": "1.0.0",
                "endpoints": hikari.H{
                    "users": "/api/v1/users",
                    "auth":  "/api/v1/auth",
                    "admin": "/api/v1/admin",
                },
            })
        })

        // Public users endpoints (read-only)
        usersGroup := v1Group.Group("/users")
        {
            usersGroup.GET("/", getUsers)
            usersGroup.GET("/:id", getUser)
        }

        // Auth endpoints
        authGroup := v1Group.Group("/auth")
        {
            authGroup.POST("/login", func(c *hikari.Context) {
                // Simplified login
                c.JSON(http.StatusOK, hikari.H{
                    "token": "Bearer valid-token",
                    "message": "Login successful",
                })
            })

            // Logout requires authentication
            authGroup.POST("/logout", func(c *hikari.Context) {
                c.JSON(http.StatusOK, hikari.H{
                    "message": "Logout successful",
                })
            }, AuthMiddleware())
        }

        // Protected endpoints (require authentication)
        protectedGroup := v1Group.Group("/protected", AuthMiddleware())
        {
            protectedGroup.GET("/profile", func(c *hikari.Context) {
                c.JSON(http.StatusOK, hikari.H{
                    "message": "Protected profile data",
                    "authenticated": c.GetBool("authenticated"),
                })
            })

            // User management (authenticated users can modify)
            userMgmtGroup := protectedGroup.Group("/users")
            {
                userMgmtGroup.POST("/", createUser)
                userMgmtGroup.PUT("/:id", updateUser)
            }
        }

        // Admin endpoints (require authentication + admin role)
        adminGroup := v1Group.Group("/admin", AuthMiddleware(), AdminMiddleware())
        {
            adminGroup.GET("/stats", func(c *hikari.Context) {
                c.JSON(http.StatusOK, hikari.H{
                    "total_users": len(users),
                    "admin_access": true,
                    "timestamp": time.Now(),
                })
            })

            // Admin user management
            adminUsersGroup := adminGroup.Group("/users")
            {
                adminUsersGroup.GET("/", func(c *hikari.Context) {
                    c.JSON(http.StatusOK, hikari.H{
                        "users": users,
                        "admin_view": true,
                    })
                })
                adminUsersGroup.DELETE("/:id", deleteUser)
            }
        }
    }

    // Static file serving
    app.GET("/static/*", func(c *hikari.Context) {
        filepath := c.Wildcard()
        c.Logger.Info("Serving static file", zap.String("file", filepath))
        c.File("./static/" + filepath)
    })

    app.ListenAndServe()
}

// Handler functions
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

This example demonstrates:

### 🏗️ Route Structure
```
/                           → Root API info
/health                     → Health check
/api/v1/                    → API v1 info
├── /users/                 → Public user operations
│   ├── GET /               → List users
│   └── GET /:id            → Get user
├── /auth/                  → Authentication
│   ├── POST /login         → Login
│   └── POST /logout        → Logout [AUTH]
├── /protected/             → Protected operations [AUTH]
│   ├── GET /profile        → User profile
│   └── /users/             → User management
│       ├── POST /          → Create user
│       └── PUT /:id        → Update user
└── /admin/                 → Admin operations [AUTH + ADMIN]
    ├── GET /stats          → System stats
    └── /users/             → Admin user management
        ├── GET /           → List all users (admin view)
        └── DELETE /:id     → Delete user
/static/*                   → Static file serving
```

### 🔧 Middleware Hierarchy
- **Global**: CORS, Content-Type
- **v1Group**: Logging middleware
- **protectedGroup**: Authentication middleware (inherits logging)
- **adminGroup**: Authentication + Admin middleware (inherits logging)

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

## � Documentation

Visit our comprehensive documentation website for detailed guides, examples, and API reference:

**🌐 [hikari-go.dev](https://gabehamasaki.github.io/hikari-docs/)**

The documentation includes:
- 🚀 **Getting Started Guide** - Quick setup and basic concepts
- 🛣️ **Routing & Middleware** - Advanced routing patterns and middleware system
- 🎯 **Context Management** - Request/response handling and data binding
- 📖 **Complete API Reference** - All methods, types, and interfaces
- 💼 **Practical Examples** - Real-world applications and use cases
- 🌍 **Multi-language Support** - Available in English and Portuguese

### Local Documentation Development

The documentation is built with Docusaurus and can be run locally:

```bash
cd docs
npm install
npm start
```

Visit `http://localhost:3000/hikari-go/` to view the documentation locally.

## �📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by popular web frameworks like Gin and Echo
- Built with ❤️ and Go
- Named after the Japanese word for "light" (光)

---

**Hikari** - Fast, lightweight, and beautiful web framework for Go 🌅
