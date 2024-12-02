package config

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

func ConnectDatastore() *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "genericissuer")
	if err != nil {
		log.Fatalf("Failed to connect to Google Cloud Datastore: %v", err)
	}
	return client
}
