package workerk8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// / Worker performs health checks on all pods in the given namespace and sends the results to a channel.
func Worker(clientset *kubernetes.Clientset, namespace string, results chan<- string) {
	// Retrieve a list of all pods in the specified namespace.
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		results <- fmt.Sprintf("Error: %v", err)
		return
	}

	// Perform health checks on each pod and send the results over the channel.
	for _, pod := range pods.Items {
		healthStatus := NotHealthyStatus
		if isPodHealthy(&pod) {
			healthStatus = HealthyStatus
		}
		results <- fmt.Sprintf(PodAndStatusAndHealth, pod.Name, pod.Status.Phase, healthStatus)
	}
}

// isPodHealthy determines if a pod is considered "healthy" based on its phase and container statuses.
func isPodHealthy(pod *corev1.Pod) bool {
	// A pod is considered healthy if it's running and all of its containers are ready.
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}
	return true
}
