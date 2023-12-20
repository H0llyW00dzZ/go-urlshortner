// Package handlers implements the HTTP or HTTPS handling logic for a URL shortener service.
// It provides handlers for creating, retrieving, updating, and deleting shortened URLs,
// leveraging Google Cloud Datastore for persistent storage. The package includes middleware
// for rate limiting and access control, ensuring that endpoints are protected against abuse
// and sensitive operations are restricted to internal services.
//
// Handlers are registered with the Gin web framework's router, forming the RESTful API of the service.
// Each handler function is tasked with handling specific HTTP or HTTPS request types, validating payloads,
// performing operations against the datastore, and crafting the HTTP or HTTPS response.
//
// Consistent and structured logging is maintained across the package using centralized logging functions,
// which aid in the systematic recording of operational events for ease of debugging and service monitoring.
//
// # Example of package usage
//
//	func main() {
//	    // Initialize a Gin router.
//	    router := gin.Default()
//
//	    // Create a new datastore client (assuming a constructor function exists).
//	    dsClient := datastore.NewClient(context.Background(), "example-project-id-0x1337")
//
//	    // Register the URL shortener's HTTP handlers with the Gin router.
//	    RegisterHandlersGin(router, dsClient)
//
//	    // Start the HTTP server on port 8080.
//	    router.Run(":8080")
//	}
//
// # Types and Variables
//
// The package defines various types to encapsulate request payloads and middleware functions:
//
//   - CreateURLPayload: Represents the JSON payload for creating a new shortened URL, containing the original URL.
//   - UpdateURLPayload: Represents the JSON payload for updating an existing shortened URL, containing the original and new URLs along with an identifier.
//   - DeleteURLPayload: Represents the JSON payload for deleting a shortened URL, containing the URL and its identifier.
//
// The following code snippets illustrate the structures of these types:
//
//	type CreateURLPayload struct {
//	    URL string `json:"url" binding:"required,url"`
//	}
//
//	type UpdateURLPayload struct {
//	    ID     string `json:"id" binding:"required"`
//	    OldURL string `json:"old_url" binding:"required,url"`
//	    NewURL string `json:"new_url" binding:"required,url"`
//	}
//
//	type DeleteURLPayload struct {
//	    ID  string `json:"id" binding:"required"`
//	    URL string `json:"url" binding:"required,url"`
//	}
//
// Middleware functions such as InternalOnly enforce access control by requiring a secret
// value in the request header, compared against an environment variable.
//
// # The package also exports several key variables
//
//   - basePath: A string representing the base path for the URL shortener's endpoints.
//   - internalSecretValue: A string used by the InternalOnly middleware to validate requests against internal services.
//   - RateLimiterStore: A sync.Map that stores rate limiters for each client IP address.
//
// # The following code snippets illustrate the declaration of these variables
//
//	var basePath string
//	var internalSecretValue string
//	var RateLimiterStore sync.Map
//
// # Handler Functions
//
// The package provides several HTTP handler functions to manage URL entities. These functions
// are designed to be registered with the Gin web framework's router and handle different
// HTTP methods and endpoints.
//
//   - getURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc:
//     Retrieves the original URL based on the short identifier provided in the request path
//     and redirects the client to it. Responds with HTTP 404 if the URL is not found, HTTP 429 if rate limit is exceeded,
//     or HTTP 500 for other errors.
//
//   - postURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc:
//     Handles the creation of a new shortened URL. It expects a JSON payload with the original
//     URL, generates a short identifier, stores the mapping, and returns the shortened URL.
//
//   - editURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc:
//     Manages the updating of an existing shortened URL. It validates the request payload,
//     verifies the existing URL, and updates it with the new URL provided.
//
//   - deleteURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc:
//     Handles the deletion of an existing shortened URL. It validates the provided ID and URL,
//     and if they match the stored entity, deletes the URL from the datastore.
//
// Each handler function utilizes the provided datastore client to interact with Google Cloud
// Datastore and leverages structured logging for operational events.
//
// # Helper Functions
//
// The package contains a variety of helper functions that support the primary handler functions.
// These helpers perform tasks such as request validation, data retrieval, rate limiting, and response generation.
//
// # Middleware
//
// The package includes middleware functions that provide additional layers of request
// processing, such as rate limiting and access control. These middleware functions are
// applied to certain handler functions to enforce security policies and request validation.
//
//   - InternalOnly():
//     A middleware function that restricts access to certain endpoints to internal services
//     only by requiring a secret value in the request header.
//
// Middleware functions are registered within the Gin router setup and are executed in the
// order they are applied to the routes.
//
// # Registering Handlers with Gin Router
//
// The handler functions are registered with the Gin router in the main application setup.
// This registration associates HTTP methods and paths with the corresponding handler functions
// and applies any necessary middleware.
//
//	func RegisterHandlersGin(router *gin.Engine, dsClient *datastore.Client) {
//	    router.GET(basePath+":id", getURLHandlerGin(dsClient))
//	    router.POST(basePath, InternalOnly(), postURLHandlerGin(dsClient))
//	    router.PUT(basePath+":id", InternalOnly(), editURLHandlerGin(dsClient))
//	    router.DELETE(basePath+":id", InternalOnly(), deleteURLHandlerGin(dsClient))
//	}
//
// The RegisterHandlersGin function is the central point for configuring the routing
// for the URL shortener service, ensuring that each endpoint is handled correctly.
//
// # Bug Fixes and Security Enhancements
//
// This section provides an overview of significant bug fixes and security enhancements
// that have been implemented in the package. The aim is to maintain transparency with
// users and to demonstrate a commitment to the security and reliability of the service.
//
//   - Version 0.3.2 (Include Latest):
//     Resolved an issue leading to Insecure Direct Object Reference (IDOR) vulnerability that was present
//     when parsing JSON payloads. Previously, malformed JSON could be used to bypass
//     payload validation, potentially allowing attackers to modify URLs without proper
//     verification. The parsing logic has been fortified to ensure that only well-formed,
//     validated JSON payloads are accepted, and any attempt to submit broken or malicious
//     JSON will be rejected.
//
// Copyright (c) 2023 by H0llyW00dzZ
package handlers
