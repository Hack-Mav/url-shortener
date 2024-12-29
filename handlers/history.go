package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Hack-Mav/url-shortener/config"
	"github.com/Hack-Mav/url-shortener/models"
	"github.com/gin-gonic/gin"

	"cloud.google.com/go/datastore"
)

// FetchHistory retrieves all ShortURL entities from Google Cloud Datastore.
func FetchHistory(dsClient *datastore.Client, cache *config.RedisClient) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Create a background context (or use one passed in if you have it).
		ctx := context.Background()

		// Prepare your query. Make sure the kind name matches what you have in Datastore.
		// Example: "ShortURL" is the datastore kind name.
		query := datastore.NewQuery("ShortURL").Order("-createdAt")

		// This slice will hold the fetched results.
		var results []models.URLMapping

		// Execute the query and populate 'results'.
		_, err := dsClient.GetAll(ctx, query, &results)
		if err != nil {
			// Handle error.
			fmt.Println("Error fetching history:", err)
			// Return an error response to the client.
			// You can customize the error message as needed.
			// For example, you can return a 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
