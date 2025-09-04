package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gabehamasaki/hikari-go/pkg/hikari"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never serialize password
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// In-memory storage
var users []User
var nextUserID = 1
var sessions = make(map[string]*User) // Simple session storage

func main() {
	app := hikari.New(":8081")

	// Initialize with admin user
	initializeUsers()

	// Global JSON middleware
	jsonMiddleware := func(next hikari.HandlerFunc) hikari.HandlerFunc {
		return func(c *hikari.Context) {
			c.SetHeader("Content-Type", "application/json")
			next(c)
		}
	}

	// Apply global middleware
	app.Use(jsonMiddleware)

	// Root welcome page
	app.GET("/", homePage)

	// API v1 group
	v1Group := app.Group("/api/v1")
	{
		// Auth group - no authentication required
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/register", register)
			authGroup.POST("/login", login)
			authGroup.POST("/logout", logout)
		}

		// User profile routes - require authentication
		profileGroup := v1Group.Group("/profile", authMiddleware)
		{
			profileGroup.GET("/", getProfile)
			profileGroup.PUT("/", updateProfile)
		}

		// User management routes - require authentication
		usersGroup := v1Group.Group("/users", authMiddleware)
		{
			usersGroup.GET("/", getUsers)
			usersGroup.GET("/:id", getUser)
			usersGroup.PUT("/:id", updateUser)
			usersGroup.DELETE("/:id", deleteUser, adminMiddleware) // Admin only
		}

		// Admin routes - require authentication and admin role
		adminGroup := v1Group.Group("/admin", authMiddleware, adminMiddleware)
		{
			adminUsersGroup := adminGroup.Group("/users")
			{
				adminUsersGroup.GET("/", adminGetUsers)
				adminUsersGroup.PATCH("/:id/activate", activateUser)
				adminUsersGroup.PATCH("/:id/deactivate", deactivateUser)
			}

			// Health and stats for admins
			adminGroup.GET("/stats", func(c *hikari.Context) {
				activeUsers := 0
				for _, user := range users { // users slice variable
					if user.Active {
						activeUsers++
					}
				}

				c.JSON(http.StatusOK, hikari.H{
					"total_users":    len(users), // users slice variable
					"active_users":   activeUsers,
					"inactive_users": len(users) - activeUsers, // users slice variable
					"sessions":       len(sessions),
				})
			})
		}

		// Health check
		v1Group.GET("/health", func(c *hikari.Context) {
			c.JSON(http.StatusOK, hikari.H{
				"status":    "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
				"service":   "user-management",
			})
		})
	}

	fmt.Println("🚀 User Management Server running on http://localhost:8081")
	fmt.Println("📋 API endpoints available at /api/v1")
	app.ListenAndServe()
}

func homePage(c *hikari.Context) {
	c.JSON(http.StatusOK, hikari.H{
		"message": "User Management API v1",
		"version": "1.0.0",
		"endpoints": hikari.H{
			"POST /api/v1/auth/register":               "Register new user",
			"POST /api/v1/auth/login":                  "Login user",
			"POST /api/v1/auth/logout":                 "Logout user",
			"GET /api/v1/users":                        "List users (authenticated)",
			"GET /api/v1/users/:id":                    "Get user by ID (authenticated)",
			"PUT /api/v1/users/:id":                    "Update user (authenticated)",
			"DELETE /api/v1/users/:id":                 "Delete user (admin only)",
			"GET /api/v1/profile":                      "Get current user profile",
			"PUT /api/v1/profile":                      "Update current user profile",
			"GET /api/v1/admin/users":                  "Admin: List all users",
			"PATCH /api/v1/admin/users/:id/activate":   "Admin: Activate user",
			"PATCH /api/v1/admin/users/:id/deactivate": "Admin: Deactivate user",
			"GET /api/v1/admin/stats":                  "Admin: Get system statistics",
			"GET /api/v1/health":                       "Health check",
		},
	})
}

