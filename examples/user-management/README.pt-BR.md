# User Management Example

Sistema completo de gerenciamento de usuários com autenticação e autorização usando Hikari-Go.

## Características

- Sistema de registro e login
- Autenticação baseada em token
- Autorização por roles (user/admin)
- Middleware de autenticação personalizado
- Endpoints protegidos
- Validação de dados
- Hash de senhas

## Como executar

```bash
cd examples/user-management
go run main.go
```

O servidor será iniciado em `http://localhost:8081`

## Usuários Padrão

### Admin
- **Username:** `admin`
- **Password:** `admin123`
- **Role:** `admin`

### User
- **Username:** `john`
- **Password:** `password123`
- **Role:** `user`

## Endpoints

### GET /
Informações sobre a API e lista de endpoints disponíveis.

### Autenticação

#### POST /auth/register
Registra um novo usuário.

**Body:**
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

#### POST /auth/login
Faz login de um usuário.

**Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "your-auth-token",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

#### POST /auth/logout
Faz logout do usuário atual.

**Headers:**
```
Authorization: Bearer your-auth-token
```

### Usuários (Requer Autenticação)

#### GET /users
Lista usuários ativos.

**Headers:**
```
Authorization: Bearer your-auth-token
```

#### GET /users/:id
Obtém informações de um usuário específico.

#### PUT /users/:id
Atualiza um usuário (usuários só podem atualizar a si mesmos, exceto admins).

**Body:**
```json
{
  "email": "newemail@example.com"
}
```

#### DELETE /users/:id (Admin Only)
Remove um usuário.

### Perfil

#### GET /profile
Obtém o perfil do usuário atual.

#### PUT /profile
Atualiza o perfil do usuário atual.

### Administração (Somente Admin)

#### GET /admin/users
Lista todos os usuários (incluindo inativos).

#### PATCH /admin/users/:id/activate
Ativa um usuário.

#### PATCH /admin/users/:id/deactivate
Desativa um usuário.

## Exemplos de Uso

### 1. Fazer Login como Admin
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. Listar Usuários (com token)
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/users
```

### 3. Registrar Novo Usuário
```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

### 4. Ver Perfil
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/profile
```

### 5. Listar Todos os Usuários (Admin)
```bash
curl -H "Authorization: Bearer ADMIN_TOKEN" \
  http://localhost:8081/admin/users
```

## Funcionalidades Demonstradas

- **Middleware Por Rota**: Middleware aplicado diretamente a rotas específicas
- **Authentication Middleware**: Verificação de tokens por rota
- **Authorization Middleware**: Controle por roles por rota
- **Password Hashing**: Uso do SHA-256 para senhas
- **Session Management**: Gerenciamento simples de sessões
- **Input Validation**: Validação de email e dados obrigatórios
- **Error Handling**: Tratamento de erros e respostas apropriadas
- **Route Protection**: Rotas protegidas por middleware de autenticação
- **Role-based Access**: Controle de acesso baseado em roles com middleware admin

## Estrutura de Segurança

- Senhas são hasheadas com SHA-256
- Tokens são gerados de forma única
- Middleware de autenticação verifica tokens
- Middleware de autorização verifica roles
- Usuários inativos não podem fazer login
- Validação de formato de email
- Senhas devem ter pelo menos 6 caracteres

## Notas de Implementação

Este é um exemplo educacional. Para produção, considere:

- Usar um sistema de hash mais seguro (bcrypt)
- Implementar JWT tokens
- Usar banco de dados real
- Adicionar rate limiting
- Implementar refresh tokens
- Adicionar logs de segurança
