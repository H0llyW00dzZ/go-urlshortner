package constant

// Component constants for structured logging.
// This is used to identify the component that is logging the message.
//
// Note: some constants are not used in the code, indicate that for future use.
const (
	ComponentNoSQL             = "datastore"
	ComponentCache             = "cache" // Currently unused.
	ComponentProjectIDENV      = "projectid"
	ComponentInternalSecretENV = "customsecretkey"
	ComponentMachineOperation  = "signal_notify"
	ComponentGopher            = "hostmachine"
)
