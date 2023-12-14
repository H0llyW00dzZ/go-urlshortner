package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/H0llyW00dzZ/ChatGPT-Next-Web-Session-Exporter/bannercli"
	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/handlers"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Ensure that any buffered log entries are flushed before exiting.
	defer func() {
		err := logmonitor.Logger.Sync()
		if err != nil {
			// Handle the error, perhaps log to stderr or a file
			fmt.Fprintf(os.Stderr, "Failed to flush log: %v\n", err)
		}
	}()

	ctx := datastore.CreateContext()

	// Get the project ID from the "DATASTORE_PROJECT_ID" environment variable
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		bannercli.PrintTypingBanner("DATASTORE_PROJECT_ID environment variable is not set.", 100*time.Millisecond)
		os.Exit(1)
	}

	datastoreClient, err := datastore.CreateDatastoreClient(ctx, projectID)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create datastore client: %v\n", err)
		bannercli.PrintTypingBanner(errorMessage, 100*time.Millisecond)
		os.Exit(1)
	}

	// Set Gin to release mode if not in development
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Use the RequestLogger middleware from the logmonitor package
	router.Use(logmonitor.RequestLogger())

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

	// Inform the user that the server is starting
	logmonitor.Logger.Info("Listening on port", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Use zap logger to log the error if initialization was successful
		logmonitor.Logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}
