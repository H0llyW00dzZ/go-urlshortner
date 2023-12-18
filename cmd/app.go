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
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize the zap logger with a development configuration.
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, constant.FailedToIntializeLoggerContextLog+" %v\n", err)
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
	datastoreClient, err := setupDatastoreClient(ctx, logger)
	if err != nil {
		handleStartupFailure(err, logger)
	}

	router := setupRouter(datastoreClient, logger)
	startServer(router, logger, datastoreClient)
}

// setupDatastoreClient creates a new Datastore client and performs a test operation to check connectivity.
func setupDatastoreClient(ctx context.Context, logger *zap.Logger) (*datastore.Client, error) {
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreConfig := datastore.NewConfig(logger, projectID)
	datastoreClient, err := datastore.CreateDatastoreClient(ctx, datastoreConfig)
	if err != nil {
		return nil, fmt.Errorf(constant.FailedToCreateDatastoreClientContextLog+" %v", err)
	}

	if err := testClientConnection(ctx, datastoreClient); err != nil {
		return nil, err
	}

	return datastoreClient, nil
}

// testClientConnection attempts to perform a test operation with the Datastore client to check connectivity.
func testClientConnection(ctx context.Context, client *datastore.Client) error {
	// Perform a test operation, such as a health check read
	// Assuming 'health_check' is a known entity for this purpose
	_, err := datastore.GetURL(ctx, client, "health_check")
	if err == datastore.ErrNotFound {
		// If the specific test entity is not found, that's okay for a health check.
		// It means the client is connected and authorized; the entity just doesn't exist.
		return nil
	} else if err != nil {
		// Any other error means there's a problem with the connection or authorization.
		return fmt.Errorf(constant.DatastoreFailedtoCheckHealthContextLog+" %v", err)
	}
	// If there's no error, the client is connected and working.
	return nil
}

// checkEnvironment checks for the presence of required environment variables.
func checkEnvironment(logger *zap.Logger) error {
	// Check for the presence of required environment variables
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf(constant.DataStoreProjectIDEnvContextLog)
	}
	return nil
}

// initializeDatastoreClient is a helper function to create a new Datastore client.
func initializeDatastoreClient(ctx context.Context, logger *zap.Logger) (*datastore.Client, error) {
	// Create and return a datastore client
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreConfig := datastore.NewConfig(logger, projectID)                     // Create a new Config instance
	datastoreClient, err := datastore.CreateDatastoreClient(ctx, datastoreConfig) // Pass the config
	if err != nil {
		return nil, fmt.Errorf(constant.FailedtoCloseDatastoreContextLog+" %v", err)
	}
	return datastoreClient, nil
}

// setupRouter creates a new Gin router and sets up the middleware.
func setupRouter(datastoreClient *datastore.Client, logger *zap.Logger) *gin.Engine {
	// Set up the router and middleware
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up the router without the default logger
	router := gin.New()
	router.Use(gin.Recovery()) // Using only the recovery middleware

	// Using custom logging middleware with zap
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
	// Create log fields using the WithComponent function to convert string constants to zapcore.Field
	logFields := logmonitor.CreateLogFields("runServer",
		logmonitor.WithComponent(constant.ComponentGopher), // Use the constant ComponentGopher for the component
	)

	// Add the server address and port to the log fields
	logFields = append(logFields,
		zap.String("address", server.Addr),
		zap.String("port", getServerPort()),
	)
	// Testing human readable logging
	// Gopher will tell info to that devops always monitor the logs
	// Log the server starting message with the common fields
	logger.Info(constant.DeployEmoji+"  "+constant.ServerStartContextLog+" "+server.Addr, logFields...)

	// Attempt to start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Add the error to the log fields for error logging
		errorLogFields := append(logFields, logmonitor.WithError(err)())
		logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.ServerFailContextLog, errorLogFields...)
		os.Exit(1)
	}
}

// waitForShutdownSignal blocks until a SIGINT or SIGTERM signal is received, then shuts down the server.
// Note: This function can be ignored if you're using a managed service, such as Google Cloud Run. In such
// environments, Google Cloud Run (on top of Knative) sends a SIGTERM signal and manages the shutdown process for you. Therefore, the
// managed service continues to handle all operational aspects, indicating that the service is running within Google Cloud Run.
func waitForShutdownSignal(server *http.Server, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// Testing human readable logging
	// Gopher will tell info to that devops always monitor the logs
	s := <-quit
	logFields := logmonitor.CreateLogFields("waitForShutdownSignal",
		logmonitor.WithComponent(constant.ComponentGopher), // Use the constant ComponentGopher for the component
		logmonitor.WithSignal(s),                           // Use the WithSignal function from logmonitor
	)
	// Log the reception of the shutdown signal.
	logger.Info(constant.SignalSatelliteEmoji+"  "+constant.SignalContextLog, logFields...)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		// Log the error using the fields and include the error message.
		logger.Fatal(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.ServerForcetoShutdownContextLog, logFields...)
	}
}

// cleanupResources gracefully closes the Datastore client and logs any errors encountered.
func cleanupResources(logger *zap.Logger, datastoreClient *datastore.Client) {
	logger.Info("Closing datastore client...")
	if err := datastore.CloseClient(datastoreClient); err != nil {
		logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedtoCloseDatastoreContextLog, zap.Error(err))
	}

	logger.Info("Server exiting")
}

func handleStartupFailure(err error, logger *zap.Logger) {
	// Log the error using the provided zap.Logger
	logFields := logmonitor.CreateLogFields("handleStartupFailure",
		logmonitor.WithComponent(constant.ComponentGopher), // Use the constant ComponentGopher for the component
		logmonitor.WithError(err),                          // Include the error here, but it will be nil if there's no error
	)
	logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.StartupFailureContextLog, logFields...)

	// Optionally, print the error using the bannercli package.
	bannercli.PrintTypingBanner(err.Error(), 100*time.Millisecond)

	// Exit the program with a non-zero status code to indicate failure
	os.Exit(1)
}
