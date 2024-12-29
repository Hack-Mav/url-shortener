package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/Hack-Mav/url-shortener/config"
	"github.com/Hack-Mav/url-shortener/models"

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

		creationTime, _ := time.Parse("2006-01-02", time.Now().String())

		// Parse the expiration date if provided
		var expiryDate time.Time
		if request.ExpiryDate != "" {
			parsedDate, err := time.Parse("2006-01-02", request.ExpiryDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration date format. Use YYYY-MM-DD."})
				return
			}
			expiryDate = parsedDate
		} else {
			// Default: Set expiration to 7 days from now
			expiryDate = time.Now().Add(7 * 24 * time.Hour)
			// expiryDate = time.Now().Add(1 * time.Minute)
		}

		// Generate short URL
		hash := sha256.Sum256([]byte(request.LongURL))
		shortID := base64.URLEncoding.EncodeToString(hash[:])[:6]

		// Store the mapping in Datastore
		ctx := context.Background()
		key := datastore.NameKey("URLMapping", shortID, nil)
		mapping := models.URLMapping{
			ShortID:    shortID,
			LongURL:    request.LongURL,
			ExpiryDate: expiryDate,
			CreatedAt:  creationTime,
		}
		_, err := dsClient.Put(ctx, key, &mapping)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL mapping"})
			return
		}

		// Cache the result in Redis
		cache.Set(ctx, shortID, request.LongURL, 0)

		c.JSON(http.StatusOK, gin.H{
			"short_url":   shortID,
			"expiry_date": expiryDate.Format("2006-01-02"),
		})
	}
}
