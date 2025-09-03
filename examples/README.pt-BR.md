# Exemplos do Hikari-Go

Esta pasta contém exemplos práticos demonstrando as funcionalidades do framework Hikari-Go. Cada exemplo é uma aplicação completa e funcional que você pode executar e estudar.

## 📋 Lista de Exemplos

### 1. [Todo App](./todo-app/)
**Porta:** `:8080`

Uma API REST completa para gerenciamento de tarefas demonstrando:
- CRUD básico (Create, Read, Update, Delete)
- Parâmetros de rota dinâmicos
- Query parameters para filtragem
- Middleware personalizado (CORS)
- Validação de dados
- Estruturas JSON organizadas

**Como executar:**
```bash
cd examples/todo-app
go run main.go
```

### 2. [User Management](./user-management/)
**Porta:** `:8081`

Sistema completo de gerenciamento de usuários com autenticação e autorização:
- Sistema de registro e login
- Autenticação baseada em tokens
- Middleware de autenticação
- Controle de acesso por roles (user/admin)
- Hash de senhas
- Validação de dados
- Endpoints protegidos

**Como executar:**
```bash
cd examples/user-management
go run main.go
```

**Usuários padrão:**
- Admin: `admin/admin123`
- User: `john/password123`

### 3. [File Upload](./file-upload/)
**Porta:** `:8082`

Sistema de upload e gerenciamento de arquivos:
- Upload de arquivo único
- Upload de múltiplos arquivos
- Download de arquivos
- Servir arquivos estáticos
- Validação de tipo e tamanho
- Listagem e remoção de arquivos
- Health check

**Como executar:**
```bash
cd examples/file-upload
go run main.go
```

**Interface de teste:** Abra `test.html` em seu navegador após iniciar o servidor.

## 🚀 Execução Rápida

Para testar todos os exemplos rapidamente:

```bash
# Terminal 1 - Todo App
cd examples/todo-app && go run main.go

# Terminal 2 - User Management
cd examples/user-management && go run main.go

# Terminal 3 - File Upload
cd examples/file-upload && go run main.go
```

**URLs de acesso:**
- Todo App: http://localhost:8080
- User Management: http://localhost:8081
- File Upload: http://localhost:8082

## 📊 Funcionalidades Demonstradas

### Routing & HTTP Methods
- `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- Parâmetros de rota (`:id`, `:name`)
- Query parameters (`?status=completed`)
- Paths dinâmicos (`/static/*`)

### Middlewares
- Middleware global (`app.Use()`)
- Middleware personalizado (CORS)
- Middleware de autenticação
- Middleware de autorização
- Chaining de middlewares

### Request/Response Handling
- JSON binding (`c.Bind()`)
- JSON responses (`c.JSON()`)
- Form data handling
- File uploads (multipart/form-data)
- Headers customizados
- Status codes personalizados

### Validação & Segurança
- Validação de entrada
- Hash de senhas
- Autenticação por tokens
- Autorização por roles
- Validação de tipos de arquivo
- Prevenção de directory traversal

### Error Handling
- Tratamento de erros HTTP
- Respostas de erro estruturadas
- Recovery middleware integrado
- Logging contextual

## 🛠️ Estrutura de Cada Exemplo

Cada aplicação de exemplo contém:

```
example-name/
├── main.go          # Código principal da aplicação
├── README.md        # Documentação específica
└── ...             # Arquivos adicionais quando necessário
```

## 📚 Como Estudar os Exemplos

1. **Comece pelo Todo App** - É o mais simples e mostra os conceitos básicos
2. **Prossiga para User Management** - Adiciona autenticação e autorização
3. **Termine com File Upload** - Demonstra manipulação de arquivos e uploads

Para cada exemplo:
1. Leia o README específico
2. Examine o código `main.go`
3. Execute a aplicação
4. Teste os endpoints com curl ou interface web
5. Experimente modificar o código

## 🔧 Requisitos

- Go 1.24.4 ou superior
- Dependências especificadas em `go.mod`:
  - `go.uber.org/zap` (logging)

## 📝 Notas de Desenvolvimento

Estes exemplos foram criados para fins educacionais e demonstram:

- **Boas práticas** de desenvolvimento web em Go
- **Padrões REST** para APIs
- **Estruturação** de aplicações web
- **Tratamento de erros** adequado
- **Documentação** clara de APIs

Para uso em produção, considere implementar:
- Banco de dados persistente
- Configuração via variáveis de ambiente
- Logs estruturados
- Testes automatizados
- Docker containers
- CI/CD pipelines

## 🤝 Contribuindo

Quer adicionar mais exemplos? Considere criar:
- WebSocket chat
- GraphQL API
- Microserviços
- Background jobs
- Template engine integration
- Database integration (PostgreSQL, MongoDB)

Cada novo exemplo deve seguir a mesma estrutura e incluir documentação completa.
