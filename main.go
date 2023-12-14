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
	// Initialize the zap logger with a development configuration.
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flush any buffered log entries

	// Pass the logger instance to other packages
	datastore.SetLogger(logger)
	logmonitor.SetLogger(logger)
	handlers.SetLogger(logger)

	if err := checkEnvironment(logger); err != nil {
		handleStartupFailure(err, logger)
	}

	ctx := datastore.CreateContext()
	datastoreClient, err := initializeDatastoreClient(ctx, logger)
	if err != nil {
		handleStartupFailure(err, logger)
	}

	router := setupRouter(datastoreClient, logger)
	startServer(router, logger)
}

func checkEnvironment(logger *zap.Logger) error {
	// Check for the presence of required environment variables
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("DATASTORE_PROJECT_ID environment variable is not set")
	}
	return nil
}

func initializeDatastoreClient(ctx context.Context, logger *zap.Logger) (*datastore.Client, error) {
	// Create and return a datastore client
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreClient, err := datastore.CreateDatastoreClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create datastore client: %v", err)
	}
	return datastoreClient, nil
}

func setupRouter(datastoreClient *datastore.Client, logger *zap.Logger) *gin.Engine {
	// Set up the router and middleware
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(logmonitor.RequestLogger(logger))
	handlers.RegisterHandlersGin(router, datastoreClient)

	return router
}

func startServer(router *gin.Engine, logger *zap.Logger) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	// Inform the devops that the server is starting
	logger.Info("Listening on address", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}

func handleStartupFailure(err error, logger *zap.Logger) {
	// Log the error using the provided zap.Logger
	logger.Error("Startup failure", zap.Error(err))

	// Optionally, print the error using the bannercli package.
	bannercli.PrintTypingBanner(err.Error(), 100*time.Millisecond)

	// Exit the program with a non-zero status code to indicate failure
	os.Exit(1)
}
