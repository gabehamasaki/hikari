# Changelog

All notable changes to the Hikari Go framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
