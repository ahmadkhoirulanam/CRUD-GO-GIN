package main

import (
	"framework/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router instance
	r := gin.Default()
	r.Static("/dist", "./dist")
	r.Static("/uploads", "./uploads")

	// Load routes from routes.go
	routes.LoadRoutes(r)

	// Start the server on port 8080
	r.Run(":8080")
}
