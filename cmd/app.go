package main

import (
	"context"
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
	defer flushLogs()

	if err := checkEnvironment(); err != nil {
		handleStartupFailure(err)
	}

	ctx := datastore.CreateContext()
	datastoreClient, err := initializeDatastoreClient(ctx)
	if err != nil {
		handleStartupFailure(err)
	}

	router := setupRouter(datastoreClient)
	startServer(router)
}

func flushLogs() {
	// Flush logs for all loggers
	if err := logmonitor.Logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush logmonitor log: %v\n", err)
	}
	if err := datastore.Logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush datastore log: %v\n", err)
	}
	if err := handlers.Logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush handlers log: %v\n", err)
	}
}

func checkEnvironment() error {
	// Check for the presence of required environment variables
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("DATASTORE_PROJECT_ID environment variable is not set")
	}
	return nil
}

func initializeDatastoreClient(ctx context.Context) (*datastore.Client, error) {
	// Create and return a datastore client
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreClient, err := datastore.CreateDatastoreClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create datastore client: %v", err)
	}
	return datastoreClient, nil
}

func setupRouter(datastoreClient *datastore.Client) *gin.Engine {
	// Set up the router and middleware
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(logmonitor.RequestLogger())
	handlers.RegisterHandlersGin(router, datastoreClient)

	return router
}

func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	// Inform the devops that the server is starting
	logmonitor.Logger.Info("Listening on address", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logmonitor.Logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}

func handleStartupFailure(err error) {
	// Handle any startup failures, print errors, and exit
	bannercli.PrintTypingBanner(err.Error(), 100*time.Millisecond)
	os.Exit(1)
}
