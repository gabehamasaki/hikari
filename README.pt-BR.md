# Hikari 🌅

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

**Hikari** (光 - "luz" em japonês) é um framework web HTTP leve, rápido e elegante para Go. Ele fornece uma base minimalista, mas poderosa, para construir aplicações web modernas e APIs com logging integrado, recuperação e capacidades de desligamento gracioso.

## ✨ Recursos

- 🚀 **Leve e Rápido** - Overhead mínimo com performance máxima
- 🛡️ **Recuperação Integrada** - Recuperação automática de pânico para evitar crashes
- 📝 **Logging Estruturado** - Logs coloridos bonitos com o logger Zap da Uber
- 🔗 **Parâmetros de Rota** - Suporte para parâmetros de rota dinâmicos (`:param`) e wildcards (`*`)
- 🧩 **Suporte a Middleware** - Sistema extensível de middleware global e por rota
- 🎯 **Baseado em Contexto** - Contexto rico com binding JSON, query params e mais
- 🛑 **Desligamento Gracioso** - Manipulação adequada de desligamento do servidor com sinais
- 📊 **Logging de Requisições** - Logging automático contextual com timing e User-Agent
- 📁 **Servidor de Arquivos** - Servir arquivos estáticos facilmente
- ⚙️ **Timeouts Configurados** - Timeouts de leitura e escrita pré-configurados (5s)

## 🚀 Início Rápido

### Instalação

```bash
go mod init seu-projeto
go get github.com/gabehamasaki/hikari-go
```

### Exemplo Básico

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
            "message": "Olá, " + c.Param("name") + "!",
            "status":  "success",
        })
    })

    app.ListenAndServe()
}
```

Execute sua aplicação:
```bash
go run main.go
```

Visite `http://localhost:8080/hello/world` para ver sua app em ação!

## 📚 Documentação

### Criando uma App

```go
app := hikari.New(":8080")
```

### Métodos HTTP

Hikari suporta todos os métodos HTTP padrão com middleware opcional por rota:

```go
// Sem middleware específico
app.GET("/users", getUsersHandler)
app.POST("/users", createUserHandler)

// Com middleware específico para a rota
app.PUT("/users/:id", updateUserHandler, authMiddleware, validationMiddleware)
app.PATCH("/users/:id", patchUserHandler, authMiddleware)
app.DELETE("/users/:id", deleteUserHandler, authMiddleware, adminMiddleware)
```

### Parâmetros de Rota

Extraia parâmetros de URLs usando a sintaxe `:param` e wildcards `*`:

```go
// Parâmetros simples
app.GET("/users/:id", func(c *hikari.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, map[string]string{"user_id": id})
})

// Múltiplos parâmetros
app.GET("/posts/:category/:id", func(c *hikari.Context) {
    category := c.Param("category")
    id := c.Param("id")
    c.JSON(http.StatusOK, map[string]string{
        "category": category,
        "post_id": id,
    })
})

// Wildcard - captura múltiplos segmentos
app.GET("/files/*", func(c *hikari.Context) {
    filepath := c.Wildcard() // Ex: "docs/api/v1/users.md"
    c.JSON(http.StatusOK, map[string]string{"file": filepath})
})

// Combinando parâmetros e wildcard
app.GET("/api/:version/*", func(c *hikari.Context) {
    version := c.Param("version")
    endpoint := c.Wildcard()
    c.JSON(http.StatusOK, map[string]string{
        "version": version,
        "endpoint": endpoint,
    })
})
```

### Métodos de Contexto

O `Context` fornece vários métodos para lidar com requisições e respostas:

#### Métodos de Resposta
```go
// Resposta JSON
c.JSON(http.StatusOK, map[string]interface{}{
    "message": "Sucesso",
    "data": data,
})

// Resposta de texto simples
c.String(http.StatusOK, "Olá, %s!", nome)

// Definir código de status
c.Status(http.StatusCreated)

// Servir arquivo estático
c.File("/path/to/file.pdf")

// Definir headers
c.SetHeader("X-Custom-Header", "value")

// Obter status atual da resposta
status := c.GetStatus()

// Obter header de resposta
contentType := c.GetHeader("Content-Type")
```

#### Métodos de Requisição
```go
// Obter parâmetro de rota
name := c.Param("name")

// Obter parâmetro wildcard
filepath := c.Wildcard()

// Obter parâmetro de query
page := c.Query("page")

// Obter valor de formulário
email := c.FormValue("email")

// Fazer bind do corpo da requisição JSON para struct
var user User
if err := c.Bind(&user); err != nil {
    c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
    return
}

// Obter método e path da requisição
method := c.Method()
path := c.Path()
```

### Middleware

Crie e use middleware personalizado - aplicável globalmente ou por rota específica:

```go
// Exemplo de middleware CORS
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

// Middleware de autenticação
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token == "" {
                c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token necessário"})
                return
            }
            next(c)
        }
    }
}

// Usar middleware globalmente (aplica a todas as rotas)
app.Use(CORSMiddleware())
app.Use(AuthMiddleware())

// Usar middleware específico por rota
app.GET("/public", publicHandler) // Sem middleware
app.GET("/protected", protectedHandler, AuthMiddleware()) // Só com auth
app.POST("/admin", adminHandler, AuthMiddleware(), AdminMiddleware()) // Múltiplos middlewares
```

