// Package datastore provides a set of tools to interact with Google Cloud Datastore.
// It abstracts datastore operations such as creating a client, saving, retrieving,
// updating, and deleting URL entities. It is designed to work with the URL shortener
// service to manage shortened URLs and their corresponding original URLs.
//
// # Types:
//
// The Client type is a wrapper around the Google Cloud Datastore client that
// provides additional methods for URL entity management. It abstracts away the
// underlying datastore implementation details.
//
//	type Client struct {
//	    *cloudDatastore.Client
//	}
//
// The URL type represents a URL entity within the datastore, with fields for the
// original URL and a unique identifier.
//
//	type URL struct {
//	    Original string `datastore:"original"` // The original URL.
//	    ID       string `datastore:"id"`       // The unique identifier for the shortened URL.
//	}
//
// # Variables:
//
// The package exposes an ErrNotFound variable, which is an error that represents
// the absence of a URL entity in the datastore.
//
//	var ErrNotFound = errors.New("no such entity")
//
// The package also includes a package-level Logger variable, which is intended to
// be used across the datastore package for consistent logging.
//
//	var Logger *zap.Logger
//
// # Handlers:
//
// The package provides functions to handle datastore operations. These functions
// include creating a new client, saving, retrieving, updating, and deleting URL
// entities.
//
//	func CreateDatastoreClient(ctx context.Context, config *Config) (*Client, error)
//	func SaveURL(ctx context.Context, client *Client, url *URL) error
//	func GetURL(ctx context.Context, client *Client, id string) (*URL, error)
//	func UpdateURL(ctx context.Context, client *Client, id string, newURL string) error
//	func DeleteURL(ctx context.Context, client *Client, id string) error
//	func CloseClient(client *Client) error
//
// # Example of package usage:
//
//	func main() {
//	    logger, _ := zap.NewDevelopment()
//	    ctx := datastore.CreateContext()
//	    config := datastore.NewConfig(logger, "my-project-id")
//	    client, err := datastore.CreateDatastoreClient(ctx, config)
//	    if err != nil {
//	        logger.Fatal("Failed to create datastore client", zap.Error(err))
//	    }
//	    defer datastore.CloseClient(client)
//
//	    // Use the client to interact with the datastore
//	    url := &datastore.URL{Original: "https://example.com", ID: "abc123"}
//	    if err := datastore.SaveURL(ctx, client, url); err != nil {
//	        logger.Fatal("Failed to save URL", zap.Error(err))
//	    }
//
//	    // Retrieve, update, or delete URLs as needed
//	    retrievedURL, err := datastore.GetURL(ctx, client, "abc123")
//	    if err != nil {
//	        logger.Error("Failed to retrieve URL", zap.Error(err))
//	    } else {
//	        logger.Info("Retrieved URL", zap.String("original", retrievedURL.Original))
//	    }
//	}
//
// The package functions are designed to be used in a concurrent environment and are
// safe for use by multiple goroutines.
//
// The CreateContext function is provided to create a new context for datastore
// operations, which can be used to control the lifetime of requests and to pass
// cancellation signals and other request-scoped values across API boundaries.
//
// It is important to close the datastore client when it is no longer needed by
// calling the CloseClient function to release any resources associated with the
// client.
//
// Copyright (c) 2023 H0llyW00dzZ
package datastore
