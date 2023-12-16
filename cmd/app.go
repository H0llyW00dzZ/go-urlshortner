// Package main/cmd is the entry point for the URL shortener service application.
// This service provides an HTTP server that handles requests for creating,
// retrieving, and editing shortened URLs. It uses the Gin web framework for
// routing and handling HTTP requests, zap for structured logging, and Google
// Cloud Datastore for storage of URL mappings.
//
// The main function initializes the necessary components such as the logger,
// the Datastore client, and the HTTP router. It also sets up the HTTP server
// and starts listening for incoming requests. The application's configuration
// is driven by environment variables, including the Datastore project ID and
// the desired port for the HTTP server.
//
// The service supports a RESTful API for managing URLs and includes middleware
// for request logging. The application is designed to be deployed as a
// containerized service, and it is capable of being scaled horizontally to
// handle high loads.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	startServer(router, logger, datastoreClient)
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

// startServer sets up and starts the HTTP server, and waits for a shutdown signal.
func startServer(router *gin.Engine, logger *zap.Logger, datastoreClient *datastore.Client) {
	server := createServer(router, logger)

	go runServer(server, logger)

	// Wait for interrupt signal to gracefully shut down the server
	waitForShutdownSignal(server, logger)

	// Close any other resources such as the datastore client
	cleanupResources(logger, datastoreClient)
}

// createServer initializes and returns a new HTTP server with the given router and logger.
func createServer(router *gin.Engine, logger *zap.Logger) *http.Server {
	port := getServerPort()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return server
}

// getServerPort retrieves the server port from the environment variable or defaults to "8080".
func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// runServer starts the server and logs any errors encountered during startup.
func runServer(server *http.Server, logger *zap.Logger) {
	logger.Info("Server is starting and Listening on address", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}

// waitForShutdownSignal blocks until a SIGINT or SIGTERM signal is received, then shuts down the server.
func waitForShutdownSignal(server *http.Server, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// Testing human readable logging
	// Gopher will tell info to that devops always monitor the logs
	s := <-quit
	logger.Info("Received shutdown signal", zap.String("signal", s.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}
}

// cleanupResources gracefully closes the Datastore client and logs any errors encountered.
func cleanupResources(logger *zap.Logger, datastoreClient *datastore.Client) {
	logger.Info("Closing datastore client...")
	if err := datastore.CloseClient(datastoreClient); err != nil {
		logger.Error("Failed to close datastore client", zap.Error(err))
	}

	logger.Info("Server exiting")
}

func handleStartupFailure(err error, logger *zap.Logger) {
	// Log the error using the provided zap.Logger
	logger.Error("Startup failure", zap.Error(err))

	// Optionally, print the error using the bannercli package.
	bannercli.PrintTypingBanner(err.Error(), 100*time.Millisecond)

	// Exit the program with a non-zero status code to indicate failure
	os.Exit(1)
}
