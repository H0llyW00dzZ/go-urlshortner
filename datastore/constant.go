package datastore

// Define datastore constants.
//
// This is used to identify the component that is logging the message.
//
// Easy Maintenance: If the error message changes, it only need to change it in one place.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	DataStoreNosuchentity         = "datastore: no such entity"
	DataStoreFailedtoCreateClient = "Failed to create client"
	DataStoreFailedtoSaveURL      = "Failed to save URL"
	DataStoreFailedtoGetURL       = "Failed to get URL"
	DataStoreFailedtoUpdateURL    = "Failed to update URL"
	DataStoreFailedtoDeleteURL    = "Failed to delete URL"
	DataStoreFailedToCloseClient  = "Failed to close client"
	DataStoreAuthInvalidToken     = "reauthentication required due to invalid token."

	// DataStoreNameKey is the name of the Kind in Datastore for URL entities.
	// Defining it here enables changing the Kind name in one place if needed.
	DataStoreNameKey = "urlz"

	// URL Info Messages
	InfoAttemptingToUpdateURLInDatastore = "Attempting to update URL in Datastore"
	InfoFailedToUpdateURLInDatastore     = "Failed to update URL in Datastore"
	InfoUpdateSuccessful                 = "URL updated successfully in the datastore"
)

// Define error object constants.
const (
	noerrortoparse        = "no error to parse"
	unexpectederrorformat = "unexpected error format"
	invalid_grant         = "invalid_grant"
	http                  = "http"
	ObjCode               = "Code"
	ObjDescription        = "Description"
	ObjDetails            = "Details"
)

// Define operation constants.
const (
	operation_CreateDatastoreClient = "CreateDatastoreClient"
)
