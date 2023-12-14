// Package datastore provides a set of functions to interact with Google Cloud Datastore.
//
// It allows for operations such as creating a new client, saving a URL entity,
// retrieving a URL entity by its ID, and closing the client connection.
// This version has been updated to use the zap logger for consistent and structured logging.
//
// Copyright (c) 2023 H0llyW00dzZ
package datastore

import (
	"context"
	"errors"

	cloudDatastore "cloud.google.com/go/datastore"
	"go.uber.org/zap"
)

// Client wraps the cloudDatastore.Client to abstract away the underlying implementation.
type Client struct {
	*cloudDatastore.Client
}

// URL represents a shortened URL with its original URL and a unique identifier.
type URL struct {
	Original string `datastore:"original"` // The original URL.
	ID       string `datastore:"id"`       // The unique identifier for the shortened URL.
}

// Logger is a package-level variable to access the zap logger throughout the datastore package.
// It is intended to be used by other functions within the package for logging purposes.
var Logger *zap.Logger

// ErrNotFound is the error returned when a requested entity is not found in the datastore.
var ErrNotFound = errors.New("datastore: no such entity")

// SetLogger sets the logger instance for the package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

func init() {
	// Initialize the zap logger with a development configuration.
	// This config is console-friendly and outputs logs in plaintext.
	config := zap.NewDevelopmentConfig()
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err) // If logger initialization fails, the application will panic.
	}
}

// CreateContext creates a new context that can be used for Datastore operations.
// It returns a non-nil, empty context.
func CreateContext() context.Context {
	return context.Background()
}

// CreateDatastoreClient creates a new client connected to Google Cloud Datastore.
// It requires a context and a projectID to initialize the connection.
// Returns a new Datastore client or an error if the connection could not be established.
func CreateDatastoreClient(ctx context.Context, projectID string) (*Client, error) {
	cloudClient, err := cloudDatastore.NewClient(ctx, projectID)
	if err != nil {
		Logger.Error("Failed to create client", zap.Error(err))
		return nil, err
	}
	return &Client{cloudClient}, nil
}

// SaveURL saves a new URL entity to Datastore under the Kind 'urlz'.
// It takes a context, a Datastore client, and a URL struct.
// If the Kind 'urlz' does not exist, Datastore will create it automatically.
// Returns an error if the URL entity could not be saved.
func SaveURL(ctx context.Context, client *Client, url *URL) error {
	key := cloudDatastore.NameKey("urlz", url.ID, nil)
	_, err := client.Put(ctx, key, url)
	if err != nil {
		// Use zap logger to log the error for consistent logging.
		Logger.Error("Failed to save URL", zap.Error(err))
		return err
	}
	return nil
}

// GetURL retrieves a URL entity by its ID from Datastore.
// It requires a context, a Datastore client, and the ID of the URL entity.
// Returns the URL entity or an error if the entity could not be retrieved.
func GetURL(ctx context.Context, dsClient *Client, id string) (*URL, error) {
	key := cloudDatastore.NameKey("urlz", id, nil)
	url := new(URL)
	err := dsClient.Get(ctx, key, url)
	if err != nil {
		if err == cloudDatastore.ErrNoSuchEntity {
			return nil, ErrNotFound
		}
		// Handle other possible errors.
		return nil, err
	}
	return url, nil
}

// CloseClient closes the Datastore client.
// It should be called to clean up resources and connections when the client is no longer needed.
// Returns an error if the client could not be closed.
func CloseClient(client *Client) error {
	err := client.Close()
	if err != nil {
		// Use zap logger to log the error for consistent logging.
		Logger.Error("Failed to close datastore client", zap.Error(err))
		return err // Now returning the error so the caller can handle it.
	}
	return nil
}
