package workerk8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Worker interacts with the Kubernetes API and sends results to a channel.
func Worker(clientset *kubernetes.Clientset, namespace string, results chan<- string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		results <- fmt.Sprintf("Error: %v", err)
		return
	}

	for _, pod := range pods.Items {
		// Perform a basic health check on the pod.
		healthStatus := NotHealthyStatus
		if isPodHealthy(&pod) {
			healthStatus = HealthyStatus
		}

		// Send the pod name, status, and health to the results channel.
		results <- fmt.Sprintf(PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
	}
}

// isPodHealthy checks if a pod is considered "healthy" (i.e., if it's running and all containers are ready).
func isPodHealthy(pod *corev1.Pod) bool {
	// Check if the pod is in the 'Running' phase.
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}

	// Check if all containers in the pod are ready.
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}

	return true
}
