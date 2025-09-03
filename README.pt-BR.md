# Hikari 🌅

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

**Hikari** (光 - "luz" em japonês) é um framework web HTTP leve, rápido e elegante para Go. Ele fornece uma base minimalista, mas poderosa, para construir aplicações web modernas e APIs com logging integrado, recuperação e capacidades de desligamento gracioso.

## ✨ Recursos

- 🚀 **Leve e Rápido** - Overhead mínimo com performance máxima
- 🛡️ **Recuperação Integrada** - Recuperação automática de pânico para evitar crashes
- 📝 **Logging Estruturado** - Logs coloridos bonitos com o logger Zap da Uber
- 🔗 **Parâmetros de Rota** - Suporte para parâmetros de rota dinâmicos (`:param`)
- 🧩 **Suporte a Middleware** - Sistema extensível de middleware
- 🎯 **Baseado em Contexto** - Contexto rico com binding JSON, query params e mais
- 🛑 **Desligamento Gracioso** - Manipulação adequada de desligamento do servidor com sinais
- 📊 **Logging de Requisições** - Logging automático de requisição/resposta com timing

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
    "github.com/gabehamasaki/hikari-go/internal/hikari"
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

Hikari suporta todos os métodos HTTP padrão:

```go
app.GET("/users", getUsersHandler)
app.POST("/users", createUserHandler)
app.PUT("/users/:id", updateUserHandler)
app.PATCH("/users/:id", patchUserHandler)
app.DELETE("/users/:id", deleteUserHandler)
```

### Parâmetros de Rota

Extraia parâmetros de URLs usando a sintaxe `:param`:

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

### Métodos de Contexto

O `Context` fornece vários métodos para lidar com requisições e respostas:

#### Métodos de Resposta
```go
// Resposta JSON
c.JSON(http.StatusOK, map[string]interface{}{
    "message": "Sucesso",
    "data": data,
})

// Definir código de status
c.Status(http.StatusCreated)

// Definir headers
c.SetHeader("X-Custom-Header", "value")
```

#### Métodos de Requisição
```go
// Obter parâmetro de rota
name := c.Param("name")

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

Crie e use middleware personalizado:

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

// Usar middleware
app.Use(CORSMiddleware())
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
Belo logging estruturado com detalhes da requisição:

```
2024-09-03 15:04:05  INFO  Request started  {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321"}
2024-09-03 15:04:05  INFO  Request completed {"method": "GET", "path": "/users/123", "remote_addr": "127.0.0.1:54321", "status": 200, "duration": "2.5ms"}
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

## 📝 Exemplo: API RESTful

```go
package main

import (
    "net/http"
    "strconv"
    "github.com/gabehamasaki/hikari-go/internal/hikari"
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

func main() {
    app := hikari.New(":8080")

    // Middleware
    app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
        return func(c *hikari.Context) {
            c.SetHeader("Content-Type", "application/json")
            next(c)
        }
    })

    // Rotas
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

    c.JSON(http.StatusCreated, map[string]interface{}{"data": newUser})
}
```

## 🛠️ Requisitos

- Go 1.21 ou superior
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
