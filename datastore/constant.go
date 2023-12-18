package datastore

// Define datastore constants.
// This is used to identify the component that is logging the message.
// Easy Maintenance: If the error message changes, it only need to change it in one place.
const (
	DataStoreNosuchentity         = "datastore: no such entity"
	DataStoreFailedtoCreateClient = "Failed to create client"
	DataStoreFailedtoSaveURL      = "Failed to save URL"
	DataStoreFailedtoGetURL       = "Failed to get URL"
	DataStoreFailedtoUpdateURL    = "Failed to update URL"
	DataStoreFailedtoDeleteURL    = "Failed to delete URL"
	DataStoreFailedToCloseClient  = "Failed to close client"

	// DataStoreNameKey is the name of the Kind in Datastore for URL entities.
	// Defining it here enables changing the Kind name in one place if needed.
	DataStoreNameKey = "urlz"

	// URL Info Messages
	InfoAttemptingToUpdateURLInDatastore = "Attempting to update URL in Datastore"
	InfoFailedToUpdateURLInDatastore     = "Failed to update URL in Datastore"
	InfoUpdateSuccessful                 = "URL updated successfully in the datastore"
)
