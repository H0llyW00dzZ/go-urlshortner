package datastore

// Define datastore constants.
// This is used to identify the component that is logging the message.
// Easy Maintenance: If the error message changes, it only need to change it in one place.
const (
	DataStoreNosuchentity = "datastore: no such entity"
	// DataStoreNameKey is the name of the Kind in Datastore for URL entities.
	// Defining it here enables changing the Kind name in one place if needed.
	DataStoreNameKey = "urlz"
)
