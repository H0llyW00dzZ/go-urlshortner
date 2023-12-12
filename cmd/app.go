package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := datastore.CreateContext()

	// Get the project ID from the "DATASTORE_PROJECT_ID" environment variable
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		fmt.Println("DATASTORE_PROJECT_ID environment variable is not set.")
		os.Exit(1)
	}

	datastoreClient, err := datastore.CreateDatastoreClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create datastore client: %v\n", err)
		os.Exit(1)
	}

	// Set Gin to release mode if not in development
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Register the handlers using Gin
	handlers.RegisterHandlersGin(router, datastoreClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Listening on port %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
