# Changelog

All notable changes to the Hikari Go framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.5] - 2025-09-06

### 🚀 Added

#### WebSocket Support System
- **Full WebSocket Integration**: Native WebSocket support with Gorilla WebSocket library
- **WebSocket Manager**: Centralized management for WebSocket connections and hubs
- **WebSocket Hubs**: Multi-hub architecture for organizing connections by topic/room
- **WebSocket Context**: Extended Context specifically for WebSocket operations with familiar API
- **Connection Lifecycle**: Complete connection management with graceful shutdown and cleanup

#### WebSocket Features
- **Hub-based Architecture**: Organize connections into separate hubs/rooms
- **Broadcasting**: Send messages to all connections in a hub
- **Direct Messaging**: Send messages to specific connections by ID
- **Connection Pooling**: Automatic connection registration and cleanup
- **Ping/Pong Handling**: Built-in keepalive mechanism with configurable intervals
- **Message Types**: Support for text and binary message types

#### WebSocket Context & API
- **WebSocketContext**: Extends standard Context with WebSocket-specific methods
  - `Send()`: Send raw bytes to the connection
  - `JSON()`: Send JSON messages with automatic marshaling
  - `String()`: Send text messages
  - `Broadcast()`: Broadcast to all connections in the hub
  - `BroadcastJSON()`: Broadcast JSON messages
  - `SendToConnection()`: Send to specific connection by ID
  - `Bind()`: Bind JSON messages to structs
- **Message Type Helpers**: `IsTextMessage()`, `IsBinaryMessage()`, `GetMessage()`
- **Connection Info**: Access to connection ID, hub name, and connection status

#### WebSocket Configuration
- **Flexible Configuration**: Comprehensive configuration options
  - Read/Write buffer sizes
  - Handshake timeout
  - Origin checking function
  - Compression support
  - Ping/Pong intervals and timeouts
- **Default Configuration**: Sensible defaults with `DefaultWebSocketConfig()`

#### Request Context Management
- **Smart Timeout Handling**: WebSocket requests bypass request timeouts automatically
- **Context Preservation**: Original HTTP request context preserved for WebSocket connections
- **Middleware Compatibility**: WebSocket routes work seamlessly with existing middleware

### 🔧 Enhanced

#### Framework Core
- **App Structure**: Added `wsManager` field to main App struct
- **Route Registration**: New `WebSocket()` method for registering WebSocket endpoints
- **Hub Management**: `GetWebSocketHub()` method for accessing hubs from application code
- **Connection Detection**: Automatic WebSocket upgrade request detection
- **Graceful Shutdown**: WebSocket connections properly closed during application shutdown

#### Context System
- **Extended Context**: WebSocketContext extends the familiar Context API
- **Seamless Integration**: WebSocket handlers use the same pattern as HTTP handlers
- **Middleware Support**: Full middleware support for WebSocket routes
- **Storage Access**: Access to Context storage system for session data

#### Connection Management
- **Thread-Safe Operations**: All WebSocket operations are thread-safe
- **Resource Cleanup**: Automatic cleanup of connections and channels
- **Error Handling**: Comprehensive error handling with structured logging
- **Connection Tracking**: Unique connection IDs and connection counting

### 🐛 Fixed
- **Request Timeout**: WebSocket requests now properly bypass request timeouts
- **Context Cancellation**: Proper context handling for long-lived WebSocket connections
- **Memory Leaks**: Automatic cleanup prevents connection and goroutine leaks
- **Concurrent Access**: Thread-safe operations on shared connection pools

### 📚 Documentation
- **WebSocket Guide**: Complete guide on implementing WebSocket functionality
- **API Reference**: Full API documentation for WebSocket-specific methods
- **Configuration Examples**: Examples of different WebSocket configurations
- **Best Practices**: Guidelines for WebSocket hub organization and connection management

### 🏗️ Code Structure
```
pkg/hikari/
├── app.go              # Enhanced with WebSocket support
├── context.go          # Base context (unchanged)
├── ws-context.go       # ✨ NEW: WebSocket context implementation
├── websocket.go        # ✨ NEW: WebSocket core functionality
├── group.go            # Route groups (existing)
├── router.go           # Enhanced routing (existing)
└── ...

examples/
├── chat-app/          # ✨ NEW: Real-time chat application
├── websocket-echo/    # ✨ NEW: Simple echo WebSocket server
├── multi-room-chat/   # ✨ NEW: Multi-room chat with hubs
└── ...
```

### 💡 Usage Examples

#### Basic WebSocket Setup
```go
app := hikari.New(":8080")
app.WithWebSocket(hikari.DefaultWebSocketConfig())

app.WebSocket("/ws/chat", "chat_room", func(c *hikari.WebSocketContext) {
    if c.IsTextMessage() {
        message := c.GetMessage()
        c.Broadcast([]byte(message))
    }
})
```

