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

	// Middleware for JSON responses (global)
	app.Use(func(next hikari.HandlerFunc) hikari.HandlerFunc {
		return func(c *hikari.Context) {
			c.SetHeader("Content-Type", "application/json")
			next(c)
		}
	})

	// Public routes
	app.GET("/", homePage)
	app.POST("/auth/register", register)
	app.POST("/auth/login", login)
	app.POST("/auth/logout", logout)

	// Protected routes (require authentication)
	app.GET("/users", getUsers, authMiddleware)
	app.GET("/users/:id", getUser, authMiddleware)
	app.PUT("/users/:id", updateUser, authMiddleware)
	app.DELETE("/users/:id", deleteUser, authMiddleware, adminMiddleware)
	app.GET("/profile", getProfile, authMiddleware)
	app.PUT("/profile", updateProfile, authMiddleware)

	// Admin only routes
	app.GET("/admin/users", adminGetUsers, authMiddleware, adminMiddleware)
	app.PATCH("/admin/users/:id/activate", activateUser, authMiddleware, adminMiddleware)
	app.PATCH("/admin/users/:id/deactivate", deactivateUser, authMiddleware, adminMiddleware)

	fmt.Println("🚀 User Management Server running on http://localhost:8081")
	app.ListenAndServe()
}

func homePage(c *hikari.Context) {
	c.JSON(http.StatusOK, hikari.H{
		"message": "User Management API",
		"version": "1.0.0",
		"endpoints": hikari.H{
			"POST /auth/register":               "Register new user",
			"POST /auth/login":                  "Login user",
			"POST /auth/logout":                 "Logout user",
			"GET /users":                        "List users (authenticated)",
			"GET /users/:id":                    "Get user by ID (authenticated)",
			"PUT /users/:id":                    "Update user (authenticated)",
			"DELETE /users/:id":                 "Delete user (admin only)",
			"GET /profile":                      "Get current user profile",
			"PUT /profile":                      "Update current user profile",
			"GET /admin/users":                  "Admin: List all users",
			"PATCH /admin/users/:id/activate":   "Admin: Activate user",
			"PATCH /admin/users/:id/deactivate": "Admin: Deactivate user",
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
