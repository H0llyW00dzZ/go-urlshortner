package constant

// Define header response for different components.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	URLmismatchContextLog                       = "URL mismatch"
	URLShorteneredContextLog                    = "URL shortened and saved"
	FailedToRetriveURLContextLog                = "Failed to retrieve URL"
	FailedToUpdateURLContextLog                 = "Failed to update URL"
	FailedToDeletedURLContextLog                = "Failed to deleted URL"
	FailedToValidateURLContextLog               = "Failed to validate URL"
	FailedToGetURLContextLog                    = "Failed to get URL"
	URLnotfoundContextLog                       = "URL not found"
	IDnotfoundContextLog                        = "ID not found"
	NoURLIDContextLog                           = "No URL and ID found"
	NoURLIDDBContextLog                         = "No URL and ID found in DB"
	URLupdateContextLog                         = "URL updated successfully"
	URLdeleteContextLog                         = "URL deleted"
	URLisNilContextLog                          = "URL is nil after GetURL call"
	URLRetriveContextLog                        = "URL retrieved successfully"
	URLDeletedSuccessfullyContextLog            = "URL deleted successfully"
	ServerStartContextLog                       = "Server is starting and Listening on address"
	ServerFailContextLog                        = "Server failed to start"
	SignalContextLog                            = "Signal received"
	ServerForcetoShutdownContextLog             = "Server is forced to shutdown:"
	DataStoreFailContextLog                     = "Datastore failed to connect"
	StartupFailedContextLog                     = "Startup failed"
	StartupFailureContextLog                    = "Startup failure"
	FailedtoCloseDatastoreContextLog            = "failed to close datastore client:"
	DatastoreFailedtoCheckHealthContextLog      = "datastore client failed health check:"
	FailedToCreateDatastoreClientContextLog     = "failed to create datastore client:"
	FailedToIntializeLoggerContextLog           = "failed to initialize logger:"
	DataStoreProjectIDEnvContextLog             = "DATASTORE_PROJECT_ID environment variable not set"
	InternelSecretEnvContextLog                 = "INTERNAL_SECRET environment variable not set"
	InvalidNewURLFormatContextLog               = "Invalid new URL format"
	MisMatchBetweenPathIDandPayloadIDContextLog = "Mismatch between path ID and payload ID"
	PathIDandPayloadIDDoesnotMatchContextLog    = "Path ID and payload ID does not match"
	InfoAttemptingToRetrieveTheCurrentURL       = "Attempting to retrieve the current URL"
	InfoFailedToRetrieveTheCurrentURL           = "Failed to retrieve the current URL for update"
	InfoOldURLDoesMatchTheCurrentURL            = "Old URL does match the current URL"
)

// Define JSON metadata for different components.
const (
	DescriptionJsonMetaData = "description"
	DetailsURLJsonMetaData  = "details_url"
)

// Define log message for middleware components.
const (
	TimestampFormat = "2006/01/02 - 15:04:05"
	RequestDetails  = "Request Details"
)