### Recursos Integrados

Hikari vem com vários recursos integrados:

#### 🛡️ Middleware de Recuperação
Recupera automaticamente de pânicos e registra o erro:

```go
// Isso é integrado e sempre habilitado
// Não é necessário adicionar middleware de recuperação manualmente
```

#### 📝 Logging de Requisições
Logging estruturado contextual com informações detalhadas da requisição:

```
2024-09-04 15:04:05  INFO  Request started  {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "user_agent": "Mozilla/5.0..."}
2024-09-04 15:04:05  INFO  Request completed {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "user_agent": "Mozilla/5.0...", "status": 200, "duration": "2.5ms"}
```

O logger é automaticamente enriquecido com informações contextuais e está disponível em handlers:

```go
app.GET("/debug", func(c *hikari.Context) {
    c.Logger.Info("Processando requisição de debug",
        zap.String("user_id", userID))
    // ... lógica do handler
})
```

#### 🛑 Desligamento Gracioso
Manipula sinais de desligamento graciosamente:

```go
// Integrado - manipula SIGINT/SIGTERM automaticamente
app.ListenAndServe()
```

## 🏗️ Estrutura do Projeto

```
seu-projeto/
├── main.go
├── go.mod
├── go.sum
└── internal/
    └── handlers/
        ├── users.go
        └── posts.go
```

## 📝 Exemplo: API RESTful Completa

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
    {ID: 1, Name: "João Silva", Email: "joao@example.com"},
    {ID: 2, Name: "Maria Santos", Email: "maria@example.com"},
}

// Middleware de autenticação simples
func AuthMiddleware() hikari.Middleware {
    return func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            token := c.Request.Header.Get("Authorization")
            if token != "Bearer valid-token" {
                c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Token inválido ou ausente"})
                return
            }
            next(c)
        }
    }
}

func main() {
    app := hikari.New(":8080")

    // Middleware global
    app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Content-Type", "application/json")
            next(c)
        }
    })

    // Rotas públicas
    app.GET("/", func(c *hikari.Context) {
        c.JSON(http.StatusOK, map[string]string{
            "message": "API Hikari funcionando!",
            "version": "1.0.0",
        })
    })

    app.GET("/users", getUsers)
    app.GET("/users/:id", getUser)

    // Rotas protegidas (com middleware específico)
    app.POST("/users", createUser, AuthMiddleware())
    app.PUT("/users/:id", updateUser, AuthMiddleware())
    app.DELETE("/users/:id", deleteUser, AuthMiddleware())

    // Rota com wildcard para servir arquivos
    app.GET("/files/*", func(c *hikari.Context) {
        filepath := c.Wildcard()
        c.Logger.Info("Servindo arquivo", zap.String("file", filepath))
        c.File("./static/" + filepath)
    })

    // Rota para resposta de texto
    app.GET("/health", func(c *hikari.Context) {
        c.String(http.StatusOK, "OK - Servidor funcionando perfeitamente!")
    })

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
        c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de usuário inválido"})
        return
    }

    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, map[string]interface{}{"data": user})
            return
        }
    }

    c.JSON(http.StatusNotFound, map[string]string{"error": "Usuário não encontrado"})
}

func createUser(c *hikari.Context) {
    var newUser User
    if err := c.Bind(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
        return
    }

    newUser.ID = len(users) + 1
    users = append(users, newUser)

    c.Logger.Info("Novo usuário criado",
        zap.Int("user_id", newUser.ID),
        zap.String("user_name", newUser.Name))

    c.JSON(http.StatusCreated, map[string]interface{}{"data": newUser})
}

func updateUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de usuário inválido"})
        return
    }

    var updatedUser User
    if err := c.Bind(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            updatedUser.ID = id
            users[i] = updatedUser
            c.JSON(http.StatusOK, map[string]interface{}{"data": updatedUser})
            return
        }
    }

    c.JSON(http.StatusNotFound, map[string]string{"error": "Usuário não encontrado"})
}

func deleteUser(c *hikari.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, map[string]string{"error": "ID de usuário inválido"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(http.StatusOK, map[string]string{"message": "Usuário removido com sucesso"})
            return
        }
    }

    c.JSON(http.StatusNotFound, map[string]string{"error": "Usuário não encontrado"})
}
```

## 🛠️ Requisitos

- Go 1.24 ou superior
- Dependências:
  - `go.uber.org/zap` - Logging estruturado
  - `go.uber.org/multierr` - Tratamento de erros

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, sinta-se à vontade para enviar um Pull Request.

1. Faça fork do projeto
2. Crie sua branch de feature (`git checkout -b feature/recurso-incrivel`)
3. Commit suas mudanças (`git commit -m 'Adiciona algum recurso incrível'`)
4. Push para a branch (`git push origin feature/recurso-incrivel`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- Inspirado por frameworks web populares como Gin e Echo
- Construído com ❤️ e Go
- Nomeado a partir da palavra japonesa para "luz" (光)

---

**Hikari** - Framework web rápido, leve e bonito para Go 🌅
