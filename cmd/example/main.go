package main

import (
	"net/http"

	"github.com/gabehamasaki/hikari-go/pkg/hikari"
)

func main() {
	app := hikari.New(":8080")

	// Recovery and Logger are now built-in and always enabled

	// Define your routes here
	app.GET("/hello/:name", func(c *hikari.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, " + c.Param("name") + "!",
			"status":  "success",
		})
	})

	app.ListenAndServe()
}
