package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gabehamasaki/hikari-go/pkg/hikari"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// In-memory storage for simplicity
var todos []Todo
var nextID = 1

func main() {
	app := hikari.New(":8080")

	// Add some sample todos
	initializeTodos()

	// Global CORS middleware
	corsMiddleware := func(next hikari.HandlerFunc) hikari.HandlerFunc {
		return func(c *hikari.Context) {
			c.SetHeader("Access-Control-Allow-Origin", "*")
			c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.SetHeader("Access-Control-Allow-Headers", "Content-Type")

			if c.Method() == "OPTIONS" {
				c.Status(http.StatusOK)
				return
			}

			next(c)
		}
	}

	// JSON Content-Type middleware
	jsonMiddleware := func(next hikari.HandlerFunc) hikari.HandlerFunc {
		return func(c *hikari.Context) {
			c.SetHeader("Content-Type", "application/json")
			next(c)
		}
	}

	// Apply global middlewares
	app.Use(corsMiddleware)
	app.Use(jsonMiddleware)

	// API v1 routes group
	v1Group := app.Group("/api/v1")
	{
		// Home page
		v1Group.GET("/", homePage)

		// Todos group
		todosGroup := v1Group.Group("/todos")
		{
			todosGroup.GET("/", getTodos)
			todosGroup.POST("/", createTodo)
			todosGroup.GET("/:id", getTodo)
			todosGroup.PUT("/:id", updateTodo)
			todosGroup.DELETE("/:id", deleteTodo)
			todosGroup.PATCH("/:id/toggle", toggleTodo)
		}

		// Health check endpoint
		v1Group.GET("/health", func(c *hikari.Context) {
			c.JSON(http.StatusOK, hikari.H{
				"status":    "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		})
	}

	// Root endpoints for backward compatibility
	app.GET("/", func(c *hikari.Context) {
		c.JSON(http.StatusOK, hikari.H{
			"message":      "Todo API - Use /api/v1 for latest version",
			"version":      "1.0.0",
			"api_endpoint": "/api/v1",
		})
	})

	app.ListenAndServe()
}

func homePage(c *hikari.Context) {
	c.JSON(http.StatusOK, hikari.H{
		"message": "Welcome to Todo API v1",
		"version": "1.0.0",
		"endpoints": hikari.H{
			"GET /api/v1/todos":              "List all todos",
			"GET /api/v1/todos/:id":          "Get todo by ID",
			"POST /api/v1/todos":             "Create new todo",
			"PUT /api/v1/todos/:id":          "Update todo",
			"DELETE /api/v1/todos/:id":       "Delete todo",
			"PATCH /api/v1/todos/:id/toggle": "Toggle todo completion",
			"GET /api/v1/health":             "Health check",
		},
	})
}

func getTodos(c *hikari.Context) {
	status := c.Query("status")
	var filteredTodos []Todo

	if status == "completed" {
		for _, todo := range todos {
			if todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	} else if status == "pending" {
		for _, todo := range todos {
			if !todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	} else {
		filteredTodos = todos
	}

	c.JSON(http.StatusOK, hikari.H{
		"todos": filteredTodos,
		"count": len(filteredTodos),
	})
}

func getTodo(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid todo ID",
		})
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "Todo not found",
	})
}

func createTodo(c *hikari.Context) {
	var newTodo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.Bind(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	if newTodo.Title == "" {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Title is required",
		})
		return
	}

	todo := Todo{
		ID:        nextID,
		Title:     newTodo.Title,
		Content:   newTodo.Content,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todos = append(todos, todo)
	nextID++

	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid todo ID",
		})
		return
	}

	var updateData struct {
		Title     *string `json:"title"`
		Content   *string `json:"content"`
		Completed *bool   `json:"completed"`
	}

	if err := c.Bind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid JSON data",
		})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			if updateData.Title != nil {
				todos[i].Title = *updateData.Title
			}
			if updateData.Content != nil {
				todos[i].Content = *updateData.Content
			}
			if updateData.Completed != nil {
				todos[i].Completed = *updateData.Completed
			}

			c.JSON(http.StatusOK, todos[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "Todo not found",
	})
}

func deleteTodo(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid todo ID",
		})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.JSON(http.StatusOK, hikari.H{
				"message": "Todo deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "Todo not found",
	})
}

func toggleTodo(c *hikari.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, hikari.H{
			"error": "Invalid todo ID",
		})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = !todos[i].Completed
			c.JSON(http.StatusOK, todos[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, hikari.H{
		"error": "Todo not found",
	})
}

func initializeTodos() {
	todos = []Todo{
		{
			ID:        1,
			Title:     "Learn Hikari-Go",
			Content:   "Study the framework documentation and create examples",
			Completed: false,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        2,
			Title:     "Build Todo API",
			Content:   "Create a REST API for managing todos",
			Completed: true,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        3,
			Title:     "Add tests",
			Content:   "Write unit tests for the API endpoints",
			Completed: false,
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
	}
	nextID = 4
}
