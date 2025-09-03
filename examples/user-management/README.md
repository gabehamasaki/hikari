# User Management Example

Complete user management system with authentication and authorization using Hikari-Go.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Features

- Registration and login system
- Token-based authentication
- Role-based authorization (user/admin)
- Custom authentication middleware
- Protected endpoints
- Data validation
- Password hashing

## How to run

```bash
cd examples/user-management
go run main.go
```

The server will start at `http://localhost:8081`

## Default Users

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
Information about the API and list of available endpoints.

### Authentication

#### POST /auth/register
Registers a new user.

**Body:**
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

#### POST /auth/login
Logs in a user.

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
Logs out the current user.

**Headers:**
```
Authorization: Bearer your-auth-token
```

### Users (Requires Authentication)

#### GET /users
Lists active users.

**Headers:**
```
Authorization: Bearer your-auth-token
```

#### GET /users/:id
Gets information about a specific user.

#### PUT /users/:id
Updates a user (users can only update themselves, except admins).

**Body:**
```json
{
  "email": "newemail@example.com"
}
```

#### DELETE /users/:id (Admin Only)
Removes a user.

### Profile

#### GET /profile
Gets the current user's profile.

#### PUT /profile
Updates the current user's profile.

### Administration (Admin Only)

#### GET /admin/users
Lists all users (including inactive ones).

#### PATCH /admin/users/:id/activate
Activates a user.

#### PATCH /admin/users/:id/deactivate
Deactivates a user.

## Usage Examples

### 1. Login as Admin
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. List Users (with token)
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/users
```

### 3. Register New User
```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

### 4. View Profile
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/profile
```

### 5. List All Users (Admin)
```bash
curl -H "Authorization: Bearer ADMIN_TOKEN" \
  http://localhost:8081/admin/users
```

## Demonstrated Features

- **Per-Route Middleware**: Middleware applied directly to specific routes
- **Authentication Middleware**: Token verification per route
- **Authorization Middleware**: Role-based control per route
- **Password Hashing**: Using SHA-256 for passwords
- **Session Management**: Simple session management
- **Input Validation**: Email and required data validation
- **Error Handling**: Error handling and appropriate responses
- **Route Protection**: Routes protected by authentication middleware
- **Role-based Access**: Role-based access control with admin middleware

## Security Structure

- Passwords are hashed with SHA-256
- Tokens are generated uniquely
- Authentication middleware verifies tokens
- Authorization middleware verifies roles
- Inactive users cannot login
- Email format validation
- Passwords must be at least 6 characters

## Implementation Notes

This is an educational example. For production, consider:

- Using a more secure hashing system (bcrypt)
- Implementing JWT tokens
- Using a real database
- Adding rate limiting
- Implementing refresh tokens
- Adding security logs
