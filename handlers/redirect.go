package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Hack-Mav/url-shortener/config"
	"github.com/Hack-Mav/url-shortener/models"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func RedirectURL(dsClient *datastore.Client, cache *config.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortID := c.Param("short_id")

		// Check Redis cache first
		longURL, err := cache.Get(config.Ctx, shortID).Result()
		if err == redis.Nil {
			// If not in cache, query Datastore
			ctx := context.Background()
			key := datastore.NameKey("URLMapping", shortID, nil)
			var mapping models.URLMapping
			err := dsClient.Get(ctx, key, &mapping)
			if err == datastore.ErrNoSuchEntity {
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Datastore error"})
				return
			}

			// Check if the URL has expired
			if time.Now().After(mapping.ExpiryDate) {
				// Remove the expired URL from the Datastore
				err := dsClient.Delete(ctx, key)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expired URL"})
					return
				}

				// Also remove from Redis cache if it exists
				cache.Del(ctx, shortID)

				c.JSON(http.StatusGone, gin.H{"error": "This URL has expired and has been removed"})
				return
			}

			// Cache the result in Redis
			cache.Set(ctx, shortID, mapping.LongURL, 0)
			longURL = mapping.LongURL
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cache error"})
			return
		}

		// Redirect to the long URL
		c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
