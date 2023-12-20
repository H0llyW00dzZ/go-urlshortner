package constant

// Define header response for different components.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	HeaderResponseError                     = "error"
	HeaderResponseInternalServerError       = "internal server error"
	HeaderResponseshortened_url             = "shortened_url"
	HeaderResponseURlUpdated                = "URL updated successfully"
	HeaderResponseInvalidRequestJSONBinding = "Invalid request - JSON binding error"
	HeaderResponseInvalidURLFormat          = "Invalid URL format"
	HeaderResponseInvalidRequestPayload     = "Invalid request payload"
	HeaderResponseInvalidRequest            = "Invalid request"
	HeaderResponseURLDeleted                = "URL deleted successfully"
	HeaderResponseIDandURLNotFound          = "ID and URL not found"
	HeaderResponseForbidden                 = "Forbidden"
	HeaderResponseFailedtoGenerateID        = "Failed to generate ID"
	HeaderResponseFailedtoSaveURL           = "Failed to save URL"
	HeaderResponseStatus                    = "status"
	HeaderResponseRateLimitExceeded         = "Too many requests, please try again later."
)

// Define header request for different components.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	HeaderRequestOldURL = "old_url"
	HeaderRequestNewURL = "new_url"
)

// Define header for different components.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	HeaderID              = "id"
	HeaderURL             = "url"
	HeaderMessage         = "message"
	HeaderSchemeHTTP      = "http"
	HeaderSchemeHTTPS     = "https"
	HeaderXProto          = "X-Forwarded-Proto"
	HeaderXinternalSecret = "X-Internal-Secret"
)

// Define gin context log for different components.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	GinContextErrLog = "errorLogged"
)
