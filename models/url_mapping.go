package models

type URLMapping struct {
	ShortID string `datastore:"short_id"`
	LongURL string `datastore:"long_url"`
}

type ShortenRequest struct {
	LongURL string `json:"long_url" binding:"required"`
}