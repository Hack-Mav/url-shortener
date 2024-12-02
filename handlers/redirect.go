package handlers

import (
	"context"
	"net/http"
	"url-shortener/config"
	"url-shortener/models"

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

			longURL = mapping.LongURL

			// Cache the result in Redis
			cache.Set(ctx, shortID, longURL, 0)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cache error"})
			return
		}

		// Redirect to the long URL
		c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
