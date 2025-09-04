# User Management Example

Complete user management system with authentication, authorization and advanced route grouping using Hikari-Go.

**Language / Idioma:** [English](README.md) | [Português Brasil](README.pt-BR.md)

## Features

- 🔐 Registration and login system
- 🎫 Token-based authentication
- 👑 Role-based authorization (user/admin)
- 🛡️ Custom authentication middleware
- 🔒 Protected endpoints with hierarchical access
- ✅ Data validation and sanitization
- 🔒 Password hashing with SHA-256
- 📊 Admin statistics and monitoring
- 🏗️ Organized route groups structure
- 🩺 Health check endpoint

## How to run

```bash
cd examples/user-management
go run main.go
```

The server will start at `http://localhost:8081`

## API Structure

The API uses a hierarchical group structure for organized access control:

```
/                    → API information
/api/v1/
├── /health          → Health check
├── /auth/           [PUBLIC]
│   ├── POST /register → Register user
│   ├── POST /login    → Login
│   └── POST /logout   → Logout
├── /profile/        [AUTH REQUIRED]
│   ├── GET /        → Get own profile
│   └── PUT /        → Update own profile
├── /users/          [AUTH REQUIRED]
│   ├── GET /        → List active users
│   ├── GET /:id     → Get user by ID
│   ├── PUT /:id     → Update user (own or admin)
│   └── DELETE /:id  → Delete user (admin only)
└── /admin/          [ADMIN REQUIRED]
    ├── GET /stats   → System statistics
    └── /users/
        ├── GET /    → List all users (including inactive)
        ├── PATCH /:id/activate   → Activate user
        └── PATCH /:id/deactivate → Deactivate user
```

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

### GET /api/v1/health
Health check endpoint for monitoring.

**Example:**
```bash
curl http://localhost:8081/api/v1/health
```

### Authentication

#### POST /api/v1/auth/register
Registers a new user.

**Body:**
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

**Example:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"secure123"}'
```

#### POST /api/v1/auth/login
Logs in a user and returns an authentication token.

**Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Example:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

#### POST /api/v1/auth/logout
Logs out the current user (invalidates token).

**Headers:**
```
Authorization: Bearer <token>
```

**Example:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/logout \
  -H "Authorization: Bearer <your-token>"
```

### Profile Management

#### GET /api/v1/profile
Returns the current user's profile information.

**Headers:**
```
Authorization: Bearer <token>
```

**Example:**
```bash
curl http://localhost:8081/api/v1/profile \
  -H "Authorization: Bearer <your-token>"
```

#### PUT /api/v1/profile
Updates the current user's profile.

**Headers:**
```
Authorization: Bearer <token>
```

**Body:**
```json
{
  "email": "newemail@example.com"
}
```

### User Management

#### GET /api/v1/users
Lists all active users (authentication required).

**Headers:**
```
Authorization: Bearer <token>
```

#### GET /api/v1/users/:id
Gets a specific user by ID (authentication required).

#### PUT /api/v1/users/:id
Updates a user. Users can only update themselves unless they are admin.

#### DELETE /api/v1/users/:id
Deletes a user (admin only).

### Admin Endpoints

#### GET /api/v1/admin/stats
Returns system statistics (admin only).

**Example Response:**
```json
{
  "total_users": 10,
  "active_users": 8,
  "inactive_users": 2,
  "sessions": 3
}
```

#### GET /api/v1/admin/users
Lists all users including inactive ones (admin only).

#### PATCH /api/v1/admin/users/:id/activate
Activates a user account (admin only).

#### PATCH /api/v1/admin/users/:id/deactivate
Deactivates a user account (admin only).

## Code Structure

### Route Groups with Middleware Hierarchy

The application demonstrates advanced route grouping with cascading middleware:

```go
// API v1 group
v1Group := app.Group("/api/v1")
{
    // Public auth group
    authGroup := v1Group.Group("/auth")
    {
        authGroup.POST("/register", register)
        authGroup.POST("/login", login)
        authGroup.POST("/logout", logout)
    }

    // Protected profile group
    profileGroup := v1Group.Group("/profile", authMiddleware)
    {
        profileGroup.GET("/", getProfile)
        profileGroup.PUT("/", updateProfile)
    }

    // Protected user management
    usersGroup := v1Group.Group("/users", authMiddleware)
    {
        usersGroup.GET("/", getUsers)
        usersGroup.DELETE("/:id", deleteUser, adminMiddleware) // Additional admin check
    }

    // Admin-only group (auth + admin middleware)
    adminGroup := v1Group.Group("/admin", authMiddleware, adminMiddleware)
    {
        adminGroup.GET("/stats", getStats)

        adminUsersGroup := adminGroup.Group("/users")
        {
            adminUsersGroup.PATCH("/:id/activate", activateUser)
            adminUsersGroup.PATCH("/:id/deactivate", deactivateUser)
        }
    }
}
```

### Middleware Stack
1. **Authentication Middleware**: Validates JWT tokens
2. **Admin Middleware**: Checks for admin role (applied after auth)
3. **JSON Middleware**: Sets content-type for responses

## Security Features

- Password hashing using SHA-256
- Token-based authentication
- Role-based access control
- Input validation and sanitization
- Protected admin endpoints

## Testing

Use the provided HTTP test file:
```
examples/user-management/requests/test-sequence.http
```

The test sequence includes:
1. User registration
2. Login scenarios (admin, user, new user)
3. Profile management
4. User management operations
5. Admin-only operations
6. Access control testing

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
