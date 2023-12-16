// Package handlers provides the HTTP handling logic for the URL shortener service.
// It includes handlers for creating, retrieving, updating, and deleting shortened URLs,
// with storage backed by Google Cloud Datastore. The package also offers middleware
// for access control, ensuring that certain operations are restricted to internal use.
//
// The handlers are designed to work with the Gin web framework and are registered
// to the Gin router, establishing the service's RESTful API. Each handler function
// is responsible for processing specific types of HTTP requests, validating input,
// interacting with the datastore, and formatting the HTTP response.
//
// Structured logging is employed throughout the package via the `logmonitor` package,
// ensuring that operational events are recorded in a consistent and searchable format.
// This facilitates debugging and monitoring of the service.
//
// Usage example:
//
//	func main() {
//	    router := gin.Default()
//	    dsClient := datastore.NewClient() // Assuming a function to create a new datastore client
//	    handlers.SetLogger(logmonitor.Logger) // Set the logger for the handlers package
//	    handlers.RegisterHandlersGin(router, dsClient)
//	    router.Run(":8080")
//	}
//
// The package defines various types to represent request payloads and middleware functions.
// The `InternalOnly` middleware function enforces access control by requiring a secret
// value in the request header, which is compared against an environment variable.
//
// Handler functions such as `getURLHandlerGin` and `postURLHandlerGin` serve as endpoints
// for fetching and storing URL mappings, respectively. The `editURLHandlerGin` and
// `deleteURLHandlerGin` functions provide the logic for updating and deleting mappings.
//
// The `RegisterHandlersGin` function is the entry point for setting up the routes and
// associating them with their handlers. It ensures that all routes are prefixed with
// a base path that can be configured via an environment variable.
//
// Copyright (c) 2023 H0llyW00dzZ
package handlers
