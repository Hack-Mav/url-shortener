package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"url-shortener/config"
	"url-shortener/models"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func ShortenURL(dsClient *datastore.Client, cache *config.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.ShortenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Generate short URL
		hash := sha256.Sum256([]byte(request.LongURL))
		shortID := base64.URLEncoding.EncodeToString(hash[:])[:6]

		// Check if short ID already exists in Datastore
		ctx := context.Background()
		key := datastore.NameKey("URLMapping", shortID, nil)
		var existingMapping models.URLMapping
		err := dsClient.Get(ctx, key, &existingMapping)
		if err == datastore.ErrNoSuchEntity {
			// Add new entry to Datastore
			mapping := models.URLMapping{ShortID: shortID, LongURL: request.LongURL}
			_, err := dsClient.Put(ctx, key, &mapping)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL mapping"})
				return
			}

			// Cache the result in Redis
			cache.Set(ctx, shortID, request.LongURL, 0)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Datastore error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"short_url": shortID})
	}
}