func register(c *hikari.Context) {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Username, email and password are required",
		})
		return
	}

	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid email format",
		})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Password must be at least 6 characters",
		})
		return
	}

	// Check if user already exists
	for _, user := range users {
		if user.Username == req.Username || user.Email == req.Email {
			c.JSON(http.StatusConflict, hikari.H{
				"error": "Username or email already exists",
			})
			return
		}
	}

	// Create new user
	user := User{
		ID:        nextUserID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashPassword(req.Password),
		Role:      "user",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, user)
	nextUserID++

	c.JSON(http.StatusCreated, hikari.H{
		"message": "User created successfully",
		"user": hikari.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func login(c *hikari.Context) {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// Find user
	var foundUser *User
	for i, user := range users {
		if user.Username == req.Username && user.Password == hashPassword(req.Password) {
			if !user.Active {
				c.JSON(http.StatusForbidden, hikari.H{
					"error": "Account is deactivated",
				})
				return
			}
			foundUser = &users[i]
			break
		}
	}

	if foundUser == nil {
		c.JSON(http.StatusUnauthorized, hikari.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Create session token (simple implementation)
	token := generateToken(foundUser.Username)
	sessions[token] = foundUser

	c.JSON(http.StatusOK, hikari.H{
		"message": "Login successful",
		"token":   token,
		"user": hikari.H{
			"id":       foundUser.ID,
			"username": foundUser.Username,
			"email":    foundUser.Email,
			"role":     foundUser.Role,
		},
	})
}

func logout(c *hikari.Context) {
	token := getTokenFromRequest(c)
	if token != "" {
		delete(sessions, token)
	}

	c.JSON(http.StatusOK, hikari.H{
		"message": "Logout successful",
	})
}

func getUsers(c *hikari.Context) {
	var publicUsers []hikari.H
	for _, user := range users {
		if user.Active {
			publicUsers = append(publicUsers, hikari.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"created_at": user.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, hikari.H{
		"users": publicUsers,
		"count": len(publicUsers),
	})
}

func getUser(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid user ID",
		})
		return
	}

	for _, user := range users {
		if user.ID == id && user.Active {
			c.JSON(http.StatusOK, hikari.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "User not found",
	})
}

func updateUser(c *hikari.Context) {
	currentUser := getCurrentUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Users can only update themselves, unless they are admin
	if currentUser.Role != "admin" && currentUser.ID != id {
		c.JSON(http.StatusForbidden, hikari.H{
			"error": "You can only update your own profile",
		})
		return
	}

	var updateData struct {
		Email *string `json:"email"`
	}

	if err := c.Bind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			if updateData.Email != nil && isValidEmail(*updateData.Email) {
				users[i].Email = *updateData.Email
			}
			users[i].UpdatedAt = time.Now()

			c.JSON(http.StatusOK, hikari.H{
				"message": "User updated successfully",
				"user": hikari.H{
					"id":       users[i].ID,
					"username": users[i].Username,
					"email":    users[i].Email,
					"role":     users[i].Role,
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "User not found",
	})
}

func deleteUser(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid user ID",
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, hikari.H{
				"message": "User deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "User not found",
	})
}

func getProfile(c *hikari.Context) {
	user := getCurrentUser(c)
	c.JSON(http.StatusOK, hikari.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"active":     user.Active,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func updateProfile(c *hikari.Context) {
	user := getCurrentUser(c)

	var updateData struct {
		Email *string `json:"email"`
	}

	if err := c.Bind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	for i, u := range users {
		if u.ID == user.ID {
			if updateData.Email != nil && isValidEmail(*updateData.Email) {
				users[i].Email = *updateData.Email
			}
			users[i].UpdatedAt = time.Now()

			c.JSON(http.StatusOK, hikari.H{
				"message": "Profile updated successfully",
				"user": hikari.H{
					"id":       users[i].ID,
					"username": users[i].Username,
					"email":    users[i].Email,
					"role":     users[i].Role,
				},
			})
			return
		}
	}
}

func adminGetUsers(c *hikari.Context) {
	var allUsers []hikari.H
	for _, user := range users {
		allUsers = append(allUsers, hikari.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"active":     user.Active,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, hikari.H{
		"users": allUsers,
		"count": len(allUsers),
	})
}

func activateUser(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid user ID",
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users[i].Active = true
			users[i].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, hikari.H{
				"message": "User activated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "User not found",
	})
}

func deactivateUser(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid user ID",
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users[i].Active = false
			users[i].UpdatedAt = time.Now()
			c.JSON(http.StatusOK, hikari.H{
				"message": "User deactivated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "User not found",
	})
}

// Middleware functions
func authMiddleware(next hikari.HandlerFunc) hikari.HandlerFunc {
	return func(c *hikari.Context) {
		token := getTokenFromRequest(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, hikari.H{
				"error": "Authentication required",
			})
			return
		}

		user, exists := sessions[token]
		if !exists || !user.Active {
			c.JSON(http.StatusUnauthorized, hikari.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Store user in context (we'll simulate this with a simple approach)
		c.Request = c.Request.WithContext(
			c.Request.Context(),
		)

		next(c)
	}
}

func adminMiddleware(next hikari.HandlerFunc) hikari.HandlerFunc {
	return func(c *hikari.Context) {
		user := getCurrentUser(c)
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, hikari.H{
				"error": "Admin access required",
			})
			return
		}
		next(c)
	}
}

// Helper functions
func getCurrentUser(c *hikari.Context) *User {
	token := getTokenFromRequest(c)
	return sessions[token]
}

func getTokenFromRequest(c *hikari.Context) string {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateToken(username string) string {
	data := fmt.Sprintf("%s_%d", username, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func initializeUsers() {
	// Create admin user
	admin := User{
		ID:        1,
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  hashPassword("admin123"), // In production, use stronger passwords
		Role:      "admin",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create regular user
	user := User{
		ID:        2,
		Username:  "john",
		Email:     "john@example.com",
		Password:  hashPassword("password123"),
		Role:      "user",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = []User{admin, user}
	nextUserID = 3
}
