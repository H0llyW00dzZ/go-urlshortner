package workerk8s

import (
	"context"
	"fmt"

	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Worker starts a worker process that retrieves all pods in a given namespace,
// performs health checks on them, and sends the results to a channel.
func Worker(ctx context.Context, clientset *kubernetes.Clientset, namespace string, results chan<- string) {
	fields := createLogFields(TaskCheckHealth, namespace)
	// Retrieve a list of pods from the namespace.
	logInfoWithEmoji(constant.InfoEmoji, "Worker started", fields...)

	pods, err := getPods(ctx, clientset, namespace)
	if err != nil {
		errMsg := fmt.Sprintf("Error retrieving pods: %v", err)
		logErrorWithEmoji(constant.ErrorEmoji, errMsg)
		results <- errMsg
		return
	}

	// Process each pod to determine its health status and send the results on the channel.
	processPods(ctx, pods, results)
	logInfoWithEmoji(constant.ModernGopherEmoji, "Worker finished processing pods", fields...)
}

// getPods fetches the list of all pods within a specific namespace.
func getPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]corev1.Pod, error) {
	// List all pods in the namespace using the provided context.
	fields := createLogFields(TaskFetchPods, namespace)
	logInfoWithEmoji(constant.ModernGopherEmoji, FetchingPods, fields...)

	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logErrorWithEmoji(constant.ModernGopherEmoji, "Failed to list pods", fields...)
		return nil, err
	}

	logInfoWithEmoji(constant.ModernGopherEmoji, PodsFetched, append(fields, zap.Int("count", len(podList.Items)))...)
	return podList.Items, nil
}

// processPods iterates over a slice of pods, performs a health check on each,
// and sends a formatted status string to the results channel.
func processPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			cancelMsg := fmt.Sprintf("Worker cancelled: %v", ctx.Err())
			logInfoWithEmoji(constant.ModernGopherEmoji, cancelMsg)
			results <- cancelMsg
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := NotHealthyStatus
			if isPodHealthy(&pod) {
				healthStatus = HealthyStatus
			}
			statusMsg := fmt.Sprintf(PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
			logInfoWithEmoji(constant.ModernGopherEmoji, PodsFetched, createLogFields(ProcessingPods, pod.Name, statusMsg)...)
			results <- statusMsg
		}
	}
}

// isPodHealthy checks if a given pod is in a running phase and all of its containers are ready.
func isPodHealthy(pod *corev1.Pod) bool {
	// Check if the pod is in the running phase.
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	// Check if all containers within the pod are ready.
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}
	return true
}
