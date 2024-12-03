package main

import (
	"url-shortener/config"
	"url-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to GCP Datastore and Redis
	dsClient := config.ConnectDatastore()
	cache := config.ConnectCache()

	// Create Gin router
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Define routes
	r.POST("/shorten", handlers.ShortenURL(dsClient, cache))
	r.GET("/:short_id", handlers.RedirectURL(dsClient, cache))

	// Start the server
	r.Run(":8080") // Default is localhost:8080
}
