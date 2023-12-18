package datastore

import (
	"context"
	"errors"

	cloudDatastore "cloud.google.com/go/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
)

// Client wraps the cloudDatastore.Client to abstract away the underlying implementation.
type Client struct {
	*cloudDatastore.Client
}

// URL represents a shortened URL with its original URL and a unique identifier.
type URL struct {
	Original string `datastore:"original"` // The original URL.
	ID       string `datastore:"id"`       // The unique identifier for the shortened URL.
}

// Config holds the configuration settings for the datastore client.
type Config struct {
	Logger    *zap.Logger
	ProjectID string
}

// Logger is a package-level variable to access the zap logger throughout the datastore package.
// It is intended to be used by other functions within the package for logging purposes.
var Logger *zap.Logger

// ErrNotFound is the error returned when a requested entity is not found in the datastore.
var ErrNotFound = errors.New(DataStoreNosuchentity)

// SetLogger sets the logger instance for the package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// NewConfig creates a new instance of Config with the given logger and project ID.
func NewConfig(logger *zap.Logger, projectID string) *Config {
	return &Config{
		Logger:    logger,
		ProjectID: projectID,
	}
}

// CreateContext creates a new context that can be used for Datastore operations.
// It returns a non-nil, empty context.
func CreateContext() context.Context {
	return context.Background()
}

// CreateDatastoreClient creates a new client connected to Google Cloud Datastore.
// It requires a context and a Config struct to initialize the connection.
// It returns a new Datastore client wrapped in a custom Client type or an error if the connection could not be established.
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
// It takes a context, a Datastore client, and a URL struct.
// If the Kind 'urlz' does not exist, Datastore will create it automatically.
// Returns an error if the URL entity could not be saved.
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
// It requires a context, a Datastore client, and the ID of the URL entity.
// Returns the URL entity or an error if the entity could not be retrieved.
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
// It requires a context, a Datastore client, the ID of the URL entity, and the new URL to update.
// Returns an error if the URL entity could not be updated.
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
// It requires a context, a Datastore client, and the ID of the URL entity.
// Returns an error if the entity could not be deleted.
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
// Returns an error if the client could not be closed.
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
