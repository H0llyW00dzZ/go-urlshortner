// Package datastore provides a set of tools to interact with Google Cloud Datastore.
// It abstracts datastore operations such as creating a client, saving, retrieving,
// updating, and deleting URL entities. It is designed to work with the URL shortener
// service to manage shortened URLs and their corresponding original URLs.
//
// The package encapsulates the complexity of direct datastore interactions and offers
// a simplified API for the rest of the application. It also integrates structured
// logging using the zap library, which allows for consistent and searchable log
// entries across the service.
//
// Usage example:
//
//	func main() {
//	    ctx := datastore.CreateContext()
//	    client, err := datastore.CreateDatastoreClient(ctx, "my-project-id")
//	    if err != nil {
//	        log.Fatalf("Failed to create datastore client: %v", err)
//	    }
//	    defer datastore.CloseClient(client)
//
//	    // Use the client to interact with the datastore
//	    url := &datastore.URL{Original: "https://example.com", ID: "abc123"}
//	    if err := datastore.SaveURL(ctx, client, url); err != nil {
//	        log.Fatalf("Failed to save URL: %v", err)
//	    }
//
//	    // Retrieve, update, or delete URLs as needed
//	}
//
// The package also defines a custom `URL` type that represents a URL entity within
// the datastore. This type is used for all operations that involve URL entities.
//
// The `Client` type wraps the Google Cloud Datastore client, providing a layer of
// abstraction that allows for mocking and testing without relying on an actual
// datastore instance.
//
// The `SetLogger` function allows for the injection of a zap logger instance, which
// is used throughout the package to log operations and errors. This ensures that
// any issues or actions are logged in a format that is consistent with the rest of
// the URL shortener service.
//
// Error handling is a critical aspect of the package. The `ErrNotFound` variable is
// returned when a requested entity is not found in the datastore, allowing calling
// code to distinguish between different types of errors and handle them accordingly.
//
// The package functions are designed to be used in a concurrent environment and are
// safe for use by multiple goroutines.
//
// The `CreateContext` function is provided to create a new context for datastore
// operations, which can be used to control the lifetime of requests and to pass
// cancellation signals and other request-scoped values across API boundaries.
//
// It is important to close the datastore client when it is no longer needed by
// calling the `CloseClient` function to release any resources associated with the
// client.
//
// Copyright (c) 2023 H0llyW00dzZ
package datastore
