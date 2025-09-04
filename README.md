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

## 🛠️ Requirements

- Go 1.24 or higher
- Dependencies:
  - `go.uber.org/zap` - Structured logging
  - `go.uber.org/multierr` - Error handling

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Hikari** - Fast, lightweight, and beautiful web framework for Go 🌅
