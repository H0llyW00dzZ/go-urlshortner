package logmonitor

// Define header response for different components.
const (
	HeaderResponseError               = "error"
	HeaderResponseInternalServerError = "internal server error"
	HeaderResponseID                  = "id"
	HeaderResponseURL                 = "url"
	HeaderResponseshortened_url       = "shortened_url"
)

// Define header request for different components.
const (
	HeaderRequestOldURL = "old_url"
	HeaderRequestNewURL = "new_url"
)
