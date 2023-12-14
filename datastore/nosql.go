// Package datastore provides a set of functions to interact with Google Cloud Datastore.
//
// It allows for operations such as creating a new client, saving a URL entity,
// retrieving a URL entity by its ID, and closing the client connection.
//
// Copyright (c) 2023 H0llyW00dzZ
package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

// URL represents a shortened URL with its original URL and a unique identifier.
type URL struct {
	Original string `datastore:"original"` // The original URL.
	ID       string `datastore:"id"`       // The unique identifier for the shortened URL.
}

// CreateContext creates a new context that can be used for Datastore operations.
// It returns a non-nil, empty context.
func CreateContext() context.Context {
	return context.Background()
}

// CreateDatastoreClient creates a new client connected to Google Cloud Datastore.
// It requires a context and a projectID to initialize the connection.
// Returns a new Datastore client or an error if the connection could not be established.
func CreateDatastoreClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return nil, err
	}
	return client, nil
}

// SaveURL saves a new URL entity to Datastore under the Kind 'urlz'.
// It takes a context, a Datastore client, and a URL struct.
// If the Kind 'urlz' does not exist, Datastore will create it automatically.
// Returns an error if the URL entity could not be saved.
func SaveURL(ctx context.Context, client *datastore.Client, url *URL) error {
	key := datastore.NameKey("urlz", url.ID, nil)
	_, err := client.Put(ctx, key, url)
	if err != nil {
		fmt.Printf("Failed to save URL: %v\n", err)
		return err
	}
	return nil
}

// GetURL retrieves a URL entity by its ID from Datastore.
// It requires a context, a Datastore client, and the ID of the URL entity.
// Returns the URL entity or an error if the entity could not be retrieved.
func GetURL(ctx context.Context, client *datastore.Client, id string) (*URL, error) {
	key := datastore.NameKey("urlz", id, nil)
	url := new(URL)
	err := client.Get(ctx, key, url)
	if err != nil {
		fmt.Printf("Failed to get URL: %v\n", err)
		return nil, err
	}
	return url, nil
}

// CloseClient closes the Datastore client.
// It should be called to clean up resources and connections when the client is no longer needed.
// Logs an error if the client could not be closed, but does not return an error.
func CloseClient(client *datastore.Client) {
	err := client.Close()
	if err != nil {
		fmt.Printf("Failed to close datastore client: %v\n", err)
	}
}
