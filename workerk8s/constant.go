package workerk8s

const (
	ErrorListingPods       = "error listing pods: %w"
	ErrorUpdatingPodLabels = "error updating pod labels: %w"
	ErrorCreatingPod       = "error creating pod: %w"
	ErrorDeletingPod       = "error deleting pod: %w"
	ErrorGettingPod        = "error getting pod: %w"
	ErrorPodNotFound       = "pod not found"
	ErrorRetrievingPods    = "Error retrieving pods: %w"
	PodAndStatus           = "Pod: %s, Status: %s"
	PodAndStatusAndHealth  = "Pod: %s, Status: %s, Health: %s"
	NotHealthyStatus       = "Not Healthy"
	HealthyStatus          = "Healthy"
	WorkerCancelled        = "Worker cancelled: %v"
	HOME                   = "HOME"
	kube                   = ".kube"
	Config                 = "config"
	errconfig              = "cannot load kubeconfig: %w"
	cannotcreatek8s        = "cannot create kubernetes client: %w"
)

const (
	FetchingPods   = "Fetching pods"
	PodsFetched    = "Pods fetched"
	ProcessingPods = "Processing pods"
)

const (
	TaskLabelKey    = "LabelKey"
	TaskCheckHealth = "CheckHealth"
	TaskGetPod      = "GetPod"
	TaskFetchPods   = "FetchPods"
	TaskProcessPod  = "ProcessPod"
)
