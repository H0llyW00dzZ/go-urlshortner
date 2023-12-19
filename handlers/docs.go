// Package handlers implements the HTTP or HTTPS handling logic for a URL shortener service.
// It provides handlers for creating, retrieving, updating, and deleting shortened URLs,
// leveraging Google Cloud Datastore for persistent storage. The package includes middleware
// for access control, which ensures that sensitive operations are restricted to internal services.
//
// Handlers are registered with the Gin web framework's router, forming the RESTful API of the service.
// Each handler function is tasked with handling specific HTTP or HTTPS request types, validating payloads,
// performing operations against the datastore, and crafting the HTTP or HTTPS response.
//
// Consistent and structured logging is maintained across the package using the `logmonitor` package,
// which aids in the systematic recording of operational events for ease of debugging and service monitoring.
//
// Example of package usage:
//
//	func main() {
//	    // Initialize a Gin router.
//	    router := gin.Default()
//
//	    // Create a new datastore client (assuming a constructor function exists).
//	    dsClient := datastore.NewClient()
//
//	    // Set the logger instance for the handlers package.
//	    handlers.SetLogger(logmonitor.NewLogger())
//
//	    // Register the URL shortener's HTTP handlers with the Gin router.
//	    handlers.RegisterHandlersGin(router, dsClient)
//
//	    // Start the HTTP server on port 8080.
//	    router.Run(":8080")
//	}
//
// The package defines various types to encapsulate request payloads and middleware functions.
// For instance, the `InternalOnly` middleware function enforces access restrictions by matching
// a secret value in the request header against a predefined environment variable.
//
// Endpoint handler functions such as `getURLHandlerGin` and `postURLHandlerGin` manage
// the retrieval and creation of URL mappings. Functions like `editURLHandlerGin` and
// `deleteURLHandlerGin` handle the updating and deletion of these mappings.
//
// The `RegisterHandlersGin` function configures the routing for the service, associating
// endpoints with their respective handler functions and applying any necessary middleware.
// It also allows for the configuration of a base path for all routes, which can be set
// through an environment variable.
//
// Copyright (c) 2023 by H0llyW00dzZ
package handlers
