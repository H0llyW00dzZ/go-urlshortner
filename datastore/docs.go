// Package datastore provides a set of tools to interact with Google Cloud Datastore.
// It abstracts datastore operations such as creating a client, saving, retrieving,
// updating, and deleting URL entities. It is designed to work with the URL shortener
// service to manage shortened URLs and their corresponding original URLs.
//
// # Types
//
//   - Client: Wraps the Google Cloud Datastore client and provides methods for URL entity management.
//   - URL: Represents a URL entity within the datastore with fields for the original URL and a unique identifier.
//   - DatastoreError: Represents structured errors from the Datastore client, including an error code, description, and optional details or URL.
//
// # Variables
//
//   - ErrNotFound: An error representing the absence of a URL entity in the datastore.
//   - Logger: A package-level variable for consistent logging. It should be set using SetLogger before using logging functions.
//
// # Handler Functions
//
// The package offers functions for datastore operations:
//   - CreateDatastoreClient: Initializes and returns a new datastore client.
//   - SaveURL: Saves a URL entity to the datastore.
//   - GetURL: Retrieves a URL entity from the datastore by ID.
//   - UpdateURL: Updates an existing URL entity in the datastore.
//   - DeleteURL: Deletes a URL entity from the datastore by ID.
//   - CloseClient: Closes the datastore client and releases resources.
//   - ParseDatastoreClientError: Parses errors from the Datastore client into a structured format.
//
// # Example Usage
//
// The following example demonstrates how to initialize a datastore client, save a URL entity,
// and retrieve it using the package's functions:
//
//	func main() {
//	    logger, _ := zap.NewDevelopment()
//	    defer logger.Sync()
//
//	    ctx := datastore.CreateContext()
//	    config := datastore.NewConfig(logger, "my-project-id")
//	    client, err := datastore.CreateDatastoreClient(ctx, config)
//	    if err != nil {
//	        logger.Fatal("Failed to create datastore client", zap.Error(err))
//	    }
//	    defer func() {
//	        if err := datastore.CloseClient(client); err != nil {
//	            logger.Error("Failed to close datastore client", zap.Error(err))
//	        }
//	    }()
//
//	    // Use the client to save a new URL entity
//	    url := &datastore.URL{Original: "https://example.com", ID: "abc123"}
//	    if err := datastore.SaveURL(ctx, client, url); err != nil {
//	        logger.Fatal("Failed to save URL", zap.Error(err))
//	    }
//
//	    // Retrieve the URL entity by ID
//	    retrievedURL, err := datastore.GetURL(ctx, client, "abc123")
//	    if err != nil {
//	        logger.Error("Failed to retrieve URL", zap.Error(err))
//	    } else {
//	        logger.Info("Retrieved URL", zap.String("original", retrievedURL.Original))
//	    }
//	}
//
// # Concurrency
//
// The package functions are designed for concurrent use and are safe for use by multiple goroutines.
//
// # Contexts
//
// The CreateContext function is provided for creating new contexts for datastore operations,
// allowing for request lifetime control and value passing across API boundaries.
//
// # Cleanup
//
// It is important to close the datastore client with CloseClient to release resources.
//
// Copyright (c) 2023 by H0llyW00dzZ
package datastore
