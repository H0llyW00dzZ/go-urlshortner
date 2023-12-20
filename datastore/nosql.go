package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cloudDatastore "cloud.google.com/go/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
)

// Client wraps the cloudDatastore.Client to abstract away the underlying implementation.
// This allows for easier mocking and testing, as well as decoupling the code from the specific datastore client used.
type Client struct {
	*cloudDatastore.Client
}

// URL represents a shortened URL with its original URL and a unique identifier.
// The struct tags specify how each field is stored in the datastore.
type URL struct {
	Original string `datastore:"original"` // The original URL.
	ID       string `datastore:"id"`       // The unique identifier for the shortened URL.
}

// Config holds the configuration settings for the datastore client.
// This includes the logger for logging operations and the project ID for Google Cloud Datastore.
type Config struct {
	Logger    *zap.Logger // The logger for logging operations within the datastore package.
	ProjectID string      // The Google Cloud project ID where the datastore is located.
}

// DatastoreError represents a structured error for the Datastore client.
// It includes details about the error code, description, and any additional details or URLs related to the error.
type DatastoreError struct {
	Code        string `json:"code"`                  // The error code.
	Description string `json:"description"`           // The human-readable error description.
	Details     string `json:"Details,omitempty"`     // Additional details about the error.
	DetailsURL  string `json:"details_url,omitempty"` // A URL with more information about the error, if available.
}

// Error implements the error interface for DatastoreError.
// This method formats the DatastoreError into a string, including the details URL if present.
func (e *DatastoreError) Error() string {
	if e.DetailsURL != "" {
		return fmt.Sprintf("Code: %s, Description: %s, Details: %s", e.Code, e.Description, e.DetailsURL)
	}
	return fmt.Sprintf("Code: %s, Description: %s", e.Code, e.Description)
}

