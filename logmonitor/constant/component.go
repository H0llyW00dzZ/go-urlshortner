package constant

// Component constants for structured logging.
// This is used to identify the component that is logging the message.
const (
	ComponentNoSQL             = "datastore"
	ComponentCache             = "cache" // Currently unused.
	ComponentProjectIDENV      = "projectid"
	ComponentInternalSecretENV = "customsecretkey"
	ComponentMachineOperation  = "signal" // Currently unused.
	ComponentGopher            = "hostmachine"
)
