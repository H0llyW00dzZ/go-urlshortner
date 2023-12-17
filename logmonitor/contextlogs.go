package logmonitor

// Define context logs for different components.
const (
	URLmissmatchContextLog                  = "URL mismatch"
	FailedToRetriveURLContextLog            = "failed to retrieve URL"
	FailedToUpdateURLContextLog             = "failed to update URL"
	NoURLContextLog                         = "No URL found"
	NoIDContextLog                          = "No ID found"
	NoURLIDContextLog                       = "No URL and ID found"
	NoURLIDDBContextLog                     = "No URL and ID found in DB"
	URLupdateContextLog                     = "URL updated"
	URLdeleteContextLog                     = "URL deleted"
	ServerStartContextLog                   = "Server is starting and Listening on address"
	ServerFailContextLog                    = "Server failed to start"
	SignalContextLog                        = "Signal received"
	ServerForcetoShutdownContextLog         = "Server is forced to shutdown:"
	DataStoreFailContextLog                 = "Datastore failed to connect"
	StartupFailedContextLog                 = "Startup failed"
	StartupFailureContextLog                = "Startup failure"
	FailedtoCloseDatastoreContextLog        = "failed to close datastore client:"
	DatastoreFailedtoCheckHealthContextLog  = "datastore client failed health check:"
	FailedToCreateDatastoreClientContextLog = "failed to create datastore client:"
	FailedToIntializeLoggerContextLog       = "failed to initialize logger:"
	DataStoreProjectIDEnvContextLog         = "DATASTORE_PROJECT_ID environment variable not set"
)
