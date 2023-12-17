package logmonitor

// Define header response for different components.
const (
	HeaderResponseError                     = "error"
	HeaderResponseInternalServerError       = "internal server error"
	HeaderResponseshortened_url             = "shortened_url"
	HeaderResponseURlUpdated                = "URL updated successfully"
	HeaderResponseInvalidRequestJSONBinding = "Invalid request - JSON binding error"
	HeaderResponseInvalidURLFormat          = "Invalid URL format"
	HeaderResponseInvalidRequestPayload     = "Invalid request payload"
	HeaderResponseINvalidRequest            = "Invalid request"
	HeaderResponseURLDeleted                = "URL deleted successfully"
	HeaderResponseIDandURLNotFound          = "ID and URL not found"
	HeaderResponseForbidden                 = "Forbidden"
	HeaderResponseFailedtoGenerateID        = "Failed to generate ID"
	HeaderResponseFailedtoSaveURL           = "Failed to save URL"
)

// Define header request for different components.
const (
	HeaderRequestOldURL = "old_url"
	HeaderRequestNewURL = "new_url"
)

// Define header for different components.
const (
	HeaderID  = "id"
	HeaderURL = "url"
)
