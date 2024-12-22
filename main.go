package main

import (
	"url-shortener/config"
	"url-shortener/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to GCP Datastore and Redis
	dsClient := config.ConnectDatastore()
	cache := config.ConnectCache()

	// Create Gin router
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins, or specify a list of allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,                                                // Allow cookies to be sent
		MaxAge:           24 * 3600,                                           // Preflight cache duration in seconds
	}))

	// Define routes
	r.POST("/shorten", handlers.ShortenURL(dsClient, cache))
	r.GET("/:short_id", handlers.RedirectURL(dsClient, cache))

	// Start the server
	r.Run(":8080") // Default is localhost:8080
}
