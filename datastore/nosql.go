package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type URL struct {
	Original string `datastore:"original"`
	ID       string `datastore:"id"`
}

// CreateContext creates a new context that can be used for Datastore operations.
func CreateContext() context.Context {
	return context.Background()
}

// CreateDatastoreClient creates a new client connected to Google Cloud Datastore.
func CreateDatastoreClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return nil, err
	}
	return client, nil
}

// SaveURL saves a new URL entity to Datastore.
func SaveURL(ctx context.Context, client *datastore.Client, url *URL) error {
	key := datastore.NameKey("URL", url.ID, nil)
	_, err := client.Put(ctx, key, url)
	if err != nil {
		fmt.Printf("Failed to save URL: %v\n", err)
		return err
	}
	return nil
}

// GetURL retrieves a URL by its ID from Datastore.
func GetURL(ctx context.Context, client *datastore.Client, id string) (*URL, error) {
	key := datastore.NameKey("URL", id, nil)
	url := new(URL)
	err := client.Get(ctx, key, url)
	if err != nil {
		fmt.Printf("Failed to get URL: %v\n", err)
		return nil, err
	}
	return url, nil
}

// CloseClient closes the Datastore client.
func CloseClient(client *datastore.Client) {
	err := client.Close()
	if err != nil {
		fmt.Printf("Failed to close datastore client: %v\n", err)
	}
}