#### Multi-Hub Chat Application
```go
// General chat
app.WebSocket("/ws/general", "general", generalChatHandler)

// VIP chat with auth middleware
app.WebSocket("/ws/vip", "vip", vipChatHandler, authMiddleware)

// Private messages
app.WebSocket("/ws/private", "private", privateChatHandler)
```

#### JSON Message Handling
```go
app.WebSocket("/ws/api", "api_hub", func(c *hikari.WebSocketContext) {
    if c.IsTextMessage() {
        var msg ChatMessage
        if err := c.Bind(&msg); err == nil {
            response := ProcessMessage(msg)
            c.JSON(response)
        }
    }
})
```

---

## [v0.1.4] - 2025-09-04

### 🚀 Added

#### Route Groups System
- **Route Groups**: Gin-like route grouping functionality with `app.Group()` method
- **Nested Groups**: Support for hierarchical route organization with unlimited nesting depth
- **Group Middleware**: Apply middleware to entire route groups with automatic inheritance
- **Pattern Prefixes**: Automatic prefix handling for grouped routes with proper normalization

#### Pattern Normalization & Validation
- **Pattern Helpers**: New helper functions for route pattern management
  - `normalizePattern()`: Removes duplicate slashes and ensures proper formatting
  - `isValidPattern()`: Validates route patterns using regex
  - `buildPattern()`: Combines prefixes with patterns safely
- **Regex Validation**: Built-in pattern validation using `duplicateSlashRegex` and `validPatternRegex`
- **Route Safety**: Automatic pattern cleanup to prevent routing conflicts

#### Enhanced Middleware System
- **Multiple Levels**: Support for global, group, and route-specific middleware
- **Middleware Inheritance**: Child groups automatically inherit parent middleware
- **Execution Order**: Predictable middleware execution chain (global → group → route → handler)
- **Conditional Middleware**: Support for conditional middleware application

### 🔧 Enhanced

#### Framework Core
- **Router Enhancement**: Updated `router.go` with pattern normalization and validation
- **Group Implementation**: New `group.go` file with complete Group struct and methods
- **HTTP Methods**: All standard HTTP methods (GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD) support groups
- **Context Handling**: Improved context management for grouped routes

#### Examples & Documentation
- **Updated Examples**: All example applications now use route groups structure
  - Todo App: `/api/v1/todos` with group-based organization
  - User Management: `/api/v1/users` and `/api/v1/auth` groups
  - File Upload: `/api/v1/files` with proper middleware inheritance
- **Professional APIs**: Examples follow REST API best practices with `/api/v1` versioning
- **Naming Conventions**: Resolved naming conflicts with "Group" suffix pattern

#### Documentation
- **Comprehensive README**: Completely rewritten documentation with:
  - Route Groups section with practical examples
  - Middleware system explanation with execution order
  - Complete RESTful API example using all new features
  - Professional API structure demonstrations
- **Example READMEs**: Updated all example documentation with new API structures
- **Code Examples**: Extensive code samples for all new functionality

### 🐛 Fixed
- **Naming Conflicts**: Resolved variable naming conflicts between group names and other variables
- **Pattern Validation**: Fixed route pattern validation issues with proper regex implementation
- **Middleware Order**: Ensured correct middleware execution order in nested groups

### 📚 Documentation
- **Route Groups Guide**: Complete guide on using route groups effectively
- **Middleware Patterns**: Common middleware examples and best practices
- **API Structure**: Professional API design patterns and conventions
- **Migration Guide**: Examples showing how to upgrade existing applications

### 🏗️ Code Structure
```
pkg/hikari/
├── app.go              # Main application struct
├── context.go          # Request context handling
├── group.go            # ✨ NEW: Route groups implementation
├── router.go           # Enhanced with pattern normalization
├── middleware.go       # Middleware system
└── ...

examples/
├── todo-app/          # Updated with route groups
├── user-management/   # Updated with route groups
├── file-upload/       # Updated with route groups
└── ...
```

### 💡 Usage Examples

#### Basic Route Groups
```go
v1Group := app.Group("/api/v1")
{
    v1Group.GET("/health", healthHandler)

    usersGroup := v1Group.Group("/users")
    {
        usersGroup.GET("/", getUsers)
        usersGroup.POST("/", createUser)
    }
}
```

#### Groups with Middleware
```go
authGroup := app.Group("/api/v1", AuthMiddleware())
{
    protectedGroup := authGroup.Group("/protected", RateLimitMiddleware())
    {
        protectedGroup.GET("/profile", getProfile)
    }
}
```

---

## [v0.1.3] - Previous Release
- Basic HTTP methods support
- Simple middleware system
- Context handling
- JSON/String responses
- File serving capabilities

---

**Full Changelog**: https://github.com/gabehamasaki/hikari-go/compare/v0.1.3...v0.1.4