// MarshalJSON ensures that the DatastoreError is marshaled correctly.
// It overrides the default JSON marshaling to include only the relevant fields.
func (e *DatastoreError) MarshalJSON() ([]byte, error) {
	type Alias DatastoreError
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

// Logger is a package-level variable to access the zap logger throughout the datastore package.
// It can be set using the SetLogger function and is used by various functions for consistent logging.
var Logger *zap.Logger

// ErrNotFound is the error returned when a requested entity is not found in the datastore.
// This error is used to signal that a URL entity with the provided ID does not exist.
var ErrNotFound = errors.New(DataStoreNosuchentity)

// SetLogger sets the logger instance for the package.
// This function configures the package-level Logger variable for use throughout the datastore package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// NewConfig creates a new instance of Config with the given logger and project ID.
// This function is used to configure the datastore client with necessary settings.
func NewConfig(logger *zap.Logger, projectID string) *Config {
	return &Config{
		Logger:    logger,
		ProjectID: projectID,
	}
}

// CreateContext creates a new context that can be used for Datastore operations.
// It returns a non-nil, empty context that can be used to carry deadlines, cancellation signals,
// and other request-scoped values across API boundaries and between processes.
func CreateContext() context.Context {
	return context.Background()
}

// CreateDatastoreClient creates a new client connected to Google Cloud Datastore.
// It initializes the connection using the provided context and configuration settings.
// The function returns a new Client instance or an error if the connection could not be established.
func CreateDatastoreClient(ctx context.Context, config *Config) (*Client, error) {
	cloudClient, err := cloudDatastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		// Create structured log fields using logmonitor
		logFields := logmonitor.CreateLogFields("CreateDatastoreClient",
			logmonitor.WithComponent(constant.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithError(err),                         // Include the error here, but it will be nil if there's no error
		)
		// Log the error with structured fields
		// Note: This logger is specifically configured for CreateDatastoreClient and is synchronized with Google Cloud Datastore's and any Google Cloud Service (e.g, Google Cloud Auth) error handling in the binary world.
		config.Logger.Error(constant.AlertEmoji+" "+DataStoreFailedtoCreateClient, logFields...)
		return nil, err
	}
	return &Client{cloudClient}, nil
}

// SaveURL saves a new URL entity to Datastore under the Kind 'urlz'.
// It uses the provided context and datastore client to save the URL struct to the datastore.
// The function returns an error if the URL entity could not be saved.
func SaveURL(ctx context.Context, client *Client, url *URL) error {
	key := cloudDatastore.NameKey(DataStoreNameKey, url.ID, nil)
	_, err := client.Put(ctx, key, url)
	if err != nil {
		// Use zap logger to log the error for consistent logging.
		logmonitor.Logger.Error(constant.AlertEmoji+" "+DataStoreFailedtoCreateClient, zap.Error(err))
		return err
	}
	return nil
}

// GetURL retrieves a URL entity by its ID from Datastore.
// It uses the provided context and datastore client to look up the URL entity by its unique identifier.
// The function returns the found URL entity or an error if the entity could not be retrieved.
func GetURL(ctx context.Context, dsClient *Client, id string) (*URL, error) {
	key := cloudDatastore.NameKey(DataStoreNameKey, id, nil)
	url := new(URL)
	err := dsClient.Get(ctx, key, url)
	if err != nil {
		if err == cloudDatastore.ErrNoSuchEntity {
			return nil, ErrNotFound
		}
		// Handle other possible errors.
		return nil, err
	}
	return url, nil
}

// UpdateURL updates an existing URL entity in Datastore with a new URL.
// It performs the update within a transaction to ensure the operation is atomic.
// The function returns an error if the URL entity could not be updated.
func UpdateURL(ctx context.Context, client *Client, id string, newURL string) error {
	key := cloudDatastore.NameKey(DataStoreNameKey, id, nil)
	// Transactionally retrieve the existing URL and update it.
	_, err := client.RunInTransaction(ctx, func(tx *cloudDatastore.Transaction) error {
		url := new(URL)
		if err := tx.Get(key, url); err != nil {
			if err == cloudDatastore.ErrNoSuchEntity {
				return ErrNotFound
			}
			return err
		}

		// Update the URL's Original field with the new URL.
		url.Original = newURL
		_, err := tx.Put(key, url)
		return err
	})

	if err != nil {
		logmonitor.Logger.Error(constant.AlertEmoji+" "+DataStoreFailedtoUpdateURL, zap.String("id", id), zap.Error(err))
		return err
	}

	return nil
}

// DeleteURL deletes a URL entity by its ID from Datastore.
// It uses the provided context and datastore client to delete the URL entity by its unique identifier.
// The function returns an error if the entity could not be deleted.
func DeleteURL(ctx context.Context, client *Client, id string) error {
	key := cloudDatastore.NameKey(DataStoreNameKey, id, nil)
	err := client.Delete(ctx, key)
	if err != nil {
		if err == cloudDatastore.ErrNoSuchEntity {
			return ErrNotFound
		}
		// Log and handle other possible errors.
		logmonitor.Logger.Error(constant.AlertEmoji+" "+DataStoreFailedtoUpdateURL, zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

// CloseClient closes the Datastore client.
// It should be called to clean up resources and connections when the client is no longer needed.
// The function returns an error if the client could not be closed.
func CloseClient(client *Client) error {
	if client == nil {
		return nil // or return an error if you expect the client to never be nil
	}
	err := client.Close()
	if err != nil {
		logmonitor.Logger.Error(constant.AlertEmoji+" "+DataStoreFailedToCloseClient, zap.Error(err))
		return err
	}
	return nil
}

// ParseDatastoreClientError parses the error from the Datastore client and returns a structured error.
// It attempts to extract meaningful information from the error returned by the datastore client
// and formats it into a DatastoreError. It returns the structured error and a parsing error, if any.
func ParseDatastoreClientError(err error) (*DatastoreError, error) {
	if err == nil {
		return nil, fmt.Errorf(noerrortoparse)
	}

	errorMessage := err.Error()
	parts := strings.Fields(errorMessage) // Use Fields to automatically handle splitting by whitespace.
	if len(parts) < 2 {
		return nil, fmt.Errorf(unexpectederrorformat)
	}

	datastoreErr := &DatastoreError{
		Code:    parts[0],                     // Assuming the code is the first part of the error message.
		Details: strings.Join(parts[1:], " "), // The rest is the details.
	}

	datastoreErr.DetailsURL = extractDetailsURL(errorMessage)
	datastoreErr = checkForSpecificError(errorMessage, datastoreErr)

	return datastoreErr, nil
}

// extractDetailsURL extracts the details URL from the error message if present.
// It looks for an "http" substring and assumes that the URL is the last part of the error message.
// The function returns the extracted URL or an empty string if no URL is found.
func extractDetailsURL(errorMessage string) string {
	if strings.Contains(errorMessage, "http") {
		parts := strings.Fields(errorMessage)
		return strings.Trim(parts[len(parts)-1], "\"") // Assuming the URL is the last part.
	}
	return ""
}

// checkForSpecificError checks for specific errors and updates the DatastoreError accordingly.
// It looks for known error patterns in the error message and sets the appropriate description
// and details in the DatastoreError. The function returns the updated DatastoreError.
func checkForSpecificError(errorMessage string, datastoreErr *DatastoreError) *DatastoreError {
	if strings.Contains(errorMessage, "invalid_grant") {
		datastoreErr.Description = DataStoreAuthInvalidToken
		datastoreErr.Details = errorMessage
	}
	return datastoreErr
}
