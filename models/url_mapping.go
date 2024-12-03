package models

import "time"

type URLMapping struct {
	ShortID    string    `datastore:"short_id"`
	LongURL    string    `datastore:"long_url,noindex"`
	ExpiryDate time.Time `datastore:"expiry_date"` // New field for expiration
}

type ShortenRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	ExpiryDate string `json:"expiry_date"` // Optional expiration date
}
