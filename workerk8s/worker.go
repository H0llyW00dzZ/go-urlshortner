package workerk8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Worker starts a worker process that retrieves all pods in a given namespace,
// performs health checks on them, and sends the results to a channel.
func Worker(ctx context.Context, clientset *kubernetes.Clientset, namespace string, results chan<- string) {
	// Retrieve a list of pods from the namespace.
	pods, err := getPods(ctx, clientset, namespace)
	if err != nil {
		// If there's an error retrieving pods, send an error message on the results channel.
		results <- fmt.Sprintf("Error retrieving pods: %v", err)
		return
	}

	// Process each pod to determine its health status and send the results on the channel.
	processPods(ctx, pods, results)
}

// getPods fetches the list of all pods within a specific namespace.
func getPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]corev1.Pod, error) {
	// List all pods in the namespace using the provided context.
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		// Return an error if the pod list cannot be retrieved.
		return nil, err
	}
	return podList.Items, nil
}

// processPods iterates over a slice of pods, performs a health check on each,
// and sends a formatted status string to the results channel.
func processPods(ctx context.Context, pods []corev1.Pod, results chan<- string) {
	for _, pod := range pods {
		select {
		case <-ctx.Done():
			// If the context is cancelled, send a cancellation message and exit the function.
			results <- fmt.Sprintf(WorkerCancelled, ctx.Err())
			return
		default:
			// Determine the health status of the pod and send the result.
			healthStatus := NotHealthyStatus
			if isPodHealthy(&pod) {
				healthStatus = HealthyStatus
			}
			results <- fmt.Sprintf(PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
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
