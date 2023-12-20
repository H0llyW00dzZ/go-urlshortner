package workerk8s

const (
	errorlistingspods      = "error listing pods: %w"
	errorupdatingpodlabels = "error updating pod labels: %w"
	errorcreatingpod       = "error creating pod: %w"
	errordeletingpod       = "error deleting pod: %w"
	errorgettingpod        = "error getting pod: %w"
	ErrorPodNotFound       = "pod not found"
	PodAndStatus           = "Pod: %s, Status: %s"
	PodAndStatusAndHealth  = "Pod: %s, Status: %s, Health: %s"
	NotHealthyStatus       = "Not Healthy"
	HealthyStatus          = "Healthy"
)
